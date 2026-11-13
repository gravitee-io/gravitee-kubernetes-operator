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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	adm "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/group"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
)

var _ = Describe("Validate create", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	It("should admit a new group the APIM dry run accepts, and persist nothing", func() {
		fixtures := fixture.Builder().
			WithContext(constants.ContextWithCredentialsFile).
			WithGroup(constants.GroupFile).
			Build()
		grp := fixtures.Group
		fixtures.Group = nil
		fixtures = fixtures.Apply()
		grp.Spec.Context = fixtures.Context.GetNamespacedName()

		Eventually(func() error {
			_, err := adm.AdmissionCtrl{}.ValidateCreate(ctx, grp)
			return err
		}, timeout, interval).Should(Succeed(), grp.Name)

		By("calling the Automation API, expecting the dry run to have created nothing")

		_, err := apim.NewClient(ctx).Env.GetGroupByHRID(ctx, refs.NewNamespacedNameFromObject(grp).HRID())
		Expect(assert.NotFoundError(err)).To(Succeed())
	})
})
