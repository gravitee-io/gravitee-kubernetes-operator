package apis

import (
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

func (d *Delegate) delete(
	api *gio.ApiDefinition,
) error {
	d.log.Info("Deleting API Definition")
	err := d.deleteConfigMap(api)

	return err
}
