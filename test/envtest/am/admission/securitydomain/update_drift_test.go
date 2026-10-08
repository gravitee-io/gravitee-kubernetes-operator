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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/amsecuritydomain"
	internal "github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/securitydomain"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var _ = Describe("Validate drift", func() {
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
		dto, err := internal.ToDomainDTO(fixtures.AMSecurityDomain)
		Expect(err).ToNot(HaveOccurred())
		remoteDesc := "remote change"
		dto.Description = &remoteDesc
		_, err = sdk.UpsertDomainWithResponse(ctx, nil, dto)
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
		dto, err := internal.ToDomainDTO(fixtures.AMSecurityDomain)
		Expect(err).ToNot(HaveOccurred())
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

var _ = Describe("Validate drift - remote fetch failure", func() {
	ctx := context.Background()
	admissionCtrl := amsecuritydomain.NewAdmissionCtrl()

	// notInAM creates the AMContext only: the domain is never sent to the AM mock.
	notInAM := func() *v1alpha1.AMSecurityDomain {
		fixtures := fixture.Builder().
			WithAMContext(constants.AMContextFile).
			WithAMSecurityDomain(constants.AMSecurityDomainBasicFile).
			Build()
		Expect(manager.Client().Create(ctx, fixtures.AMContext)).To(Succeed())
		return fixtures.AMSecurityDomain
	}

	It("should apply fetch-failure policy when a synced domain no longer exists remotely", func() {
		sd := notInAM()
		dto, err := internal.ToDomainDTO(sd)
		Expect(err).ToNot(HaveOccurred())
		sd.Status.Key = dto.Key

		newSD := sd.DeepCopy()
		newSD.Spec.Description = new("changed")

		_, err = admissionCtrl.ValidateUpdate(ctx, sd, newSD)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("not found during drift detection"))
	})

	It("should admit the finalizer save of a domain not created in AM yet", func() {
		sd := notInAM()

		newSD := sd.DeepCopy()
		controllerutil.AddFinalizer(newSD, core.AMSecurityDomainFinalizer)

		_, err := admissionCtrl.ValidateUpdate(ctx, sd, newSD)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should admit a spec fix of a domain AM never accepted", func() {
		sd := notInAM()

		newSD := sd.DeepCopy()
		newSD.Spec.Description = new("fixed")

		_, err := admissionCtrl.ValidateUpdate(ctx, sd, newSD)
		Expect(err).ToNot(HaveOccurred())
	})
})
