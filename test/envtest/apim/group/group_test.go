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

package apimgroup

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	admission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/group"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Group on the APIM SDK", func() {
	ctx := context.Background()

	It("creates the group in APIM by its HRID and reports its id and member count", func() {
		grp := groupFixtures(constants.GroupFile).Apply().Group

		Expect(assert.GroupCompleted(grp)).To(Succeed())
		Expect(k8s.IsAutomationAPIManaged(grp)).To(BeTrue())

		status, remote := remoteGroup(ctx, grp)
		Expect(status).To(Equal(http.StatusOK))
		Expect(remote.Name).To(Equal(grp.Spec.Name))
		Expect(remote.Id).ToNot(BeNil())
		Expect(grp.Status.ID).To(Equal(*remote.Id))
		Expect(grp.Status.OrgID).To(Equal(apim.MockOrg))
		Expect(grp.Status.EnvID).To(Equal(apim.MockEnv))
		if remote.MemberCount != nil {
			Expect(grp.Status.Members).To(Equal(uint(*remote.MemberCount)))
		}
	})

	It("updates the group in APIM", func() {
		grp := groupFixtures(constants.GroupFile).Apply().Group
		Expect(assert.GroupCompleted(grp)).To(Succeed())

		updated := grp.DeepCopy()
		updated.Spec.Name = grp.Spec.Name + "-renamed"
		Expect(manager.Client().Update(ctx, updated)).To(Succeed())

		Eventually(ctx, func(g Gomega) {
			status, remote := remoteGroup(ctx, grp)
			g.Expect(status).To(Equal(http.StatusOK))
			g.Expect(remote.Name).To(Equal(updated.Spec.Name))
		}, constants.EventualTimeout, constants.Interval).Should(Succeed())
	})

	It("deletes the group from APIM", func() {
		grp := groupFixtures(constants.GroupFile).Apply().Group
		Expect(assert.GroupCompleted(grp)).To(Succeed())

		Expect(manager.Client().Delete(ctx, grp.DeepCopy())).To(Succeed())
		Eventually(ctx, func() error {
			return assert.Deleted(ctx, "Group", grp)
		}, constants.EventualTimeout, constants.Interval).Should(Succeed())

		status, _ := remoteGroup(ctx, grp)
		Expect(status).To(Equal(http.StatusNotFound))
	})

	It("admits a new group the APIM dry run accepts", func() {
		objects := groupFixtures(constants.GroupFile)
		grp := objects.Group
		objects.Group = nil
		objects.Apply()
		grp.Spec.Context = objects.Context.GetNamespacedName()

		_, err := admission.AdmissionCtrl{}.ValidateCreate(ctx, grp)
		Expect(err).ToNot(HaveOccurred())

		By("expecting the dry run to persist nothing")
		status, _ := remoteGroup(ctx, grp)
		Expect(status).To(Equal(http.StatusNotFound))
	})

	It("rejects an update when the group was changed in APIM", func() {
		grp := groupFixtures(constants.GroupFile).Apply().Group
		Expect(assert.GroupCompleted(grp)).To(Succeed())

		By("renaming the group in APIM, outside the operator")
		resp, err := apim.NewMockSDKClient().UpsertGroupWithResponse(ctx, nil, sdk.Group{
			Hrid: refs.NewNamespacedNameFromObject(grp).HRID(),
			Name: "renamed-in-apim",
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode()).To(Equal(http.StatusOK))

		By("expecting an update that does not realign with APIM to be rejected")
		updated := grp.DeepCopy()
		updated.Spec.NotifyMembers = !grp.Spec.NotifyMembers
		_, err = admission.AdmissionCtrl{}.ValidateUpdate(ctx, grp, updated)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("renamed-in-apim"))
	})

	It("leaves members to APIM when the ignore-members annotation is set", func() {
		objects := groupFixtures(constants.GroupFile)
		objects.Group.SetAnnotations(map[string]string{core.IgnoreMembersAnnotation: "true"})
		grp := objects.Apply().Group
		Expect(assert.GroupCompleted(grp)).To(Succeed())

		By("expecting the new group to get none of the CRD members")
		status, remote := remoteGroup(ctx, grp)
		Expect(status).To(Equal(http.StatusOK))
		Expect(remote.Members).To(BeEmpty())

		By("adding a member in APIM, outside the operator")
		outside := remote.Spec()
		outside.Members = []sdk.GroupMember{{
			Source: "memory", SourceId: "added-in-apim", Roles: map[string]string{"API": "USER"},
		}}
		resp, err := apim.NewMockSDKClient().UpsertGroupWithResponse(ctx, nil, outside)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode()).To(Equal(http.StatusOK))

		By("expecting an update to be admitted although the members differ")
		updated := grp.DeepCopy()
		updated.Spec.Name = grp.Spec.Name + "-renamed"
		_, err = admission.AdmissionCtrl{}.ValidateUpdate(ctx, grp, updated)
		Expect(err).ToNot(HaveOccurred())
		Expect(manager.Client().Update(ctx, updated)).To(Succeed())

		By("expecting the reconcile to rename the group and keep the members set in APIM")
		Eventually(ctx, func(g Gomega) {
			status, remote := remoteGroup(ctx, grp)
			g.Expect(status).To(Equal(http.StatusOK))
			g.Expect(remote.Name).To(Equal(updated.Spec.Name))
			g.Expect(remote.Members).To(Equal(outside.Members))
		}, constants.EventualTimeout, constants.Interval).Should(Succeed())
	})
})
