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

package apim

import (
	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/group"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func roles(api string) *sdk.GroupDefaultMemberRoles {
	return &sdk.GroupDefaultMemberRoles{Api: &api}
}

var _ = Describe("Group Drift detection", func() {

	DescribeTable("equivalent values",
		func(crd, remote any) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
		},
		Entry("empty struct",
			sdk.Group{},
			sdk.Group{},
		),
		Entry("equal struct",
			completeGroup(),
			completeGroup(),
		),
		Entry("notifyMembers ignored and empty members equivalent to nil",
			sdk.Group{
				Hrid: "my-group",
				Name: "My Group",
			},
			sdk.Group{
				Hrid:          "my-group",
				Name:          "My Group",
				NotifyMembers: new(true),
				Members:       make([]sdk.GroupMember, 0),
			},
		),
		Entry("crd-only members ignored, paired through their identity",
			sdk.Group{
				Name: "My Group",
				Members: []sdk.GroupMember{
					{Source: "gravitee", SourceId: "admin"},
					{Source: "memory", SourceId: "local-only"},
				},
			},
			sdk.Group{
				Name: "My Group",
				Members: []sdk.GroupMember{
					{Source: "gravitee", SourceId: "admin"},
				},
			},
		),
		Entry("crd-only member listed first is ignored",
			sdk.Group{
				Name: "My Group",
				Members: []sdk.GroupMember{
					{Source: "memory", SourceId: "local-only"},
					{Source: "gravitee", SourceId: "admin"},
				},
			},
			sdk.Group{
				Name: "My Group",
				Members: []sdk.GroupMember{
					{Source: "gravitee", SourceId: "admin"},
				},
			},
		),
		Entry("undeclared default member roles accept any remote defaults",
			sdk.Group{Name: "My Group"},
			sdk.Group{Name: "My Group", DefaultMemberRoles: roles("USER")},
		),
		Entry("empty default member roles match a group without defaults",
			sdk.Group{Name: "My Group", DefaultMemberRoles: &sdk.GroupDefaultMemberRoles{}},
			sdk.Group{Name: "My Group"},
		),
		Entry("empty default member roles match the empty object APIM answers",
			sdk.Group{Name: "My Group", DefaultMemberRoles: &sdk.GroupDefaultMemberRoles{}},
			sdk.Group{Name: "My Group", DefaultMemberRoles: &sdk.GroupDefaultMemberRoles{}},
		),
		Entry("undeclared default member roles accept the empty object APIM answers",
			sdk.Group{Name: "My Group"},
			sdk.Group{Name: "My Group", DefaultMemberRoles: &sdk.GroupDefaultMemberRoles{}},
		),
		Entry("crd members against empty remote are ignored",
			sdk.Group{
				Name: "e2e-group-no-roles",
				Members: []sdk.GroupMember{
					{Source: "gravitee", SourceId: "admin"},
				},
			},
			sdk.Group{
				Name:    "e2e-group-no-roles",
				Members: []sdk.GroupMember{},
			},
		),
		Entry("member roles defaulted to {} by the CRD match a remote member without roles",
			sdk.Group{
				Name:    "My Group",
				Members: []sdk.GroupMember{{Source: "gravitee", SourceId: "admin", Roles: map[string]string{}}},
			},
			sdk.Group{
				Name:    "My Group",
				Members: []sdk.GroupMember{{Source: "gravitee", SourceId: "admin"}},
			},
		),
	)

	DescribeTable("drifted values",
		func(crd, remote sdk.Group) {
			result := drift.DetectWithNamespace(crd, remote, "")
			Expect(result.DriftDetected()).To(BeTrue())
		},
		Entry("changed role",
			sdk.Group{Name: "My Group", DefaultMemberRoles: roles("USER")},
			sdk.Group{Name: "My Group", DefaultMemberRoles: roles("OWNER")},
		),
		Entry("scope added in APIM",
			sdk.Group{Name: "My Group", DefaultMemberRoles: roles("USER")},
			sdk.Group{Name: "My Group", DefaultMemberRoles: &sdk.GroupDefaultMemberRoles{
				Api: new("USER"), Application: new("USER"),
			}},
		),
		Entry("removed in APIM",
			sdk.Group{Name: "My Group", DefaultMemberRoles: roles("USER")},
			sdk.Group{Name: "My Group"},
		),
		Entry("removed in APIM, which answers an empty object",
			sdk.Group{Name: "My Group", DefaultMemberRoles: roles("USER")},
			sdk.Group{Name: "My Group", DefaultMemberRoles: &sdk.GroupDefaultMemberRoles{}},
		),
		Entry("added in APIM while the CRD declares none",
			sdk.Group{Name: "My Group", DefaultMemberRoles: &sdk.GroupDefaultMemberRoles{}},
			sdk.Group{Name: "My Group", DefaultMemberRoles: roles("USER")},
		),
		Entry("role of a shared member changed in APIM",
			sdk.Group{
				Name:    "My Group",
				Members: []sdk.GroupMember{{Source: "gravitee", SourceId: "admin", Roles: map[string]string{"API": "USER"}}},
			},
			sdk.Group{
				Name:    "My Group",
				Members: []sdk.GroupMember{{Source: "gravitee", SourceId: "admin", Roles: map[string]string{"API": "OWNER"}}},
			},
		),
		Entry("member removed in APIM is not hidden by a crd-only member",
			sdk.Group{
				Name: "My Group",
				Members: []sdk.GroupMember{
					{Source: "gravitee", SourceId: "admin", Roles: map[string]string{"API": "USER"}},
				},
			},
			sdk.Group{
				Name: "My Group",
				Members: []sdk.GroupMember{
					{Source: "gravitee", SourceId: "admin", Roles: map[string]string{"API": "OWNER"}},
					{Source: "memory", SourceId: "remote-only"},
				},
			},
		),
	)

	Describe("State turned into the payload type by Spec()", func() {
		It("matches the CRD it was synced from, state-only fields dropped", func() {
			crd, err := model.ToGroup(&v1alpha1.Group{
				ObjectMeta: metav1.ObjectMeta{Name: "developers", Namespace: "default"},
				Spec: v1alpha1.GroupSpec{Type: &group.Type{
					HRID:          "default-developers",
					Name:          "Developers",
					NotifyMembers: true,
					Members: []group.Member{{
						Source: "memory", SourceID: "api1",
						Roles: map[group.RoleScope]string{group.APIRoleScope: "OWNER"},
					}},
				}},
			})
			Expect(err).ToNot(HaveOccurred())

			state := parseJSON[sdk.GroupState](`{
				"id": "4d1c5f0a-0000-0000-0000-000000000001",
				"organizationId": "DEFAULT",
				"environmentId": "DEFAULT",
				"hrid": "default-developers",
				"name": "Developers",
				"notifyMembers": true,
				"memberCount": 1,
				"errors": {"warning": ["something defaulted"]},
				"defaultMemberRoles": {},
				"members": [{"source": "memory", "sourceId": "api1", "roles": {"API": "OWNER"}}]
			}`)
			expectNoDrift(drift.DetectWithNamespace(crd, state.Spec(), ""))
		})
	})

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(completeGroup(), completeGroup())
		})
	})
})

func completeGroup() sdk.Group {
	GinkgoHelper()
	return loadFixture[sdk.Group]("group_dto.json")
}
