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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
)

var _ = Describe("Create", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	It("should create the group in APIM by its HRID", func() {
		fixtures := fixture.Builder().
			WithContext(constants.ContextWithCredentialsFile).
			WithGroup(constants.GroupFile).
			Build().
			Apply()

		By("expecting group status to be completed and managed by the Automation API")

		Expect(assert.GroupCompleted(fixtures.Group)).To(Succeed())
		Expect(assert.GroupAccepted(fixtures.Group)).To(Succeed())
		Expect(assert.ManagedByAutomationAPI(fixtures.Group)).To(Succeed())

		By("calling the Automation API, expecting the group to match the CRD and its status")

		apim := apim.NewClient(ctx)
		hrid := refs.NewNamespacedNameFromObject(fixtures.Group).HRID()

		Eventually(func(g Gomega) {
			remote, err := apim.Env.GetGroupByHRID(ctx, hrid)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(remote.Name).To(Equal(fixtures.Group.Spec.Name))
			g.Expect(remote.Id).ToNot(BeNil())
			g.Expect(fixtures.Group.Status.ID).To(Equal(*remote.Id))
			g.Expect(remote.Members).To(HaveLen(len(fixtures.Group.Spec.Members)))
		}, timeout, interval).Should(Succeed(), fixtures.Group.Name)

		By("expecting group event to have been emitted")

		assert.EventsEmitted(fixtures.Group, "UpdateStarted", "UpdateSucceeded")
	})
})
