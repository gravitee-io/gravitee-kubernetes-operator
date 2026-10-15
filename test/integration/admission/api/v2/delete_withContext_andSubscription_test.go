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

package v2

import (
	"context"
	"fmt"

	adm "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v2"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"

	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"
)

var _ = Describe("Validate delete", labels.WithContext, func() {
	ctx := context.Background()
	admissionCtrl := adm.AdmissionCtrl{}

	fixtures := fixture.
		Builder().
		WithApplication(constants.ApplicationWithClientIDFile).
		WithAPI(constants.ApiWithJWTPlan).
		WithContext(constants.ContextWithCredentialsFile).
		WithSubscription(constants.SubscriptionFile).
		Build()

	clientID := random.GetName()
	fixtures.Application.Spec.Settings.App.ClientID = &clientID
	fixtures.Subscription.Spec.API.Name = fixtures.API.Name
	fixtures.Subscription.Spec.API.Kind = core.CRDApiDefinitionResource
	fixtures.Subscription.Namespace = constants.Namespace

	fixtures.Apply()

	referencedError := func() error {
		return errors.NewSeveref(
			"cannot delete [%s] because it is referenced in 1 subscriptions. "+
				"Subscriptions must be deleted before the API definition. "+
				"You can review the subscriptions using the following command: "+
				"kubectl get subscriptions.gravitee.io -A "+
				"-o jsonpath='{.items[?(@.spec.api.name==\"%s\")].metadata.name}'",
			fixtures.API.GetRef(), fixtures.API.Name,
		)
	}

	It("should fail with subscription", func() {
		Eventually(func() error {
			Expect(manager.GetLatest(ctx, fixtures.API)).ToNot(HaveOccurred())
			_, err := admissionCtrl.ValidateDelete(ctx, fixtures.API)
			return assert.Equals("error", referencedError(), err)
		}, constants.EventualTimeout, constants.Interval).Should(Succeed())
	})

	It("should fail with subscription when the subscription counter is lost", func() {
		By("setting the counter to 0, as a write conflict on subscription creation does")

		Eventually(func() error {
			if err := manager.GetLatest(ctx, fixtures.API); err != nil {
				return err
			}
			if fixtures.API.Status.SubscriptionCount == 0 {
				return nil
			}
			fixtures.API.Status.SubscriptionCount = 0
			if err := manager.Client().Status().Update(ctx, fixtures.API); err != nil {
				return err
			}
			return fmt.Errorf("subscription counter not seen at 0 yet")
		}, constants.EventualTimeout, constants.Interval).Should(Succeed())

		// validated against the object seen at 0: a reconcile may write the stale counter back
		_, err := admissionCtrl.ValidateDelete(ctx, fixtures.API)
		Expect(assert.Equals("error", referencedError(), err)).To(Succeed())
	})
})
