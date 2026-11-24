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

package group

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/group"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	adm "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/group"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Update with ignore-members", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	It("should leave the group members to APIM", func() {
		fixtures := fixture.Builder().
			WithContext(constants.ContextWithCredentialsFile).
			WithGroup(constants.GroupFile).
			Build()
		fixtures.Group.SetAnnotations(map[string]string{core.IgnoreMembersAnnotation: "true"})
		fixtures = fixtures.Apply()

		Expect(assert.GroupCompleted(fixtures.Group)).To(Succeed())

		apim := apim.NewClient(ctx)
		hrid := refs.NewNamespacedNameFromObject(fixtures.Group).HRID()

		By("expecting the new group to get none of the CRD members")

		Eventually(func(g Gomega) {
			remote, err := apim.Env.GetGroupByHRID(ctx, hrid)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(remote.Members).To(BeEmpty())
		}, timeout, interval).Should(Succeed(), fixtures.Group.Name)

		By("setting the members in APIM, outside the operator")

		outside := fixtures.Group.DeepCopy()
		outside.SetAnnotations(nil)
		outside.PopulateIDs(nil, true)
		outside.Spec.Members = []group.Member{{
			Source:   "memory",
			SourceID: "application1",
			Roles:    map[group.RoleScope]string{group.APIRoleScope: "USER"},
		}}
		_, err := apim.Env.ImportGroup(ctx, outside)
		Expect(err).ToNot(HaveOccurred())

		By("expecting an update to be admitted although the members differ")

		updated := fixtures.Group.DeepCopy()
		updated.Spec.Name += "-renamed"

		Eventually(func() error {
			_, err := adm.AdmissionCtrl{}.ValidateUpdate(ctx, fixtures.Group, updated)
			return err
		}, timeout, interval).Should(Succeed(), fixtures.Group.Name)

		Eventually(func() error {
			return manager.UpdateSafely(ctx, updated)
		}, timeout, interval).Should(Succeed())

		By("expecting the reconcile to rename the group and keep the members set in APIM")

		Eventually(func(g Gomega) {
			remote, err := apim.Env.GetGroupByHRID(ctx, hrid)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(remote.Name).To(Equal(updated.Spec.Name))
			g.Expect(remote.Members).To(HaveLen(1))
			g.Expect(remote.Members[0].SourceId).To(Equal("application1"))
		}, timeout, interval).Should(Succeed(), fixtures.Group.Name)
	})
})
