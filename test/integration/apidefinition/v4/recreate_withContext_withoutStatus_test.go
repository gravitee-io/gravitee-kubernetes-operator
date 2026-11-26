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

	v4model "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	adm "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/apidefinition"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Create", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	It("should re-attach the API GKO 4.11 created when the CR is re-created without status", func() {
		fixtures := fixture.Builder().
			WithApplication(constants.Application).
			WithAPIv4(constants.ApiV4WithApiKeyPlanFile).
			WithContext(constants.ContextWithSecretFile).
			Build()
		api := fixtures.APIv4
		fixtures.APIv4 = nil
		fixtures = fixtures.Apply()

		By("defaulting the CR the way the API server does for the webhook and the controller")

		defaulted := api.DeepCopy()
		Expect(manager.Client().Create(ctx, defaulted, crclient.DryRunAll)).To(Succeed())
		api.Spec = defaulted.Spec

		client := apim.NewClient(ctx)

		By("creating the API in APIM the way GKO 4.11 did")

		legacy := api.DeepCopy()
		Expect(apidefinition.ResolveV4Refs(ctx, legacy, false)).To(Succeed())
		legacy.Spec.DefinitionContext = v4model.NewDefaultKubernetesContext().MergeWith(legacy.Spec.DefinitionContext)
		apiID := legacy.LegacyID(client.Context)
		legacy.Spec.ID = apiID
		legacy.PopulateIDs(client.Context, false)
		planID := (*legacy.Spec.Plans)[planKey].ID
		_, err := client.APIs.ImportV4(legacy)
		Expect(err).ToNot(HaveOccurred())

		By("subscribing an application to its plan")

		subscription, err := client.Subscriptions.Subscribe(apiID, fixtures.Application.Status.ID, planID)
		Expect(err).ToNot(HaveOccurred())

		Eventually(func() (string, error) {
			return subscriptionStatusInAPIM(client, apiID, subscription.ID)
		}, timeout, interval).Should(Equal("ACCEPTED"))

		By("admitting the CR re-created without status")

		_, err = adm.AdmissionCtrl{}.ValidateCreate(ctx, api)
		Expect(err).ToNot(HaveOccurred())

		By("re-creating the CR without status")

		Expect(manager.Client().Create(ctx, api)).To(Succeed())

		Eventually(func() error {
			if getErr := manager.GetLatest(ctx, api); getErr != nil {
				return getErr
			}
			if completedErr := assert.ApiV4Completed(api); completedErr != nil {
				return completedErr
			}
			return assert.ApiV4Accepted(api)
		}, timeout, interval).Should(Succeed())

		By("expecting the CR to point at the API, plan and subscription GKO 4.11 created")

		Expect(api.Status.ID).To(Equal(apiID))
		Expect(api.Status.Plans).To(HaveKeyWithValue(planKey, planID))
		Expect(k8s.IsAutomationAPIManaged(api)).To(BeFalse())
		Expect(subscriptionStatusInAPIM(client, apiID, subscription.ID)).To(Equal("ACCEPTED"))
	})
})
