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

package v2

import (
	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
)

func toProxy(
	listeners []*v4.Listener, endpointGroups []*v4.EndpointGroup, analytics *v4.Analytics,
) *v2.Proxy {
	proxy := &v2.Proxy{}
	proxy.Groups = toEndpointGroups(endpointGroups)
	proxy.VirtualHosts = toVirtualHosts(listeners)

	if len(proxy.Groups) == 0 && len(proxy.VirtualHosts) == 0 {
		return nil
	}

	if analytics != nil {
		proxy.Logging = toLogging(analytics)
	}

	return proxy
}

func toVirtualHosts(listeners []*v4.Listener) []*v2.VirtualHost {
	var virtualHosts []*v2.VirtualHost
	for _, listener := range listeners {
		for _, path := range getPaths(listener) {
			virtualHosts = append(virtualHosts, toVirtualHost(path))
		}
	}
	return virtualHosts
}

func getPaths(listener *v4.Listener) []*v4.Path {
	impl := listener.GetSlice("paths")
	if impl == nil {
		return nil
	}

	paths := make([]*v4.Path, 0)
	for _, p := range impl {
		pImpl := utils.ToGenericStringMap(p)
		path := v4.NewPath(pImpl.GetString("host"), pImpl.GetString("path"))
		paths = append(paths, path)
	}

	return paths
}

func toVirtualHost(path *v4.Path) *v2.VirtualHost {
	return v2.NewVirtualHost(path.Host, path.Path)
}
