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
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/group"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
)

var _ = Describe("Group payload", func() {
	sent := func(payload sdk.Group) map[string]any {
		GinkgoHelper()
		data, err := json.Marshal(payload)
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

		payload, err := model.ToGroup(grp)
		Expect(err).ToNot(HaveOccurred())

		Expect(sent(payload)).To(Equal(map[string]any{
			"hrid":          "default-developers",
			"name":          "Developers",
			"notifyMembers": false,
			"members": []any{map[string]any{
				"source": "memory", "sourceId": "api1", "roles": map[string]any{"API": "OWNER"},
			}},
		}))
	})

	It("leaves default member roles unset when the CRD does not declare them", func() {
		payload, err := model.ToGroup(developers())
		Expect(err).ToNot(HaveOccurred())
		Expect(payload.DefaultMemberRoles).To(BeNil())
	})

	It("does not send the CRD id, which is not part of the Automation API payload", func() {
		grp := developers()
		grp.Spec.ID = "4d1c5f0a-0000-0000-0000-000000000001"
		payload, err := model.ToGroup(grp)
		Expect(err).ToNot(HaveOccurred())
		Expect(sent(payload)).ToNot(HaveKey("id"))
	})
})

var _ = Describe("Group status", func() {
	It("defaults a missing member count to zero", func() {
		Expect(model.ToGroupStatus(sdk.GroupState{}).Members).To(BeZero())
	})

	It("never reports a negative member count", func() {
		Expect(model.ToGroupStatus(sdk.GroupState{MemberCount: new(int64(-1))}).Members).To(BeZero())
	})
})
