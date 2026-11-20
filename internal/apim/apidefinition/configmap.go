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

package apidefinition

import (
	"context"

	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
)

const (
	definitionKey        = "definition"
	definitionVersionKey = "apiDefinitionVersion"
	managedByKey         = "managed-by"
	gioTypeKey           = "gio-type"
	orgKey               = "organizationId"
	envKey               = "environmentId"
	// DefaultOrgID and DefaultEnvID are written when the API has no management context.
	DefaultOrgID = "DEFAULT"
	DefaultEnvID = "DEFAULT"
)

// GatewayConfigMap is the ConfigMap a gateway reads an API definition from when the API syncs from Kubernetes.
// It lives in the API's namespace, under the API's name, owned by the API.
type GatewayConfigMap struct {
	Name      string
	Namespace string
	Owner     metav1.OwnerReference
	// OrgID and EnvID default to DEFAULT when empty.
	OrgID string
	EnvID string
	// DefinitionVersion is written as apiDefinitionVersion when set (v4 APIs).
	DefinitionVersion string
	// Definition is the JSON definition the gateway deploys.
	Definition []byte
}

// SaveConfigMap creates or updates the gateway ConfigMap.
func SaveConfigMap(ctx context.Context, target GatewayConfigMap) error {
	cm := &v1.ConfigMap{}
	cm.SetOwnerReferences([]metav1.OwnerReference{target.Owner})
	cm.Namespace = target.Namespace
	cm.Name = target.Name
	cm.CreationTimestamp = metav1.Now()
	cm.Labels = map[string]string{
		managedByKey: core.CRDGroup,
		gioTypeKey:   core.CRDApiDefinitionResource + "." + core.CRDGroup,
	}

	cm.Data = map[string]string{orgKey: DefaultOrgID, envKey: DefaultEnvID}
	if target.OrgID != "" {
		cm.Data[orgKey] = target.OrgID
		cm.Data[envKey] = target.EnvID
	}
	if target.DefinitionVersion != "" {
		cm.Data[definitionVersionKey] = target.DefinitionVersion
	}
	cm.Data[definitionKey] = string(target.Definition)

	current := &v1.ConfigMap{}
	err := k8s.GetClient().Get(ctx, types.NamespacedName{Name: cm.Name, Namespace: cm.Namespace}, current)
	if kerrors.IsNotFound(err) {
		log.Debug(ctx, "Creating config map for API", "namespace", cm.Namespace, "name", cm.Name)
		return k8s.GetClient().Create(ctx, cm)
	}
	if err != nil {
		return err
	}

	log.Debug(ctx, "Updating config map for API", "namespace", cm.Namespace, "name", cm.Name)
	return k8s.GetClient().Update(ctx, cm)
}

// DeleteConfigMap deletes the gateway ConfigMap of the API named name in namespace, if any.
func DeleteConfigMap(ctx context.Context, name, namespace string) error {
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	log.Debug(ctx, "Deleting any config map associated to the API", "namespace", namespace, "name", name)
	return client.IgnoreNotFound(k8s.GetClient().Delete(ctx, configMap))
}

// DeleteOwnedConfigMap deletes the gateway ConfigMap of an API that no longer syncs from Kubernetes. It reads the
// ConfigMap through the client first, so an API that never had one costs no delete call, and leaves alone a
// ConfigMap of the same name that the API does not own.
func DeleteOwnedConfigMap(ctx context.Context, name, namespace string, owner types.UID) error {
	current := &v1.ConfigMap{}
	err := k8s.GetClient().Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, current)
	if err != nil {
		return client.IgnoreNotFound(err)
	}
	for _, ref := range current.GetOwnerReferences() {
		if ref.UID == owner {
			log.Debug(ctx, "Deleting config map following a switch to sync from management",
				"namespace", namespace, "name", name)
			return client.IgnoreNotFound(k8s.GetClient().Delete(ctx, current))
		}
	}
	return nil
}
