package apis

import gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"

func (d *Delegate) Handle(api *gio.ApiDefinition) error {
	if api.IsBeingDeleted() {
		return d.Delete(api)
	}

	if api.IsBeingCreated() {
		return d.Create(api)
	}

	if api.IsBeingUpdated() {
		return d.Update(api)
	}
	return nil
}
