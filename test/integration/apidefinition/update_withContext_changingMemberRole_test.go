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

package apidefinition

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
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

var _ = Describe("Update", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	It("should change the role of an API member", func() {
		fixtures := fixture.
			Builder().
			WithAPI(constants.ApiWithMembersAndGroups).
			WithContext(constants.ContextWithCredentialsFile).
			Build()

		By("initializing a service account in current organization")

		apim := apim.NewClient(ctx)

		saName := random.GetName()

		Expect(apim.Org.CreateUser(model.NewServiceAccount(saName))).To(Succeed())

		By("applying the API with created service account as members")

		primaryOwner := base.NewMemoryMember("admin", "PRIMARY_OWNER")

		saMember := base.NewGraviteeMember(saName, "REVIEWER")

		fixtures.API.Spec.Members = []*base.Member{saMember}

		fixtures = fixtures.Apply()

		By("checking that exported API has two members")

		Eventually(func() error {
			export, err := apim.Export.V2Api(fixtures.API.Status.ID)
			if err != nil {
				return err
			}
			return assert.Equals("members", []*base.Member{primaryOwner, saMember}, export.Spec.Members)
		}, timeout, interval).Should(Succeed(), fixtures.API.Name)

		By("changing role of API member from REVIEWER to USER")
		saMember.Role = "USER"
		expectedMembers := []*base.Member{primaryOwner, saMember}
		fixtures.API.Spec.Members = expectedMembers

		Eventually(func() error {
			return manager.UpdateSafely(ctx, fixtures.API)
		}, timeout, interval).Should(Succeed(), fixtures.API.Name)

		By("checking that exported API has one member left")

		Eventually(func() error {
			export, err := apim.Export.V2Api(fixtures.API.Status.ID)
			if err != nil {
				return err
			}
			return assert.Equals("members", expectedMembers, export.Spec.Members)
		}, timeout, interval).Should(Succeed(), fixtures.API.Name)
	})
})
