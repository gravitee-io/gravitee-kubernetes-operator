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
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/catalogmcpserver"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimcatalog "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/catalogmcpserver"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	xErrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
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

func headerAuth() *catalogmcpserver.Auth {
	return &catalogmcpserver.Auth{
		Type:   catalogmcpserver.AuthTypeHeader,
		Header: &catalogmcpserver.HeaderAuth{Name: "Authorization", Value: "Bearer ghp"},
	}
}

func toCatalogMcpServerDTO(srv *v1alpha1.CatalogMcpServer) sdk.AimCatalogMcpServer {
	GinkgoHelper()
	dto, err := apimcatalog.ToDTO(srv)
	Expect(err).ToNot(HaveOccurred())
	return dto
}

func marshalJSON(v any) string {
	GinkgoHelper()
	raw, err := json.Marshal(v)
	Expect(err).ToNot(HaveOccurred())
	return string(raw)
}

// catalogMcpServerState is what the platform answers for the github server: the payload without its
// credentials, the ids, the findings and what it discovered.
const catalogMcpServerState = `{
	"hrid": "gravitee-github",
	"entityId": "mcp-server.github",
	"connection": {"endpoint": "https://api.githubcopilot.com/mcp/", "transport": "HTTP",
		"auth": {"type": "HEADER", "name": "Authorization"}},
	"id": "8c21f0a2", "organizationId": "DEFAULT", "environmentId": "DEFAULT",
	"authPolicyRef": "MCPServer::\"github\"",
	"errors": {"warning": ["something to know"]},
	"lastSyncedAt": "2026-09-23T09:12:04Z",
	"protocolVersion": "2025-03-26",
	"serverInfo": {"name": "github-mcp", "version": "1.4.0"},
	"tools": [{"name": "create_issue", "entityId": "mcp-tool.github.create_issue",
		"authPolicyRef": "MCPTool::\"github.create_issue\""}],
	"prompts": [{"name": "summarize_pr", "entityId": "mcp-prompt.github.summarize_pr"}],
	"resources": [{"uri": "github://repos", "name": "repositories", "entityId": "mcp-resource.github.repositories"}]
}`

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
			dto := toCatalogMcpServerDTO(newCatalogMcpServer(nil))
			Expect(dto.Hrid).To(Equal("gravitee-github"))
			Expect(dto.EntityId).To(Equal("mcp-server.github"))
			Expect(dto.Connection.Endpoint).To(Equal("https://api.githubcopilot.com/mcp/"))
		})

		It("defaults the transport to HTTP", func() {
			dto := toCatalogMcpServerDTO(newCatalogMcpServer(nil))
			Expect(dto.Connection.Transport).To(Equal(sdk.AimCatalogMcpServerTransportHTTP))
		})

		DescribeTable("sends the auth flat, with its type",
			func(auth *catalogmcpserver.Auth, expected string) {
				dto := toCatalogMcpServerDTO(newCatalogMcpServer(auth))
				Expect(marshalJSON(dto.Connection.Auth)).To(MatchJSON(expected))
			},
			Entry("an explicit NONE when the spec has none", nil, `{"type": "NONE"}`),
			Entry("NONE", &catalogmcpserver.Auth{Type: catalogmcpserver.AuthTypeNone}, `{"type": "NONE"}`),
			Entry("HEADER", headerAuth(), `{"type": "HEADER", "name": "Authorization", "value": "Bearer ghp"}`),
			Entry("OAUTH2", &catalogmcpserver.Auth{
				Type: catalogmcpserver.AuthTypeOAuth2,
				OAuth2: &catalogmcpserver.OAuth2Auth{
					ClientID:     "gravitee-catalog",
					ClientSecret: "s3cr3t",
					TokenURL:     "https://auth.example.com/oauth/token",
					Scope:        new("mcp:read"),
				},
			}, `{"type": "OAUTH2", "clientId": "gravitee-catalog", "clientSecret": "s3cr3t",
				"tokenUrl": "https://auth.example.com/oauth/token", "scope": "mcp:read"}`),
			Entry("OAUTH2 without a scope", &catalogmcpserver.Auth{
				Type: catalogmcpserver.AuthTypeOAuth2,
				OAuth2: &catalogmcpserver.OAuth2Auth{
					ClientID:     "gravitee-catalog",
					ClientSecret: "s3cr3t",
					TokenURL:     "https://auth.example.com/oauth/token",
				},
			}, `{"type": "OAUTH2", "clientId": "gravitee-catalog", "clientSecret": "s3cr3t",
				"tokenUrl": "https://auth.example.com/oauth/token"}`),
		)

		DescribeTable("refuses an auth it cannot send",
			func(auth *catalogmcpserver.Auth) {
				_, err := apimcatalog.ToDTO(newCatalogMcpServer(auth))
				Expect(err).To(HaveOccurred())
			},
			Entry("HEADER without its header", &catalogmcpserver.Auth{Type: catalogmcpserver.AuthTypeHeader}),
			Entry("OAUTH2 without its oauth2", &catalogmcpserver.Auth{Type: catalogmcpserver.AuthTypeOAuth2}),
			Entry("an unknown type", &catalogmcpserver.Auth{Type: "BASIC"}),
		)

		It("serialises the payload with the wire field names", func() {
			dto := toCatalogMcpServerDTO(newCatalogMcpServer(headerAuth()))
			Expect(marshalJSON(dto)).To(MatchJSON(`{
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

	Describe("status", func() {
		It("copies the ids, the findings, what the platform discovered and the policy references", func() {
			state := sdk.AimCatalogMcpServerState{}
			Expect(json.Unmarshal([]byte(catalogMcpServerState), &state)).To(Succeed())
			srv := newCatalogMcpServer(nil)

			Expect(apimcatalog.UpdateStatus(context.Background(), srv, apimcatalog.Response{State: state})).To(Succeed())

			Expect(srv.Status.ID).To(Equal("8c21f0a2"))
			Expect(srv.Status.OrgID).To(Equal("DEFAULT"))
			Expect(srv.Status.EnvID).To(Equal("DEFAULT"))
			Expect(srv.Status.HRID).To(Equal("gravitee-github"))
			Expect(srv.Status.Errors.Warning).To(ConsistOf("something to know"))
			Expect(srv.Status.Errors.Severe).To(BeEmpty())
			Expect(srv.Status.AuthPolicyRef).To(Equal(`MCPServer::"github"`))
			Expect(srv.Status.Discovered).To(Equal(catalogmcpserver.Discovered{
				LastSyncedAt:    "2026-09-23T09:12:04Z",
				ProtocolVersion: "2025-03-26",
				ServerInfo:      &catalogmcpserver.ServerInfo{Name: "github-mcp", Version: "1.4.0"},
				Tools: []catalogmcpserver.Tool{{
					Name:          "create_issue",
					EntityID:      "mcp-tool.github.create_issue",
					AuthPolicyRef: `MCPTool::"github.create_issue"`,
				}},
				Prompts: []catalogmcpserver.Prompt{{Name: "summarize_pr", EntityID: "mcp-prompt.github.summarize_pr"}},
				Resources: []catalogmcpserver.Resource{{
					URI: "github://repos", Name: "repositories", EntityID: "mcp-resource.github.repositories",
				}},
			}))
			Expect(k8s.IsAutomationAPIManaged(srv)).To(BeTrue())
		})

		It("leaves out what the state leaves out", func() {
			srv := newCatalogMcpServer(nil)
			srv.Status.AuthPolicyRef = `MCPServer::"stale"`

			resp := apimcatalog.Response{State: sdk.AimCatalogMcpServerState{Id: new("8c21f0a2")}}
			Expect(apimcatalog.UpdateStatus(context.Background(), srv, resp)).To(Succeed())

			Expect(srv.Status.ID).To(Equal("8c21f0a2"))
			Expect(srv.Status.AuthPolicyRef).To(BeEmpty())
			Expect(srv.Status.Errors.Severe).To(BeNil())
			Expect(srv.Status.Errors.Warning).To(BeNil())
			Expect(srv.Status.Discovered).To(Equal(catalogmcpserver.Discovered{}))
		})
	})

	Describe("calls", func() {
		const collection = "/automation/organizations/org1/environments/env1/aim/catalog/mcp-servers"

		var (
			ctx    context.Context
			fake   *fakeAutomationAPI
			client *apim.SDKClient
		)

		BeforeEach(func() {
			ctx = context.Background()
			fake = newFakeAutomationAPI()
			fake.response = catalogMcpServerState
			DeferCleanup(fake.Close)
			var err error
			client, err = apim.NewSDKClient(ctx, managementContext(fake.URL, bearer("token")))
			Expect(err).ToNot(HaveOccurred())
		})

		It("puts the payload on the catalog collection and answers the state", func() {
			resp, err := apimcatalog.Upsert(ctx, client, toCatalogMcpServerDTO(newCatalogMcpServer(headerAuth())))
			Expect(err).ToNot(HaveOccurred())

			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodPut))
			Expect(req.Path).To(Equal(collection))
			Expect(req.Query.Get("dryRun")).To(Equal("false"))
			Expect(marshalJSON(req.Body)).To(MatchJSON(`{
				"hrid": "gravitee-github",
				"entityId": "mcp-server.github",
				"connection": {
					"endpoint": "https://api.githubcopilot.com/mcp/",
					"transport": "HTTP",
					"auth": {"type": "HEADER", "name": "Authorization", "value": "Bearer ghp"}
				}
			}`))
			Expect(*resp.State.Id).To(Equal("8c21f0a2"))
			Expect(resp.GetOrgID()).To(Equal("DEFAULT"))
			Expect(resp.GetEnvID()).To(Equal("DEFAULT"))
		})

		It("dry runs and returns the findings of the state", func() {
			fake.response = `{"hrid": "gravitee-github", "entityId": "mcp-server.github",
				"connection": {"endpoint": "https://api.githubcopilot.com/mcp/", "transport": "HTTP"},
				"errors": {"severe": ["entityId already used"], "warning": ["slow upstream"]}}`

			findings, err := apimcatalog.DryRun(ctx, client, toCatalogMcpServerDTO(newCatalogMcpServer(nil)))
			Expect(err).ToNot(HaveOccurred())

			Expect(fake.lastRequest().Query.Get("dryRun")).To(Equal("true"))
			Expect(findings.Severe).To(HaveLen(1))
			Expect(findings.Severe[0].Message).To(Equal("entityId already used"))
			Expect(findings.Warning).To(HaveLen(1))
			Expect(findings.Warning[0].Message).To(Equal("slow upstream"))
		})

		It("returns the findings of a refused dry run", func() {
			fake.status = http.StatusBadRequest
			fake.errorResponse = `{"entityId": "mcp-server.github",
				"errors": {"severe": ["upstream unreachable"]}}`

			findings, err := apimcatalog.DryRun(ctx, client, toCatalogMcpServerDTO(newCatalogMcpServer(nil)))
			Expect(err).ToNot(HaveOccurred())

			Expect(findings.Severe).To(HaveLen(1))
			Expect(findings.Severe[0].Message).To(Equal("upstream unreachable"))
		})

		It("returns a 400 that carries no findings as the error", func() {
			fake.status = http.StatusBadRequest
			fake.errorResponse = `{"message": "hrid is required", "http_status": 400}`

			_, err := apimcatalog.DryRun(ctx, client, toCatalogMcpServerDTO(newCatalogMcpServer(nil)))

			Expect(xErrors.IsBadRequest(err)).To(BeTrue())
		})

		It("gets the remote as the payload a PUT sends", func() {
			remote, err := apimcatalog.GetRemote(ctx, client, toCatalogMcpServerDTO(newCatalogMcpServer(nil)))
			Expect(err).ToNot(HaveOccurred())

			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodGet))
			Expect(req.Path).To(Equal(collection + "/gravitee-github"))
			Expect(remote).To(Equal(sdk.AimCatalogMcpServer{
				Hrid:     "gravitee-github",
				EntityId: "mcp-server.github",
				Connection: sdk.AimCatalogMcpServerConnection{
					Endpoint:  "https://api.githubcopilot.com/mcp/",
					Transport: sdk.AimCatalogMcpServerTransportHTTP,
					Auth: &sdk.AimCatalogMcpServerAuth{Header: &sdk.AimCatalogMcpServerHeaderAuth{
						Type: sdk.AimCatalogMcpServerHeaderAuthTypeHEADER,
						Name: "Authorization",
					}},
				},
			}))
		})

		It("reports a remote that is gone as not found", func() {
			fake.status = http.StatusNotFound

			_, err := apimcatalog.GetRemote(ctx, client, toCatalogMcpServerDTO(newCatalogMcpServer(nil)))

			Expect(xErrors.IsNotFound(err)).To(BeTrue())
		})

		It("deletes by hrid and treats a server already gone as deleted", func() {
			dto := toCatalogMcpServerDTO(newCatalogMcpServer(nil))

			fake.status = http.StatusNoContent
			Expect(apimcatalog.Delete(ctx, client, dto)).To(Succeed())
			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodDelete))
			Expect(req.Path).To(Equal(collection + "/gravitee-github"))

			fake.status = http.StatusNotFound
			Expect(apimcatalog.Delete(ctx, client, dto)).To(Succeed())
		})
	})

	Describe("context reference", func() {
		It("reports an unset context as not found, so a delete releases the finalizer", func() {
			srv := newCatalogMcpServer(nil)
			srv.Spec.Context = nil

			err := apim.ResolveContextRef(context.Background(), srv, srv.Namespace)

			Expect(apierrors.IsNotFound(err)).To(BeTrue())
			Expect(err).To(MatchError("a management context reference (spec.contextRef) is required"))
		})
	})

	Describe("a refused apply", func() {
		refusedWith := func(statusCode int, body string) error {
			return xErrors.ServerError{StatusCode: statusCode, Body: body}
		}

		It("records the findings in the status and stays recoverable", func() {
			srv := newCatalogMcpServer(nil)

			err := apimcatalog.OnUpsertError(context.Background(), srv,
				refusedWith(http.StatusBadRequest, `{"errors": {"severe": ["upstream unreachable"]}}`))

			Expect(err).To(MatchError(ContainSubstring("catalog mcp server [github] was refused by the automation api")))
			Expect(xErrors.IsRecoverable(err)).To(BeTrue())
			Expect(srv.Status.Errors.Severe).To(ConsistOf("upstream unreachable"))
		})

		It("stays recoverable when the 400 carries no findings", func() {
			srv := newCatalogMcpServer(nil)

			err := apimcatalog.OnUpsertError(context.Background(), srv,
				refusedWith(http.StatusBadRequest, `{"message": "hrid is required", "http_status": 400}`))

			Expect(xErrors.IsRecoverable(err)).To(BeTrue())
			Expect(srv.Status.Errors.Severe).To(BeEmpty())
		})

		It("leaves any other failure to the lifecycle", func() {
			err := apimcatalog.OnUpsertError(context.Background(), newCatalogMcpServer(nil),
				refusedWith(http.StatusInternalServerError, `{"errors": {"severe": ["boom"]}}`))

			Expect(err).ToNot(HaveOccurred())
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

			findings, refused := model.AutomationRefusal(err)

			Expect(refused).To(BeTrue())
			Expect(findings.Severe).To(ConsistOf("upstream unreachable", "entityId already used"))
			Expect(findings.Warning).To(ConsistOf("slow upstream"))
		})

		It("reads the findings through a wrapped error", func() {
			err := fmt.Errorf("apply failed: %w",
				refusedWith(http.StatusBadRequest, `{"errors": {"severe": ["upstream unreachable"]}}`))

			findings, refused := model.AutomationRefusal(err)

			Expect(refused).To(BeTrue())
			Expect(findings.Severe).To(ConsistOf("upstream unreachable"))
		})

		DescribeTable("is not a refusal",
			func(err error) {
				_, refused := model.AutomationRefusal(err)
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
