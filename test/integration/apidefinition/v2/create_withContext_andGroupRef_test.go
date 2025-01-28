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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/group"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
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

	It("should add group reference to API", func() {
		fixtures := fixture.
			Builder().
			WithContext(constants.ContextWithCredentialsFile).
			WithGroup(constants.GroupFile).
			WithAPI(constants.ApiWithContextFile).
			Build()

		By("initializing a service account in current organization")

		apim := apim.NewClient(ctx)
		saName := random.GetName()
		Expect(apim.Org.CreateUser(model.NewServiceAccount(saName))).To(Succeed())

		By("adding the sa to the Group")

		groupMember := base.NewGraviteeGroupMember(saName, "OWNER")
		fixtures.Group.Spec.Members = []group.Member{groupMember}

		fixtures.API.Spec.GroupRefs = []refs.NamespacedName{refs.NewNamespacedName(
			fixtures.Group.Namespace,
			fixtures.Group.Name,
		)}

		fixtures = fixtures.Apply()

		By("checking that group has been assigned in exported API")

		expectedGroups := []string{fixtures.Group.Spec.Name}

		Eventually(func() error {
			apiExport, err := apim.Export.V2Api(fixtures.API.Status.ID)
			if err != nil {
				return err
			}

			apiGroups := apiExport.Spec.Groups

			return assert.SliceEqualsSorted(
				"groups",
				expectedGroups, apiGroups,
				strings.Compare,
			)
		}, timeout, interval).Should(Succeed(), fixtures.API.Name)
	})
})
