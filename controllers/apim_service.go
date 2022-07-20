package controllers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	graviteeiov1alpha1 "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	uuid "github.com/satori/go.uuid"
	v1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ApiDefinitionReconciler) createApiDefinition(
	ctx context.Context,
	apiDefinition *graviteeiov1alpha1.ApiDefinition,
	orgId string,
	envId string,
) error {
	log := logr.FromContextOrDiscard(ctx)

	// Plan is not required from the CRD, but is expected by the Gateway, so we must create at least one
	createDefaultPlan(apiDefinition, log)

	// Ensure that IDs have been generated
	generateIds(apiDefinition)

	// Import the API Definition:
	// * Create the ConfigMap for the Gateway
	// * Call the Management API if a ManagementContext is defined
	err := r.importApiDefinition(ctx, apiDefinition, orgId, envId)
	if err != nil {
		log.Error(err, "Unexpected error while importing the API Definition")
		return err
	}

	apiDefinition.Status.ApiID = apiDefinition.Spec.CrossId

	err = r.Status().Update(ctx, apiDefinition)
	if err != nil {
		log.Error(err, "Unexpected error while updating status")
		return err
	}

	return nil
}

func (r *ApiDefinitionReconciler) updateApiDefinition(
	ctx context.Context,
	apiDefinition *graviteeiov1alpha1.ApiDefinition,
	orgId string,
	envId string,
) error {
	log := logr.FromContextOrDiscard(ctx)

	// Plan is not required from the CRD, but is expected by the Gateway, so we must create at least one
	createDefaultPlan(apiDefinition, log)

	// Ensure that IDs have been generated
	generateIds(apiDefinition)

	return r.importApiDefinition(ctx, apiDefinition, orgId, envId)
}

func (r *ApiDefinitionReconciler) importApiDefinition(
	ctx context.Context,
	apiDefinition *graviteeiov1alpha1.ApiDefinition,
	orgId string,
	envId string,
) error {
	log := logr.FromContextOrDiscard(ctx)

	// Define the API definition context
	apiDefinition.Spec.DefinitionContext = &model.DefinitionContext{
		Origin: "kubernetes",
		// Could also be api_definition_only in a near future
		Mode: "fully_managed",
	}

	// Marshal the APIDefinition to JSON
	apiJson, err := json.Marshal(apiDefinition.Spec)

	if err != nil {
		log.Error(err, "Unable to generate json api definition for api '%s' (%s). %s",
			apiDefinition.Name, apiDefinition.Spec.Id)
		return err
	}

	updated, err := r.updateConfigMap(ctx, apiDefinition, orgId, envId, apiJson, log)

	if err != nil {
		log.Error(err, "Unable to create or update ConfigMap for API '%s' (%s). %s",
			apiDefinition.Name, apiDefinition.Spec.Id)
		return err
	}

	if updated {
		err = r.importToManagementApi(ctx, apiDefinition, orgId, envId, apiJson, log)
		if err != nil {
			log.Error(err, "Unable to import API to the Management API '%s' (%s). %s",
				apiDefinition.Name, apiDefinition.Spec.Id)
			return err
		}
	}

	return nil
}

// This function is applied to all ingresses which are using the ApiDefinition template
// As per Kubernetes Finalizers (https://kubernetes.io/docs/concepts/overview/working-with-objects/finalizers/)
func (r *ApiDefinitionReconciler) importApiDefinitionTemplate(ctx context.Context, apiDefinition *graviteeiov1alpha1.ApiDefinition, namespace string) (ctrl.Result, error) {

	// We are first looking if the template is in deletion phase, the Kubernetes API marks the object for
	// deletion by populating .metadata.deletionTimestamp
	if !apiDefinition.DeletionTimestamp.IsZero() {
		if util.ContainsFinalizer(apiDefinition, keys.ApiDefinitionTemplateFinalizer) {
			ingressList := netv1.IngressList{}

			// Retrieves the ingresses from the namespace
			err := r.List(ctx, &ingressList, client.InNamespace(namespace))
			if err != nil {
				if !kerrors.IsNotFound(err) {
					return ctrl.Result{}, err
				}
			}

			var ingresses []string

			for _, ingress := range ingressList.Items {
				if ingress.GetAnnotations()[keys.IngressTemplateAnnotation] == apiDefinition.Name {
					ingresses = append(ingresses, ingress.GetName())
				}
			}

			// There are existing ingresses wich to the ApiDefinition template, re-schedule deletion
			if len(ingresses) > 0 {
				return ctrl.Result{RequeueAfter: time.Second * RequeueAfterTime},
					fmt.Errorf("can not delete %s %v depends on it", apiDefinition.Name, ingresses)
			}

			util.RemoveFinalizer(apiDefinition, keys.ApiDefinitionTemplateFinalizer)

			return ctrl.Result{}, r.Update(ctx, apiDefinition)
		}

		return ctrl.Result{}, nil
	}

	// Adding or updating a new ApiDefinition template
	// If it is a creation, adding the Finalizers to keep track of the deletion
	if !util.ContainsFinalizer(apiDefinition, keys.ApiDefinitionTemplateFinalizer) {
		util.AddFinalizer(apiDefinition, keys.ApiDefinitionTemplateFinalizer)

		return ctrl.Result{}, r.Update(ctx, apiDefinition)
	}

	ingressList := netv1.IngressList{}

	// Listing ingresses from the same namespace
	err := r.List(ctx, &ingressList, client.InNamespace(namespace))

	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Then look if one of the ingress is refering to the ApiDefinition template
	for _, ingress := range ingressList.Items {
		if ingress.GetAnnotations()[keys.IngressTemplateAnnotation] == apiDefinition.Name {
			// Notify the gravitee ingress controller of an ingress to update (see ingress_controller.go)
		}
	}

	return ctrl.Result{}, nil
}

