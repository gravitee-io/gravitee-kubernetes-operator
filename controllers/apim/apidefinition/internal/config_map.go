// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

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

func (d *Delegate) updateConfigMap(api *gio.ApiDefinition, context *DelegateContext) error {
	if api.Spec.State == model.StateStopped {
		if err := d.deleteConfigMap(api.Namespace, api.Name); err != nil {
			d.log.Error(err, "Unable to delete ConfigMap from API definition")
			return err
		}
	} else {
		if err := d.saveConfigMap(api, context); err != nil {
			d.log.Error(err, "Unable to create or update ConfigMap from API definition")
			return err
		}
	}

	return nil
}

func (d *Delegate) saveConfigMap(
	apiDefinition *gio.ApiDefinition, context *DelegateContext,
) error {
	if apiDefinition.Spec.State == model.StateStopped {
		return nil
	}

	// Create config map with some specific metadata that will be used to check changes across 'Update' events.
	cm := &v1.ConfigMap{}

	// Set OwnerReference on config map to be able to delete it when API is deleted.
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
		"definitionVersion": apiDefinition.ResourceVersion,
	}

	spec := apiDefinition.Spec.DeepCopy()

	if context != nil && context.hasManagement() {
		status := apiDefinition.Status.Contexts[context.Location]
		cm.Data["organizationId"] = status.OrgID
		cm.Data["environmentId"] = status.EnvID
		spec.ID = status.ID
	} else {
		spec.ID = apiDefinition.GetID()
	}

	jsonSpec, err := json.Marshal(spec)
	if err != nil {
		return err
	}

	cm.Data["definition"] = string(jsonSpec)

	currentApiDefinition := &v1.ConfigMap{}
	err = d.k8s.Get(d.ctx, types.NamespacedName{Name: cm.Name, Namespace: cm.Namespace}, currentApiDefinition)

	if err == nil {
		if currentApiDefinition.Data["definitionVersion"] != apiDefinition.ResourceVersion {
			d.log.Info("Updating ConfigMap", "id", apiDefinition.Spec.ID)
			// Only update the config map if resource version has changed (means api definition has changed).
			err = d.k8s.Update(d.ctx, cm)
		} else {
			d.log.Info("No change detected on api. Skipped.", "id", apiDefinition.Spec.ID)
			return nil
		}
	} else {
		d.log.Info("Creating config map for api.", "id", apiDefinition.Spec.ID, "name", apiDefinition.Name)
		err = d.k8s.Create(d.ctx, cm)
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
	err := d.k8s.Delete(d.ctx, configMap)

	if errors.IsNotFound(err) {
		return nil
	}

	return err
}
