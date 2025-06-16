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

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"k8s.io/apimachinery/pkg/util/sets"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func buildListeners(ctx context.Context, route *gwAPIv1.HTTPRoute) ([]*v4.GenericListener, error) {
	listener := v4.NewHTTPListener()
	hostnames := k8s.ResolveRouteHostnames(ctx, route)
	paths, err := getPaths(ctx, route, hostnames)
	if err != nil {
		return nil, err
	}
	listener.Paths = paths

	return []*v4.GenericListener{
		v4.ToGenericListener(listener),
	}, nil
}

func getPaths(ctx context.Context, route *gwAPIv1.HTTPRoute, hostnames []string) ([]*v4.Path, error) {
	if len(hostnames) == 0 {
		paths, err := getPathsWithoutHostnames(ctx, route)
		if err != nil {
			return nil, err
		}

		return paths, nil
	}
	return getPathsWithHostnames(route, hostnames), nil
}

func getPathsWithoutHostnames(ctx context.Context, route *gwAPIv1.HTTPRoute) ([]*v4.Path, error) {
	paths := sets.New[*v4.Path]()
	routePaths := extractPaths(route)

	for _, path := range routePaths {
		if len(route.Spec.ParentRefs) == 0 {
			paths.Insert(v4.NewPath("", path))
		} else {
			for _, parentRef := range route.Spec.ParentRefs {
				host, err := resolveParentRefHostname(ctx, route, parentRef)
				if err != nil {
					return nil, err
				}
				paths.Insert(v4.NewPath(host, path))
			}
		}
	}

	return paths.UnsortedList(), nil
}

func resolveParentRefHostname(ctx context.Context, route *gwAPIv1.HTTPRoute,
	parentRef gwAPIv1.ParentReference) (string, error) {
	if parentRef.SectionName != nil && k8s.IsGatewayKind(parentRef) {
		gateway, err := k8s.ResolveGateway(ctx, route.ObjectMeta, parentRef)
		if err != nil {
			return "", err
		}

		for _, l := range gateway.Spec.Listeners {
			if l.Name == *parentRef.SectionName && l.Hostname != nil {
				return string(*l.Hostname), nil
			}
		}
	}

	return "", nil
}

func getPathsWithHostnames(route *gwAPIv1.HTTPRoute, hostnames []string) []*v4.Path {
	paths := []*v4.Path{}
	routePaths := extractPaths(route)
	for _, hostname := range hostnames {
		for _, path := range routePaths {
			paths = append(paths, v4.NewPath(hostname, path))
		}
	}
	return paths
}

func extractPaths(route *gwAPIv1.HTTPRoute) []string {
	paths := sets.New[string]()
	for _, rule := range route.Spec.Rules {
		for _, match := range rule.Matches {
			if *match.Path.Type != gwAPIv1.PathMatchRegularExpression {
				paths.Insert(addTrailingSlash(*match.Path.Value))
			} else {
				paths.Insert(rootPath)
			}
		}
	}
	return paths.UnsortedList()
}
