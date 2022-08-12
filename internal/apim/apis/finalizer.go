package apis

import (
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (d *Delegate) finalizer(
	apiDefinition *gio.ApiDefinition,
) error {
	util.AddFinalizer(apiDefinition, keys.ApiDefinitionDeletionFinalizer)

	return d.k8sClient.Update(d.ctx, apiDefinition)
}
