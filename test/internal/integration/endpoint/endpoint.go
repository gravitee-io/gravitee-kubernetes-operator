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

package endpoint

import (
	"fmt"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	. "github.com/onsi/gomega"
)

func FromStrURL(strURL string) http.URL {
	url, err := http.NewURL(strURL)
	Expect(err).ToNot(HaveOccurred())
	return url
}

func ForV2(api *v1alpha1.ApiDefinition) http.URL {
	path := api.Spec.Proxy.VirtualHosts[0].Path
	return FromStrURL(constants.GatewayUrl).WithPath(path)
}

func ForV4Proxy(l v4.Listener) http.URL {
	url, err := resolveHTTPListener(l)
	Expect(err).ToNot(HaveOccurred())
	return url
}

func resolveHTTPListener(l v4.Listener) (http.URL, error) {
	switch t := l.(type) {
	case *v4.HttpListener:
		return FromStrURL(constants.GatewayUrl).WithPath(t.Paths[0].Path), nil
	case *v4.GenericListener:
		return resolveHTTPListener(t.ToListener())
	default:
		return http.URL{}, fmt.Errorf("unknown proxy type %s", t)
	}
}
