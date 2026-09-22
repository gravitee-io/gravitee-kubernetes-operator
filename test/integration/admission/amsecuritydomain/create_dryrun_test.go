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
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate create", labels.WithContext, func() {
	ctx := context.Background()
	admissionCtrl := amsecuritydomain.NewAdmissionCtrl()

	It("should pass dryrun validation", func() {
		fixtures := fixture.Builder().
			WithAMContext(constants.AMContextFile).
			WithAMSecurityDomain(constants.AMSecurityDomainBasicFile).
			Build()
		fixtures.Apply()

		_, err := admissionCtrl.ValidateCreate(ctx, fixtures.AMSecurityDomain)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should fail dryrun validation when server rejects", func() {
		By("restarting mock server with dryRunReject=true")
		mockServer.Stop()
		rejectServer := am.Start(true)
		defer func() {
			rejectServer.Stop()
			mockServer = am.Start(false)
		}()

		fixtures := fixture.Builder().
			WithAMContext(constants.AMContextFile).
			WithAMSecurityDomain(constants.AMSecurityDomainBasicFile).
			Build()
		fixtures.Apply()

		_, err := admissionCtrl.ValidateCreate(ctx, fixtures.AMSecurityDomain)
		Expect(err).To(HaveOccurred())
	})

	It("should fail when AMContext does not exist", func() {
		fixtures := fixture.Builder().
			WithAMSecurityDomain(constants.AMSecurityDomainBasicFile).
			Build()

		_, err := admissionCtrl.ValidateCreate(ctx, fixtures.AMSecurityDomain)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("could not resolve context"))
	})
})
