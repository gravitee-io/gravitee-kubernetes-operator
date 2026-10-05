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
	"github.com/gravitee-io/gravitee-automation-sdk/common/pkg/unstructured"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/identityprovider"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Identity provider drift detection", func() {
	It("does not drift on the passwords AM returns masked", func() {
		expectNoDrift(drift.DetectWithNamespace(identityProviderForDrift(), remoteIdentityProvider(), ""))
	})

	It("detects drift on a configuration value AM does not mask", func() {
		remote := remoteIdentityProvider()
		remote.Configuration = unstructured.StringifiedFrom(map[string]any{
			"passwordEncoder": "BCrypt",
			"users": []any{map[string]any{
				"username": "jdoe",
				"email":    "jane.doe@example.com",
				"password": "********",
			}},
		})

		expectDrift(drift.DetectWithNamespace(identityProviderForDrift(), remote, ""), `configuration:
  users:
    email: "john.doe@example.com" != "jane.doe@example.com"`)
	})

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			// system is false on a regular identity provider, a zero value this check rejects.
			crd := identityProviderForDrift()
			crd.System = new(true)
			remote := remoteIdentityProvider()
			remote.System = new(true)
			expectedEquivalentNotHavingAnyZeroValue(crd, remote)
		})
	})
})

func identityProviderForDrift() sdk.IdentityProvider {
	GinkgoHelper()
	crd := loadFixture[v1alpha1.AMIdentityProvider]("identityprovider_crd.json")
	dto, err := identityprovider.ToIdentityProviderDTOForDrift(&crd)
	Expect(err).ToNot(HaveOccurred())
	return dto.IdentityProvider
}

func remoteIdentityProvider() sdk.IdentityProvider {
	GinkgoHelper()
	return loadFixture[sdk.IdentityProvider]("identityprovider_dto.json")
}
