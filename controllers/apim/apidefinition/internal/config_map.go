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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1beta1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
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

func (d *Delegate) saveConfigMap(
	apiDefinition *v1beta1.ApiDefinition,
) error {
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

	definition := apiDefinition.ToGatewayDefinition()

	cm.Data = map[string]string{
		definitionVersionKey: apiDefinition.ResourceVersion,
	}

	if d.apim != nil {
		cm.Data[orgKey] = d.apim.OrgID()
		cm.Data[envKey] = d.apim.EnvID()
	} else {
		cm.Data[orgKey] = defaultOrgId
		cm.Data[envKey] = defaultEnvId
	}

	jsonDefinition, err := json.Marshal(definition)
	if err != nil {
		return err
	}

	cm.Data[definitionKey] = string(jsonDefinition)

	currentApiDefinition := &v1.ConfigMap{}

	lookupKey := types.NamespacedName{Name: cm.Name, Namespace: cm.Namespace}

	err = d.k8s.Get(d.ctx, lookupKey, currentApiDefinition)
	if errors.IsNotFound(err) {
		log.Debug(d.ctx, "Storing API definition to config map")
		return d.k8s.Create(d.ctx, cm)
	}

	if err != nil {
		return err
	}

	// Only update the config map if resource version has changed (means api definition has changed).
	if currentApiDefinition.Data[definitionVersionKey] != apiDefinition.ResourceVersion {
		log.Debug(d.ctx, "Updating API definition in config map")
		return d.k8s.Update(d.ctx, cm)
	}

	log.Debug(d.ctx, "Skipping config map update as API definition has not changed")
	return nil
}

func (d *Delegate) deleteConfigMap(api *v1beta1.ApiDefinition) error {
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      api.Name,
			Namespace: api.Namespace,
		},
	}

	log.Debug(d.ctx, "Deleting config map associated to API (if exists)")
	return client.IgnoreNotFound(d.k8s.Delete(d.ctx, configMap))
}
