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

package v4

import (
	"context"

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

	It("should create a V4 API with a single group", func() {
		fixtures := fixture.
			Builder().
			WithAPIv4(constants.ApiV4).
			WithContext(constants.ContextWithCredentialsFile).
			Build()

		By("create group in APIM")
		groupName := random.GetName()
		apim := apim.NewClient(ctx)
		Expect(apim.Env.CreateGroup(&model.Group{Name: groupName})).To(Succeed())

		By("applying the API with created group")
		fixtures.APIv4.Spec.Groups = []string{groupName}
		fixtures = fixtures.Apply()

		By("verifying that group was created")
		expectedGroups := []string{groupName}

		Eventually(func() error {
			apiExport, err := apim.Export.V4Api(fixtures.APIv4.Status.ID)
			if err != nil {
				return err
			}

			exportedGroups := apiExport.Spec.Groups

			return assert.Equals("groups", expectedGroups, exportedGroups)
		}, timeout, interval).Should(Succeed(), fixtures.APIv4.Name)
	})
})
