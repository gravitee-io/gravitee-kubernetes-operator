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
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Update", labels.WithContext, func() {
	timeout := constants.EventualTimeout

	interval := constants.Interval

	ctx := context.Background()

	It("should update IDP group mapping condition", func() {
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

		originalCondition := fixtures.IDPGroupMapping.Spec.Condition

		By("updating IDP group mapping condition")

		updated := fixtures.IDPGroupMapping.DeepCopy()
		updated.Spec.Condition = "{#jsonPath(#profile, '$.roles') matches 'senior-developer'}"

		Expect(manager.UpdateSafely(ctx, updated)).To(Succeed())

		By("expecting IDP group mapping to be updated")

		Eventually(func() error {
			err := manager.GetLatest(ctx, fixtures.IDPGroupMapping)
			if err != nil {
				return err
			}
			return assert.Equals("condition", updated.Spec.Condition, fixtures.IDPGroupMapping.Spec.Condition)
		}, timeout, interval).Should(Succeed())

		By("verifying condition has changed")

		Expect(fixtures.IDPGroupMapping.Spec.Condition).NotTo(Equal(originalCondition))
		Expect(fixtures.IDPGroupMapping.Spec.Condition).To(Equal(updated.Spec.Condition))

		By("calling rest API, expecting condition to be updated in IDP configuration")

		apimClient := apim.NewClient(ctx)

		Eventually(func() error {
			idpConfig, err := apimClient.Configuration.GetIDPConfiguration(fixtures.IDPGroupMapping.Spec.IDPID)
			if err != nil {
				return err
			}

			// Verify the old condition no longer exists
			for _, gm := range idpConfig.GroupMappings {
				if gm.Condition == originalCondition {
					return assert.Equals("old condition should not exist", false, true)
				}
			}

			// Verify the new condition exists
			found := false
			for _, gm := range idpConfig.GroupMappings {
				if gm.Condition == updated.Spec.Condition {
					found = true
					break
				}
			}

			if !found {
				return assert.Equals("new condition found", true, false)
			}
			return nil
		}, timeout, interval).Should(Succeed(), fixtures.IDPGroupMapping.Name)
	})

	It("should update IDP group mapping groups", func() {
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

		originalGroupsCount := len(fixtures.IDPGroupMapping.Spec.Groups)

		By("updating IDP group mapping to add more groups")

		updated := fixtures.IDPGroupMapping.DeepCopy()
		updated.Spec.Groups = append(updated.Spec.Groups, "api-publishers")

		Expect(manager.UpdateSafely(ctx, updated)).To(Succeed())

		By("expecting IDP group mapping to have updated groups")

		Eventually(func() error {
			err := manager.GetLatest(ctx, fixtures.IDPGroupMapping)
			if err != nil {
				return err
			}
			if len(fixtures.IDPGroupMapping.Spec.Groups) != len(updated.Spec.Groups) {
				return assert.Equals("groups count", len(updated.Spec.Groups), len(fixtures.IDPGroupMapping.Spec.Groups))
			}
			return nil
		}, timeout, interval).Should(Succeed())

		By("verifying groups have been updated")

		Expect(len(fixtures.IDPGroupMapping.Spec.Groups)).To(Equal(originalGroupsCount + 1))

		By("calling rest API, expecting groups to be updated in IDP configuration")

		apimClient := apim.NewClient(ctx)

		Eventually(func() error {
			idpConfig, err := apimClient.Configuration.GetIDPConfiguration(fixtures.IDPGroupMapping.Spec.IDPID)
			if err != nil {
				return err
			}

			// Verify the group mapping has the updated number of groups
			found := false
			for _, gm := range idpConfig.GroupMappings {
				if gm.Condition == fixtures.IDPGroupMapping.Spec.Condition {
					if len(gm.Groups) == len(updated.Spec.Groups) {
						found = true
						break
					}
				}
			}

			if !found {
				return assert.Equals("updated groups count found", true, false)
			}
			return nil
		}, timeout, interval).Should(Succeed(), fixtures.IDPGroupMapping.Name)
	})
})
