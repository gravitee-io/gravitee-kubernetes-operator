package delegate

import (
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	managementapierror "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi/clienterror"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"errors"
)

func (d *Delegate) Delete(
	apiDefinition *gio.ApiDefinition,
) error {
	// Do nothing if finalizer is already removed
	if !util.ContainsFinalizer(apiDefinition, keys.ApiDefinitionDeletionFinalizer) {
		return nil
	}

	var apiNotFoundError *managementapierror.ApiNotFoundError
	if d.IsConnectedToManagementApi() {
		d.log.Info("Delete API definition into Management API")
		err := d.apimClient.DeleteApi(apiDefinition.Status.ID)
		if errors.As(err, &apiNotFoundError) {
			d.log.Info("The API has already been deleted", "id", apiDefinition.Status.ID)
		}
		if err != nil && !errors.As(err, &apiNotFoundError) {
			d.log.Error(err, "Unable to delete API definition into Management API")
			return err
		}
	}

	// Remove finalizer when API definition is fully deleted
	util.RemoveFinalizer(apiDefinition, keys.ApiDefinitionDeletionFinalizer)

	return d.k8sClient.Update(d.ctx, apiDefinition)
}
