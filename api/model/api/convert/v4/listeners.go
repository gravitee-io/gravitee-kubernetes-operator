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

package v4

import (
	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
)

func toListeners(proxy *v2.Proxy) []*v4.Listener {
	if proxy == nil || len(proxy.VirtualHosts) == 0 {
		return nil
	}

	listener := v4.NewHttpListenerBase()
	paths := make([]interface{}, 0)
	for _, vHost := range proxy.VirtualHosts {
		paths = append(paths, v4.NewPath(vHost.Host, vHost.Path))
	}
	listener.Put("paths", paths)
	return []*v4.Listener{listener}
}
