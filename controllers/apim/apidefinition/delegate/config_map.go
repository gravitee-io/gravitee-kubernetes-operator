package delegate

import (
	"encoding/json"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (d *Delegate) saveConfigMap(
	apiDefinition *gio.ApiDefinition,
) error {
	if apiDefinition.Spec.State == model.StateStopped {
		return nil
	}

	apiJson, err := json.Marshal(apiDefinition.Spec)
	if err != nil {
		d.log.Error(err, "Unable to marshall API definition as JSON")
		return err
	}

	// Create configmap with some specific metadata that will be used to check changes across 'Update' events.
	cm := &v1.ConfigMap{}

	// Set OwnerReference on configmap to be able to delete it when API is deleted.
	// 📝 ConfigMap should be in same namespace as ApiDefinition.
	newOwnerReferences := []metav1.OwnerReference{
		{
			Kind:       apiDefinition.Kind,
			Name:       apiDefinition.Name,
			APIVersion: apiDefinition.APIVersion,
			UID:        apiDefinition.UID,
		},
	}
	cm.SetOwnerReferences(newOwnerReferences)

	cm.Namespace = apiDefinition.Namespace
	cm.Name = apiDefinition.Name
	cm.CreationTimestamp = metav1.Now()
	cm.Labels = map[string]string{
		"managed-by": keys.CrdGroup,
		"gio-type":   keys.CrdApiDefinitionResource + "." + keys.CrdGroup,
	}

	cm.Data = map[string]string{
		"definition":        string(apiJson),
		"definitionVersion": apiDefinition.ResourceVersion,
	}

	if d.managementContext != nil {
		cm.Data["organizationId"] = d.managementContext.Spec.OrgId
		cm.Data["environmentId"] = d.managementContext.Spec.EnvId
	}

	currentApiDefinition := &v1.ConfigMap{}
	err = d.k8sClient.Get(d.ctx, types.NamespacedName{Name: cm.Name, Namespace: cm.Namespace}, currentApiDefinition)

	if err == nil {
		if currentApiDefinition.Data["definitionVersion"] != apiDefinition.ResourceVersion {
			d.log.Info("Updating ConfigMap", "id", apiDefinition.Spec.Id)
			// Only update the config map if resource version has changed (means api definition has changed).
			err = d.k8sClient.Update(d.ctx, cm)
		} else {
			d.log.Info("No change detected on api. Skipped.", "id", apiDefinition.Spec.Id)
			return nil
		}
	} else {
		d.log.Info("Creating config map for api.", "id", apiDefinition.Spec.Id, "name", apiDefinition.Name)
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
