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
		var apiContextFixture *gio.ApiContext
		var apiDefinitionFixture *gio.ApiDefinition

		var savedApiDefinition *gio.ApiDefinition

		var apiLookupKey types.NamespacedName
		var contextLookupKey types.NamespacedName

		BeforeEach(func() {
			By("Create a management context to synchronize with the REST API")

			fixtureGenerator := internal.NewFixtureGenerator()

			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Api:      internal.BasicApiFile,
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
			apiContextFixture = apiContext
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

			apiContextBad := apiContextFixture.DeepCopy()
			apiContextBad.Spec.Management.Auth.SecretRef = nil
			apiContextBad.Spec.Management.Auth.BearerToken = "bad-token"

			Eventually(func() error {
				update := new(gio.ApiContext)
				if err := k8sClient.Get(ctx, contextLookupKey, update); err != nil {
					return err
				}
				apiContextBad.Spec.DeepCopyInto(&update.Spec)
				return k8sClient.Update(ctx, update)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Update the API definition")

			apiDefinition := savedApiDefinition.DeepCopy()

			Eventually(func() error {
				return k8sClient.Get(ctx, apiLookupKey, apiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			apiDefinition.Spec.Name = "new-name"

			Eventually(func() error {
				update := new(gio.ApiDefinition)
				if err := k8sClient.Get(ctx, apiLookupKey, update); err != nil {
					return err
				}
				apiDefinition.Spec.DeepCopyInto(&update.Spec)
				return k8sClient.Update(ctx, update)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Check API definition processing status")

			Eventually(func() error {
				if err := k8sClient.Get(ctx, apiLookupKey, savedApiDefinition); err != nil {
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

			apiContextRight := apiContextBad.DeepCopy()
			apiContextRight.Spec = apiContextFixture.Spec

			Eventually(func() error {
				update := new(gio.ApiContext)
				if err := k8sClient.Get(ctx, contextLookupKey, update); err != nil {
					return err
				}
				apiContextRight.Spec.DeepCopyInto(&update.Spec)
				return k8sClient.Update(ctx, update)
			}, timeout, interval).ShouldNot(HaveOccurred())

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

			apiContextBad := apiContextFixture.DeepCopy()
			apiContextBad.Spec.Management.BaseUrl = "http://bad-url:8083"

			Eventually(func() error {
				update := new(gio.ApiContext)
				if err := k8sClient.Get(ctx, contextLookupKey, update); err != nil {
					return err
				}
				apiContextBad.Spec.DeepCopyInto(&update.Spec)
				return k8sClient.Update(ctx, update)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Update the API definition")

			apiDefinition := savedApiDefinition.DeepCopy()
			apiDefinition.Spec.Name = "new-name"

			Eventually(func() error {
				update := new(gio.ApiDefinition)
				if err := k8sClient.Get(ctx, apiLookupKey, update); err != nil {
					return err
				}
				apiDefinition.Spec.DeepCopyInto(&update.Spec)
				return k8sClient.Update(ctx, update)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Check API definition processing status")

			Eventually(func() error {
				if err := k8sClient.Get(ctx, apiLookupKey, savedApiDefinition); err != nil {
					return err
				}

				return internal.AssertStatusContextMatches(savedApiDefinition, contextLookupKey, &gio.StatusContext{
					Status:  gio.ProcessingStatusFailed,
					EnvID:   "DEFAULT",
					OrgID:   "DEFAULT",
					CrossID: apiDefinition.GetOrGenerateCrossID(),
					ID:      apiDefinition.PickID(contextLookupKey.String()),
					State:   "STARTED",
				})
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Set right BaseUrl in ManagementContext")

			apiContextRight := apiContextBad.DeepCopy()
			apiContextRight.Spec = apiContextFixture.Spec

			Eventually(func() error {
				update := new(gio.ApiContext)
				if err := k8sClient.Get(ctx, contextLookupKey, update); err != nil {
					return err
				}
				apiContextRight.Spec.DeepCopyInto(&update.Spec)
				return k8sClient.Update(ctx, update)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Check API definition processing status")

			Eventually(func() error {
				if err := k8sClient.Get(ctx, apiLookupKey, savedApiDefinition); err != nil {
					return err
				}

				return internal.AssertStatusContextMatches(savedApiDefinition, contextLookupKey, &gio.StatusContext{
					Status:  gio.ProcessingStatusCompleted,
					EnvID:   "DEFAULT",
					OrgID:   "DEFAULT",
					ID:      apiDefinition.PickID(contextLookupKey.String()),
					CrossID: apiDefinition.GetOrGenerateCrossID(),
					State:   "STARTED",
				})
			}, timeout, interval).ShouldNot(HaveOccurred())
		})
	})
})
