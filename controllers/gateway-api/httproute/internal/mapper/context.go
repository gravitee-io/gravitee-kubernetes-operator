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

package mapper

import (
	"context"
	"fmt"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type gatewayCacheKey struct{}

// GatewayCache stores fetched gateways keyed by a unique string representation of the parent reference.
type GatewayCache map[string]*gwAPIv1.Gateway

func WithGatewayCache(ctx context.Context, cache GatewayCache) context.Context {
	return context.WithValue(ctx, gatewayCacheKey{}, cache)
}

func getGatewayCache(ctx context.Context) GatewayCache {
	if cache, ok := ctx.Value(gatewayCacheKey{}).(GatewayCache); ok {
		return cache
	}
	return nil
}

func gatewayKey(routeMeta metaV1.ObjectMeta, ref gwAPIv1.ParentReference) string {
	ns := ref.Namespace
	if ns == nil {
		routeNS := gwAPIv1.Namespace(routeMeta.Namespace)
		ns = &routeNS
	}
	key := string(*ns) + "/" + string(ref.Name)
	if ref.SectionName != nil {
		key += "/" + string(*ref.SectionName)
	}
	if ref.Port != nil {
		key += "/" + fmt.Sprintf("%d", *ref.Port)
	}
	return key
}

func resolveGateway(
	ctx context.Context,
	routeMeta metaV1.ObjectMeta,
	parentRef gwAPIv1.ParentReference,
	resolveFn func(context.Context, metaV1.ObjectMeta, gwAPIv1.ParentReference) (*gwAPIv1.Gateway, error),
) (*gwAPIv1.Gateway, error) {
	cache := getGatewayCache(ctx)
	if cache == nil {
		return resolveFn(ctx, routeMeta, parentRef)
	}

	key := gatewayKey(routeMeta, parentRef)
	if cachedGw, found := cache[key]; found {
		return cachedGw, nil
	}

	gw, err := resolveFn(ctx, routeMeta, parentRef)
	if err == nil {
		cache[key] = gw
	}
	return gw, err
}
