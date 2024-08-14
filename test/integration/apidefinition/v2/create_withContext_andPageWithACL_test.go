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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
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

	It("should import page with ACLs in APIM", func() {
		fixtures := fixture.
			Builder().
			WithAPI(constants.ApiWithMarkdownPage).
			WithContext(constants.ContextWithCredentialsFile).
			Build()

		groupName := random.GetName()

		acl := []base.AccessControl{
			{
				ReferenceID:   groupName,
				ReferenceType: "GROUP",
			},
		}

		fixtures.API.Spec.Pages["markdown"].AccessControls = acl

		By("initializing a group in APIM")

		apim := apim.NewClient(ctx)

		Expect(apim.Env.CreateGroup(&model.Group{Name: groupName})).To(Succeed())

		By("applying resources")

		fixtures = fixtures.Apply()

		By("checking that exported API has page with ACL")

		Eventually(func() error {
			api, xErr := apim.Export.V2Api(fixtures.API.Status.ID)
			if xErr != nil {
				return xErr
			}

			page := api.Spec.Pages["hello-markdown"]

			if err := assert.NotNil("page", page); err != nil {
				return err
			}

			return assert.Equals("accessControls", page.AccessControls, acl)
		}, timeout, interval).Should(Succeed(), fixtures.API.Name)
	})
})
