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
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/gateway-api/httproute/internal/mapper"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/uuid"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const (
	configMapPrefix    = "httproute-"
	maxK8sNameLength   = 253
	hashSuffixLength   = 8
	definitionKey      = "definition"
	apiDefVersionKey   = "apiDefinitionVersion"
	apiDefVersionValue = "4.0.0"
	orgKey             = "organizationId"
	envKey             = "environmentId"
	defaultOrgID       = "DEFAULT"
	defaultEnvID       = "DEFAULT"
	managedByKey       = "managed-by"
	gioTypeKey         = "gio-type"
)

func Program(ctx context.Context, route *gwAPIv1.HTTPRoute) error {
	api, err := mapper.Map(ctx, route)
	if err != nil {
		return err
	}
	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		if err := k8s.GetClient().Get(ctx, client.ObjectKeyFromObject(api), api); client.IgnoreNotFound(err) != nil {
			return err
		}

		return k8s.CreateOrUpdate(ctx, api, func() error {
			api.SetOwnerReferences(getOwnerReferences(route))
			spec, err := mapper.MapSpec(ctx, route)
			if err != nil {
				return err
			}
			api.Spec = spec
			return nil
		})
	})
}

// ProgramConfigMap creates a ConfigMap directly from the HTTPRoute, bypassing
// the intermediate ApiV4Definition CR. This frees the CR name so a user-managed
// ApiV4Definition with the same name can coexist.
func ProgramConfigMap(ctx context.Context, route *gwAPIv1.HTTPRoute) error {
	api, err := mapper.Map(ctx, route)
	if err != nil {
		return err
	}

	api.Spec.ID = uuid.FromStrings(string(route.UID))
	api.Spec.CrossID = uuid.FromStrings(route.Namespace, route.Name)

	if api.Spec.Plans != nil {
		for key, plan := range *api.Spec.Plans {
			plan.ID = uuid.FromStrings(api.Spec.ID, key)
		}
	}

	gwDef := api.Spec.Api.ToGatewayDefinition()
	jsonSpec, err := json.Marshal(gwDef)
	if err != nil {
		return err
	}

	cm := &v1.ConfigMap{
		ObjectMeta: metaV1.ObjectMeta{
			Name:            buildConfigMapName(route),
			Namespace:       route.Namespace,
			OwnerReferences: getOwnerReferences(route),
			Labels: map[string]string{
				managedByKey: core.CRDGroup,
				gioTypeKey:   core.CRDApiDefinitionResource + "." + core.CRDGroup,
			},
		},
		Data: map[string]string{
			definitionKey:    string(jsonSpec),
			apiDefVersionKey: apiDefVersionValue,
			orgKey:           defaultOrgID,
			envKey:           defaultEnvID,
		},
	}

	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		existing := &v1.ConfigMap{}
		err := k8s.GetClient().Get(ctx, client.ObjectKeyFromObject(cm), existing)
		if errors.IsNotFound(err) {
			return k8s.GetClient().Create(ctx, cm)
		}
		if err != nil {
			return err
		}
		existing.Labels = cm.Labels
		existing.Data = cm.Data
		existing.SetOwnerReferences(cm.GetOwnerReferences())
		return k8s.GetClient().Update(ctx, existing)
	})
}

func buildConfigMapName(route *gwAPIv1.HTTPRoute) string {
	name := configMapPrefix + route.Name
	if len(name) <= maxK8sNameLength {
		return name
	}
	h := sha256.Sum256([]byte(route.Name))
	suffix := hex.EncodeToString(h[:])[:hashSuffixLength]
	truncateAt := maxK8sNameLength - len(configMapPrefix) - 1 - hashSuffixLength
	return fmt.Sprintf("%s%s-%s", configMapPrefix, route.Name[:truncateAt], suffix)
}

func getOwnerReferences(httpRoute *gwAPIv1.HTTPRoute) []metaV1.OwnerReference {
	kind := httpRoute.GetObjectKind().GroupVersionKind().Kind
	version := httpRoute.GetObjectKind().GroupVersionKind().GroupVersion().String()
	return []metaV1.OwnerReference{
		{
			Kind:       kind,
			APIVersion: version,
			Name:       httpRoute.GetName(),
			UID:        httpRoute.GetUID(),
		},
	}
}
