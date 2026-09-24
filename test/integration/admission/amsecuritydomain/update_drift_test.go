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

package amsecuritydomain

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/amsecuritydomain"
	internal "github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/securitydomain"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate drift", labels.WithContext, func() {
	ctx := context.Background()
	admissionCtrl := amsecuritydomain.NewAdmissionCtrl()

	It("should not drift when only the CRD changes", func() {
		fixtures := fixture.Builder().
			WithAMContext(constants.AMContextFile).
			WithAMSecurityDomain(constants.AMSecurityDomainBasicFile).
			Build()
		fixtures.Apply()

		By("updating the local CRD description")
		newSD := fixtures.AMSecurityDomain.DeepCopy()
		desc := "updated description"
		newSD.Spec.Description = &desc

		_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.AMSecurityDomain, newSD)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should detect drift when remote is modified", func() {
		fixtures := fixture.Builder().
			WithAMContext(constants.AMContextFile).
			WithAMSecurityDomain(constants.AMSecurityDomainBasicFile).
			Build()
		fixtures.Apply()

		By("modifying the remote domain via the AM mock")
		sdk := am.NewSDKClient()
		dto := internal.ToDomainDTO(fixtures.AMSecurityDomain)
		remoteDesc := "remote change"
		dto.Description = &remoteDesc
		_, err := sdk.UpsertDomainWithResponse(ctx, nil, dto)
		Expect(err).ToNot(HaveOccurred())

		By("changing the local CRD description")
		newSD := fixtures.AMSecurityDomain.DeepCopy()
		localDesc := "local CRD change"
		newSD.Spec.Description = &localDesc

		Eventually(func() error {
			_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.AMSecurityDomain, newSD)
			return assert.DriftDetected(`description: "local CRD change" != "remote change"`, err)
		}, constants.EventualTimeout, constants.Interval).Should(Succeed())
	})

	It("should not drift when CRD realigns with remote", func() {
		fixtures := fixture.Builder().
			WithAMContext(constants.AMContextFile).
			WithAMSecurityDomain(constants.AMSecurityDomainBasicFile).
			Build()
		fixtures.Apply()

		By("modifying the remote domain via the AM mock")
		sdk := am.NewSDKClient()
		dto := internal.ToDomainDTO(fixtures.AMSecurityDomain)
		remoteDesc := "remote value"
		dto.Description = &remoteDesc

		resp, err := sdk.UpsertDomainWithResponse(ctx, nil, dto)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.JSON200).ToNot(BeNil())

		By("updating the CRD to match the remote")
		newSD := fixtures.AMSecurityDomain.DeepCopy()
		newSD.Spec.Description = &remoteDesc

		Eventually(func() error {
			_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.AMSecurityDomain, newSD)
			return err
		}, constants.EventualTimeout, constants.Interval).Should(Succeed())
	})
})

var _ = Describe("Validate drift - remote fetch failure", labels.WithContext, func() {
	ctx := context.Background()

	It("should apply fetch-failure policy when domain does not exist remotely", func() {
		admissionCtrl := amsecuritydomain.NewAdmissionCtrl()

		fixtures := fixture.Builder().
			WithAMContext(constants.AMContextFile).
			WithAMSecurityDomain(constants.AMSecurityDomainBasicFile).
			Build()

		newSD := fixtures.AMSecurityDomain.DeepCopy()
		desc := "changed"
		newSD.Spec.Description = &desc

		_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.AMSecurityDomain, newSD)
		if err != nil {
			Expect(err.Error()).To(ContainSubstring("not found during drift detection"))
		}
	})
})

