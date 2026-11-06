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

	DescribeTable("compares an oidc setting against the defaults AM fills in",
		func(old, updated *bool, remote bool, expectDrift bool) {
			merged := drift.Merge(
				drift.DetectWithNamespace(domainWithDynamicRegistration(old), remoteDomainWithDynamicRegistration(remote), ""),
				drift.DetectWithNamespace(domainWithDynamicRegistration(updated), remoteDomainWithDynamicRegistration(remote), ""),
			)

			Expect(merged.DriftDetected()).To(Equal(expectDrift), merged.String())
		},
		Entry("a setting left unset stays unset", nil, nil, false, false),
		Entry("the first change to a setting AM defaulted is a CRD change", nil, new(true), false, false),
		Entry("a setting first set to the value AM holds is no change", nil, new(false), false, false),
		Entry("a CRD change to a declared setting is accepted", new(false), new(true), false, false),
		Entry("an AM change to a declared setting is drift", new(true), new(true), false, true),
		Entry("unsetting a declared setting AM changed is drift", new(true), nil, false, true),
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

// domainWithDynamicRegistration is a domain whose spec sets no oidc unless dynamicRegistration is set.
func domainWithDynamicRegistration(dynamicRegistration *bool) sdk.Domain {
	domain := sdk.Domain{Name: "cs"}
	if dynamicRegistration != nil {
		domain.Oidc = &sdk.AutomationOidcSettings{
			ClientRegistrationSettings: &sdk.AutomationClientRegistrationSettings{
				DynamicClientRegistrationEnabled: dynamicRegistration,
			},
		}
	}
	return domain
}

// remoteDomainWithDynamicRegistration is the domain AM returns: every client registration
// setting filled in, false by default.
func remoteDomainWithDynamicRegistration(dynamicRegistration bool) sdk.Domain {
	return sdk.Domain{
		Name: "cs",
		Oidc: &sdk.AutomationOidcSettings{
			ClientRegistrationSettings: &sdk.AutomationClientRegistrationSettings{
				AllowHttpSchemeRedirectUri:               new(false),
				AllowLocalhostRedirectUri:                new(false),
				AllowRedirectUriParamsExpressionLanguage: new(false),
				AllowWildCardRedirectUri:                 new(false),
				AllowedScopesEnabled:                     new(false),
				ClientTemplateEnabled:                    new(false),
				DynamicClientRegistrationEnabled:         new(dynamicRegistration),
				OpenDynamicClientRegistrationEnabled:     new(false),
			},
		},
	}
}
