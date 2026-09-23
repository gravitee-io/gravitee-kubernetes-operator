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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Catalog MCP server Drift detection", func() {

	DescribeTable("equivalent values",
		func(crd, remote any) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
		},
		Entry("empty struct",
			model.CatalogMcpServerDTO{},
			model.CatalogMcpServerDTO{},
		),
		Entry("equal struct",
			completeCatalogMcpServerDTO(),
			completeCatalogMcpServerDTO(),
		),
		Entry("HEADER auth: the platform never returns the header value",
			model.ToCatalogMcpServerDTO(headerCatalogMcpServerCRD()),
			withoutCredentials(model.ToCatalogMcpServerDTO(headerCatalogMcpServerCRD())),
		),
		Entry("OAUTH2 auth: the platform never returns the client secret",
			model.ToCatalogMcpServerDTO(oauth2CatalogMcpServerCRD()),
			withoutCredentials(model.ToCatalogMcpServerDTO(oauth2CatalogMcpServerCRD())),
		),
		Entry("no auth in the spec reads back as NONE",
			model.ToCatalogMcpServerDTO(noAuthCatalogMcpServerCRD()),
			withAuth(withoutCredentials(model.ToCatalogMcpServerDTO(noAuthCatalogMcpServerCRD())),
				&model.CatalogMcpServerAuthDTO{Type: "NONE"}),
		),
		Entry("empty description equivalent to nil",
			withDescription(model.ToCatalogMcpServerDTO(noAuthCatalogMcpServerCRD()), nil),
			withDescription(model.ToCatalogMcpServerDTO(noAuthCatalogMcpServerCRD()), new("")),
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
			model.ToCatalogMcpServerDTO(headerCatalogMcpServerCRD()),
			withEndpoint(withoutCredentials(model.ToCatalogMcpServerDTO(headerCatalogMcpServerCRD())),
				"https://mcp.elsewhere.example.com/"),
		),
		Entry("auth type changed on the platform",
			model.ToCatalogMcpServerDTO(headerCatalogMcpServerCRD()),
			withAuth(model.ToCatalogMcpServerDTO(headerCatalogMcpServerCRD()),
				&model.CatalogMcpServerAuthDTO{Type: "NONE"}),
		),
	)

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(completeCatalogMcpServerDTO(), completeCatalogMcpServerDTO())
		})
	})
})

func completeCatalogMcpServerDTO() model.CatalogMcpServerDTO {
	GinkgoHelper()
	return loadFixture[model.CatalogMcpServerDTO]("catalogmcpserver_full_dto.json")
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
func withoutCredentials(dto model.CatalogMcpServerDTO) model.CatalogMcpServerDTO {
	if dto.Connection.Auth != nil {
		auth := *dto.Connection.Auth
		auth.Value = nil
		auth.ClientSecret = nil
		dto.Connection.Auth = &auth
	}
	return dto
}

func withAuth(dto model.CatalogMcpServerDTO, auth *model.CatalogMcpServerAuthDTO) model.CatalogMcpServerDTO {
	dto.Connection.Auth = auth
	return dto
}

func withDescription(dto model.CatalogMcpServerDTO, description *string) model.CatalogMcpServerDTO {
	dto.Description = description
	return dto
}

func withEndpoint(dto model.CatalogMcpServerDTO, endpoint string) model.CatalogMcpServerDTO {
	dto.Connection.Endpoint = endpoint
	return dto
}

func withHRID(dto model.CatalogMcpServerDTO, hrid string) model.CatalogMcpServerDTO {
	dto.HRID = hrid
	return dto
}
