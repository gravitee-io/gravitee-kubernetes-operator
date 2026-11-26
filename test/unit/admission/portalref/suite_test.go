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

package portalref_test

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/rest"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

func TestPortalRefAdmission(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PortalListing, Documentation and PortalLink portalRef admission tests suite")
}

// platformCalls counts the requests the APIM stub receives: a spec that expects the dry run to be
// skipped asserts none arrived.
var platformCalls atomic.Int32

const crds = "/apis/gravitee.io/v1alpha1/namespaces/travel/"

// The dynamic client resolves the dev-ctx context, pointing at the APIM stub, the API "orders" and
// two portals: "unsynced", which GKO has not written to APIM yet, and "synced", which it has.
var _ = BeforeSuite(func() {
	platform := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		platformCalls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	DeferCleanup(platform.Close)

	objects := map[string]string{
		crds + "managementcontexts/dev-ctx": `{"apiVersion":"gravitee.io/v1alpha1","kind":"ManagementContext",` +
			`"metadata":{"name":"dev-ctx","namespace":"travel"},` +
			`"spec":{"baseUrl":"` + platform.URL + `","organizationId":"DEFAULT","environmentId":"DEFAULT",` +
			`"auth":{"bearerToken":"test"}}}`,
		crds + "apiv4definitions/orders": `{"apiVersion":"gravitee.io/v1alpha1","kind":"ApiV4Definition",` +
			`"metadata":{"name":"orders","namespace":"travel"},` +
			`"spec":{"contextRef":{"name":"dev-ctx"},"name":"orders","version":"1.0","type":"PROXY"}}`,
		crds + "portals/unsynced": portal("unsynced", `{}`),
		crds + "portals/synced":   portal("synced", `{"id":"f9326ba3-0da0-3d4f-a7d5-f4a729b9e3fc"}`),
	}

	cluster := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if body, ok := objects[r.URL.Path]; ok {
			_, _ = w.Write([]byte(body))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"kind":"Status","apiVersion":"v1","status":"Failure",` +
			`"reason":"NotFound","code":404,"message":"not found"}`))
	}))
	DeferCleanup(cluster.Close)
	dynamic.UseConfig(&rest.Config{Host: cluster.URL})
})

func portal(name, status string) string {
	return `{"apiVersion":"gravitee.io/v1alpha1","kind":"Portal",` +
		`"metadata":{"name":"` + name + `","namespace":"travel"},` +
		`"spec":{"contextRef":{"name":"dev-ctx"},"name":"Travel portal"},"status":` + status + `}`
}
