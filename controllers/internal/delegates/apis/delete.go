package apis

import (
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (d *Delegate) Delete(
	api *gio.ApiDefinition,
) error {
	d.log.Info("Deleting API Definition")
	err := d.deleteConfigMap(api)

	return err
}

func (d *Delegate) deleteConfigMap(
	api *gio.ApiDefinition,
) error {
	configMap := &v1.ConfigMap{}

	d.log.Info("Deleting ConfigMap associated to API")
	err := d.cli.Get(d.ctx, types.NamespacedName{Name: api.Name, Namespace: api.Namespace}, configMap)

	if err != nil {
		err = d.cli.Delete(d.ctx, configMap)
	}

	return err
}
