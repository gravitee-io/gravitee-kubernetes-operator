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

package sharedpolicygroups

import (
	"context"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/utils/ptr"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	"k8s.io/apimachinery/pkg/api/errors"
)

var _ = Describe("Delete", labels.WithContext, func() {
	timeout := constants.EventualTimeout

	interval := constants.Interval

	ctx := context.Background()

	It("should delete only with no more API V4 references", func() {
		fixtures1 := fixture.Builder().
			WithContext(constants.ContextWithCredentialsFile).
			WithSharedPolicyGroups(constants.SharedPolicyGroupsFile).
			Build().
			Apply()

		fixtures2 := fixture.Builder().
			WithAPIv4(constants.ApiV4).
			Build()

		fixtures2.APIv4.Spec.Context = fixtures1.Context.GetNamespacedName()
		fixtures2.APIv4.Spec.Flows = []*v4.Flow{
			{
				Enabled: ptr.To(true),
				Request: []*v4.FlowStep{
					{
						SharedPolicyGroup: &refs.NamespacedName{
							Name: fixtures1.SharedPolicyGroup.Name,
						},
					},
				},
			},
		}

		fixtures2.Apply()

		By("deleting SharedPolicyGroup context")

		Expect(manager.Client().Delete(ctx, fixtures1.SharedPolicyGroup)).To(Succeed())

		By("expecting to still find SharedPolicyGroup")

		checkUntil := constants.ConsistentTimeout
		Consistently(func() error {
			kErr := manager.GetLatest(ctx, fixtures1.SharedPolicyGroup)
			return kErr
		}, checkUntil, interval).Should(Succeed())

		By("deleting the API V4 definition")

		Expect(manager.Client().Delete(ctx, fixtures2.APIv4)).To(Succeed())

		By("expecting SharedPolicyGroup to have been deleted")

		Eventually(func() error {
			kErr := manager.GetLatest(ctx, fixtures1.SharedPolicyGroup)
			if errors.IsNotFound(kErr) {
				return nil
			}
			return assert.Equals("error", "[NOT FOUND]", kErr)
		}, timeout, interval).Should(Succeed())
	})
})
