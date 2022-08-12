package apis

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (d *Delegate) saveConfigMap(
	api *gio.ApiDefinition,
	apiJson []byte,
) error {
	if api.Spec.State == model.StateStopped {
		return nil
	}

	// Create configmap with some specific metadata that will be used to check changes across 'Update' events.
	cm := &v1.ConfigMap{}

	// Set OwnerReference on configmap to be able to delete it when API is deleted.
	// 📝 ConfigMap should be in same namespace as ApiDefinition.
	newOwnerReferences := []metav1.OwnerReference{
		{
			Kind:       api.Kind,
			Name:       api.Name,
			APIVersion: api.APIVersion,
			UID:        api.UID,
		},
	}
	cm.SetOwnerReferences(newOwnerReferences)

	cm.Namespace = api.Namespace
	cm.Name = api.Name
	cm.CreationTimestamp = metav1.Now()
	cm.Labels = map[string]string{
		"managed-by": keys.CrdGroup,
		"gio-type":   keys.CrdApiDefinitionResource + "." + keys.CrdGroup,
	}

	cm.Data = map[string]string{
		"definition":        string(apiJson),
		"definitionVersion": api.ResourceVersion,
	}

	if d.managementContext != nil {
		cm.Data["organizationId"] = d.managementContext.Spec.OrgId
		cm.Data["environmentId"] = d.managementContext.Spec.EnvId
	}

	currentApiDefinition := &v1.ConfigMap{}
	err := d.k8sClient.Get(d.ctx, types.NamespacedName{Name: cm.Name, Namespace: cm.Namespace}, currentApiDefinition)

	if err == nil {
		if currentApiDefinition.Data["definitionVersion"] != api.ResourceVersion {
			d.log.Info("Updating ConfigMap", "id", api.Spec.Id)
			// Only update the config map if resource version has changed (means api definition has changed).
			err = d.k8sClient.Update(d.ctx, cm)
		} else {
			d.log.Info("No change detected on api. Skipped.", "id", api.Spec.Id)
			return nil
		}
	} else {
		d.log.Info("Creating config map for api.", "id", api.Spec.Id, "name", api.Name)
		err = d.k8sClient.Create(d.ctx, cm)
	}
	return err
}

func (d *Delegate) deleteConfigMap(apiNamespace string, apiName string) error {
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      apiName,
			Namespace: apiNamespace,
		},
	}

	d.log.Info("Deleting ConfigMap associated to API if exist")
	err := d.k8sClient.Delete(d.ctx, configMap)

	if errors.IsNotFound(err) {
		return nil
	}

	return err
}
