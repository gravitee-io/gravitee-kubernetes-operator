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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/group"
	adm "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/group"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate update", labels.WithContext, func() {
	interval := constants.Interval
	timeout := constants.EventualTimeout
	ctx := context.Background()

	admissionCtrl := adm.AdmissionCtrl{}

	DescribeTable("should return severe error with primary owner role",
		func(roleScope group.RoleScope) {
			fixtures := fixture.
				Builder().
				WithContext(constants.ContextWithCredentialsFile).
				WithGroup(constants.GroupFile).
				Build()

			By("initializing a service account in current organization")

			apim := apim.NewClient(ctx)
			saName := random.GetName()

			Expect(apim.Org.CreateUser(model.NewServiceAccount(saName))).To(Succeed())

			By("adding the sa to the Group")

			groupMember := base.NewGraviteeGroupMember(saName, "OWNER")
			fixtures.Group.Spec.Members = []group.Member{groupMember}

			fixtures = fixtures.Apply()
			obj := fixtures.Group
			newObj := obj.DeepCopy()

			By("adding the sa role to be primary owner for scope " + string(roleScope))

			newObj.Spec.Members[0].Roles[roleScope] = "PRIMARY_OWNER"

			Eventually(func() error {
				_, err := admissionCtrl.ValidateUpdate(ctx, obj, newObj)
				return assert.Equals(
					"error",
					errors.NewSevere("setting a member with the primary owner role is not allowed"),
					err,
				)
			}, timeout, interval).Should(Succeed(), fixtures.Group.Name)

			Expect(manager.Delete(ctx, fixtures.Group)).To(Succeed())
		},
		Entry("on API role scope", group.APIRoleScope),
		Entry("on application role scope", group.ApplicationRoleScope),
		Entry("on integration role scope", group.APIRoleScope),
	)
})
