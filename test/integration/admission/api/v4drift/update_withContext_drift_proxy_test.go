package v4

import (
	"context"

	admissionv4 "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/onsi/gomega"

	. "github.com/onsi/ginkgo/v2"
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

var _ = FDescribe("Validate drift for PROXY APIs", labels.WithContext, func() {
	ctx := context.Background()
	admissionCtrl := admissionv4.AdmissionCtrl{}

	It("should not drift on a simple update with minimal fields", func() {
		fixtures := fixture.Builder().
			WithAPIv4(constants.ApiV4WithContextFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("changing the API description")
		newAPI := fixtures.APIv4.DeepCopy()
		setDescription(newAPI, updatedDescription)

		_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.APIv4, newAPI)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	It("should detect drift with minimal fields", func() {
		fixtures := fixture.Builder().
			WithAPIv4(constants.ApiV4WithContextFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("changing the remote API description")
		newAPI := fixtures.APIv4.DeepCopy()

		By("changing the CRD API description")
		validateDescriptionDrift(ctx, admissionCtrl, fixtures.APIv4, newAPI, fixtures.Context)
	})

	It("should not drift on a simple update with all fields", func() {
		fixtures := fixture.Builder().
			WithContext(constants.ContextWithCredentialsFile).
			WithGroup(constants.GroupFile).
			WithNotification(constants.NotificationWithGroupFile).
			WithResource(constants.ApiResourceCacheRedisFile).
			WithSharedPolicyGroups(constants.SharedPolicyGroupsFile).
			WithAPIv4(constants.ApiV4ProxyDriftFullFile).
			Build()

		fixtures.APIv4.Spec.Resources[0].Ref.Name = fixtures.Resource.Name
		fixtures.APIv4.Spec.Flows[0].Response[0].SharedPolicyGroup.Name = fixtures.SharedPolicyGroup.Name
		fixtures.Apply()

		By("changing the API description")
		newAPI := fixtures.APIv4.DeepCopy()
		setDescription(newAPI, updatedDescription)

		_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.APIv4, newAPI)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	It("should detect drift with all fields", func() {
		fixtures := fixture.Builder().
			WithContext(constants.ContextWithCredentialsFile).
			WithGroup(constants.GroupFile).
			WithNotification(constants.NotificationWithGroupFile).
			WithResource(constants.ApiResourceCacheRedisFile).
			WithSharedPolicyGroups(constants.SharedPolicyGroupsFile).
			WithAPIv4(constants.ApiV4ProxyDriftFullFile).
			Build()

		fixtures.APIv4.Spec.Resources[0].Ref.Name = fixtures.Resource.Name
		fixtures.APIv4.Spec.Flows[0].Response[0].SharedPolicyGroup.Name = fixtures.SharedPolicyGroup.Name
		fixtures.Apply()

		By("changing the remote API description")
		newAPI := fixtures.APIv4.DeepCopy()

		By("changing the CRD API description")
		validateDescriptionDrift(ctx, admissionCtrl, fixtures.APIv4, newAPI, fixtures.Context)
	})
})
