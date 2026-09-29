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

var _ = Describe("Update", func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval
	ctx := context.Background()

	It("should update security domain in AM", func() {
		fixtures := fixture.Builder().
			AddSecret(constants.AMContextSecretFile).
			WithAMContext(constants.AMContextFile).
			WithAMSecurityDomain(constants.AMSecurityDomainBasicFile).
			Build().
			Apply()

		By("expecting security domain to be accepted")

		Expect(assert.AMSecurityDomainAccepted(fixtures.AMSecurityDomain)).To(Succeed())

		By("updating security domain name")

		updated := fixtures.AMSecurityDomain.DeepCopy()
		updated.Spec.Name += "-updated"

		Expect(manager.UpdateSafely(ctx, updated)).To(Succeed())

		By("expecting updated security domain to be accepted")

		Eventually(ctx, func() error {
			if err := manager.GetLatest(ctx, updated); err != nil {
				return err
			}
			return assert.AMSecurityDomainAccepted(updated)
		}, timeout, interval).Should(Succeed(), updated.Name)
	})
})
