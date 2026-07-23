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
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Delete", labels.WithContext, func() {
	timeout := constants.EventualTimeout

	interval := constants.Interval

	ctx := context.Background()

	It("should delete IDP group mapping", func() {
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

		// Store the original condition for verification later
		originalCondition := fixtures.IDPGroupMapping.Spec.Condition
		originalIDPID := fixtures.IDPGroupMapping.Spec.IDPID

		By("deleting IDP group mapping")

		Expect(manager.Client().Delete(ctx, fixtures.IDPGroupMapping.DeepCopy())).To(Succeed())

		By("expecting IDP group mapping to be deleted")

		Eventually(func() error {
			err := manager.Client().Get(ctx, types.NamespacedName{
				Namespace: fixtures.IDPGroupMapping.Namespace,
				Name:      fixtures.IDPGroupMapping.Name,
			}, new(v1alpha1.IDPGroupMapping))
			if !errors.IsNotFound(err) {
				return assert.Equals("error", "[NOT FOUND]", err)
			}
			return nil
		}, timeout, interval).Should(Succeed())

		By("calling rest API, expecting group mapping to be removed from IDP configuration")

		apimClient := apim.NewClient(ctx)

		Eventually(func() error {
			idpConfig, err := apimClient.Configuration.GetIDPConfiguration(originalIDPID)
			if err != nil {
				return err
			}

			// Verify the group mapping no longer exists in the IDP configuration
			for _, gm := range idpConfig.GroupMappings {
				if gm.Condition == originalCondition {
					return assert.Equals("group mapping should not exist", false, true)
				}
			}
			return nil
		}, timeout, interval).Should(Succeed())
	})
})
