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
	"net/http"
	"net/http/httptest"
	"strings"
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

// acceptingContextPath is where the cluster stub serves ManagementContext gravitee/accepting-ctx,
// whose APIM accepts every PUT and answers the state acceptedStates holds for its collection.
const acceptingContextPath = "/apis/gravitee.io/v1alpha1/namespaces/gravitee/managementcontexts/accepting-ctx"

// acceptedStates are the states the accepting APIM answers, by collection path suffix.
var acceptedStates = map[string]string{
	"/aim/catalog/mcp-servers": `{"hrid":"gravitee-github","id":"7a2e","organizationId":"DEFAULT",` +
		`"environmentId":"DEFAULT","authPolicyRef":"MCPServer::\"github\"","tools":[{"name":"create_issue",` +
		`"entityId":"mcp-tool.github.create_issue","authPolicyRef":"MCPTool::\"github.create_issue\""}]}`,
	"/aim/mcp-proxies": `{"hrid":"gravitee-support-studio","id":"9b3d","organizationId":"DEFAULT",` +
		`"environmentId":"DEFAULT","tools":[{"server":"gravitee-github-mcp","tool":"search_issues",` +
		`"alias":"gh_search","entityId":"mcp-tool.github.search_issues",` +
		`"authPolicyRef":"MCPTool::\"github.search_issues\""}]}`,
}

// studioOnlyContextPath is where the cluster stub serves ManagementContext gravitee/studio-only-ctx,
// whose APIM accepts every PUT and answers studioOnlyStates: an McpProxy state with studio.tools and
// no top-level tools list, as a platform that predates it does.
const studioOnlyContextPath = "/apis/gravitee.io/v1alpha1/namespaces/gravitee/managementcontexts/studio-only-ctx"

var studioOnlyStates = map[string]string{
	"/aim/mcp-proxies": `{"hrid":"gravitee-support-studio","id":"9b3d","organizationId":"DEFAULT",` +
		`"environmentId":"DEFAULT","studio":{"tools":[{"server":"gravitee-github-mcp","tool":"search_issues",` +
		`"alias":"gh_search","entityId":"mcp-tool.github.search_issues"}],"upstreamAuth":[],"enableFGA":true}}`,
}

// answering accepts every PUT and answers the state states holds for its collection path suffix.
func answering(states map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodPut {
			for suffix, state := range states {
				if strings.HasSuffix(r.URL.Path, suffix) {
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

// The dynamic client resolves ManagementContexts and falls back to the ambient kubeconfig, which
// reaches a developer's cluster and exits where there is none (CI). Every lookup answers 404 here,
// but gravitee/accepting-ctx and gravitee/studio-only-ctx.
var _ = BeforeSuite(func() {
	accepting := httptest.NewServer(answering(acceptedStates))
	DeferCleanup(accepting.Close)
	studioOnly := httptest.NewServer(answering(studioOnlyStates))
	DeferCleanup(studioOnly.Close)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case acceptingContextPath:
			_, _ = w.Write([]byte(contextBody("accepting-ctx", accepting.URL)))
			return
		case studioOnlyContextPath:
			_, _ = w.Write([]byte(contextBody("studio-only-ctx", studioOnly.URL)))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"kind":"Status","apiVersion":"v1","status":"Failure",` +
			`"reason":"NotFound","code":404,"message":"not found"}`))
	}))
	DeferCleanup(server.Close)
	dynamic.UseConfig(&rest.Config{Host: server.URL})
})
