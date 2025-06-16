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

	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/gateway-api/httproute/internal/mapper"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func Program(ctx context.Context, route *gwAPIv1.HTTPRoute) error {
	api, err := mapper.Map(ctx, route)
	if err != nil {
		return err
	}
	api.SetOwnerReferences(getOwnerReferences(route))
	return k8s.CreateOrUpdate(ctx, api, func() error {
		spec, err := mapper.MapSpec(ctx, route)
		if err != nil {
			return err
		}
		api.Spec = spec
		return nil
	})
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
