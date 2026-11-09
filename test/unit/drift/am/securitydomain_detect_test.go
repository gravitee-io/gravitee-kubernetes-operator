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

package am

import (
	"github.com/gravitee-io/gravitee-automation-sdk/am-sdk/v2/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/securitydomain"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Security domain drift detection", func() {
	It("does not drift between the CRD mapping and the AM response", func() {
		expectNoDrift(drift.DetectWithNamespace(securityDomainForDrift(), remoteSecurityDomain(), ""))
	})

	It("does not drift on the data plane AM resolves when the CRD leaves it unset", func() {
		crd := securityDomainForDrift()
		crd.DataPlaneId = nil

		expectNoDrift(drift.DetectWithNamespace(crd, remoteSecurityDomain(), ""))
	})

	It("does not drift on a default allowed algorithm only AM holds", func() {
		remote := remoteSecurityDomain()
		spiffe := remote.Oidc.WorkloadIdentitySettings
		spiffe.DefaultAllowedAlgorithms = append(spiffe.DefaultAllowedAlgorithms, "ES256")

		expectNoDrift(drift.DetectWithNamespace(securityDomainForDrift(), remote, ""))
	})

	It("detects drift on a declared default allowed algorithm AM does not hold", func() {
		crd := securityDomainForDrift()
		spiffe := crd.Oidc.WorkloadIdentitySettings
		spiffe.DefaultAllowedAlgorithms = append(spiffe.DefaultAllowedAlgorithms, "ES256")

		result := drift.DetectWithNamespace(crd, remoteSecurityDomain(), "")
		Expect(result.DriftDetected()).To(BeTrue())
		Expect(result.String()).To(ContainSubstring("defaultAllowedAlgorithms"))
	})

	It("detects drift when AM adds a tag", func() {
		remote := remoteSecurityDomain()
		remote.Tags = append(remote.Tags, "prod")

		expectDrift(drift.DetectWithNamespace(securityDomainForDrift(), remote, ""), `tags: (2 unchanged)
  <unchanged>
              !=
                 + "prod"`)
	})

	It("detects drift when AM drops a CORS method", func() {
		remote := remoteSecurityDomain()
		remote.CorsSettings.AllowedMethods = []string{"DELETE", "POST", "GET"}

		expectDrift(drift.DetectWithNamespace(securityDomainForDrift(), remote, ""), `corsSettings:
  allowedMethods: (3 unchanged)
    + "PUT"
            !=
               <unchanged>`)
	})

	DescribeTable("compares oidc settings against the defaults AM fills in",
		func(old, updated func(*sdk.Domain), expectDrift bool) {
			remote := domainWithOidcDefaults()
			merged := drift.Merge(
				drift.DetectWithNamespace(domainSpec(old), remote, ""),
				drift.DetectWithNamespace(domainSpec(updated), remote, ""),
			)

			Expect(merged.DriftDetected()).To(Equal(expectDrift), merged.String())
		},
		Entry("no oidc on either version is no change", noChange, noChange, false),
		Entry("the first oidc setting is a CRD change", noChange, withDynamicRegistration(true), false),
		Entry("a later setting AM defaulted is a CRD change",
			withDynamicRegistration(false), with(withDynamicRegistration(false), withPlainFapi), false),
		Entry("a CIBA number set over AM's default is a CRD change",
			withCibaDisabled, with(withCibaDisabled, withAuthReqExpiry(300)), false),
		Entry("an AM change to a declared setting is drift",
			withDynamicRegistration(true), withDynamicRegistration(true), true),
		Entry("an untagged field set over a console change is drift", noChange, withDescription, true),
	)

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(securityDomainForDrift(), remoteSecurityDomain())
		})
	})
})

func securityDomainForDrift() sdk.Domain {
	GinkgoHelper()
	crd := loadFixture[v1alpha1.AMSecurityDomain]("securitydomain_crd.json")
	dto, err := securitydomain.ToDomainDTO(&crd)
	Expect(err).ToNot(HaveOccurred())
	return dto
}

func remoteSecurityDomain() sdk.Domain {
	GinkgoHelper()
	return loadFixture[sdk.Domain]("securitydomain_dto.json")
}

// domainSpec is a domain whose spec declares only what change sets.
func domainSpec(change func(*sdk.Domain)) sdk.Domain {
	domain := sdk.Domain{Name: "cs"}
	change(&domain)
	return domain
}

// domainWithOidcDefaults is the domain AM returns for a spec without oidc, after a console user set a description.
func domainWithOidcDefaults() sdk.Domain {
	return sdk.Domain{
		Name:        "cs",
		Description: new("set in the console"),
		Oidc: &sdk.AutomationOidcSettings{
			RedirectUriStrictMatching: new(false),
			ClientRegistrationSettings: &sdk.AutomationClientRegistrationSettings{
				AllowHttpSchemeRedirectUri:               new(false),
				AllowLocalhostRedirectUri:                new(false),
				AllowRedirectUriParamsExpressionLanguage: new(false),
				AllowWildCardRedirectUri:                 new(false),
				AllowedScopesEnabled:                     new(false),
				ClientTemplateEnabled:                    new(false),
				DynamicClientRegistrationEnabled:         new(false),
				OpenDynamicClientRegistrationEnabled:     new(false),
			},
			SecurityProfileSettings: &sdk.SecurityProfileSettings{
				EnableFapiBrazil: new(false),
				EnablePlainFapi:  new(false),
			},
			CibaSettings: &sdk.AutomationCIBASettings{
				Enabled:              new(false),
				AuthReqExpiry:        new(int32(600)),
				TokenReqInterval:     new(int32(5)),
				BindingMessageLength: new(int32(256)),
			},
			WorkloadIdentitySettings: &sdk.SpiffeDomainSettings{
				Enabled:               new(false),
				MaxJwtLifetimeSeconds: new(int32(300)),
				ClockSkewSeconds:      new(int32(30)),
			},
		},
	}
}

func noChange(*sdk.Domain) {}

func with(changes ...func(*sdk.Domain)) func(*sdk.Domain) {
	return func(domain *sdk.Domain) {
		for _, change := range changes {
			change(domain)
		}
	}
}

func oidc(domain *sdk.Domain) *sdk.AutomationOidcSettings {
	if domain.Oidc == nil {
		domain.Oidc = &sdk.AutomationOidcSettings{}
	}
	return domain.Oidc
}

func withDynamicRegistration(enabled bool) func(*sdk.Domain) {
	return func(domain *sdk.Domain) {
		oidc(domain).ClientRegistrationSettings = &sdk.AutomationClientRegistrationSettings{
			DynamicClientRegistrationEnabled: new(enabled),
		}
	}
}

func withPlainFapi(domain *sdk.Domain) {
	oidc(domain).SecurityProfileSettings = &sdk.SecurityProfileSettings{EnablePlainFapi: new(true)}
}

func withCibaDisabled(domain *sdk.Domain) {
	oidc(domain).CibaSettings = &sdk.AutomationCIBASettings{Enabled: new(false)}
}

func withAuthReqExpiry(seconds int32) func(*sdk.Domain) {
	return func(domain *sdk.Domain) {
		oidc(domain).CibaSettings.AuthReqExpiry = new(seconds)
	}
}

func withDescription(domain *sdk.Domain) {
	domain.Description = new("set in the resource")
}
