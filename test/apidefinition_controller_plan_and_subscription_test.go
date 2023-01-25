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

/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	apimModel "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
)

var _ = Describe("Checking ApiKey plan and subscription", Ordered, func() {

	httpClient := http.Client{Timeout: 5 * time.Second}

	Context("Checking ApiKey plan and subscription", Ordered, func() {
		var apiDefinitionFixture *gio.ApiDefinition

		var savedApiDefinition *gio.ApiDefinition

		var apiLookupKey types.NamespacedName
		var contextLookupKey types.NamespacedName

		var gatewayEndpoint string
		var apim *internal.APIM

		BeforeAll(func() {
			By("Create a management context to synchronize with the REST API")

			fixtureGenerator := internal.NewFixtureGenerator()

			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Api:      internal.ApiWithApiKeyPlanFile,
				Contexts: []string{internal.ContextWithSecretFile},
			})

			Expect(err).ToNot(HaveOccurred())

			apiContext := &fixtures.Contexts[0]
			Expect(k8sClient.Create(ctx, apiContext)).Should(Succeed())

			contextLookupKey = types.NamespacedName{Name: apiContext.Name, Namespace: namespace}
			Eventually(func() error {
				return k8sClient.Get(ctx, contextLookupKey, apiContext)
			}, timeout, interval).Should(Succeed())

			By("Create an API definition resource stared by default")

			apiDefinition := fixtures.Api
			Expect(k8sClient.Create(ctx, apiDefinition)).Should(Succeed())

			apiDefinitionFixture = apiDefinition
			apiLookupKey = types.NamespacedName{Name: apiDefinitionFixture.Name, Namespace: namespace}

			By("Expect the API Definition is Ready")

			savedApiDefinition = new(gio.ApiDefinition)
			Eventually(func() error {
				if err = k8sClient.Get(ctx, apiLookupKey, savedApiDefinition); err != nil {
					return err
				}
				return internal.AssertStatusContextIsSet(savedApiDefinition)
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
				contextLookupKey,
				func(mgmtApi *apimModel.ApiEntity) string { return mgmtApi.Plans[0].Id },
			)

			By("Call gateway with subscription api key")
			Eventually(func() error {
				res, callErr := getWithGioApiKey(&httpClient, gatewayEndpoint, apiKey)
				return internal.AssertNoErrorAndHTTPStatus(callErr, res, http.StatusOK)
			}, timeout, interval).ShouldNot(HaveOccurred())

		})

		It("Should update ApiDefinition resource", func() {

			By("Update ApiDefinition path & name")

			updatedApiDefinition := savedApiDefinition.DeepCopy()

			expectedPath := updatedApiDefinition.Spec.Proxy.VirtualHosts[0].Path + "-updated"
			expectedName := updatedApiDefinition.Spec.Name + "-updated"
			updatedApiDefinition.Spec.Proxy.VirtualHosts[0].Path = expectedPath
			updatedApiDefinition.Spec.Name = expectedName

			Eventually(func() error {
				update := new(gio.ApiDefinition)
				if err := k8sClient.Get(ctx, apiLookupKey, update); err != nil {
					return err
				}
				updatedApiDefinition.Spec.DeepCopyInto(&update.Spec)
				return k8sClient.Update(ctx, update)
			}, timeout, interval).ShouldNot(HaveOccurred())

			// Wait for the ApiDefinition to be updated
			Eventually(func() error {
				expectedGeneration := savedApiDefinition.Status.ObservedGeneration + 1
				k8sErr := k8sClient.Get(ctx, apiLookupKey, updatedApiDefinition)
				return internal.AssertNoErrorAndObservedGenerationEquals(
					k8sErr, updatedApiDefinition, expectedGeneration,
				)
			}, timeout, interval).ShouldNot(HaveOccurred())

			// Update savedApiDefinition & global var with last Get
			savedApiDefinition = updatedApiDefinition.DeepCopy()
			gatewayEndpoint = internal.GatewayUrl + savedApiDefinition.Spec.Proxy.VirtualHosts[0].Path

			By("Update ApiDefinition add ApiKey plan")

			secondPlanCrossId := uuid.FromStrings("second-plan-cross-id")
			updatedApiDefinition.Spec.Plans = append(savedApiDefinition.Spec.Plans, &model.Plan{
				CrossId:  secondPlanCrossId,
				Name:     "G.K.O. Second ApiKey",
				Security: "API_KEY",
				Status:   "PUBLISHED",
			})

			Eventually(func() error {
				update := new(gio.ApiDefinition)
				if err := k8sClient.Get(ctx, apiLookupKey, update); err != nil {
					return err
				}
				updatedApiDefinition.Spec.DeepCopyInto(&update.Spec)
				return k8sClient.Update(ctx, update)
			}, timeout, interval).ShouldNot(HaveOccurred())

			// Wait for the ApiDefinition to be updated
			Eventually(func() error {
				expectedGeneration := savedApiDefinition.Status.ObservedGeneration + 1
				k8sErr := k8sClient.Get(ctx, apiLookupKey, updatedApiDefinition)
				return internal.AssertNoErrorAndObservedGenerationEquals(
					k8sErr, updatedApiDefinition, expectedGeneration,
				)
			}, timeout, interval).ShouldNot(HaveOccurred())

			// Update savedApiDefinition & global var with last Get
			savedApiDefinition = updatedApiDefinition.DeepCopy()

			apiKey := createSubscriptionAndGetApiKey(
				apim,
				savedApiDefinition,
				contextLookupKey,
				func(mgmtApi *apimModel.ApiEntity) string {
					for _, plan := range mgmtApi.Plans {
						if plan.CrossId == secondPlanCrossId {
							return plan.Id
						}
					}
					return ""
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
				_, apiErr := apim.APIs.GetByID(internal.GetStatusId(savedApiDefinition, contextLookupKey))
				return errors.IgnoreNotFound(apiErr)
			}, timeout, interval).ShouldNot(HaveOccurred())
		})
	})

	Context("Checking Api with no plan", Ordered, func() {
		var apiDefinitionFixture *gio.ApiDefinition
		var savedApiDefinition *gio.ApiDefinition
		var apiLookupKey types.NamespacedName
		var contextLookupKey types.NamespacedName
		var gatewayEndpoint string
		var apim *internal.APIM

		BeforeAll(func() {
			By("Create a management context to synchronize with the REST API")

			fixtureGenerator := internal.NewFixtureGenerator()

			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Api:      internal.ApiWithContextNoPlanFile,
				Contexts: []string{internal.ContextWithSecretFile},
			})

			Expect(err).ToNot(HaveOccurred())

			apiContext := &fixtures.Contexts[0]
			Expect(k8sClient.Create(ctx, apiContext)).Should(Succeed())

			contextLookupKey = types.NamespacedName{Name: apiContext.Name, Namespace: namespace}
			Eventually(func() error {
				return k8sClient.Get(ctx, contextLookupKey, apiContext)
			}, timeout, interval).Should(Succeed())

			By("Create an API definition resource stared by default")

			apiDefinition := fixtures.Api
			Expect(k8sClient.Create(ctx, apiDefinition)).Should(Succeed())

			apiDefinitionFixture = apiDefinition
			apiLookupKey = types.NamespacedName{Name: apiDefinitionFixture.Name, Namespace: namespace}

			By("Expect the API Definition is Ready")

			savedApiDefinition = new(gio.ApiDefinition)
			Eventually(func() error {
				if err = k8sClient.Get(ctx, apiLookupKey, savedApiDefinition); err != nil {
					return err
				}
				return internal.AssertStatusContextIsSet(savedApiDefinition)
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
			updatedApiDefinition.Spec.Plans = append(savedApiDefinition.Spec.Plans, &model.Plan{
				CrossId:  keyLessPlanCrossId,
				Name:     "G.K.O. KeyLess Plan",
				Security: "KEY_LESS",
				Status:   "PUBLISHED",
			})

			Eventually(func() error {
				update := new(gio.ApiDefinition)
				if err := k8sClient.Get(ctx, apiLookupKey, update); err != nil {
					return err
				}
				updatedApiDefinition.Spec.DeepCopyInto(&update.Spec)
				return k8sClient.Update(ctx, update)
			}, timeout, interval).ShouldNot(HaveOccurred())

			// Wait for the ApiDefinition to be updated
			Eventually(func() error {
				expectedGeneration := savedApiDefinition.Status.ObservedGeneration + 1
				k8sErr := k8sClient.Get(ctx, apiLookupKey, updatedApiDefinition)
				return internal.AssertNoErrorAndObservedGenerationEquals(
					k8sErr, updatedApiDefinition, expectedGeneration,
				)
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
				_, apiErr := apim.APIs.GetByID(internal.GetStatusId(savedApiDefinition, contextLookupKey))
				return errors.IgnoreNotFound(apiErr)
			}, timeout, interval).ShouldNot(HaveOccurred())
		})
	})
})

func createSubscriptionAndGetApiKey(
	apim *internal.APIM,
	createdApiDefinition *gio.ApiDefinition,
	contextLocation types.NamespacedName,
	planSelector func(*apimModel.ApiEntity) string,
) string {
	// Get first active application
	mgmtApplications, mgmtErr := apim.Applications.Search("", "ACTIVE")
	Expect(mgmtErr).ToNot(HaveOccurred())
	defaultApplication := mgmtApplications[0]

	// Get Api description with plan
	mgmtApi, mgmtErr := apim.APIs.GetByID(internal.GetStatusId(createdApiDefinition, contextLocation))
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
