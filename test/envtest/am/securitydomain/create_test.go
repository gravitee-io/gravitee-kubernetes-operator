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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Create", func() {
	ctx := context.Background()

	It("should create a basic security domain", func() {
		fixtures := fixture.Builder().
			AddSecret(constants.AMContextSecretFile).
			WithAMContext(constants.AMContextFile).
			WithAMSecurityDomain(constants.AMSecurityDomainBasicFile).
			Build().
			Apply()

		By("expecting security domain to be accepted")

		Expect(assert.AMSecurityDomainAccepted(fixtures.AMSecurityDomain)).To(Succeed())

		By("expecting status to contain the domain ID")

		Expect(fixtures.AMSecurityDomain.Status.ID).ToNot(BeEmpty())
		Expect(fixtures.AMSecurityDomain.Status.OrgID).To(Equal("DEFAULT"))
		Expect(fixtures.AMSecurityDomain.Status.EnvID).To(Equal("DEFAULT"))
	})

	It("should create a security domain with full settings", func() {
		fixtures := fixture.Builder().
			AddSecret(constants.AMContextSecretFile).
			WithAMContext(constants.AMContextFile).
			WithAMSecurityDomain(constants.AMSecurityDomainFullFile).
			Build().
			Apply()

		By("expecting security domain to be accepted")

		Expect(assert.AMSecurityDomainAccepted(fixtures.AMSecurityDomain)).To(Succeed())

		By("expecting status to contain the domain ID")

		Expect(fixtures.AMSecurityDomain.Status.ID).ToNot(BeEmpty())
	})

	It("should fail when AMContext has bad token", func() {
		fixtures := fixture.Builder().
			WithAMContext(constants.AMContextBadTokenFile).
			WithAMSecurityDomain(constants.AMSecurityDomainBasicFile).
			Build()

		Expect(manager.Client().Create(ctx, fixtures.AMContext)).ToNot(HaveOccurred())
		Expect(manager.Client().Create(ctx, fixtures.AMSecurityDomain)).ToNot(HaveOccurred())

		By("expecting security domain to fail")

		Eventually(ctx, func() error {
			if err := manager.GetLatest(ctx, fixtures.AMSecurityDomain); err != nil {
				return err
			}
			return assert.AMSecurityDomainFailed(fixtures.AMSecurityDomain)
		}, constants.EventualTimeout, constants.Interval).Should(Succeed(), fixtures.AMSecurityDomain.Name)
	})

	It("should fail when AMContext does not exist", func() {
		fixtures := fixture.Builder().
			WithAMSecurityDomain(constants.AMSecurityDomainBasicFile).
			Build()

		Expect(manager.Client().Create(ctx, fixtures.AMSecurityDomain)).ToNot(HaveOccurred())

		By("expecting security domain to fail with unresolved ref")

		Eventually(ctx, func() error {
			if err := manager.GetLatest(ctx, fixtures.AMSecurityDomain); err != nil {
				return err
			}
			return assert.AMSecurityDomainFailed(fixtures.AMSecurityDomain)
		}, constants.EventualTimeout, constants.Interval).Should(Succeed(), fixtures.AMSecurityDomain.Name)
	})

	It("should fail when AMContext is unreachable", func() {
		fixtures := fixture.Builder().
			WithAMContext(constants.AMContextUnreachableFile).
			WithAMSecurityDomain(constants.AMSecurityDomainBasicFile).
			Build()

		Expect(manager.Client().Create(ctx, fixtures.AMContext)).ToNot(HaveOccurred())
		Expect(manager.Client().Create(ctx, fixtures.AMSecurityDomain)).ToNot(HaveOccurred())

		By("expecting security domain to fail")

		Eventually(ctx, func() error {
			if err := manager.GetLatest(ctx, fixtures.AMSecurityDomain); err != nil {
				return err
			}
			return assert.AMSecurityDomainFailed(fixtures.AMSecurityDomain)
		}, constants.EventualTimeout, constants.Interval).Should(Succeed(), fixtures.AMSecurityDomain.Name)
	})
})
