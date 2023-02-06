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
	"fmt"
	"net/http"
	"strings"
	"time"

	model "github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("An API template", func() {
	var fixturesGenerator *internal.FixtureGenerator

	BeforeEach(func() {
		fixturesGenerator = internal.NewFixtureGenerator()
	})

	Context("with one context holding values", func() {

		When("not synced with an APIM instance", func() {

			It("can create an API, using context values as a template context", func() {
				fixtures, err := fixturesGenerator.NewFixtures(internal.FixtureFiles{
					Api: internal.ApiWithContextFile,
					Contexts: []string{
						internal.ContextWithSecretFile,
					},
				}, func(fixtures *internal.Fixtures) {
					fixtures.Contexts[0].Spec.Management = nil
				})

				Expect(err).ToNot(HaveOccurred())

				By("creating an API context with values and no APIM instance defined")

				apiContext := &fixtures.Contexts[0]
				Expect(k8sClient.Create(ctx, apiContext)).Should(Succeed())

				Eventually(func() error {
					return k8sClient.Get(ctx, apiContext.GetNamespacedName().ToK8sType(), apiContext)
				}, timeout, interval).Should(Succeed())

				By("creating an API using this context")

				apiDefinition := fixtures.Api
				Expect(k8sClient.Create(ctx, apiDefinition)).Should(Succeed())

				Eventually(func() error {
					return k8sClient.Get(ctx, apiDefinition.GetNamespacedName().ToK8sType(), apiDefinition)
				}, timeout, interval).Should(Succeed())

				By("calling successfully the API using the context path templated with the context values")

				httpClient := http.Client{Timeout: 5 * time.Second}
				endpoint := internal.GatewayUrl + fixturesGenerator.AddSuffix("/context-dev")
				Eventually(func() error {
					res, callErr := httpClient.Get(endpoint)
					return internal.AssertNoErrorAndHTTPStatus(callErr, res, http.StatusOK)
				}, timeout, interval).ShouldNot(HaveOccurred())
			})
		})

		When("synced with an APIM instance", func() {

			It("can create an API, using contexts value as a template context", func() {
				fixtures, err := fixturesGenerator.NewFixtures(internal.FixtureFiles{
					Api: internal.ApiWithContextFile,
					Contexts: []string{
						internal.ContextWithSecretFile,
					},
				})
				Expect(err).ToNot(HaveOccurred())

				By("creating an API context with values and an APIM instance defined")

				apiContext := &fixtures.Contexts[0]
				Expect(k8sClient.Create(ctx, apiContext)).Should(Succeed())

				Eventually(func() error {
					return k8sClient.Get(ctx, apiContext.GetNamespacedName().ToK8sType(), apiContext)
				}, timeout, interval).Should(Succeed())

				By("creating an API using this context")

				apiDefinition := fixtures.Api
				Expect(k8sClient.Create(ctx, apiDefinition)).Should(Succeed())

				Eventually(func() error {
					if getErr := k8sClient.Get(ctx, apiDefinition.GetNamespacedName().ToK8sType(), apiDefinition); getErr != nil {
						return getErr
					}
					return internal.AssertStatusContextIsSet(apiDefinition)
				}, timeout, interval).Should(Succeed())

				By("getting the API from APIM and see that its name has been templated with the context values")

				Eventually(func() error {
					lookupKey := apiContext.GetNamespacedName().ToK8sType()
					id := internal.GetStatusId(apiDefinition, lookupKey)
					Expect(id).ToNot(BeEmpty())

					apim, apimErr := internal.NewAPIM(ctx)
					Expect(apimErr).ToNot(HaveOccurred())

					api, getErr := apim.APIs.GetByID(id)
					if getErr != nil {
						return getErr
					}

					expectedName := fixturesGenerator.AddSuffix("Context (dev)")
					return internal.AssertEquals("API Name", expectedName, api.Name)
				}, timeout, interval).ShouldNot(HaveOccurred())

				By("calling successfully the API using the context path templated with the context values")

				httpClient := http.Client{Timeout: 5 * time.Second}
				endpoint := internal.GatewayUrl + fixturesGenerator.AddSuffix("/context-dev")
				Eventually(func() error {
					res, callErr := httpClient.Get(endpoint)
					return internal.AssertNoErrorAndHTTPStatus(callErr, res, http.StatusOK)
				}, timeout, interval).ShouldNot(HaveOccurred())
			})
		})
	})

	Context("with two contexts holding values", func() {

		When("not synced with an APIM instance", func() {

			It("can create two APIs, using contexts values as a template context", func() {
				fixtures, err := fixturesGenerator.NewFixtures(internal.FixtureFiles{
					Api: internal.ApiWithContextFile,
					Contexts: []string{
						internal.ContextWithSecretFile,
					},
				}, func(fixtures *internal.Fixtures) {

					fixtures.Contexts[0].Spec.Management = nil

					devContext := &fixtures.Contexts[0]

					stagingContext := devContext.DeepCopy()
					stagingContext.Name = strings.Replace(devContext.Name, "dev", "staging", 1)
					stagingContext.Spec.Values = map[string]string{
						"env": "staging",
					}

					fixtures.Contexts = []gio.ManagementContext{
						*stagingContext,
						*devContext,
					}

					fixtures.Api.Spec.Contexts = []model.NamespacedName{
						stagingContext.GetNamespacedName(),
						devContext.GetNamespacedName(),
					}
				})

				Expect(err).ToNot(HaveOccurred())

				By("creating two API contexts with values and no APIM instance defined")

				contexts := fixtures.Contexts
				for i := range fixtures.Contexts {
					context := &contexts[i]
					Expect(k8sClient.Create(ctx, context)).Should(Succeed())
					Eventually(func() error {
						return k8sClient.Get(ctx, context.GetNamespacedName().ToK8sType(), new(gio.ManagementContext))
					}, timeout, interval).Should(Succeed())
				}

				By("creating an API using this contexts")

				apiDefinition := fixtures.Api
				Expect(k8sClient.Create(ctx, apiDefinition)).Should(Succeed())

				Eventually(func() error {
					getErr := k8sClient.Get(ctx, apiDefinition.GetNamespacedName().ToK8sType(), apiDefinition)
					if getErr != nil {
						return getErr
					}
					return internal.AssertStatusContextsLen(apiDefinition, 2)
				}, timeout, interval).Should(Succeed())

				By("calling successfully the API using the contexts paths templated with the contexts values")

				httpClient := http.Client{Timeout: 5 * time.Second}

				stagingEndpoint := internal.GatewayUrl + fixturesGenerator.AddSuffix("/context-staging")
				Eventually(func() error {
					res, callErr := httpClient.Get(stagingEndpoint)
					return internal.AssertNoErrorAndHTTPStatus(callErr, res, http.StatusOK)
				}, timeout, interval).ShouldNot(HaveOccurred())

				devEndpoint := internal.GatewayUrl + fixturesGenerator.AddSuffix("/context-dev")
				Eventually(func() error {
					res, callErr := httpClient.Get(devEndpoint)
					return internal.AssertNoErrorAndHTTPStatus(callErr, res, http.StatusOK)
				}, timeout, interval).ShouldNot(HaveOccurred())
			})
		})

		When("both contexts synced with the same APIM instance", func() {
			fmt.Fprintf(GinkgoWriter, `
				When multiple contexts are set to point on the same instance,
				each context must be synced with one dedicated environment to
				allow multiple API creations in APIM.

				This is because we compute the cross ID in a reproducible manner
				using the API namespace and name, which would result in updating
				the same API for each context instead of creating a new one if we 
				are targeting the same env.

				As multi-env is only supported through cockpit, we do not cover this
				in tests yet.
			`)
		})

	})
})
