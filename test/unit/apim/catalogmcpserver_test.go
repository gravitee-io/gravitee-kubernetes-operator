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
package apim_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/catalogmcpserver"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	xErrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

func newCatalogMcpServer(auth *catalogmcpserver.Auth) *v1alpha1.CatalogMcpServer {
	return &v1alpha1.CatalogMcpServer{
		ObjectMeta: metav1.ObjectMeta{Name: "github", Namespace: "gravitee"},
		Spec: v1alpha1.CatalogMcpServerSpec{
			Context: &refs.NamespacedName{Name: "dev-ctx"},
			Type: catalogmcpserver.Type{
				EntityID: "mcp-server.github",
				Connection: catalogmcpserver.Connection{
					Endpoint: "https://api.githubcopilot.com/mcp/",
					Auth:     auth,
				},
			},
		},
	}
}

var _ = Describe("Catalog MCP server", func() {
	Describe("spec hash", func() {
		It("changes with the entity id and the connection", func() {
			base := newCatalogMcpServer(nil)
			otherID := newCatalogMcpServer(nil)
			otherID.Spec.EntityID = "mcp-server.gitlab"
			otherEndpoint := newCatalogMcpServer(nil)
			otherEndpoint.Spec.Connection.Endpoint = "https://mcp.gitlab.example.com/"

			Expect(base.Spec.Hash()).To(Equal(newCatalogMcpServer(nil).Spec.Hash()))
			Expect(base.Spec.Hash()).ToNot(Equal(otherID.Spec.Hash()))
			Expect(base.Spec.Hash()).ToNot(Equal(otherEndpoint.Spec.Hash()))
		})
	})

	Describe("wire mapping", func() {
		It("derives the hrid from the namespace and name", func() {
			dto := model.ToCatalogMcpServerDTO(newCatalogMcpServer(nil))
			Expect(dto.HRID).To(Equal("gravitee-github"))
			Expect(dto.EntityID).To(Equal("mcp-server.github"))
			Expect(dto.Connection.Endpoint).To(Equal("https://api.githubcopilot.com/mcp/"))
		})

		It("defaults the transport to HTTP", func() {
			dto := model.ToCatalogMcpServerDTO(newCatalogMcpServer(nil))
			Expect(dto.Connection.Transport).To(Equal("HTTP"))
		})

		It("sends an explicit NONE auth when the spec has none", func() {
			dto := model.ToCatalogMcpServerDTO(newCatalogMcpServer(nil))
			Expect(dto.Connection.Auth).ToNot(BeNil())
			Expect(dto.Connection.Auth.Type).To(Equal("NONE"))
			Expect(dto.Connection.Auth.Name).To(BeNil())
			Expect(dto.Connection.Auth.ClientID).To(BeNil())
		})

		It("flattens a HEADER auth", func() {
			dto := model.ToCatalogMcpServerDTO(newCatalogMcpServer(&catalogmcpserver.Auth{
				Type:   catalogmcpserver.AuthTypeHeader,
				Header: &catalogmcpserver.HeaderAuth{Name: "Authorization", Value: "Bearer ghp"},
			}))
			Expect(dto.Connection.Auth.Type).To(Equal("HEADER"))
			Expect(*dto.Connection.Auth.Name).To(Equal("Authorization"))
			Expect(*dto.Connection.Auth.Value).To(Equal("Bearer ghp"))
			Expect(dto.Connection.Auth.ClientID).To(BeNil())
			Expect(dto.Connection.Auth.TokenURL).To(BeNil())
		})

		It("flattens an OAUTH2 auth", func() {
			dto := model.ToCatalogMcpServerDTO(newCatalogMcpServer(&catalogmcpserver.Auth{
				Type: catalogmcpserver.AuthTypeOAuth2,
				OAuth2: &catalogmcpserver.OAuth2Auth{
					ClientID:     "gravitee-catalog",
					ClientSecret: "s3cr3t",
					TokenURL:     "https://auth.example.com/oauth/token",
					Scope:        new("mcp:read"),
				},
			}))
			Expect(dto.Connection.Auth.Type).To(Equal("OAUTH2"))
			Expect(*dto.Connection.Auth.ClientID).To(Equal("gravitee-catalog"))
			Expect(*dto.Connection.Auth.ClientSecret).To(Equal("s3cr3t"))
			Expect(*dto.Connection.Auth.TokenURL).To(Equal("https://auth.example.com/oauth/token"))
			Expect(*dto.Connection.Auth.Scope).To(Equal("mcp:read"))
			Expect(dto.Connection.Auth.Name).To(BeNil())
			Expect(dto.Connection.Auth.Value).To(BeNil())
		})

		It("serialises the payload with the wire field names", func() {
			dto := model.ToCatalogMcpServerDTO(newCatalogMcpServer(&catalogmcpserver.Auth{
				Type:   catalogmcpserver.AuthTypeHeader,
				Header: &catalogmcpserver.HeaderAuth{Name: "Authorization", Value: "Bearer ghp"},
			}))
			raw, err := json.Marshal(dto)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(raw)).To(MatchJSON(`{
				"hrid": "gravitee-github",
				"entityId": "mcp-server.github",
				"connection": {
					"endpoint": "https://api.githubcopilot.com/mcp/",
					"transport": "HTTP",
					"auth": {"type": "HEADER", "name": "Authorization", "value": "Bearer ghp"}
				}
			}`))
		})
	})

	Describe("state decoding", func() {
		It("decodes the payload, the ids and what the platform discovered", func() {
			body := `{
				"hrid": "gravitee-github",
				"entityId": "mcp-server.github",
				"connection": {"endpoint": "https://api.githubcopilot.com/mcp/", "transport": "HTTP", "auth": {"type": "HEADER", "name": "Authorization"}},
				"id": "8c21f0a2", "organizationId": "DEFAULT", "environmentId": "DEFAULT",
				"errors": {"warning": ["something to know"]},
				"lastSyncedAt": "2026-09-23T09:12:04Z",
				"protocolVersion": "2025-03-26",
				"serverInfo": {"name": "github-mcp", "version": "1.4.0"},
				"tools": [{"name": "create_issue", "entityId": "mcp-tool.github.create_issue"}],
				"prompts": [{"name": "summarize_pr", "entityId": "mcp-prompt.github.summarize_pr"}],
				"resources": [{"uri": "github://repos", "name": "repositories", "entityId": "mcp-resource.github.repositories"}]
			}`
			state := new(model.CatalogMcpServerState)
			Expect(json.Unmarshal([]byte(body), state)).To(Succeed())

			Expect(state.HRID).To(Equal("gravitee-github"))
			Expect(state.EntityID).To(Equal("mcp-server.github"))
			Expect(state.Connection.Auth.Value).To(BeNil())
			Expect(state.ID).To(Equal("8c21f0a2"))
			Expect(state.OrgID).To(Equal("DEFAULT"))
			Expect(state.Errors.Warning).To(ConsistOf("something to know"))
			Expect(state.LastSyncedAt).To(Equal("2026-09-23T09:12:04Z"))
			Expect(state.ProtocolVersion).To(Equal("2025-03-26"))
			Expect(state.ServerInfo.Name).To(Equal("github-mcp"))
			Expect(state.Tools).To(HaveLen(1))
			Expect(state.Tools[0].EntityID).To(Equal("mcp-tool.github.create_issue"))
			Expect(state.Prompts[0].Name).To(Equal("summarize_pr"))
			Expect(state.Resources[0].URI).To(Equal("github://repos"))
		})
	})

	Describe("refusal findings", func() {
		refusedWith := func(statusCode int, body string) error {
			return xErrors.ServerError{StatusCode: statusCode, Body: body}
		}

		It("reads the findings a refused apply carries in errors.severe", func() {
			err := refusedWith(http.StatusBadRequest, `{
				"entityId": "mcp-server.github",
				"connection": {"endpoint": "https://api.githubcopilot.com/mcp/", "transport": "HTTP"},
				"errors": {"severe": ["upstream unreachable", "entityId already used"], "warning": ["slow upstream"]}
			}`)

			findings, refused := model.CatalogMcpServerRefusal(err)

			Expect(refused).To(BeTrue())
			Expect(findings.Severe).To(ConsistOf("upstream unreachable", "entityId already used"))
			Expect(findings.Warning).To(ConsistOf("slow upstream"))
		})

		It("reads the findings through a wrapped error", func() {
			err := fmt.Errorf("apply failed: %w",
				refusedWith(http.StatusBadRequest, `{"errors": {"severe": ["upstream unreachable"]}}`))

			findings, refused := model.CatalogMcpServerRefusal(err)

			Expect(refused).To(BeTrue())
			Expect(findings.Severe).To(ConsistOf("upstream unreachable"))
		})

		DescribeTable("is not a refusal",
			func(err error) {
				_, refused := model.CatalogMcpServerRefusal(err)
				Expect(refused).To(BeFalse())
			},
			Entry("a 400 without severe findings",
				refusedWith(http.StatusBadRequest, `{"errors": {"warning": ["slow upstream"]}}`)),
			Entry("a 400 in the platform error shape, a missing hrid for instance",
				refusedWith(http.StatusBadRequest, `{"message": "hrid is required", "http_status": 400}`)),
			Entry("a 400 whose body is not JSON",
				refusedWith(http.StatusBadRequest, `Bad Request`)),
			Entry("another status carrying a state body",
				refusedWith(http.StatusInternalServerError, `{"errors": {"severe": ["boom"]}}`)),
			Entry("an error that is not a server answer",
				fmt.Errorf("connection refused")),
		)
	})
})
