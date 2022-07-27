package context

import (
	"k8s.io/apimachinery/pkg/types"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

func (d *Delegate) Get(
	api *gio.ApiDefinition,
) (*gio.ManagementContext, error) {
	contextRef := api.Spec.Context

	apimContext := new(gio.ManagementContext)
	ns := types.NamespacedName{Name: contextRef.Name, Namespace: contextRef.Namespace}

	d.log.Info("Looking for context from", "namespace", contextRef.Namespace, "name", contextRef.Name)

	err := d.cli.Get(d.ctx, ns, apimContext)

	if err != nil {
		return nil, err
	}

	return apimContext, nil
}
