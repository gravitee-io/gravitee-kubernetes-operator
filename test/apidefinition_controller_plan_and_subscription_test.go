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

package test

import (
	"net/http"
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/uuid"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model/api"
)

var _ = Describe("Checking ApiKey plan and subscription", Ordered, func() {

	httpClient := http.Client{Timeout: 5 * time.Second}

	Context("with plan and subscription", Ordered, func() {
		var apiDefinitionFixture *v1alpha1.ApiDefinition

		var savedApiDefinition *v1alpha1.ApiDefinition

		var apiLookupKey types.NamespacedName
		var contextLookupKey types.NamespacedName

		var gatewayEndpoint string
		var apim *internal.APIM

		BeforeAll(func() {
			By("Create a management context to synchronize with the REST API")

			fixtureGenerator := internal.NewFixtureGenerator()

			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Api:     internal.ApiWithApiKeyPlanFile,
				Context: internal.ClusterContextFile,
			})

			Expect(err).ToNot(HaveOccurred())

			managementContext := fixtures.Context
			Expect(k8sClient.Create(ctx, managementContext)).Should(Succeed())

			contextLookupKey = types.NamespacedName{Name: managementContext.Name, Namespace: namespace}
			Eventually(func() error {
				return k8sClient.Get(ctx, contextLookupKey, managementContext)
			}, timeout, interval).Should(Succeed())

			By("Creating an API definition resource started by default")

			apiDefinition := fixtures.Api
			Expect(k8sClient.Create(ctx, apiDefinition)).Should(Succeed())

			apiDefinitionFixture = apiDefinition
			apiLookupKey = types.NamespacedName{Name: apiDefinitionFixture.Name, Namespace: namespace}

			By("Expecting the API Definition to be Ready")

			savedApiDefinition = new(v1alpha1.ApiDefinition)
			Eventually(func() error {
				err = k8sClient.Get(ctx, apiLookupKey, savedApiDefinition)
				return internal.AssertNoErrorAndStatusCompleted(err, savedApiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())
			gatewayEndpoint = internal.GatewayUrl + savedApiDefinition.Spec.Proxy.VirtualHosts[0].Path

			apim, err = internal.NewAPIM(ctx)
			Expect(err).ToNot(HaveOccurred())

		})

		It("Should return unauthorize without subscription", func() {
			Eventually(func() error {
				res, callErr := httpClient.Get(gatewayEndpoint)
				return internal.AssertNoErrorAndHTTPStatus(callErr, res, http.StatusUnauthorized)
			}, timeout, interval).ShouldNot(HaveOccurred())
		})

		It("Should return success with subscription", func() {
			By("Create a subscription and get api key")
			apiKey := createSubscriptionAndGetApiKey(
				apim,
				savedApiDefinition,
				func(mgmtApi *api.Entity) string {
					plans, err := apim.Plans.ListByAPI(mgmtApi.ID)
					Expect(err).ToNot(HaveOccurred())
					Expect(len(plans.Data)).To(Equal(1))
					return plans.Data[0].ID
				},
			)

			By("Call gateway with subscription api key")
			Eventually(func() error {
				res, callErr := getWithGioApiKey(&httpClient, gatewayEndpoint, apiKey)
				return internal.AssertNoErrorAndHTTPStatus(callErr, res, http.StatusOK)
			}, timeout, interval).ShouldNot(HaveOccurred())

		})

		It("should validate API key", func() {
			// Update savedApiDefinition & global var with last Get
			gatewayEndpoint = internal.GatewayUrl + savedApiDefinition.Spec.Proxy.VirtualHosts[0].Path

			By("Updating ApiDefinition to add ApiKey plan")

			newApiKeyPlan := v2.NewPlan(
				base.
					NewPlan("G.K.O. API Key Plan - 2", "").
					WithStatus(base.PublishedPlanStatus),
			).WithSecurity("API_KEY")

			updatedApiDefinition := savedApiDefinition.DeepCopy()
			updatedApiDefinition.Spec.Plans = append(savedApiDefinition.Spec.Plans,
				newApiKeyPlan,
			)

			Eventually(func() error {
				update := new(v1alpha1.ApiDefinition)
				if err := k8sClient.Get(ctx, apiLookupKey, update); err != nil {
					return err
				}
				updatedApiDefinition.Spec.DeepCopyInto(&update.Spec)
				return k8sClient.Update(ctx, update)
			}, timeout, interval).ShouldNot(HaveOccurred())

			// Wait for the ApiDefinition to be updated
			Eventually(func() error {
				k8sErr := k8sClient.Get(ctx, apiLookupKey, updatedApiDefinition)
				return internal.AssertNoErrorAndStatusCompleted(k8sErr, updatedApiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Subscribing to the API Key plan")

			apiKey := createSubscriptionAndGetApiKey(
				apim,
				updatedApiDefinition,
				func(mgmtApi *api.Entity) string {
					plans, err := apim.Plans.ListByAPI(mgmtApi.ID)
					Expect(err).ToNot(HaveOccurred())
					Expect(len(plans.Data)).To(Equal(2))
					planID := ""
					for _, plan := range plans.Data {
						if plan.Name == newApiKeyPlan.Name {
							planID = plan.ID
						}
					}
					Expect(planID).ToNot(BeEmpty())
					return planID
				},
			)

			By("Call gateway with subscription api key of second plan")
			Eventually(func() error {
				res, callErr := getWithGioApiKey(&httpClient, gatewayEndpoint, apiKey)
				return internal.AssertNoErrorAndHTTPStatus(callErr, res, http.StatusOK)
			}, timeout, interval).ShouldNot(HaveOccurred())
		})

		It("Should delete ApiDefinition resource", func() {
			By("Delete the ApiDefinition resource")
			err := k8sClient.Delete(ctx, apiDefinitionFixture)
			Expect(err).ToNot(HaveOccurred())

			By("Call deleted API definition URL and expect 404")
			Eventually(func() error {
				res, callErr := httpClient.Get(gatewayEndpoint)
				return internal.AssertNoErrorAndHTTPStatus(callErr, res, http.StatusNotFound)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Get the API definition from ManagementApi and expect deleted state")

			Eventually(func() error {
				_, apiErr := apim.APIs.GetByID(savedApiDefinition.Status.ID)
				return errors.IgnoreNotFound(apiErr)
			}, timeout, interval).ShouldNot(HaveOccurred())
		})
	})

	Context("with no plan", Ordered, func() {
		var apiDefinitionFixture *v1alpha1.ApiDefinition
		var savedApiDefinition *v1alpha1.ApiDefinition
		var apiLookupKey types.NamespacedName
		var contextLookupKey types.NamespacedName
		var gatewayEndpoint string
		var apim *internal.APIM

		BeforeAll(func() {
			By("Create a management context to synchronize with the REST API")

			fixtureGenerator := internal.NewFixtureGenerator()

			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Api:     internal.ApiWithContextNoPlanFile,
				Context: internal.ClusterContextFile,
			})

			Expect(err).ToNot(HaveOccurred())

			managementContext := fixtures.Context
			Expect(k8sClient.Create(ctx, managementContext)).Should(Succeed())

			contextLookupKey = types.NamespacedName{Name: managementContext.Name, Namespace: namespace}
			Eventually(func() error {
				return k8sClient.Get(ctx, contextLookupKey, managementContext)
			}, timeout, interval).Should(Succeed())

			By("Create an API definition resource stared by default")

			apiDefinition := fixtures.Api
			Expect(k8sClient.Create(ctx, apiDefinition)).Should(Succeed())

			apiDefinitionFixture = apiDefinition
			apiLookupKey = types.NamespacedName{Name: apiDefinitionFixture.Name, Namespace: namespace}

			By("Expecting the API Definition to be Ready")

			savedApiDefinition = new(v1alpha1.ApiDefinition)
			Eventually(func() error {
				err = k8sClient.Get(ctx, apiLookupKey, savedApiDefinition)
				return internal.AssertNoErrorAndStatusCompleted(err, savedApiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			gatewayEndpoint = internal.GatewayUrl + savedApiDefinition.Spec.Proxy.VirtualHosts[0].Path

			apim, err = internal.NewAPIM(ctx)
			Expect(err).ToNot(HaveOccurred())

			By("Should have no plan in the API definition")
			Expect(0, len(savedApiDefinition.Spec.Plans))
		})

		It("Should return NotFound Exception without plan", func() {
			Eventually(func() error {
				res, callErr := httpClient.Get(gatewayEndpoint)
				return internal.AssertNoErrorAndHTTPStatus(callErr, res, http.StatusNotFound)
			}, timeout, interval).ShouldNot(HaveOccurred())
		})

		It("Should be reachable using KeyLess plan", func() {
			By("Update ApiDefinition add KeyLess plan")
			updatedApiDefinition := savedApiDefinition.DeepCopy()
			keyLessPlanCrossId := uuid.FromStrings("key-less-plan-cross-id")
			updatedApiDefinition.Spec.Plans = append(savedApiDefinition.Spec.Plans,
				v2.NewPlan(
					base.NewPlan("G.K.O. KeyLess Plan", "").
						WithStatus(base.PublishedPlanStatus).
						WithCrossID(keyLessPlanCrossId),
				).WithSecurity("KEY_LESS"),
			)

			Eventually(func() error {
				update := new(v1alpha1.ApiDefinition)
				if err := k8sClient.Get(ctx, apiLookupKey, update); err != nil {
					return err
				}
				updatedApiDefinition.Spec.DeepCopyInto(&update.Spec)
				return k8sClient.Update(ctx, update)
			}, timeout, interval).ShouldNot(HaveOccurred())

			// Wait for the ApiDefinition to be updated
			Eventually(func() error {
				k8sErr := k8sClient.Get(ctx, apiLookupKey, updatedApiDefinition)
				return internal.AssertNoErrorAndStatusCompleted(k8sErr, updatedApiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			// Update savedApiDefinition & global var with last Get
			savedApiDefinition = updatedApiDefinition.DeepCopy()

			By("Call gateway with subscription api key of second plan")
			Eventually(func() error {
				res, callErr := httpClient.Get(gatewayEndpoint)
				return internal.AssertNoErrorAndHTTPStatus(callErr, res, http.StatusOK)
			}, timeout, interval).ShouldNot(HaveOccurred())
		})

		It("Should delete ApiDefinition resource", func() {
			By("Delete the ApiDefinition resource")
			err := k8sClient.Delete(ctx, apiDefinitionFixture)
			Expect(err).ToNot(HaveOccurred())

			By("Call deleted API definition URL and expect 404")
			Eventually(func() error {
				res, callErr := httpClient.Get(gatewayEndpoint)
				return internal.AssertNoErrorAndHTTPStatus(callErr, res, http.StatusNotFound)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Get the API definition from ManagementApi and expect deleted state")

			Eventually(func() error {
				_, apiErr := apim.APIs.GetByID(savedApiDefinition.Status.ID)
				return errors.IgnoreNotFound(apiErr)
			}, timeout, interval).ShouldNot(HaveOccurred())
		})
	})
})

func createSubscriptionAndGetApiKey(
	apim *internal.APIM,
	createdApiDefinition *v1alpha1.ApiDefinition,
	planSelector func(*api.Entity) string,
) string {
	// Get first active application
	mgmtApplications, mgmtErr := apim.Applications.Search("", "ACTIVE")
	Expect(mgmtErr).ToNot(HaveOccurred())
	defaultApplication := mgmtApplications[0]

	// Get Api description with plan
	mgmtApi, mgmtErr := apim.APIs.GetByID(createdApiDefinition.Status.ID)
	Expect(mgmtErr).ToNot(HaveOccurred())

	planId := planSelector(mgmtApi)

	// Create subscription
	mgmtSubscription, mgmtErr := apim.Subscriptions.Subscribe(mgmtApi.ID, defaultApplication.Id, planId)
	Expect(mgmtErr).ToNot(HaveOccurred())

	// Get subscription api keys
	mgmtSubscriptionApiKeys, mgmtErr := apim.Subscriptions.GetApiKeys(mgmtApi.ID, mgmtSubscription.Id)
	Expect(mgmtErr).ToNot(HaveOccurred())

	return mgmtSubscriptionApiKeys[0].Key
}

func getWithGioApiKey(httpClient *http.Client, gatewayEndpoint string, apiKey string) (*http.Response, error) {
	req, callErr := http.NewRequest("GET", gatewayEndpoint, nil)
	if callErr != nil {
		return nil, callErr
	}
	req.Header.Set("X-Gravitee-Api-Key", apiKey)

	res, callErr := httpClient.Do(req)
	if callErr != nil {
		return nil, callErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	return res, nil
}
