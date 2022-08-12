package apis

import (
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

func (d *Delegate) Handle(apiDefinition *gio.ApiDefinition) error {
	log := d.log.WithValues("name", apiDefinition.Name)

	if !apiDefinition.HasFinalizer() {
		log.Info("Add Finalizer to API definition")
		return d.finalizer(apiDefinition)
	}

	if apiDefinition.IsBeingDeleted() {
		log.Info("Deleting API definition")
		return d.delete(apiDefinition)
	}

	if apiDefinition.IsBeingCreated() {
		log.Info("Creating API definition")
		return d.create(apiDefinition)
	}

	if apiDefinition.IsBeingUpdated() {
		log.Info("Updating API definition")
		return d.update(apiDefinition)
	}
	return nil
}
