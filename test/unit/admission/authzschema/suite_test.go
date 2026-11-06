// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package authzschema_test

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

func TestAuthzSchemaAdmission(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AuthzSchema admission hook tests suite")
}

// platformCalls counts the requests the APIM stub receives: a spec that expects the dry run to be
// skipped asserts none arrived.
var platformCalls atomic.Int32

// The dynamic client resolves ManagementContexts: it answers the dev-ctx context of the travel
// namespace, pointing at the APIM stub, and 404 for anything else.
var _ = BeforeSuite(func() {
	platform := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		platformCalls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	DeferCleanup(platform.Close)

	cluster := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/apis/gravitee.io/v1alpha1/namespaces/travel/managementcontexts/dev-ctx" {
			_, _ = w.Write([]byte(`{"apiVersion":"gravitee.io/v1alpha1","kind":"ManagementContext",` +
				`"metadata":{"name":"dev-ctx","namespace":"travel"},` +
				`"spec":{"baseUrl":"` + platform.URL + `","organizationId":"DEFAULT","environmentId":"DEFAULT",` +
				`"auth":{"bearerToken":"test"}}}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"kind":"Status","apiVersion":"v1","status":"Failure",` +
			`"reason":"NotFound","code":404,"message":"not found"}`))
	}))
	DeferCleanup(cluster.Close)
	dynamic.UseConfig(&rest.Config{Host: cluster.URL})
})
