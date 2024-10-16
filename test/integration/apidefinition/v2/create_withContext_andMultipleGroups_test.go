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

package v2

import (
	"context"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	It("should create an API with multiple groups and then remove one group", func() {
		fixtures3GroupsAPI := fixture.
			Builder().
			WithAPI(constants.ApiWithMembersAndGroups).
			WithContext(constants.ContextWithCredentialsFile).
			Build()

		By("create groups in APIM")
		groupName1 := random.GetName()
		groupName2 := random.GetName()
		groupName3 := random.GetName()
		apim := apim.NewClient(ctx)
		Expect(apim.Env.CreateGroup(&model.Group{Name: groupName1})).To(Succeed())
		Expect(apim.Env.CreateGroup(&model.Group{Name: groupName2})).To(Succeed())
		Expect(apim.Env.CreateGroup(&model.Group{Name: groupName3})).To(Succeed())

		By("applying the API with created group")
		fixtures3GroupsAPI.API.Spec.Groups = []string{groupName1, groupName2, groupName3}
		fixtures3GroupsAPI = fixtures3GroupsAPI.Apply()

		By("verifying that groups were created")
		expected3Groups := []string{groupName1, groupName2, groupName3}
		Eventually(func() error {
			apiExport, err := apim.Export.V2Api(fixtures3GroupsAPI.API.Status.ID)
			if err != nil {
				return err
			}

			exported3Groups := apiExport.Spec.Groups
			return assert.SliceEqualsSorted("groups", expected3Groups, exported3Groups, strings.Compare)
		}, timeout, interval).Should(Succeed(), fixtures3GroupsAPI.API.Name)

		By("removing one group")
		fixtures2GroupsAPI := fixture.
			Builder().
			WithAPI(constants.ApiWithMembersAndGroups).
			WithContext(constants.ContextWithCredentialsFile).
			Build()

		fixtures2GroupsAPI.API.Spec.Groups = []string{groupName1, groupName3}
		fixtures2GroupsAPI = fixtures2GroupsAPI.Apply()

		By("verifying that groups were updated")
		expected2Groups := []string{groupName1, groupName3}
		Eventually(func() error {
			updatedApiExport, err := apim.Export.V2Api(fixtures2GroupsAPI.API.Status.ID)
			if err != nil {
				return err
			}

			exported2Groups := updatedApiExport.Spec.Groups
			return assert.SliceEqualsSorted("groups", expected2Groups, exported2Groups, strings.Compare)
		}, timeout, interval).Should(Succeed(), fixtures2GroupsAPI.API.Name)
	})
})
