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
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
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
// whose APIM refuses every PUT with the ownership finding of the resource's kind.
const refusingContextPath = "/apis/gravitee.io/v1alpha1/namespaces/gravitee/managementcontexts/refusing-ctx"

// acceptingContextPath is where the cluster stub serves ManagementContext gravitee/accepting-ctx,
// whose APIM accepts every PUT and answers the state acceptedStates holds for its collection.
const acceptingContextPath = "/apis/gravitee.io/v1alpha1/namespaces/gravitee/managementcontexts/accepting-ctx"

// acceptedStates are the states the accepting APIM answers, by collection path suffix.
var acceptedStates = map[string]string{
	"/authz/policies": `{"hrid":"gravitee-hotel-booking-read","id":"5f1c","organizationId":"DEFAULT",` +
		`"environmentId":"DEFAULT","policyStatus":"DEPLOYED"}`,
}

// lastReceived holds the last body an accepting or the cascade APIM received, by collection path
// suffix.
var lastReceived sync.Map

// answering accepts every PUT, recording its body in lastReceived, and answers the state states
// holds for its collection path suffix.
func answering(states map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodPut {
			for suffix, state := range states {
				if strings.HasSuffix(r.URL.Path, suffix) {
					body, _ := io.ReadAll(r.Body)
					lastReceived.Store(suffix, string(body))
					_, _ = w.Write([]byte(state))
					return
				}
			}
		}
		w.WriteHeader(http.StatusNotFound)
	}
}

func contextBody(name, baseURL string) string {
	return `{"apiVersion":"gravitee.io/v1alpha1","kind":"ManagementContext",` +
		`"metadata":{"name":"` + name + `","namespace":"gravitee"},` +
		`"spec":{"baseUrl":"` + baseURL + `","organizationId":"DEFAULT","environmentId":"DEFAULT",` +
		`"auth":{"bearerToken":"test"}}}`
}

// cascadeContextPath is where the cluster stub serves ManagementContext gravitee/cascade-ctx, whose
// APIM holds no parent anymore: it refuses a RESOURCE declared under one, as the platform does once
// that parent is deleted, and accepts anything else.
const cascadeContextPath = "/apis/gravitee.io/v1alpha1/namespaces/gravitee/managementcontexts/cascade-ctx"

// missingParentFinding is the platform's finding for a RESOURCE declared under parent, which does not exist.
func missingParentFinding(parent string) string {
	return "parent " + parent + " does not exist: a RESOURCE is deleted with its parent, " +
		"so declare the parent first or drop it from parents"
}

// cascadePlatform answers an AuthzEntity PUT as a platform whose parents have all been deleted,
// recording its body in lastReceived.
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
	body, err := io.ReadAll(r.Body)
	if r.Method != http.MethodPut || err != nil || json.Unmarshal(body, &sent) != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if strings.HasSuffix(r.URL.Path, "/authz/entities") {
		lastReceived.Store("/authz/entities", string(body))
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
// but gravitee/refusing-ctx, gravitee/accepting-ctx and gravitee/cascade-ctx.
var _ = BeforeSuite(func() {
	// The 400 the refusing APIM answers: an entityId another source holds.
	refusingPlatformBody := `{"errors":{"severe":[` + mustJSON(ownershipFinding) + `]}}`
	refusingPolicyBody := `{"errors":{"severe":[` + mustJSON(policyOwnershipFinding) + `]}}`
	platform := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		if strings.HasSuffix(r.URL.Path, "/authz/policies") {
			_, _ = w.Write([]byte(refusingPolicyBody))
			return
		}
		_, _ = w.Write([]byte(refusingPlatformBody))
	}))
	DeferCleanup(platform.Close)
	cascade := httptest.NewServer(http.HandlerFunc(cascadePlatform))
	DeferCleanup(cascade.Close)

	accepting := httptest.NewServer(answering(acceptedStates))
	DeferCleanup(accepting.Close)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case refusingContextPath:
			_, _ = w.Write([]byte(contextBody("refusing-ctx", platform.URL)))
			return
		case acceptingContextPath:
			_, _ = w.Write([]byte(contextBody("accepting-ctx", accepting.URL)))
			return
		case cascadeContextPath:
			_, _ = w.Write([]byte(contextBody("cascade-ctx", cascade.URL)))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"kind":"Status","apiVersion":"v1","status":"Failure",` +
			`"reason":"NotFound","code":404,"message":"not found"}`))
	}))
	DeferCleanup(server.Close)
	dynamic.UseConfig(&rest.Config{Host: server.URL})
})
