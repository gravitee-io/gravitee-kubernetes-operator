package context

import (
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

func (d *Delegate) Get(
	api *gio.ApiDefinition,
) (*gio.ManagementContext, error) {
	contextRef := api.Spec.Context

	if contextRef == nil {
		group := gio.GroupVersion.Group
		resource := keys.CrdManagementContextResource
		ref := schema.GroupResource{Group: group, Resource: resource}
		return nil, kerrors.NewNotFound(ref, "")
	}

	apimContext := new(gio.ManagementContext)
	ns := types.NamespacedName{Name: contextRef.Name, Namespace: contextRef.Namespace}

	d.log.Info("Looking for context from", "namespace", contextRef.Namespace, "name", contextRef.Name)

	err := d.cli.Get(d.ctx, ns, apimContext)

	if err != nil {
		return nil, err
	}

	return apimContext, nil
}
