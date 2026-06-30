package group

import (
	"context"

	admission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/group"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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

var _ = Describe("Validate drift", labels.WithContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := admission.AdmissionCtrl{}

	It("should not drift on a simple update", func() {
		fixtures := fixture.
			Builder().
			WithGroup(constants.GroupFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("changing the group name")
		newGroup := fixtures.Group.DeepCopy()
		newGroup.Spec.Name = "updated name"

		Eventually(func() error {
			_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Group, newGroup)
			return err
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	It("should detect drift", func() {
		fixtures := fixture.
			Builder().
			WithGroup(constants.GroupFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("changing the remote group name")

		newGroup := fixtures.Group.DeepCopy()
		newGroup.PopulateIDs(nil, true)

		newGroup.Spec.Name = "remote updated name"

		apim := apim.NewClient(ctx)

		_, err := apim.Env.ImportGroup(newGroup)
		Expect(err).ToNot(HaveOccurred())

		By("changing the CRD group name")

		newGroup.Spec.Name = "local CRD name"

		By("checking that group validation returns error")

		Eventually(func() error {
			_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Group, newGroup)
			return assert.DriftDetected("name: \"local CRD name\" != \"remote updated name\"", err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})
})
