package apim_test

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

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/rest"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	//+kubebuilder:scaffold:imports
)

func TestAPIM(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "APIM unit tests suite")
}

// refusingContextPath is where the cluster stub serves ManagementContext gravitee/refusing-ctx,
// whose APIM refuses every PUT with the ownership finding.
const refusingContextPath = "/apis/gravitee.io/v1alpha1/namespaces/gravitee/managementcontexts/refusing-ctx"

// cascadeContextPath is where the cluster stub serves ManagementContext gravitee/cascade-ctx, whose
// APIM holds no parent anymore: it refuses a RESOURCE declared under one, as the platform does once
// that parent is deleted, and accepts anything else.
const cascadeContextPath = "/apis/gravitee.io/v1alpha1/namespaces/gravitee/managementcontexts/cascade-ctx"

// missingParentFinding is the platform's finding for a RESOURCE declared under parent, which does not exist.
func missingParentFinding(parent string) string {
	return "parent " + parent + " does not exist: a RESOURCE is deleted with its parent, " +
		"so declare the parent first or drop it from parents"
}

// cascadePlatform answers an AuthzEntity PUT as a platform whose parents have all been deleted.
func cascadePlatform(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sent := struct {
		HRID    string `json:"hrid"`
		Kind    string `json:"kind"`
		Parents []struct {
			EntityType string `json:"entityType"`
			EntityID   string `json:"entityId"`
		} `json:"parents"`
	}{}
	if r.Method != http.MethodPut || json.NewDecoder(r.Body).Decode(&sent) != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if sent.Kind == "RESOURCE" && len(sent.Parents) > 0 {
		parent := sent.Parents[0].EntityType + `::"` + sent.Parents[0].EntityID + `"`
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"errors":{"severe":[` + mustJSON(missingParentFinding(parent)) + `]}}`))
		return
	}
	_, _ = w.Write([]byte(`{"hrid":"` + sent.HRID + `","id":"5e7b","organizationId":"DEFAULT","environmentId":"DEFAULT"}`))
}

// The dynamic client resolves ManagementContexts and falls back to the ambient kubeconfig, which
// reaches a developer's cluster and exits where there is none (CI). Every lookup answers 404 here,
// but gravitee/refusing-ctx and gravitee/cascade-ctx.
var _ = BeforeSuite(func() {
	// The 400 the refusing APIM answers: an entityId another source holds.
	refusingPlatformBody := `{"errors":{"severe":[` + mustJSON(ownershipFinding) + `]}}`
	platform := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(refusingPlatformBody))
	}))
	DeferCleanup(platform.Close)
	cascade := httptest.NewServer(http.HandlerFunc(cascadePlatform))
	DeferCleanup(cascade.Close)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == refusingContextPath {
			_, _ = w.Write([]byte(`{"apiVersion":"gravitee.io/v1alpha1","kind":"ManagementContext",` +
				`"metadata":{"name":"refusing-ctx","namespace":"gravitee"},` +
				`"spec":{"baseUrl":"` + platform.URL + `","organizationId":"DEFAULT","environmentId":"DEFAULT",` +
				`"auth":{"bearerToken":"test"}}}`))
			return
		}
		if r.URL.Path == cascadeContextPath {
			_, _ = w.Write([]byte(`{"apiVersion":"gravitee.io/v1alpha1","kind":"ManagementContext",` +
				`"metadata":{"name":"cascade-ctx","namespace":"gravitee"},` +
				`"spec":{"baseUrl":"` + cascade.URL + `","organizationId":"DEFAULT","environmentId":"DEFAULT",` +
				`"auth":{"bearerToken":"test"}}}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"kind":"Status","apiVersion":"v1","status":"Failure",` +
			`"reason":"NotFound","code":404,"message":"not found"}`))
	}))
	DeferCleanup(server.Close)
	dynamic.UseConfig(&rest.Config{Host: server.URL})
})
