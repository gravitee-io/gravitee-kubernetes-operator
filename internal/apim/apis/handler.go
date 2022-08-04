package apis

import gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"

func (d *Delegate) Handle(api *gio.ApiDefinition) error {
	log := d.log.WithValues("name", api.Name)
	if api.IsBeingDeleted() {
		log.Info("Updating API definition")
		return d.delete(api)
	}

	if api.IsBeingCreated() {
		log.Info("Creating API definition")
		return d.create(api)
	}

	if api.IsBeingUpdated() {
		log.Info("Updating API definition")
		return d.update(api)
	}
	return nil
}