// Add a default keyless plan to the api definition if no plan is defined.
func createDefaultPlan(apiDefinition *graviteeiov1alpha1.ApiDefinition, log logr.Logger) {
	plans := apiDefinition.Spec.Plans

	if len(plans) == 0 {
		log.Info("Define default plan for API")
		apiDefinition.Spec.Plans = []*model.Plan{
			{
				Name:     "Free",
				Security: "KEY_LESS",
				Status:   "PUBLISHED",
			},
		}
	}
}

func (r *ApiDefinitionReconciler) updateConfigMap(
	ctx context.Context,
	apiDefinition *graviteeiov1alpha1.ApiDefinition,
	orgId string,
	envId string,
	apiJson []byte,
	log logr.Logger,
) (bool, error) {
	// Create configmap with some specific metadata that will be used to check changes across 'Update' events.
	cm := &v1.ConfigMap{}

	cm.Namespace = apiDefinition.Namespace
	cm.Name = apiDefinition.Name
	cm.CreationTimestamp = metav1.Now()
	cm.Labels = map[string]string{
		"managed-by": keys.CrdGroup,
		"gio-type":   keys.CrdApiDefinitionResource + "." + keys.CrdGroup,
	}

	cm.Data = map[string]string{
		"definition":        string(apiJson),
		"definitionVersion": apiDefinition.ResourceVersion,
		"organizationId":    orgId,
		"environmentId":     envId,
	}

	currentapiDefinition := &v1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{Name: cm.Name, Namespace: cm.Namespace}, currentapiDefinition)

	if err == nil {
		if currentapiDefinition.Data["definitionVersion"] != apiDefinition.ResourceVersion {
			log.Info("Updating ConfigMap", "id", apiDefinition.Spec.Id)
			// Only update the confimap if resource version has changed (means api definition has changed).
			err = r.Update(ctx, cm)
		} else {
			log.Info("No change detected on api. Skipped.", "id", apiDefinition.Spec.Id)
			return false, nil
		}
	} else {
		log.Info("Creating configmap for api.", "id", apiDefinition.Spec.Id, "name", apiDefinition.Name)
		err = r.Create(ctx, cm)
	}
	return true, err
}

