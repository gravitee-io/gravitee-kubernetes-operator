package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/internal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	v1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const requestTimeout = 5

func (r *ApiDefinitionReconciler) createApiDefinition(
	ctx context.Context,
	apiDefinition *gio.ApiDefinition,
) error {
	log := log.FromContext(ctx)

	// Plan is not required from the CRD, but is expected by the Gateway, so we must create at least one
	createDefaultPlan(ctx, apiDefinition)

	apimCtx, err := internal.GetApimContext(ctx, r.Client, apiDefinition)
	if client.IgnoreNotFound(err) != nil {
		return fmt.Errorf("An error has occured while trying to find a management context %s", err)
	}

	// Ensure that IDs have been generated
	internal.GenerateIds(apimCtx, apiDefinition)
	internal.SetDeployedAt(apiDefinition)

	// Import the API Definition:
	// * Create the ConfigMap for the Gateway
	// * Call the Management API if a ManagementContext is defined
	err = r.importApiDefinition(ctx, apiDefinition)
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
	apiDefinition *gio.ApiDefinition,
) error {
	// Plan is not required from the CRD, but is expected by the Gateway, so we must create at least one
	createDefaultPlan(ctx, apiDefinition)

	apimCtx, err := internal.GetApimContext(ctx, r.Client, apiDefinition)
	if client.IgnoreNotFound(err) != nil {
		return fmt.Errorf("An error has occured while trying to find a management context %s", err)
	}

	// Ensure that IDs have been generated
	internal.GenerateIds(apimCtx, apiDefinition)
	internal.SetDeployedAt(apiDefinition)

	return r.importApiDefinition(ctx, apiDefinition)
}

