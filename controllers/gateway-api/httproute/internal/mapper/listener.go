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

func buildListeners(ctx context.Context, route *gwAPIv1.HTTPRoute) []*v4.GenericListener {
	listener := v4.NewHTTPListener()
	hostnames := k8s.ResolveRouteHostnames(ctx, route)
	listener.Paths = getPaths(route, hostnames)
	return []*v4.GenericListener{
		v4.ToGenericListener(listener),
	}
}

func getPaths(route *gwAPIv1.HTTPRoute, hostnames []string) []*v4.Path {
	if len(hostnames) == 0 {
		return getPathsWithoutHostnames(route)
	}
	return getPathsWithHostnames(route, hostnames)
}

func getPathsWithoutHostnames(route *gwAPIv1.HTTPRoute) []*v4.Path {
	paths := []*v4.Path{}
	routePaths := extractPaths(route)
	for _, path := range routePaths {
		paths = append(paths, v4.NewPath("", path))
	}
	return paths
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
