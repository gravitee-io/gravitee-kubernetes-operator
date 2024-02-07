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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	definitionVersionKey = "definitionVersion"
	definitionKey        = "definition"
	managedByKey         = "managed-by"
	gioTypeKey           = "gio-type"
	orgKey               = "organizationId"
	envKey               = "environmentId"
	defaultEnvId         = "DEFAULT"
	defaultOrgId         = "DEFAULT"
)

func (d *Delegate) updateConfigMap(api *gio.ApiDefinition) error {
	if api.Spec.State == base.StateStopped {
		if err := d.deleteConfigMap(api); err != nil {
			d.log.Error(err, "Unable to delete ConfigMap from API definition")
			return err
		}
	} else {
		if err := d.saveConfigMap(api); err != nil {
			d.log.Error(err, "Unable to create or update ConfigMap from API definition")
			return err
		}
	}

	return nil
}

func (d *Delegate) saveConfigMap(
	apiDefinition *gio.ApiDefinition,
) error {
	if apiDefinition.Spec.State == base.StateStopped {
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
		managedByKey: keys.CrdGroup,
		gioTypeKey:   keys.CrdApiDefinitionResource + "." + keys.CrdGroup,
	}

	cm.Data = map[string]string{
		definitionVersionKey: apiDefinition.ResourceVersion,
	}

	spec := &(apiDefinition.Spec)

	if d.apim != nil {
		cm.Data[orgKey] = d.apim.OrgID()
		cm.Data[envKey] = d.apim.EnvID()
	} else {
		cm.Data[orgKey] = defaultOrgId
		cm.Data[envKey] = defaultEnvId
	}

	if spec.ID == "" {
		spec.ID = string(apiDefinition.UID)
	}

	jsonSpec, err := json.Marshal(spec)
	if err != nil {
		return err
	}

	cm.Data[definitionKey] = string(jsonSpec)

	currentApiDefinition := &v1.ConfigMap{}

	lookupKey := types.NamespacedName{Name: cm.Name, Namespace: cm.Namespace}

	err = d.k8s.Get(d.ctx, lookupKey, currentApiDefinition)
	if errors.IsNotFound(err) {
		d.log.Info("Creating config map for API.", "id", apiDefinition.Spec.ID, "name", apiDefinition.Name)
		return d.k8s.Create(d.ctx, cm)
	}

	if err != nil {
		return err
	}

	// Only update the config map if resource version has changed (means api definition has changed).
	if currentApiDefinition.Data[definitionVersionKey] != apiDefinition.ResourceVersion {
		d.log.Info("Updating ConfigMap", "id", apiDefinition.Spec.ID)
		return d.k8s.Update(d.ctx, cm)
	}

	d.log.Info("No change detected on API. Skipped.", "id", apiDefinition.Spec.ID)
	return nil
}

func (d *Delegate) deleteConfigMap(api *gio.ApiDefinition) error {
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      api.Name,
			Namespace: api.Namespace,
		},
	}

	d.log.Info("Deleting Config Map associated to API if exists")
	return client.IgnoreNotFound(d.k8s.Delete(d.ctx, configMap))
}