func (r *ApiDefinitionReconciler) importApiDefinition(
	ctx context.Context,
	apiDefinition *gio.ApiDefinition,
) error {
	log := log.FromContext(ctx)

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

	updated, err := r.updateConfigMap(ctx, apiDefinition, apiJson)

	if err != nil {
		log.Error(err, "Unable to create or update ConfigMap for API '%s' (%s). %s",
			apiDefinition.Name, apiDefinition.Spec.Id)
		return err
	}

	if updated {
		err = r.importToManagementApi(ctx, apiDefinition, apiJson)
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
func (r *ApiDefinitionReconciler) importApiDefinitionTemplate(
	ctx context.Context,
	apiDefinition *gio.ApiDefinition,
	namespace string,
) (ctrl.Result, error) {
	// We are first looking if the template is in deletion phase, the Kubernetes API marks the object for
	// deletion by populating .metadata.deletionTimestamp
	if !apiDefinition.DeletionTimestamp.IsZero() {
		if !util.ContainsFinalizer(apiDefinition, keys.ApiDefinitionTemplateFinalizer) {
			return ctrl.Result{}, nil
		}

		ingressList := netv1.IngressList{}

		// Retrieves the ingresses from the namespace
		err := r.List(ctx, &ingressList, client.InNamespace(namespace))
		if err != nil && !kerrors.IsNotFound(err) {
			return ctrl.Result{}, err
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

	return ctrl.Result{}, nil
}

// Add a default keyless plan to the api definition if no plan is defined.
func createDefaultPlan(ctx context.Context, api *gio.ApiDefinition) {
	log := log.FromContext(ctx)

	plans := api.Spec.Plans

	if len(plans) == 0 {
		log.Info("Define default plan for API")
		api.Spec.Plans = []*model.Plan{
			{
				Name:     "GKO DEFAULT",
				Security: "KEY_LESS",
				Status:   "PUBLISHED",
			},
		}
	}
}

func (r *ApiDefinitionReconciler) updateConfigMap(
	ctx context.Context,
	apiDefinition *gio.ApiDefinition,
	apiJson []byte,
) (bool, error) {
	log := log.FromContext(ctx)

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
	}

	apimCtx, err := internal.GetApimContext(ctx, r.Client, apiDefinition)
	if client.IgnoreNotFound(err) != nil {
		log.Error(err, "An error has occured trying to find an APIM context")
	}

	if apimCtx != nil {
		cm.Data["organizationId"] = apimCtx.Spec.OrgId
		cm.Data["environmentId"] = apimCtx.Spec.EnvId
	}

	currentapiDefinition := &v1.ConfigMap{}
	err = r.Get(ctx, types.NamespacedName{Name: cm.Name, Namespace: cm.Namespace}, currentapiDefinition)

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
	apiDefinition *gio.ApiDefinition,
	apiJson []byte,
) error {
	apiId := apiDefinition.Status.ApiID
	apiName := apiDefinition.Spec.Name

	log := log.FromContext(ctx).WithValues("apiId", apiId).WithValues("API.Name", apiName, "API.ID", apiId)

	if apiDefinition.Spec.Context == nil {
		log.Info("No management context associated to the API, skipping import to Management API")
		return nil
	}

	apimCtx, err := internal.GetApimContext(ctx, r.Client, apiDefinition)

	if err != nil {
		return err
	}

	// Call management side to push api also.
	client := http.Client{Timeout: requestTimeout * time.Second}

	findApiResp, findApiErr := r.findApisByCrossId(ctx, apimCtx, apiId, client)

	if findApiErr != nil {
		return findApiErr
	}

	if findApiResp.Body != nil {
		defer findApiResp.Body.Close()
	}

	if findApiResp.StatusCode != http.StatusOK {
		// TODO parse response body as a map and log
		return fmt.Errorf("an error as occured trying to find API %s, HTTP Status: %d ", apiId, findApiResp.StatusCode)
	}

	body, readErr := ioutil.ReadAll(findApiResp.Body)
	if readErr != nil {
		log.Error(readErr, "Unable to read apis response body")
	}

	// If the API does not exist (ie. 404) it should be a POST
	importHttpMethod := http.MethodPut
	var result []interface{}
	err = json.Unmarshal(body, &result)

	if err != nil {
		log.Error(err, "Unable to marshal API definition")
		return err
	}

	if len(result) == 0 {
		log.Info("No match found for API, switching to creation mode")
		importHttpMethod = http.MethodPost
	}

	importResp, importErr := r.importApi(ctx, importHttpMethod, apimCtx, apiJson, client)

	if importErr != nil {
		log.Error(importErr, "Unable to import the api into the Management API")
		return importErr
	}

	if importResp.Body != nil {
		defer importResp.Body.Close()
	}

	if importResp.StatusCode < 200 || importResp.StatusCode > 299 {
		log.Error(nil, "Unable to import the api into the Management API")
		return fmt.Errorf("management has returned a %d code", importResp.StatusCode)
	}

	log.Info("Api has been pushed to the Management API")
	return nil
}

func (r *ApiDefinitionReconciler) importApi(
	ctx context.Context,
	importHttpMethod string,
	apimCtx *gio.ManagementContext,
	apiJson []byte,
	client http.Client,
) (*http.Response, error) {
	url := internal.BuildApimUrl(apimCtx, "/apis/import?definitionVersion=2.0.0")
	req, err := http.NewRequestWithContext(ctx, importHttpMethod, url, bytes.NewBuffer(apiJson))

	if err != nil {
		return nil, fmt.Errorf("unable to import the api into the Management API")
	}

	req.Header.Add("Content-Type", "application/json")
	internal.SetApimAuth(apimCtx, req)
	resp, err := client.Do(req)
	return resp, err
}

func (r *ApiDefinitionReconciler) findApisByCrossId(
	ctx context.Context,
	apimCtx *gio.ManagementContext,
	apiId string,
	client http.Client,
) (*http.Response, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		internal.BuildApimUrl(apimCtx, "/apis?crossId="+apiId),
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("an error as occured while trying to create new findApisByCrossId request")
	}

	internal.SetApimAuth(apimCtx, req)
	resp, err := client.Do(req)
	return resp, err
}

func (r *ApiDefinitionReconciler) deleteApiDefinition(
	ctx context.Context,
	apiDefinition gio.ApiDefinition,
) error {
	r.Log.Info("Deleting API Definition")
	err := r.deleteApiDefinitionConfigMap(ctx, apiDefinition)

	return err
}

func (r *ApiDefinitionReconciler) deleteApiDefinitionConfigMap(
	ctx context.Context,
	apiDefinition gio.ApiDefinition,
) error {
	configMap := &v1.ConfigMap{}

	r.Log.Info("Deleting ConfigMap associated to API")
	err := r.Get(ctx, types.NamespacedName{Name: apiDefinition.Name, Namespace: apiDefinition.Namespace}, configMap)

	if err != nil {
		err = r.Delete(ctx, configMap)
	}

	return err
}
