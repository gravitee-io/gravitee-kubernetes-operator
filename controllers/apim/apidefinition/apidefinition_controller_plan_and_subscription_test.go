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
package apidefinition

import (
	"context"
	"errors"
	"net/http"
	"time"

	clientError "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi/clienterror"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	managementapi "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi"
	managementapimodel "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test"
)

var _ = Describe("Checking ApiKey plan and subscription", Ordered, func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		namespace = "default"

		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	ctx := context.Background()
	httpClient := http.Client{Timeout: 5 * time.Second}

	Context("Checking ApiKey plan and subscription", Ordered, func() {
		var managementContextFixture *gio.ManagementContext
		var apiDefinitionFixture *gio.ApiDefinition

		var savedApiDefinition *gio.ApiDefinition

		var apiLookupKey types.NamespacedName
		var contextLookupKey types.NamespacedName

		var gatewayEndpoint string
		var mgmtClient *managementapi.Client

		BeforeAll(func() {
			By("Create a management context to synchronize with the REST API")
			managementContext, err := test.NewManagementContext(
				"../../../config/samples/context/dev/managementcontext_credentials.yaml")
			Expect(err).ToNot(HaveOccurred())
			Expect(k8sClient.Create(ctx, managementContext)).Should(Succeed())

			By("Create an API definition resource stared by default")
			apiDefinition, err := test.NewApiDefinition("../../../config/samples/apim/apikey-example-with-ctx.yml")
			Expect(err).ToNot(HaveOccurred())
			Expect(k8sClient.Create(ctx, apiDefinition)).Should(Succeed())

			apiDefinitionFixture = apiDefinition
			managementContextFixture = managementContext
			apiLookupKey = types.NamespacedName{Name: apiDefinitionFixture.Name, Namespace: namespace}
			contextLookupKey = types.NamespacedName{Name: managementContextFixture.Name, Namespace: namespace}

			By("Expect the API Definition is Ready")
			savedApiDefinition = new(gio.ApiDefinition)
			Eventually(func() bool {
				k8sErr := k8sClient.Get(ctx, apiLookupKey, savedApiDefinition)
				return k8sErr == nil && savedApiDefinition.Status.CrossID != ""
			}, timeout, interval).Should(BeTrue())

			gatewayEndpoint = test.GatewayUrl + savedApiDefinition.Spec.Proxy.VirtualHosts[0].Path
			mgmtClient = managementapi.NewClient(ctx, managementContextFixture, httpClient)

		})

		AfterAll(func() {
			err := k8sClient.Delete(ctx, apiDefinitionFixture)
			if err != nil {
				// wait deleted only if not already deleted
				Eventually(func() error {
					return k8sClient.Get(ctx, apiLookupKey, apiDefinitionFixture)
				}, timeout, interval).ShouldNot(Succeed())
			}

			err = k8sClient.Delete(ctx, managementContextFixture)
			if err != nil {
				// wait deleted only if not already deleted
				Eventually(func() error {
					return k8sClient.Get(ctx, contextLookupKey, managementContextFixture)
				}, timeout, interval).ShouldNot(Succeed())
			}
		})

		It("Should return unauthorize without subscription", func() {
			Eventually(func() bool {
				res, callErr := httpClient.Get(gatewayEndpoint)
				return callErr == nil && res.StatusCode == 401
			}, timeout, interval).Should(BeTrue())
		})

		It("Should return success with subscription", func() {
			By("Create a subscription and get api key")
			apiKey := createSubscriptionAndGetApiKey(
				mgmtClient,
				savedApiDefinition,
				func(mgmtApi *managementapimodel.ApiEntity) string { return mgmtApi.Plans[0].Id },
			)

			By("Call gateway with subscription api key")
			Eventually(func() bool {
				res, callErr := getWithGioApiKey(&httpClient, gatewayEndpoint, apiKey)
				return callErr == nil && res.StatusCode == 200
			}, timeout, interval).Should(BeTrue())
		})

		It("Should update ApiDefinition resource", func() {

			By("Update ApiDefinition path & name")

			updatedApiDefinition := savedApiDefinition.DeepCopy()

			expectedPath := updatedApiDefinition.Spec.Proxy.VirtualHosts[0].Path + "-updated"
			expectedName := updatedApiDefinition.Spec.Name + "-updated"
			updatedApiDefinition.Spec.Proxy.VirtualHosts[0].Path = expectedPath
			updatedApiDefinition.Spec.Name = expectedName

			err := k8sClient.Update(ctx, updatedApiDefinition)
			Expect(err).ToNot(HaveOccurred())

			// Wait for the ApiDefinition to be updated
			Eventually(func() bool {
				k8sErr := k8sClient.Get(ctx, apiLookupKey, updatedApiDefinition)
				return k8sErr == nil && updatedApiDefinition.Status.Generation == savedApiDefinition.Status.Generation+1
			}, timeout, interval).Should(BeTrue())

			// Update savedApiDefinition & global var with last Get
			savedApiDefinition = updatedApiDefinition.DeepCopy()
			gatewayEndpoint = test.GatewayUrl + savedApiDefinition.Spec.Proxy.VirtualHosts[0].Path

			By("Update ApiDefinition add ApiKey plan")

			secondPlanCrossId := utils.ToUUID("second-plan-cross-id")
			updatedApiDefinition.Spec.Plans = append(savedApiDefinition.Spec.Plans, &model.Plan{
				CrossId:  secondPlanCrossId,
				Name:     "G.K.O. Second ApiKey",
				Security: "API_KEY",
				Status:   "PUBLISHED",
			})

			err = k8sClient.Update(ctx, updatedApiDefinition)
			Expect(err).ToNot(HaveOccurred())

			// Wait for the ApiDefinition to be updated
			Eventually(func() bool {
				k8sErr := k8sClient.Get(ctx, apiLookupKey, updatedApiDefinition)
				return k8sErr == nil &&
					updatedApiDefinition.Status.Generation == savedApiDefinition.Status.Generation+1
			}, timeout, interval).Should(BeTrue())

			// Update savedApiDefinition & global var with last Get
			savedApiDefinition = updatedApiDefinition.DeepCopy()

			apiKey := createSubscriptionAndGetApiKey(
				mgmtClient,
				savedApiDefinition,
				func(mgmtApi *managementapimodel.ApiEntity) string {
					for _, plan := range mgmtApi.Plans {
						if plan.CrossId == secondPlanCrossId {
							return plan.Id
						}
					}
					return ""
				},
			)

			By("Call gateway with subscription api key of second plan")
			Eventually(func() bool {
				res, callErr := getWithGioApiKey(&httpClient, gatewayEndpoint, apiKey)
				return callErr == nil && res.StatusCode == 200
			}, timeout, interval).Should(BeTrue())
		})

		It("Should delete ApiDefinition resource", func() {
			By("Delete the ApiDefinition resource")
			err := k8sClient.Delete(ctx, apiDefinitionFixture)
			Expect(err).ToNot(HaveOccurred())

			By("Call deleted API definition URL and expect 404")
			Eventually(func() bool {
				res, callErr := httpClient.Get(gatewayEndpoint)
				return callErr == nil && res.StatusCode == 404
			}, timeout, interval).Should(BeTrue())

			By("Get the API definition from ManagementApi and expect deleted state")

			Eventually(func() bool {
				_, apiErr := mgmtClient.GetApiById(savedApiDefinition.Status.ID)
				var apiNotFoundError *clientError.ApiNotFoundError
				return apiErr != nil && errors.As(apiErr, &apiNotFoundError)
			}, timeout, interval).Should(BeTrue())
		})
	})
})

func createSubscriptionAndGetApiKey(
	mgmtClient *managementapi.Client,
	createdApiDefinition *gio.ApiDefinition,
	planSelector func(*managementapimodel.ApiEntity) string,
) string {
	// Get first active application
	mgmtApplications, mgmtErr := mgmtClient.SearchApplications("", "ACTIVE")
	Expect(mgmtErr).ToNot(HaveOccurred())
	defaultApplication := mgmtApplications[0]

	// Get Api description with plan
	mgmtApi, mgmtErr := mgmtClient.GetApiById(createdApiDefinition.Status.ID)
	Expect(mgmtErr).ToNot(HaveOccurred())

	planId := planSelector(mgmtApi)

	// Create subscription
	mgmtSubscription, mgmtErr := mgmtClient.SubscribeToPlan(mgmtApi.Id, defaultApplication.Id, planId)
	Expect(mgmtErr).ToNot(HaveOccurred())

	// Get subscription api keys
	mgmtSubscriptionApiKeys, mgmtErr := mgmtClient.GetSubscriptionApiKey(mgmtApi.Id, mgmtSubscription.Id)
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
