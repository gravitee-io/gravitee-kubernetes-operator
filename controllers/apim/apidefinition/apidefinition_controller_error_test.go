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
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	k8sErr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test"
)

var _ = Describe("Checking NoneRecoverable && Recoverable error", Label("DisableSmokeExpect"), func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		namespace = "default"

		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	ctx := context.Background()

	Context("With basic ApiDefinition & ManagementContext", func() {
		var managementContextFixture *gio.ManagementContext
		var apiDefinitionFixture *gio.ApiDefinition

		var savedApiDefinition *gio.ApiDefinition

		var apiLookupKey types.NamespacedName
		var contextLookupKey types.NamespacedName

		BeforeEach(func() {
			By("Create a management context to synchronize with the REST API")
			managementContext, err := test.NewManagementContext(
				"../../../config/samples/context/dev/managementcontext_credentials.yaml")
			Expect(err).ToNot(HaveOccurred())
			Expect(k8sClient.Create(ctx, managementContext)).Should(Succeed())

			By("Create an API definition resource stared by default")
			apiDefinition, err := test.NewApiDefinition("../../../config/samples/apim/basic-example-with-ctx.yml")
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
		})

		AfterEach(func() {
			err := k8sClient.Delete(ctx, apiDefinitionFixture)
			if !k8sErr.IsNotFound(err) {
				// wait deleted only if not already deleted
				Eventually(func() error {
					return k8sClient.Get(ctx, apiLookupKey, apiDefinitionFixture)
				}, timeout, interval).ShouldNot(Succeed())
			}

			err = k8sClient.Delete(ctx, managementContextFixture)
			if !k8sErr.IsNotFound(err) {
				// wait deleted only if not already deleted
				Eventually(func() error {
					return k8sClient.Get(ctx, contextLookupKey, managementContextFixture)
				}, timeout, interval).ShouldNot(Succeed())
			}
		})

		It("Should not requeue reconcile with 401 error", func() {

			By("Set bad credentials in ManagementContext")
			managementContextBad := managementContextFixture.DeepCopy()
			managementContextBad.Spec.Auth.Credentials.Username = "bad-username"

			err := k8sClient.Update(ctx, managementContextBad)
			Expect(err).ToNot(HaveOccurred())

			By("Update the API definition")
			apiDefinition := savedApiDefinition.DeepCopy()
			apiDefinition.Spec.Name = "new-name"

			err = k8sClient.Update(ctx, apiDefinition)
			Expect(err).ToNot(HaveOccurred())

			By("Check API definition processing status")
			Eventually(func() bool {
				k8sErr := k8sClient.Get(ctx, apiLookupKey, savedApiDefinition)
				return k8sErr == nil && savedApiDefinition.Status.ProcessingStatus == gio.ProcessingStatusFailed
			}, timeout, interval).Should(BeTrue())

			By("Set right credentials in ManagementContext")
			managementContextRight := managementContextBad.DeepCopy()
			managementContextRight.Spec = managementContextFixture.Spec

			err = k8sClient.Update(ctx, managementContextRight)
			Expect(err).ToNot(HaveOccurred())

		})

		It("Should requeue reconcile with bad ManagementContext BaseUrl", func() {

			By("Set bad BaseUrl in ManagementContext")
			managementContextBad := managementContextFixture.DeepCopy()
			managementContextBad.Spec.BaseUrl = "http://bad-url:8083"

			err := k8sClient.Update(ctx, managementContextBad)
			Expect(err).ToNot(HaveOccurred())

			By("Update the API definition")
			apiDefinition := savedApiDefinition.DeepCopy()
			apiDefinition.Spec.Name = "new-name"

			err = k8sClient.Update(ctx, apiDefinition)
			Expect(err).ToNot(HaveOccurred())

			By("Check API definition processing status")
			Eventually(func() bool {
				k8sErr := k8sClient.Get(ctx, apiLookupKey, savedApiDefinition)
				return k8sErr == nil && savedApiDefinition.Status.ProcessingStatus == gio.ProcessingStatusReconciling
			}, timeout, interval).Should(BeTrue())

			By("Set right BaseUrl in ManagementContext")
			managementContextRight := managementContextBad.DeepCopy()
			managementContextRight.Spec = managementContextFixture.Spec

			err = k8sClient.Update(ctx, managementContextRight)
			Expect(err).ToNot(HaveOccurred())

			// TODO: fix it
			// By("Check API definition processing status")
			// Eventually(func() bool {
			// 	k8sErr := k8sClient.Get(ctx, apiLookupKey, savedApiDefinition)
			// 	fmt.Println(savedApiDefinition.Status.ProcessingStatus)
			// 	return k8sErr == nil && savedApiDefinition.Status.ProcessingStatus == gio.ProcessingStatusCompleted
			// }, timeout, interval).Should(BeTrue())
		})
	})
})
