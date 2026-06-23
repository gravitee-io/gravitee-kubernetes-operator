package drift

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Application Drift detection", func() {

	DescribeTable("equivalent values",
		func(crd, api any) {
			expectNoDrift(drift.Detect(crd, api))
		},
		Entry("empty struct",
			model.ApplicationDTO{},
			model.ApplicationDTO{},
		),
		Entry("equal struct",
			completeApplication(),
			completeApplication(),
		),
		Entry("equal struct",
			completeApplication(),
			completeSpec().ToDTO(),
		),
		Entry("equivalent struct",
			model.ApplicationDTO{},
			model.ApplicationDTO{
				Status:        "ACTIVE",
				Groups:        make([]string, 0),
				Members:       make([]application.Member, 0),
				Metadata:      make([]application.Metadata, 0),
				Settings:      &application.Setting{},
				NotifyMembers: ptr(false),
			},
		),
		Entry("client certificate equivalent",
			model.ApplicationDTO{
				Settings: &application.Setting{
					TLS: &application.TLSSettings{
						ClientCertificate: `-----BEGIN CERTIFICATE-----
-----END CERTIFICATE-----`,
					},
				},
			},
			model.ApplicationDTO{
				Settings: &application.Setting{
					TLS: &application.TLSSettings{
						ClientCertificate: `-----BEGIN CERTIFICATE-----
-----END CERTIFICATE-----
`,
					},
				},
			},
		),
		Entry("listed certificates equivalent",
			model.ApplicationDTO{
				Settings: &application.Setting{
					TLS: &application.TLSSettings{
						ClientCertificates: []application.ClientCertificate{
							{
								Content: `-----BEGIN CERTIFICATE-----
-----END CERTIFICATE-----`,
							},
						},
					},
				},
			},
			model.ApplicationDTO{
				Settings: &application.Setting{
					TLS: &application.TLSSettings{
						ClientCertificates: []application.ClientCertificate{
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
				Settings: &application.Setting{
					TLS: &application.TLSSettings{
						ClientCertificates: []application.ClientCertificate{
							{
								Ref: &application.CertificateRef{
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
				Settings: &application.Setting{
					TLS: &application.TLSSettings{
						ClientCertificates: []application.ClientCertificate{
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
				Settings: &application.Setting{
					TLS: &application.TLSSettings{
						ClientCertificates: []application.ClientCertificate{
							{
								StartsAt: "2023-07-25T02:43:16+03:00",
								EndsAt:   "2023-08-25T23:43:16Z",
							},
						},
					},
				},
			},
			model.ApplicationDTO{
				Settings: &application.Setting{
					TLS: &application.TLSSettings{
						ClientCertificates: []application.ClientCertificate{
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
})

func completeApplication() model.ApplicationDTO {
	GinkgoHelper()
	return model.ApplicationDTO{
		ID:          "123456",
		HRID:        "my-app",
		Name:        "My Application",
		Status:      "ACTIVE",
		Description: "This is my application",
		Settings: &application.Setting{
			App: &application.SimpleSettings{
				AppType:  "Web",
				ClientID: ptr("my-client"),
			},
			Oauth: &application.OAuthClientSettings{
				ApplicationType: "M2M",
				GrantTypes:      []application.GrantType{application.GrantTypeClientCredentials},
				RedirectUris:    []string{"https://my-app"},
			},
			TLS: &application.TLSSettings{
				ClientCertificate: "CERTIFICATE",
				ClientCertificates: []application.ClientCertificate{
					{
						Name:     "Foo",
						Content:  "CERTIFICATE",
						StartsAt: "2023-07-25T20:43:16-03:00",
						EndsAt:   "2024-07-25T20:43:16-03:00",
						Encoded:  false,
					},
				},
			},
		},
		Background: "None",
		Domain:     "gravitee.io",
		Groups:     []string{"group1", "group2"},
		PictureURL: "None",
		Picture:    "None",
	}
}

func completeSpec() *v1alpha1.Application {
	GinkgoHelper()
	return &v1alpha1.Application{
		Spec: v1alpha1.ApplicationSpec{
			Application: application.Application{
				ID:          "123456",
				HRID:        "my-app",
				Name:        "My Application",
				Description: "This is my application",
				Settings: &application.Setting{
					App: &application.SimpleSettings{
						AppType:  "Web",
						ClientID: ptr("my-client"),
					},
					Oauth: &application.OAuthClientSettings{
						ApplicationType: "M2M",
						GrantTypes:      []application.GrantType{application.GrantTypeClientCredentials},
						RedirectUris:    []string{"https://my-app"},
					},
					TLS: &application.TLSSettings{
						ClientCertificate: "CERTIFICATE",
						ClientCertificates: []application.ClientCertificate{
							{
								Name:     "Foo",
								Content:  "CERTIFICATE",
								StartsAt: "2023-07-25T20:43:16-03:00",
								EndsAt:   "2024-07-25T20:43:16-03:00",
								Encoded:  false,
							},
						},
					},
				},
				Background: ptr("None"),
				Domain:     ptr("gravitee.io"),
				Groups:     []string{"group1", "group2"},
				PictureURL: ptr("None"),
				Picture:    ptr("None"),
			}},
	}
}
