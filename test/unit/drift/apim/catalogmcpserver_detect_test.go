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
package apim

import (
	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	apimcatalog "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/catalogmcpserver"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var catalogNoAuth = &sdk.AimCatalogMcpServerAuth{
	None: &sdk.AimCatalogMcpServerNoAuth{Type: sdk.AimCatalogMcpServerNoAuthTypeNONE},
}

var _ = Describe("Catalog MCP server Drift detection", func() {

	DescribeTable("equivalent values",
		func(crd, remote any) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
		},
		Entry("empty struct",
			sdk.AimCatalogMcpServer{},
			sdk.AimCatalogMcpServer{},
		),
		Entry("equal struct",
			completeCatalogMcpServerDTO(),
			completeCatalogMcpServerDTO(),
		),
		Entry("HEADER auth: the platform never returns the header value",
			catalogMcpServerDTO(headerCatalogMcpServerCRD()),
			withoutCredentials(catalogMcpServerDTO(headerCatalogMcpServerCRD())),
		),
		Entry("OAUTH2 auth: the platform never returns the client secret",
			catalogMcpServerDTO(oauth2CatalogMcpServerCRD()),
			withoutCredentials(catalogMcpServerDTO(oauth2CatalogMcpServerCRD())),
		),
		Entry("OAUTH2 auth: empty scope equivalent to nil",
			withScope(catalogMcpServerDTO(oauth2CatalogMcpServerCRD()), nil),
			withScope(catalogMcpServerDTO(oauth2CatalogMcpServerCRD()), new("")),
		),
		Entry("no auth in the spec reads back as NONE",
			catalogMcpServerDTO(noAuthCatalogMcpServerCRD()),
			withAuth(catalogMcpServerDTO(noAuthCatalogMcpServerCRD()), catalogNoAuth),
		),
		Entry("empty description equivalent to nil",
			withDescription(catalogMcpServerDTO(noAuthCatalogMcpServerCRD()), nil),
			withDescription(catalogMcpServerDTO(noAuthCatalogMcpServerCRD()), new("")),
		),
		Entry("hrid is identity, not payload",
			withHRID(completeCatalogMcpServerDTO(), "default-drift-full"),
			withHRID(completeCatalogMcpServerDTO(), "something-else"),
		),
	)

	DescribeTable("drifted values",
		func(crd, remote any) {
			result := drift.DetectWithNamespace(crd, remote, "")
			Expect(result.DriftDetected()).To(BeTrue())
		},
		Entry("endpoint changed on the platform",
			catalogMcpServerDTO(headerCatalogMcpServerCRD()),
			withEndpoint(withoutCredentials(catalogMcpServerDTO(headerCatalogMcpServerCRD())),
				"https://mcp.elsewhere.example.com/"),
		),
		Entry("auth type changed on the platform",
			catalogMcpServerDTO(headerCatalogMcpServerCRD()),
			withAuth(catalogMcpServerDTO(headerCatalogMcpServerCRD()), catalogNoAuth),
		),
		Entry("OAUTH2 scope changed on the platform",
			catalogMcpServerDTO(oauth2CatalogMcpServerCRD()),
			withScope(withoutCredentials(catalogMcpServerDTO(oauth2CatalogMcpServerCRD())), new("mcp:write")),
		),
	)

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(completeCatalogMcpServerDTO(), completeCatalogMcpServerDTO())
		})
	})
})

// completeCatalogMcpServerDTO sets every field, every auth variant included. JSON carries one variant of
// the union, so the fixture holds the OAUTH2 one and the others are set here.
func completeCatalogMcpServerDTO() sdk.AimCatalogMcpServer {
	GinkgoHelper()
	dto := loadFixture[sdk.AimCatalogMcpServer]("catalogmcpserver_full_dto.json")
	auth := *dto.Connection.Auth
	auth.None = catalogNoAuth.None
	auth.Header = &sdk.AimCatalogMcpServerHeaderAuth{
		Type:  sdk.AimCatalogMcpServerHeaderAuthTypeHEADER,
		Name:  "Authorization",
		Value: new("Bearer token"),
	}
	dto.Connection.Auth = &auth
	return dto
}

func catalogMcpServerDTO(srv *v1alpha1.CatalogMcpServer) sdk.AimCatalogMcpServer {
	GinkgoHelper()
	dto, err := apimcatalog.ToDTO(srv)
	Expect(err).ToNot(HaveOccurred())
	return dto
}

func headerCatalogMcpServerCRD() *v1alpha1.CatalogMcpServer {
	GinkgoHelper()
	fixture := loadFixture[v1alpha1.CatalogMcpServer]("catalogmcpserver_header_crd.json")
	return &fixture
}

func oauth2CatalogMcpServerCRD() *v1alpha1.CatalogMcpServer {
	GinkgoHelper()
	fixture := loadFixture[v1alpha1.CatalogMcpServer]("catalogmcpserver_oauth2_crd.json")
	return &fixture
}

func noAuthCatalogMcpServerCRD() *v1alpha1.CatalogMcpServer {
	GinkgoHelper()
	fixture := loadFixture[v1alpha1.CatalogMcpServer]("catalogmcpserver_header_crd.json")
	fixture.Spec.Connection.Auth = nil
	return &fixture
}

// withoutCredentials is what the platform answers: the payload with every credential removed.
func withoutCredentials(dto sdk.AimCatalogMcpServer) sdk.AimCatalogMcpServer {
	if dto.Connection.Auth == nil {
		return dto
	}
	auth := *dto.Connection.Auth
	if auth.Header != nil {
		header := *auth.Header
		header.Value = nil
		auth.Header = &header
	}
	if auth.OAuth2 != nil {
		oauth2 := *auth.OAuth2
		oauth2.ClientSecret = nil
		auth.OAuth2 = &oauth2
	}
	dto.Connection.Auth = &auth
	return dto
}

func withAuth(dto sdk.AimCatalogMcpServer, auth *sdk.AimCatalogMcpServerAuth) sdk.AimCatalogMcpServer {
	dto.Connection.Auth = auth
	return dto
}

func withScope(dto sdk.AimCatalogMcpServer, scope *string) sdk.AimCatalogMcpServer {
	auth := *dto.Connection.Auth
	oauth2 := *auth.OAuth2
	oauth2.Scope = scope
	auth.OAuth2 = &oauth2
	dto.Connection.Auth = &auth
	return dto
}

func withDescription(dto sdk.AimCatalogMcpServer, description *string) sdk.AimCatalogMcpServer {
	dto.Description = description
	return dto
}

func withEndpoint(dto sdk.AimCatalogMcpServer, endpoint string) sdk.AimCatalogMcpServer {
	dto.Connection.Endpoint = endpoint
	return dto
}

func withHRID(dto sdk.AimCatalogMcpServer, hrid string) sdk.AimCatalogMcpServer {
	dto.Hrid = hrid
	return dto
}
