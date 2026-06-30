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

package drift

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Application Drift detection", func() {

	DescribeTable("equivalent values",
		func(crd, remote any) {
			expectNoDrift(drift.Detect(crd, remote))
		},
		Entry("empty struct",
			model.ApplicationDTO{},
			model.ApplicationDTO{},
		),
		Entry("equal struct",
			completeApplicationDTO(),
			completeApplicationDTO(),
		),
		Entry("equal struct",
			completeApplicationDTO(),
			model.ToApplicationDTO(completeApplicationCRD().Spec),
		),
		Entry("equivalent struct",
			model.ApplicationDTO{},
			model.ApplicationDTO{
				ID:            "123456",
				HRID:          "my-app",
				Status:        "ACTIVE",
				Groups:        make([]string, 0),
				Members:       make([]model.ApplicationMemberDTO, 0),
				Metadata:      make([]model.ApplicationMetadataDTO, 0),
				Settings:      &model.ApplicationSettingsDTO{},
				NotifyMembers: ptr(false),
			},
		),
		Entry("client certificate equivalent",
			model.ApplicationDTO{
				Settings: &model.ApplicationSettingsDTO{
					TLS: &model.ApplicationTLSSettingsDTO{
						ClientCertificate: `-----BEGIN CERTIFICATE-----
-----END CERTIFICATE-----`,
					},
				},
			},
			model.ApplicationDTO{
				Settings: &model.ApplicationSettingsDTO{
					TLS: &model.ApplicationTLSSettingsDTO{
						ClientCertificate: `-----BEGIN CERTIFICATE-----
-----END CERTIFICATE-----
`,
					},
				},
			},
		),
		Entry("listed certificates equivalent",
			model.ApplicationDTO{
				Settings: &model.ApplicationSettingsDTO{
					TLS: &model.ApplicationTLSSettingsDTO{
						ClientCertificates: []model.ApplicationClientCertificateDTO{
							{
								Content: `-----BEGIN CERTIFICATE-----
-----END CERTIFICATE-----`,
							},
						},
					},
				},
			},
			model.ApplicationDTO{
				Settings: &model.ApplicationSettingsDTO{
					TLS: &model.ApplicationTLSSettingsDTO{
						ClientCertificates: []model.ApplicationClientCertificateDTO{
							{
								Content: `-----BEGIN CERTIFICATE-----
-----END CERTIFICATE-----
`,
							},
						},
					},
				},
			},
		),
		Entry("listed certificates equivalent, ref ignored",
			model.ApplicationDTO{
				Settings: &model.ApplicationSettingsDTO{
					TLS: &model.ApplicationTLSSettingsDTO{
						ClientCertificates: []model.ApplicationClientCertificateDTO{
							{
								Ref: &model.ApplicationCertificateRefDTO{
									Kind: "configmaps",
									Name: "foo",
									Key:  "cert",
								},
								Content: `-----BEGIN CERTIFICATE-----
-----END CERTIFICATE-----`,
							},
						},
					},
				},
			},
			model.ApplicationDTO{
				Settings: &model.ApplicationSettingsDTO{
					TLS: &model.ApplicationTLSSettingsDTO{
						ClientCertificates: []model.ApplicationClientCertificateDTO{
							{
								Content: `-----BEGIN CERTIFICATE-----
-----END CERTIFICATE-----
`,
							},
						},
					},
				},
			},
		),
		Entry("listed certificates equivalent, start/end date",
			model.ApplicationDTO{
				Settings: &model.ApplicationSettingsDTO{
					TLS: &model.ApplicationTLSSettingsDTO{
						ClientCertificates: []model.ApplicationClientCertificateDTO{
							{
								StartsAt: "2023-07-25T02:43:16+03:00",
								EndsAt:   "2023-08-25T23:43:16Z",
							},
						},
					},
				},
			},
			model.ApplicationDTO{
				Settings: &model.ApplicationSettingsDTO{
					TLS: &model.ApplicationTLSSettingsDTO{
						ClientCertificates: []model.ApplicationClientCertificateDTO{
							{
								StartsAt: "2023-07-24T20:43:16-03:00",
								EndsAt:   "2023-08-25T23:43:16-00:00",
							},
						},
					},
				},
			},
		),
	)

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(completeApplicationDTO(), completeApplicationDTO())
		})
	})
})

func completeApplicationDTO() model.ApplicationDTO {
	GinkgoHelper()
	return loadFixture[model.ApplicationDTO]("application_dto.json")
}

func completeApplicationCRD() *v1alpha1.Application {
	GinkgoHelper()
	fixture := loadFixture[v1alpha1.Application]("application_crd.json")
	return &fixture
}
