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
	"context"
	"encoding/json"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
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
	defaultEnvID         = "DEFAULT"
	defaultOrgID         = "DEFAULT"
)

func updateConfigMap(
	ctx context.Context,
	api *v1alpha1.ApiDefinition,
) error {
	if api.Spec.State == base.StateStopped {
		if err := deleteConfigMap(ctx, api); err != nil {
			log.Error(ctx, err, "Unable to delete config map for API definition", log.KeyValues(api)...)
			return err
		}
	} else {
		if err := saveConfigMap(ctx, api); err != nil {
			log.Error(ctx, err, "Unable to create or update config map for API definition", log.KeyValues(api)...)
			return err
		}
	}

	return nil
}

func saveConfigMap(
	ctx context.Context,
	apiDefinition core.ApiDefinitionObject,
) error {
	// Create config map with some specific metadata that will be used to check changes across 'Update' events.
	cm := &v1.ConfigMap{}

	// Set OwnerReference on config map to be able to delete it when API is deleted.
	// 📝 ConfigMap should be in same namespace as ApiDefinition.
	newOwnerReferences := []metav1.OwnerReference{
		{
			Kind:       apiDefinition.GetObjectKind().GroupVersionKind().Kind,
			Name:       apiDefinition.GetName(),
			APIVersion: apiDefinition.GetObjectKind().GroupVersionKind().GroupVersion().String(),
			UID:        apiDefinition.GetUID(),
		},
	}
	cm.SetOwnerReferences(newOwnerReferences)

	cm.Namespace = apiDefinition.GetNamespace()
	cm.Name = apiDefinition.GetName()

	cm.CreationTimestamp = metav1.Now()
	cm.Labels = map[string]string{
		managedByKey: core.CRDGroup,
		gioTypeKey:   core.CRDApiDefinitionResource + "." + core.CRDGroup,
	}

	cm.Data = map[string]string{
		definitionVersionKey: apiDefinition.GetResourceVersion(),
	}

	if apiDefinition.GetOrgID() != "" {
		cm.Data[orgKey] = apiDefinition.GetOrgID()
		cm.Data[envKey] = apiDefinition.GetEnvID()
	} else {
		cm.Data[orgKey] = defaultOrgID
		cm.Data[envKey] = defaultEnvID
	}

	var payload any
	switch t := apiDefinition.(type) {
	case *v1alpha1.ApiDefinition:
		spec := &(t.Spec)
		if spec.ID == "" {
			spec.ID = string(t.UID)
		}
		payload = spec
	case *v1alpha1.ApiV4Definition:
		cm.Data["apiDefinitionVersion"] = "4.0.0"
		payload = t.ToGatewayDefinition()
	}

	jsonSpec, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	cm.Data[definitionKey] = string(jsonSpec)

	currentApiDefinition := &v1.ConfigMap{}

	lookupKey := types.NamespacedName{Name: cm.Name, Namespace: cm.Namespace}

	err = k8s.GetClient().Get(ctx, lookupKey, currentApiDefinition)
	if errors.IsNotFound(err) {
		log.Debug(ctx, "Creating config map for API", log.KeyValues(apiDefinition)...)
		return k8s.GetClient().Create(ctx, cm)
	}

	if err != nil {
		return err
	}

	if currentApiDefinition.Data[definitionVersionKey] != apiDefinition.GetResourceVersion() {
		log.Debug(ctx, "Updating config map", log.KeyValues(apiDefinition)...)
		return k8s.GetClient().Update(ctx, cm)
	}

	return nil
}

func deleteConfigMap(ctx context.Context, api client.Object) error {
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      api.GetName(),
			Namespace: api.GetNamespace(),
		},
	}

	log.Debug(ctx, "Deleting any config map associated to the API", log.KeyValues(api)...)
	return client.IgnoreNotFound(k8s.GetClient().Delete(ctx, configMap))
}
