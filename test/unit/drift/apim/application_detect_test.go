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

package apim

import (
	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	apimapp "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Application Drift detection", func() {

	DescribeTable("equivalent values",
		func(crd, remote any) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, "namespace"))
		},
		Entry("empty struct",
			apimapp.DTO{},
			apimapp.DTO{},
		),
		Entry("equal struct",
			completeApplicationDTO(),
			completeApplicationDTO(),
		),
		Entry("equal struct",
			completeApplicationDTO(),
			toApplicationDTO(completeApplicationCRD()),
		),
		Entry("crd-only members ignored",
			apimapp.DTO{Application: sdk.Application{
				Members: []sdk.Member{
					{Source: "memory", SourceId: "admin", Role: "ADMIN"},
					{Source: "memory", SourceId: "local-only", Role: "USER"},
				},
			}},
			apimapp.DTO{Application: sdk.Application{
				Members: []sdk.Member{
					{Source: "memory", SourceId: "admin", Role: "ADMIN"},
				},
			}},
		),
		Entry("equivalent struct",
			apimapp.DTO{Application: sdk.Application{
				Members: nil,
				Groups:  []string{"foo", "bar"},
			}},
			apimapp.DTO{Application: sdk.Application{
				Hrid:          "my-app",
				Status:        new(sdk.ApplicationSpecStatus("ACTIVE")),
				Groups:        []string{"namespace-bar"},
				Members:       make([]sdk.Member, 0),
				Metadata:      make([]sdk.Metadata, 0),
				Settings:      &sdk.ApplicationSettings{},
				NotifyMembers: new(false),
			}},
		),
		Entry("client certificate equivalent",
			apimapp.DTO{Application: sdk.Application{
				Settings: &sdk.ApplicationSettings{
					Tls: &sdk.ApplicationTLSSettings{
						ClientCertificate: new(`-----BEGIN CERTIFICATE-----
-----END CERTIFICATE-----`),
					},
				},
			}},
			apimapp.DTO{Application: sdk.Application{
				Settings: &sdk.ApplicationSettings{
					Tls: &sdk.ApplicationTLSSettings{
						ClientCertificate: new(`-----BEGIN CERTIFICATE-----
-----END CERTIFICATE-----
`),
					},
				},
			}},
		),
		Entry("listed certificates equivalent",
			apimapp.DTO{Application: sdk.Application{
				Settings: &sdk.ApplicationSettings{
					Tls: &sdk.ApplicationTLSSettings{
						ClientCertificates: []sdk.ClientCertificate{
							{
								Content: `-----BEGIN CERTIFICATE-----
-----END CERTIFICATE-----`,
							},
						},
					},
				},
			}},
			apimapp.DTO{Application: sdk.Application{
				Settings: &sdk.ApplicationSettings{
					Tls: &sdk.ApplicationTLSSettings{
						ClientCertificates: []sdk.ClientCertificate{
							{
								Content: `-----BEGIN CERTIFICATE-----
-----END CERTIFICATE-----
`,
							},
						},
					},
				},
			}},
		),
		Entry("listed certificates equivalent, ref ignored",
			toApplicationDTO(&v1alpha1.Application{Spec: v1alpha1.ApplicationSpec{Application: application.Application{
				Settings: &application.Setting{
					TLS: &application.TLSSettings{
						ClientCertificates: []application.ClientCertificate{
							{
								Content: `-----BEGIN CERTIFICATE-----
-----END CERTIFICATE-----`,
								Ref: &application.CertificateRef{Name: "client-cert", Key: "tls.crt"},
							},
						},
					},
				},
			}}}),
			apimapp.DTO{Application: sdk.Application{
				Settings: &sdk.ApplicationSettings{
					Tls: &sdk.ApplicationTLSSettings{
						ClientCertificates: []sdk.ClientCertificate{
							{
								Content: `-----BEGIN CERTIFICATE-----
-----END CERTIFICATE-----
`,
							},
						},
					},
				},
			}},
		),
		Entry("expired cert omitted from remote is not drift",
			apimapp.DTO{Application: sdk.Application{
				Settings: &sdk.ApplicationSettings{
					Tls: &sdk.ApplicationTLSSettings{
						ClientCertificates: []sdk.ClientCertificate{
							{Name: "expired", Content: "CERT-A", EndsAt: new("2000-01-01T00:00:00Z")},
							{Name: "active", Content: "CERT-B", EndsAt: new("2099-01-01T00:00:00Z")},
						},
					},
				},
			}},
			apimapp.DTO{Application: sdk.Application{
				Settings: &sdk.ApplicationSettings{
					Tls: &sdk.ApplicationTLSSettings{
						ClientCertificates: []sdk.ClientCertificate{
							{Name: "active", Content: "CERT-B", EndsAt: new("2099-01-01T00:00:00Z")},
						},
					},
				},
			}},
		),
		Entry("scheduled cert omitted from remote is not drift",
			apimapp.DTO{Application: sdk.Application{
				Settings: &sdk.ApplicationSettings{
					Tls: &sdk.ApplicationTLSSettings{
						ClientCertificates: []sdk.ClientCertificate{
							{Name: "scheduled", Content: "CERT-A", StartsAt: new("2099-01-01T00:00:00Z")},
						},
					},
				},
			}},
			apimapp.DTO{Application: sdk.Application{
				Settings: &sdk.ApplicationSettings{
					Tls: &sdk.ApplicationTLSSettings{
						ClientCertificates: []sdk.ClientCertificate{},
					},
				},
			}},
		),
		Entry("listed certificates equivalent, start/end date",
			apimapp.DTO{Application: sdk.Application{
				Settings: &sdk.ApplicationSettings{
					Tls: &sdk.ApplicationTLSSettings{
						ClientCertificates: []sdk.ClientCertificate{
							{
								StartsAt: new("2023-07-25T02:43:16+03:00"),
								EndsAt:   new("2099-08-25T23:43:16Z"),
							},
						},
					},
				},
			}},
			apimapp.DTO{Application: sdk.Application{
				Settings: &sdk.ApplicationSettings{
					Tls: &sdk.ApplicationTLSSettings{
						ClientCertificates: []sdk.ClientCertificate{
							{
								StartsAt: new("2023-07-24T20:43:16-03:00"),
								EndsAt:   new("2099-08-25T23:43:16-00:00"),
							},
						},
					},
				},
			}},
		),
	)

	It("detects a console-side change of the OAuth additional client metadata", func() {
		withMetadata := func(softwareID string) apimapp.DTO {
			return apimapp.DTO{Application: sdk.Application{Settings: &sdk.ApplicationSettings{
				Oauth: &sdk.ApplicationOAuthClientSettings{
					AdditionalClientMetadata: map[string]string{"software_id": softwareID},
				},
			}}}
		}

		result := drift.DetectWithNamespace(withMetadata("crd-template"), withMetadata("console-template"), "")

		Expect(result.DriftDetected()).To(BeTrue())
	})

	It("treats absent and empty OAuth additional client metadata as equivalent", func() {
		crd := apimapp.DTO{Application: sdk.Application{Settings: &sdk.ApplicationSettings{
			Oauth: &sdk.ApplicationOAuthClientSettings{},
		}}}
		remote := apimapp.DTO{Application: sdk.Application{Settings: &sdk.ApplicationSettings{
			Oauth: &sdk.ApplicationOAuthClientSettings{AdditionalClientMetadata: map[string]string{}},
		}}}

		expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
	})

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(completeApplicationDTO(), completeApplicationDTO())
		})
	})

	It("refuses a group APIM added when neither version of the application sets groups", func() {
		remote := apimapp.DTO{Application: sdk.Application{Description: "before", Groups: []string{"dev"}}}
		merged := drift.Merge(
			drift.DetectWithNamespace(apimapp.DTO{Application: sdk.Application{Description: "before"}}, remote, "namespace"),
			drift.DetectWithNamespace(apimapp.DTO{Application: sdk.Application{Description: "after"}}, remote, "namespace"),
		)

		Expect(merged.DriftDetected()).To(BeTrue(), merged.String())
	})
})

func completeApplicationDTO() apimapp.DTO {
	GinkgoHelper()
	return loadFixture[apimapp.DTO]("application_dto.json")
}

func toApplicationDTO(app *v1alpha1.Application) apimapp.DTO {
	GinkgoHelper()
	dto, err := apimapp.ToDTO(app)
	Expect(err).ToNot(HaveOccurred())
	return dto
}

func completeApplicationCRD() *v1alpha1.Application {
	GinkgoHelper()
	fixture := loadFixture[v1alpha1.Application]("application_crd.json")
	return &fixture
}
