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

package idpgroupmapping

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
)

var _ = Describe("Create", labels.WithContext, func() {
	timeout := constants.EventualTimeout

	interval := constants.Interval

	ctx := context.Background()

	It("should create IDP group mapping", func() {
		fixtures := fixture.Builder().
			AddSecret(constants.ContextSecretFile).
			WithGroup(constants.GroupFile).
			WithIDPGroupMapping(constants.IDPGroupMappingFile).
			WithContext(constants.ContextWithSecretFile).
			Build().
			Apply()

		By("expecting IDP group mapping status to be completed")

		Expect(assert.IDPGroupMappingCompleted(fixtures.IDPGroupMapping)).To(Succeed())
		Expect(assert.IDPGroupMappingAccepted(fixtures.IDPGroupMapping)).To(Succeed())

		By("verifying organization and environment IDs are set")

		Expect(fixtures.IDPGroupMapping.Status.OrgID).NotTo(BeEmpty())
		Expect(fixtures.IDPGroupMapping.Status.EnvID).NotTo(BeEmpty())

		By("verifying group names have been resolved")

		Expect(fixtures.IDPGroupMapping.Status.Groups).NotTo(BeEmpty())
		Expect(len(fixtures.IDPGroupMapping.Status.Groups)).To(Equal(len(fixtures.IDPGroupMapping.Spec.Groups)))

		By("calling rest API, expecting to find IDP group mapping in IDP configuration")

		apimClient := apim.NewClient(ctx)

		Eventually(func() error {
			idpConfig, err := apimClient.Configuration.GetIDPConfiguration(fixtures.IDPGroupMapping.Spec.IDPID)
			if err != nil {
				return err
			}

			// Verify the group mapping exists in the IDP configuration
			found := false
			for _, gm := range idpConfig.GroupMappings {
				if gm.Condition == fixtures.IDPGroupMapping.Spec.Condition {
					if len(gm.Groups) == len(fixtures.IDPGroupMapping.Status.Groups) {
						found = true
						break
					}
				}
			}

			if !found {
				return assert.Equals("group mapping found", true, false)
			}
			return nil
		}, timeout, interval).Should(Succeed(), fixtures.IDPGroupMapping.Name)
	})

	It("should create IDP group mapping with multiple groups", func() {
		fixtures := fixture.Builder().
			AddSecret(constants.ContextSecretFile).
			WithGroup(constants.GroupFile).
			WithIDPGroupMapping(constants.IDPGroupMappingMultipleGroups).
			WithContext(constants.ContextWithSecretFile).
			Build().
			Apply()

		By("expecting IDP group mapping status to be completed")

		Expect(assert.IDPGroupMappingCompleted(fixtures.IDPGroupMapping)).To(Succeed())
		Expect(assert.IDPGroupMappingAccepted(fixtures.IDPGroupMapping)).To(Succeed())

		By("verifying multiple groups are mapped")

		Expect(fixtures.IDPGroupMapping.Status.Groups).NotTo(BeEmpty())
		Expect(len(fixtures.IDPGroupMapping.Status.Groups)).To(Equal(2))

		By("calling rest API, expecting to find IDP group mapping with multiple groups")

		apimClient := apim.NewClient(ctx)

		Eventually(func() error {
			idpConfig, err := apimClient.Configuration.GetIDPConfiguration(fixtures.IDPGroupMapping.Spec.IDPID)
			if err != nil {
				return err
			}

			// Verify the group mapping exists with multiple groups
			found := false
			for _, gm := range idpConfig.GroupMappings {
				if gm.Condition == fixtures.IDPGroupMapping.Spec.Condition {
					if len(gm.Groups) == 2 {
						found = true
						break
					}
				}
			}

			if !found {
				return assert.Equals("group mapping with multiple groups found", true, false)
			}
			return nil
		}, timeout, interval).Should(Succeed(), fixtures.IDPGroupMapping.Name)
	})
})
