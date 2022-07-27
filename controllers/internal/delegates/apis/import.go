package apis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	netv1 "k8s.io/api/networking/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

func (d *Delegate) importToManagementApi(
	api *gio.ApiDefinition,
	apimCtx *gio.ManagementContext,
	apiJson []byte,
) error {
	apiId := api.Status.ApiID
	apiName := api.Spec.Name

	log := d.log.WithValues("apiId", apiId).WithValues("api.name", apiName, "api.crossId", apiId)

	if apimCtx == nil {
		log.Info("No management context associated to the API, skipping import to Management API")
		return nil
	}

	findApiResp, findApiErr := d.findByCrossId(apimCtx, apiId)

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

	err := json.Unmarshal(body, &result)
	if err != nil {
		log.Error(err, "Unable to marshal API definition")
		return err
	}

	if len(result) == 0 {
		log.Info("No match found for API, switching to creation mode", "crossId", apiId)
		importHttpMethod = http.MethodPost
	}

	importResp, importErr := d.doImport(importHttpMethod, apimCtx, apiJson)

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

func (d *Delegate) doImport(
	importHttpMethod string,
	apimCtx *gio.ManagementContext,
	apiJson []byte,
) (*http.Response, error) {
	url := apimCtx.Spec.BuildUrl("/apis/import?definitionVersion=2.0.0")
	req, err := http.NewRequestWithContext(d.ctx, importHttpMethod, url, bytes.NewBuffer(apiJson))

	if err != nil {
		return nil, fmt.Errorf("unable to import the api into the Management API")
	}

	req.Header.Add("Content-Type", "application/json")
	apimCtx.Spec.Authenticate(req)
	resp, err := d.http.Do(req)
	return resp, err
}

// This function is applied to all ingresses which are using the ApiDefinition template
// As per Kubernetes Finalizers (https://kubernetes.io/docs/concepts/overview/working-with-objects/finalizers/)
// First return value defines if we should requeue or not.
func (d *Delegate) ImportApiDefinitionTemplate(
	apiDefinition *gio.ApiDefinition,
	namespace string,
) (bool, error) {
	// We are first looking if the template is in deletion phase, the Kubernetes API marks the object for
	// deletion by populating .metadata.deletionTimestamp
	if !apiDefinition.DeletionTimestamp.IsZero() {
		if !util.ContainsFinalizer(apiDefinition, keys.ApiDefinitionTemplateFinalizer) {
			return false, nil
		}

		ingressList := netv1.IngressList{}

		// Retrieves the ingresses from the namespace
		err := d.cli.List(d.ctx, &ingressList, client.InNamespace(namespace))
		if err != nil && !kerrors.IsNotFound(err) {
			return false, err
		}

		var ingresses []string

		for _, ingress := range ingressList.Items {
			if ingress.GetAnnotations()[keys.IngressTemplateAnnotation] == apiDefinition.Name {
				ingresses = append(ingresses, ingress.GetName())
			}
		}

		// There are existing ingresses wich to the ApiDefinition template, re-schedule deletion
		if len(ingresses) > 0 {
			err = fmt.Errorf("can not delete %s %v depends on it", apiDefinition.Name, ingresses)
			return true, err
		}

		util.RemoveFinalizer(apiDefinition, keys.ApiDefinitionTemplateFinalizer)

		return false, d.cli.Update(d.ctx, apiDefinition)
	}

	// Adding or updating a new ApiDefinition template
	// If it is a creation, adding the Finalizers to keep track of the deletion
	if !util.ContainsFinalizer(apiDefinition, keys.ApiDefinitionTemplateFinalizer) {
		util.AddFinalizer(apiDefinition, keys.ApiDefinitionTemplateFinalizer)

		return false, d.cli.Update(d.ctx, apiDefinition)
	}

	ingressList := netv1.IngressList{}

	// Listing ingresses from the same namespace
	err := d.cli.List(d.ctx, &ingressList, client.InNamespace(namespace))

	if err != nil {
		return false, client.IgnoreNotFound(err)
	}

	return false, nil
}
