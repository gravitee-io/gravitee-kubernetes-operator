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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	admissionv4 "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// GKO-3452: plans are sent in the order of plans.<key>.order, which APIM keeps.
var _ = Describe("Validate drift on the plan order", labels.WithContext, func() {
	ctx := context.Background()
	admissionCtrl := admissionv4.AdmissionCtrl{}

	It("should keep plan orders that are the reverse of the plan keys, without drift", func() {
		api := applyAPIWithPlanOrders(1, 2)

		By("checking APIM lists the plans in the order of the spec")
		Expect(remotePlanHRIDs(ctx, api)).To(Equal([]string{"KeyLess", "ApiKey"}))

		By("changing the API description after the first sync")
		newAPI := api.DeepCopy()
		setDescription(newAPI, updatedDescription)

		_, err := admissionCtrl.ValidateUpdate(ctx, api, newAPI)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should reorder existing plans when their orders are swapped, without drift", func() {
		api := applyAPIWithPlanOrders(1, 2)
		Expect(remotePlanHRIDs(ctx, api)).To(Equal([]string{"KeyLess", "ApiKey"}))

		By("swapping the plan orders in the spec")
		reordered := api.DeepCopy()
		plans := *reordered.Spec.Plans
		plans["KeyLess"].Order = new(2)
		plans["ApiKey"].Order = new(1)
		Expect(manager.UpdateSafely(ctx, reordered)).To(Succeed())

		By("checking APIM lists the plans in the new order")
		Eventually(func() ([]string, error) {
			return remotePlanHRIDs(ctx, reordered)
		}, constants.EventualTimeout, constants.Interval).Should(Equal([]string{"ApiKey", "KeyLess"}))

		By("changing the API description after the reorder")
		newAPI := reordered.DeepCopy()
		setDescription(newAPI, updatedDescription)

		_, err := admissionCtrl.ValidateUpdate(ctx, reordered, newAPI)
		Expect(err).ToNot(HaveOccurred())
	})
})

// applyAPIWithPlanOrders applies an API with a keyless plan and an API key plan, set to the given orders.
func applyAPIWithPlanOrders(keyLessOrder, apiKeyOrder int) *v1alpha1.ApiV4Definition {
	fixtures := fixture.Builder().
		WithAPIv4(constants.ApiV4WithContextFile).
		WithContext(constants.ContextWithCredentialsFile).
		Build()
	plans := *fixtures.APIv4.Spec.Plans
	plans["KeyLess"].Plan = &base.Plan{Order: new(keyLessOrder)}
	apiKeySecurity := v4.NewPlanSecurity("API_KEY")
	apiKey := v4.NewPlan().WithSecurity(&apiKeySecurity)
	apiKey.Name = "API key plan"
	apiKey.Order = new(apiKeyOrder)
	plans["ApiKey"] = apiKey
	fixtures.Apply()
	Expect(assert.ApiV4Completed(fixtures.APIv4)).To(Succeed())
	return fixtures.APIv4
}

func remotePlanHRIDs(ctx context.Context, api *v1alpha1.ApiV4Definition) ([]string, error) {
	remote, err := apim.NewClient(ctx).APIs.GetV4ByHRID(refs.NewNamespacedNameFromObject(api).HRID())
	if err != nil {
		return nil, err
	}
	return planHRIDs(remote.Plans), nil
}

func planHRIDs(plans model.APIV4PlansDTO) []string {
	hrids := make([]string, 0, len(plans))
	for _, plan := range plans {
		hrids = append(hrids, plan.HRID)
	}
	return hrids
}
