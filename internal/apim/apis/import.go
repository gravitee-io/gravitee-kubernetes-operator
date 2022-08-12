package apis

import (
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	netv1 "k8s.io/api/networking/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

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
		err := d.k8sClient.List(d.ctx, &ingressList, client.InNamespace(namespace))
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

		return false, d.k8sClient.Update(d.ctx, apiDefinition)
	}

	// Adding or updating a new ApiDefinition template
	// If it is a creation, adding the Finalizers to keep track of the deletion
	if !util.ContainsFinalizer(apiDefinition, keys.ApiDefinitionTemplateFinalizer) {
		util.AddFinalizer(apiDefinition, keys.ApiDefinitionTemplateFinalizer)

		return false, d.k8sClient.Update(d.ctx, apiDefinition)
	}

	ingressList := netv1.IngressList{}

	// Listing ingresses from the same namespace
	err := d.k8sClient.List(d.ctx, &ingressList, client.InNamespace(namespace))

	if err != nil {
		return false, client.IgnoreNotFound(err)
	}

	return false, nil
}
