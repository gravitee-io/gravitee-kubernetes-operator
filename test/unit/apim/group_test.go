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

package apim_test

import (
	"context"
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/group"
	apimgroup "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/group"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

var _ = Describe("Group payload", func() {
	sent := func(dto apimgroup.DTO) map[string]any {
		GinkgoHelper()
		data, err := json.Marshal(dto)
		Expect(err).ToNot(HaveOccurred())
		var result map[string]any
		Expect(json.Unmarshal(data, &result)).To(Succeed())
		return result
	}

	It("carries the CRD fields and applies no SDK default", func() {
		grp := developers()
		grp.Spec.NotifyMembers = false
		grp.Spec.Members = []group.Member{{
			Source: "memory", SourceID: "api1",
			Roles: map[group.RoleScope]string{group.APIRoleScope: "OWNER"},
		}}

		Expect(sent(toDTO(grp))).To(Equal(map[string]any{
			"hrid":          "default-developers",
			"name":          "Developers",
			"notifyMembers": false,
			"members": []any{map[string]any{
				"source": "memory", "sourceId": "api1", "roles": map[string]any{"API": "OWNER"},
			}},
		}))
	})

	It("leaves default member roles unset when the CRD does not declare them", func() {
		Expect(toDTO(developers()).DefaultMemberRoles).To(BeNil())
	})

	It("does not send the CRD id, which is not part of the Automation API payload", func() {
		grp := developers()
		grp.Spec.ID = "4d1c5f0a-0000-0000-0000-000000000001"
		Expect(sent(toDTO(grp))).ToNot(HaveKey("id"))
	})

	It("sets the namespace-name hrid of a new group, without changing the CR", func() {
		grp := developers()
		grp.Spec.HRID = ""
		Expect(toDTO(grp).Identity()).To(Equal("default-developers"))
		Expect(grp.Spec.HRID).To(BeEmpty())
	})

	It("addresses a legacy group by the UUID in its status", func() {
		grp := developers()
		grp.Status.ID = "4d1c5f0a-0000-0000-0000-000000000001"
		Expect(toDTO(grp).Identity()).To(Equal(grp.Status.ID))
	})

	It("addresses a group by the UUID set in its spec when it has no status yet", func() {
		grp := developers()
		grp.Spec.ID = "4d1c5f0a-0000-0000-0000-000000000001"
		dto := toDTO(grp)
		Expect(dto.Identity()).To(Equal(grp.Spec.ID))
		Expect(dto.HridContainsUUID).To(BeTrue())
	})

	It("prefers the status UUID to the spec one", func() {
		grp := developers()
		grp.Spec.ID = "4d1c5f0a-0000-0000-0000-000000000002"
		grp.Status.ID = "4d1c5f0a-0000-0000-0000-000000000001"
		Expect(toDTO(grp).Identity()).To(Equal(grp.Status.ID))
	})

	It("addresses an Automation API managed group by its hrid, even with a status id", func() {
		grp := developers()
		grp.Status.ID = "4d1c5f0a-0000-0000-0000-000000000001"
		k8s.AddAutomationAPIManagedCondition(grp)
		Expect(toDTO(grp).Identity()).To(Equal("default-developers"))
	})

	It("keeps the members in the payload under ignore-members, and flags it", func() {
		grp := developers()
		grp.Spec.Members = []group.Member{{Source: "memory", SourceID: "api1"}}
		grp.Annotations = map[string]string{core.IgnoreMembersAnnotation: "true"}
		dto := toDTO(grp)
		Expect(dto.Members).To(HaveLen(1))
		Expect(dto.IgnoreMembers).To(BeTrue())
	})

	It("flags a legacy group known only by its spec id as unadopted", func() {
		grp := developers()
		grp.Spec.ID = "4d1c5f0a-0000-0000-0000-000000000001"
		Expect(toDTO(grp).Unadopted).To(BeTrue())
	})

	It("does not flag a group synced by UUID or by hrid as unadopted", func() {
		synced := developers()
		synced.Spec.ID = "4d1c5f0a-0000-0000-0000-000000000001"
		synced.Status.ID = synced.Spec.ID
		Expect(toDTO(synced).Unadopted).To(BeFalse())

		managed := developers()
		managed.Spec.ID = "4d1c5f0a-0000-0000-0000-000000000001"
		k8s.AddAutomationAPIManagedCondition(managed)
		Expect(toDTO(managed).Unadopted).To(BeFalse())

		Expect(toDTO(developers()).Unadopted).To(BeFalse())
	})
})

var _ = Describe("Group status", func() {
	ctx := context.Background()
	id := "4d1c5f0a-0000-0000-0000-000000000001"

	It("copies the id, member count, org and env", func() {
		grp := developers()
		resp := apimgroup.Response{
			State: sdk.GroupState{Id: &id, MemberCount: new(int64(2))},
			OrgID: "org1", EnvID: "env1",
		}
		Expect(apimgroup.UpdateStatus(ctx, grp, resp)).To(Succeed())
		Expect(grp.Status.ID).To(Equal(id))
		Expect(grp.Status.Members).To(Equal(uint(2)))
		Expect(grp.Status.OrgID).To(Equal("org1"))
		Expect(grp.Status.EnvID).To(Equal("env1"))
	})

	It("clears the deprecated processing status", func() {
		grp := developers()
		grp.Status.ProcessingStatus = core.ProcessingStatusCompleted
		Expect(apimgroup.UpdateStatus(ctx, grp, apimgroup.Response{})).To(Succeed())
		Expect(grp.Status.ProcessingStatus).To(BeEmpty())
	})

	It("defaults a missing member count to zero", func() {
		grp := developers()
		grp.Status.Members = 3
		Expect(apimgroup.UpdateStatus(ctx, grp, apimgroup.Response{})).To(Succeed())
		Expect(grp.Status.Members).To(BeZero())
	})

	It("never reports a negative member count", func() {
		grp := developers()
		grp.Status.Members = 3
		resp := apimgroup.Response{State: sdk.GroupState{MemberCount: new(int64(-1))}}
		Expect(apimgroup.UpdateStatus(ctx, grp, resp)).To(Succeed())
		Expect(grp.Status.Members).To(BeZero())
	})

	It("marks a group synced by hrid as Automation API managed", func() {
		grp := developers()
		Expect(apimgroup.UpdateStatus(ctx, grp, apimgroup.Response{State: sdk.GroupState{Id: &id}})).To(Succeed())
		Expect(k8s.IsAutomationAPIManaged(grp)).To(BeTrue())
	})

	It("leaves a legacy group, synced by UUID, unmanaged by the Automation API", func() {
		grp := developers()
		grp.Status.ID = id
		resp := apimgroup.Response{State: sdk.GroupState{Id: &id}, HridContainsUUID: true}
		Expect(apimgroup.UpdateStatus(ctx, grp, resp)).To(Succeed())
		Expect(k8s.IsAutomationAPIManaged(grp)).To(BeFalse())
	})

	It("leaves a group synced by the UUID in its spec, with no status yet, unmanaged", func() {
		grp := developers()
		grp.Spec.ID = id
		dto := toDTO(grp)
		resp := apimgroup.Response{State: sdk.GroupState{Id: &id}, HridContainsUUID: dto.HridContainsUUID}
		Expect(apimgroup.UpdateStatus(ctx, grp, resp)).To(Succeed())
		Expect(k8s.IsAutomationAPIManaged(grp)).To(BeFalse())
		Expect(grp.Status.ID).To(Equal(id))
	})
})
