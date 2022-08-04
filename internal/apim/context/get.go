package context

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
)

func (d *Delegate) Get(
	contextRef *model.ContextRef,
) (*gio.ManagementContext, error) {
	apimContext := new(gio.ManagementContext)
	ns := types.NamespacedName{Name: contextRef.Name, Namespace: contextRef.Namespace}

	d.log.Info("Looking for context from", "namespace", contextRef.Namespace, "name", contextRef.Name)

	err := d.cli.Get(d.ctx, ns, apimContext)

	if err != nil {
		return nil, err
	}

	return apimContext, nil
}
