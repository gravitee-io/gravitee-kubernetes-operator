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
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal"
)

var _ = Describe("Checking NoneRecoverable && Recoverable error", Label("DisableSmokeExpect"), func() {

	Context("With basic ApiDefinition & ManagementContext", func() {
		var managementContextFixture *gio.ApiContext
		var apiDefinitionFixture *gio.ApiDefinition

		var savedApiDefinition *gio.ApiDefinition

		var apiLookupKey types.NamespacedName
		var contextLookupKey types.NamespacedName

		BeforeEach(func() {
			By("Create a management context to synchronize with the REST API")

			fixtureGenerator := internal.NewFixtureGenerator()

			apiWithContext, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Api:     internal.BasicApiFile,
				Context: internal.ContextWithSecretFile,
			})

			Expect(err).ToNot(HaveOccurred())

			managementContext := apiWithContext.Context
			Expect(k8sClient.Create(ctx, managementContext)).Should(Succeed())
			contextLookupKey = types.NamespacedName{Name: managementContext.Name, Namespace: namespace}

			Eventually(func() error {
				return k8sClient.Get(ctx, contextLookupKey, managementContext)
			}, timeout, interval).Should(Succeed())

			By("Create an API definition resource stared by default")

			apiDefinition := apiWithContext.Api
			Expect(k8sClient.Create(ctx, apiDefinition)).Should(Succeed())

			apiDefinitionFixture = apiDefinition
			managementContextFixture = managementContext
			apiLookupKey = types.NamespacedName{Name: apiDefinitionFixture.Name, Namespace: namespace}

			By("Expect the API Definition is Ready")
			savedApiDefinition = new(gio.ApiDefinition)
			Eventually(func() error {
				if err = k8sClient.Get(ctx, apiLookupKey, savedApiDefinition); err != nil {
					return err
				}
				return internal.AssertStatusContextIsSet(savedApiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())
		})

		It("Should not requeue reconcile with 401 error", func() {

			By("Set bad credentials in ManagementContext")
			managementContextBad := managementContextFixture.DeepCopy()
			managementContextBad.Spec.Management.Auth.SecretRef = nil
			managementContextBad.Spec.Management.Auth.BearerToken = "bad-token"

			err := k8sClient.Update(ctx, managementContextBad)
			Expect(err).ToNot(HaveOccurred())

			By("Update the API definition")
			apiDefinition := savedApiDefinition.DeepCopy()

			Eventually(func() error {
				return k8sClient.Get(ctx, apiLookupKey, apiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			apiDefinition.Spec.Name = "new-name"

			err = k8sClient.Update(ctx, apiDefinition)
			Expect(err).ToNot(HaveOccurred())

			By("Check API definition processing status")
			Eventually(func() error {
				if err = k8sClient.Get(ctx, apiLookupKey, savedApiDefinition); err != nil {
					return err
				}
				context := internal.GetStatusContext(savedApiDefinition, contextLookupKey)
				if context == nil {
					return fmt.Errorf("context not found")
				}
				if context.Status != gio.ProcessingStatusFailed {
					return fmt.Errorf("expected status %s, got %s", gio.ProcessingStatusFailed, context.Status)
				}
				return nil
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Check events")
			Expect(getEventsReason(apiDefinitionFixture)).Should(ContainElements([]string{"UpdateStarted", "UpdateFailed"}))

			By("Set right credentials in ManagementContext")
			managementContextRight := managementContextBad.DeepCopy()
			managementContextRight.Spec = managementContextFixture.Spec

			err = k8sClient.Update(ctx, managementContextRight)
			Expect(err).ToNot(HaveOccurred())

			By("Check that API definition has been reconciled on ManagementContext update")

			apim, err := internal.NewAPIM(ctx)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() error {
				api, cliErr := apim.APIs.GetByID(internal.GetStatusId(savedApiDefinition, contextLookupKey))
				if cliErr != nil {
					return cliErr
				}

				if api.Name != "new-name" {
					return fmt.Errorf("expected name %s, got %s", "new-name", api.Name)
				}

				context := internal.GetStatusContext(savedApiDefinition, contextLookupKey)

				if context.ID != api.ID {
					return fmt.Errorf("expected id %s, got %s", api.ID, context.ID)
				}

				return nil
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Check events")
			Expect(getEventsReason(apiDefinitionFixture)).Should(ContainElements([]string{"UpdateSucceeded"}))
		})

		It("Should requeue reconcile with bad ManagementContext BaseUrl", func() {

			By("Set bad BaseUrl in ManagementContext")
			managementContextBad := managementContextFixture.DeepCopy()
			managementContextBad.Spec.Management.BaseUrl = "http://bad-url:8083"

			err := k8sClient.Update(ctx, managementContextBad)
			Expect(err).ToNot(HaveOccurred())

			By("Update the API definition")
			apiDefinition := savedApiDefinition.DeepCopy()
			apiDefinition.Spec.Name = "new-name"

			err = k8sClient.Update(ctx, apiDefinition)
			Expect(err).ToNot(HaveOccurred())

			By("Check API definition processing status")
			Eventually(func() error {
				if err = k8sClient.Get(ctx, apiLookupKey, savedApiDefinition); err != nil {
					return err
				}

				return internal.AssertStatusContextMatches(savedApiDefinition, contextLookupKey, &gio.StatusContext{
					Status:  gio.ProcessingStatusFailed,
					EnvID:   "DEFAULT",
					OrgID:   "DEFAULT",
					CrossID: apiDefinition.GetOrGenerateCrossID(),
					ID:      apiDefinition.GetID(),
					State:   "STARTED",
				})
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Set right BaseUrl in ManagementContext")
			managementContextRight := managementContextBad.DeepCopy()
			managementContextRight.Spec = managementContextFixture.Spec

			err = k8sClient.Update(ctx, managementContextRight)
			Expect(err).ToNot(HaveOccurred())

			By("Check API definition processing status")
			Eventually(func() error {
				if err = k8sClient.Get(ctx, apiLookupKey, savedApiDefinition); err != nil {
					return err
				}

				return internal.AssertStatusContextMatches(savedApiDefinition, contextLookupKey, &gio.StatusContext{
					Status:  gio.ProcessingStatusCompleted,
					EnvID:   "DEFAULT",
					OrgID:   "DEFAULT",
					ID:      apiDefinition.GetID(),
					CrossID: apiDefinition.GetOrGenerateCrossID(),
					State:   "STARTED",
				})
			}, timeout, interval).ShouldNot(HaveOccurred())
		})
	})
})
