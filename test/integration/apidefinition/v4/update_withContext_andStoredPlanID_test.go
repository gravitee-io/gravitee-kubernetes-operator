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
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/apidefinition"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

const planKey = "API_KEY"

// subscriptionStatus is the part of the Management API v2 subscription read this test needs.
type subscriptionStatus struct {
	Status string `json:"status"`
}

var _ = Describe("Update", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	It("should keep the stored plan id and its subscription", func() {
		fixtures := fixture.Builder().
			WithApplication(constants.Application).
			WithAPIv4(constants.ApiV4WithApiKeyPlanFile).
			WithContext(constants.ContextWithSecretFile).
			Build().
			Apply()

		Expect(assert.ApiV4Completed(fixtures.APIv4)).To(Succeed())
		Expect(assert.ApiV4Accepted(fixtures.APIv4)).To(Succeed())

		client := apim.NewClient(ctx)
		apiHRID := refs.NewNamespacedNameFromObject(fixtures.APIv4).HRID()
		storedPlanID := uuid.NewV4String()

		By("storing a plan id that APIM does not derive from the plan HRID")

		Eventually(func() error {
			if err := manager.GetLatest(ctx, fixtures.APIv4); err != nil {
				return err
			}
			if fixtures.APIv4.Status.Plans[planKey] == storedPlanID {
				return nil
			}
			if fixtures.APIv4.Status.Plans == nil {
				fixtures.APIv4.Status.Plans = map[string]string{}
			}
			fixtures.APIv4.Status.Plans[planKey] = storedPlanID
			if err := manager.Client().Status().Update(ctx, fixtures.APIv4); err != nil {
				return err
			}
			return fmt.Errorf("stored plan id not seen yet")
		}, timeout, interval).Should(Succeed())

		existing := fixtures.APIv4.DeepCopy()
		Expect(apidefinition.PrepareV4SpecForAutomation(ctx, existing, false)).To(Succeed())
		existing.PopulateIDs(client.Context, true)
		plan := (*existing.Spec.Plans)[planKey]
		if plan.Plan == nil {
			plan.Plan = &base.Plan{}
		}
		plan.ID = storedPlanID
		_, err := client.APIs.ImportV4(existing)
		Expect(err).ToNot(HaveOccurred())

		Expect(planIDInAPIM(client, apiHRID)).To(Equal(storedPlanID))

		By("subscribing an application to the stored plan")

		apiFromAPIM, err := client.APIs.GetV4ByHRID(apiHRID)
		Expect(err).ToNot(HaveOccurred())
		subscription, err := client.Subscriptions.Subscribe(
			apiFromAPIM.ID, fixtures.Application.Status.ID, storedPlanID,
		)
		Expect(err).ToNot(HaveOccurred())

		Eventually(func() (string, error) {
			return subscriptionStatusInAPIM(client, apiFromAPIM.ID, subscription.ID)
		}, timeout, interval).Should(Equal("ACCEPTED"))

		By("updating the API with a different plan id in its spec")

		description := "updated with a stored plan id"
		updated := fixtures.APIv4.DeepCopy()
		updated.Spec.Description = &description
		specPlan := (*updated.Spec.Plans)[planKey]
		if specPlan.Plan == nil {
			specPlan.Plan = &base.Plan{}
		}
		specPlan.ID = uuid.NewV4String()
		Expect(manager.UpdateSafely(ctx, updated)).To(Succeed())

		Eventually(func() error {
			api, getErr := client.APIs.GetV4ByHRID(apiHRID)
			if getErr != nil {
				return getErr
			}
			return assert.Equals("description", &description, api.Description)
		}, timeout, interval).Should(Succeed())

		By("expecting the stored plan and its subscription to be kept")

		Expect(planIDInAPIM(client, apiHRID)).To(Equal(storedPlanID))
		Expect(subscriptionStatusInAPIM(client, apiFromAPIM.ID, subscription.ID)).To(Equal("ACCEPTED"))
	})
})

func planIDInAPIM(client *apim.APIM, apiHRID string) (string, error) {
	api, err := client.APIs.GetV4ByHRID(apiHRID)
	if err != nil {
		return "", err
	}
	for _, plan := range api.Plans {
		if plan.HRID == planKey {
			return plan.ID, nil
		}
	}
	return "", fmt.Errorf("plan %s not found in API %s", planKey, apiHRID)
}

func subscriptionStatusInAPIM(client *apim.APIM, apiID, subscriptionID string) (string, error) {
	url := client.APIs.EnvV2Target("apis").WithPath(apiID).WithPath("subscriptions").WithPath(subscriptionID)
	sub := new(subscriptionStatus)
	if err := client.APIs.HTTP.Get(url.String(), sub); err != nil {
		return "", err
	}
	return sub.Status, nil
}
