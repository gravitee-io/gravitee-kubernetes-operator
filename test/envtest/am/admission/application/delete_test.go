// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package amapplication

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/amcertificate"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/amidentityprovider"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Validate delete of a certificate or identity provider an application uses (D2)", func() {
	ctx := context.Background()
	certificateCtrl := amcertificate.NewAdmissionCtrl()
	identityProviderCtrl := amidentityprovider.NewAdmissionCtrl()

	review := "You can review those by running the following command: " +
		"kubectl get amapplications.gravitee.io -n " + constants.Namespace

	It("should reject deleting a certificate an application references, until the application is deleted", func() {
		fixtures := withDomain(constants.AMApplicationMinimalFile).Apply()
		cert, app := fixtures.AMCertificate, fixtures.AMApplication

		Eventually(func() error {
			_, err := certificateCtrl.ValidateDelete(ctx, cert)
			return err
		}, constants.EventualTimeout, constants.Interval).Should(MatchError(And(
			ContainSubstring("cannot be deleted because 1 AM applications use it"), ContainSubstring(review))))

		By("deleting the application")

		Expect(manager.Client().Delete(ctx, app.DeepCopy())).To(Succeed())
		Eventually(ctx, func() error {
			return assert.Deleted(ctx, "AMApplication", app)
		}, constants.EventualTimeout, constants.Interval).Should(Succeed(), app.Name)

		Eventually(func() error {
			_, err := certificateCtrl.ValidateDelete(ctx, cert)
			return err
		}, constants.EventualTimeout, constants.Interval).Should(Succeed())
	})

	It("should reject deleting a certificate an application uses by key", func() {
		fixtures := withDomain(constants.AMApplicationKeysFile)
		fixtures.AMApplication.Spec.Certificate = new(refs.NewNamespacedNameFromObject(fixtures.AMCertificate).HRID())
		fixtures.Apply()

		Eventually(func() error {
			_, err := certificateCtrl.ValidateDelete(ctx, fixtures.AMCertificate)
			return err
		}, constants.EventualTimeout, constants.Interval).Should(MatchError(
			ContainSubstring("cannot be deleted because 1 AM applications use it")))
	})

	It("should reject deleting an identity provider an application references, until the application is deleted", func() {
		fixtures := withDomain(constants.AMApplicationFullFile).Apply()
		idp, app := fixtures.AMIdentityProvider, fixtures.AMApplication

		Eventually(func() error {
			_, err := identityProviderCtrl.ValidateDelete(ctx, idp)
			return err
		}, constants.EventualTimeout, constants.Interval).Should(MatchError(And(
			ContainSubstring("cannot be deleted because 1 AM applications use it"), ContainSubstring(review))))

		By("deleting the application")

		Expect(manager.Client().Delete(ctx, app.DeepCopy())).To(Succeed())
		Eventually(ctx, func() error {
			return assert.Deleted(ctx, "AMApplication", app)
		}, constants.EventualTimeout, constants.Interval).Should(Succeed(), app.Name)

		Eventually(func() error {
			_, err := identityProviderCtrl.ValidateDelete(ctx, idp)
			return err
		}, constants.EventualTimeout, constants.Interval).Should(Succeed())
	})

	It("should reject deleting an identity provider an application uses by key", func() {
		fixtures := withDomain(constants.AMApplicationKeysFile)
		fixtures.AMApplication.Spec.IdentityProviders[0].Identity =
			new(refs.NewNamespacedNameFromObject(fixtures.AMIdentityProvider).HRID())
		fixtures.Apply()

		Eventually(func() error {
			_, err := identityProviderCtrl.ValidateDelete(ctx, fixtures.AMIdentityProvider)
			return err
		}, constants.EventualTimeout, constants.Interval).Should(MatchError(
			ContainSubstring("cannot be deleted because 1 AM applications use it")))
	})
})