func (r *ApiDefinitionReconciler) importToManagementApi(
	ctx context.Context,
	apiDefinition *graviteeiov1alpha1.ApiDefinition,
	orgId string,
	envId string,
	apiJson []byte,
	log logr.Logger,
) error {
	const timeout = 5
	apiId := apiDefinition.Status.ApiID
	apiName := apiDefinition.Spec.Name

	if apiDefinition.Spec.Context != nil {
		mgmtContextInst, err := getManagementContext(ctx, r.Client, log, apiDefinition)

		if err != nil {
			return err
		}

		// Call management side to push api also.
		client := http.Client{Timeout: timeout * time.Second}

		// Do reconciliation with the Management API
		request, err := http.NewRequest(
			http.MethodGet,
			mgmtContextInst.Spec.BaseUrl+"/management/organizations/"+orgId+"/environments/"+envId+"/apis?crossId="+apiId,
			nil,
		)
		setRequestAuth(request, mgmtContextInst)
		response, err := client.Do(request)

		// If the API does not exist (ie. 404) it should be a POST
		importHttpMethod := http.MethodPut

		if err != nil {
			log.Error(err, "Error")
		}

		if response.Body != nil {
			defer response.Body.Close()
		}

		if response.StatusCode != http.StatusOK {
			// TODO parse response body as a map and log
			return fmt.Errorf("an error as occured trying to find API %s, HTTP Status: %d ", apiId, response.StatusCode)
		}

		body, readErr := ioutil.ReadAll(response.Body)
		if readErr != nil {
			log.Error(readErr, "Error")
		}

		var result []interface{}
		err = json.Unmarshal(body, &result)

		if err != nil {
			log.Error(err, "Unable to marshal API definition")
			return err
		} else if len(result) == 0 {
			log.Info("No match found for API, switching to creation mode", "apiId", apiId)
			importHttpMethod = http.MethodPost
		}

		request, err = http.NewRequestWithContext(
			ctx, importHttpMethod, mgmtContextInst.Spec.BaseUrl+"/management/organizations/"+orgId+"/environments/"+
				envId+"/apis/import?definitionVersion=2.0.0", bytes.NewBuffer(apiJson),
		)

		if err != nil {
			log.Error(err, "Unable to import the api into the Management API")
			return err
		}

		request.Header.Add("Content-Type", "application/json")
		setRequestAuth(request, mgmtContextInst)
		response, err = client.Do(request)

		if err != nil {
			log.Error(err, "Unable to import the api into the Management API", apiName, apiId, err)
			return err
		}

		if response.StatusCode < 200 || response.StatusCode > 299 {
			log.Error(nil, "Unable to import the api into the Management API", apiName, apiId)
			return fmt.Errorf("management has returned a %d code", response.StatusCode)
		}

		if response.Body != nil {
			defer response.Body.Close()
		}
		log.Info("Api has been pushed to the Management API", apiName, apiId)
	} else {
		log.Info("No management context associated to the API, skipping import to Management API")
	}
	return nil
}

func (r *ApiDefinitionReconciler) deleteApiDefinition(
	ctx context.Context,
	apiDefinition graviteeiov1alpha1.ApiDefinition,
) error {
	r.Log.Info("Deleting API Definition")
	err := r.deleteApiDefinitionConfigMap(ctx, apiDefinition)

	return err
}

func (r *ApiDefinitionReconciler) deleteApiDefinitionConfigMap(
	ctx context.Context,
	apiDefinition graviteeiov1alpha1.ApiDefinition,
) error {
	configMap := &v1.ConfigMap{}

	r.Log.Info("Deleting ConfigMap associated to API")
	err := r.Get(ctx, types.NamespacedName{Name: apiDefinition.Name, Namespace: apiDefinition.Namespace}, configMap)

	if err != nil {
		err = r.Delete(ctx, configMap)
	}

	return err
}

func setRequestAuth(request *http.Request, managementContext graviteeiov1alpha1.ManagementContext) {
	if managementContext.Spec.Auth != nil {
		bearerToken := managementContext.Spec.Auth.BearerToken
		if bearerToken != "" {
			request.Header.Add("Authorization", "Bearer "+bearerToken)
		} else if managementContext.Spec.Auth.Credentials != nil {
			username := managementContext.Spec.Auth.Credentials.Username
			password := managementContext.Spec.Auth.Credentials.Password
			if username != "" {
				request.SetBasicAuth(username, password)
			}
		}
	}
}

// This function is used to generate all the IDs needed for communicating with the Management API
// It doesn't override IDs if these one have been defined.
func generateIds(apiDefinition *graviteeiov1alpha1.ApiDefinition) {
	// If a CrossID is defined at the API level, reuse it.
	// If not, just generate a new CrossID
	if apiDefinition.Spec.CrossId == "" {
		// The ID of the API will be based on the API Name and Namespace to ensure consistency
		apiDefinition.Spec.CrossId =
			toUUID(types.NamespacedName{Namespace: apiDefinition.Namespace, Name: apiDefinition.Name}.String())
	}

	plans := apiDefinition.Spec.Plans

	for i, plan := range plans {
		if plan.CrossId == "" {
			plan.CrossId = toUUID(apiDefinition.Spec.CrossId + fmt.Sprint(i))
		}
		plan.Status = "PUBLISHED"
	}

	//TODO: manage metadata
}

func toUUID(decoded string) string {
	encoded := base64.RawStdEncoding.EncodeToString([]byte(decoded))
	return uuid.NewV3(uuid.NamespaceURL, encoded).String()
}
