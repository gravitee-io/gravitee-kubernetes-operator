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

package amidentityprovider

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-automation-sdk/common/pkg/unstructured"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/amidentityprovider"
	internal "github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/identityprovider"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
)

var _ = Describe("Validate drift", func() {
	ctx := context.Background()
	admissionCtrl := amidentityprovider.NewAdmissionCtrl()

	// inAM applies the domain only, edits the identity provider spec with local, and stores the identity provider
	// in the AM mock as AM returns it: the spec, with the fields remote overrides.
	inAM := func(
		local func(*v1alpha1.AMIdentityProviderSpec), remote func(*internal.IdentityProvider),
	) *v1alpha1.AMIdentityProvider {
		fixtures := withDomain()
		idp := fixtures.AMIdentityProvider
		local(&idp.Spec)
		fixtures.AMIdentityProvider = nil
		fixtures.Apply()

		dto, err := internal.ToIdentityProviderDTO(idp)
		Expect(err).ToNot(HaveOccurred())
		remote(&dto)
		resp, err := am.NewSDKClient().UpsertIdentityProviderWithResponse(ctx, dto.DomainKey, nil, dto.IdentityProvider)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.JSON200).ToNot(BeNil(), string(resp.Body))
		return idp
	}
	validateUpdate := func(idp *v1alpha1.AMIdentityProvider) func() error {
		return func() error {
			updated := idp.DeepCopy()
			updated.Annotations = map[string]string{"updated": "true"}
			_, err := admissionCtrl.ValidateUpdate(ctx, idp, updated)
			return err
		}
	}

	// Left unset, name, type and configuration never drift: a field the CRD does not set is not compared.
	// Set, they are ignored by AM for the system identity provider, which returns its own.
	It("should not drift on the name, type and configuration AM replaces for the system identity provider", func() {
		idp := inAM(func(spec *v1alpha1.AMIdentityProviderSpec) {
			spec.System = new(true)
		}, func(remote *internal.IdentityProvider) {
			remote.Name = new("Default Identity Provider")
			remote.Type = new("gravitee-am-idp")
			remote.Configuration = unstructured.StringifiedFrom(map[string]any{"passwordEncoder": "BCrypt"})
		})

		Consistently(validateUpdate(idp), constants.ConsistentTimeout, constants.Interval).Should(Succeed())
	})

	It("should detect drift on the name of a regular identity provider", func() {
		idp := inAM(func(*v1alpha1.AMIdentityProviderSpec) {}, func(remote *internal.IdentityProvider) {
			remote.Name = new("Renamed in AM")
		})

		Eventually(validateUpdate(idp), constants.EventualTimeout, constants.Interval).
			Should(MatchError(And(ContainSubstring("drift detected"), ContainSubstring("Renamed in AM"))))
	})
})
