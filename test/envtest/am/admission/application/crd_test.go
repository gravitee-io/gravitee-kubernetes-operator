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
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appmodel "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/am/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"
)

// The CRD rules (schema and CEL) are evaluated by the API server, not by the webhook handler.
var _ = Describe("CRD validation", func() {
	ctx := context.Background()

	create := func(obj client.Object) error {
		err := manager.Client().Create(ctx, obj)
		if err == nil {
			DeferCleanup(func() { Expect(client.IgnoreNotFound(manager.Client().Delete(ctx, obj))).To(Succeed()) })
		}
		return err
	}
	newApplication := func() *v1alpha1.AMApplication {
		app := &v1alpha1.AMApplication{
			ObjectMeta: metaV1.ObjectMeta{Name: random.GetName(), Namespace: constants.Namespace},
		}
		app.Spec.DomainRef = refs.NamespacedName{Name: "some-domain"}
		app.Spec.Name = "Application"
		app.Spec.Type = "native"
		app.Spec.Certificate = new("signing")
		return app
	}

	It("should require exactly one of certificate or certificateRef", func() {
		neither := newApplication()
		neither.Spec.Certificate = nil
		Expect(create(neither)).To(MatchError(ContainSubstring("exactly one of certificate or certificateRef must be set")))

		both := newApplication()
		both.Spec.CertificateRef = &refs.NamespacedName{Name: "signing"}
		Expect(create(both)).To(MatchError(ContainSubstring("exactly one of certificate or certificateRef must be set")))
	})

	It("should require exactly one of identity or identityRef on an identity provider", func() {
		neither := newApplication()
		neither.Spec.IdentityProviders = []appmodel.AutomationApplicationIdentityProvider{{}}
		Expect(create(neither)).To(MatchError(ContainSubstring("exactly one of identity or identityRef must be set")))

		both := newApplication()
		both.Spec.IdentityProviders = []appmodel.AutomationApplicationIdentityProvider{
			{Identity: new("users"), IdentityRef: &refs.NamespacedName{Name: "users"}},
		}
		Expect(create(both)).To(MatchError(ContainSubstring("exactly one of identity or identityRef must be set")))
	})

	It("should reject empty ref names", func() {
		app := newApplication()
		app.Spec.DomainRef.Name = ""
		Expect(create(app)).To(MatchError(ContainSubstring("domainRef.name must not be empty")))

		app = newApplication()
		app.Spec.Certificate, app.Spec.CertificateRef = nil, &refs.NamespacedName{}
		Expect(create(app)).To(MatchError(ContainSubstring("certificateRef.name must not be empty")))

		app = newApplication()
		app.Spec.IdentityProviders = []appmodel.AutomationApplicationIdentityProvider{
			{IdentityRef: &refs.NamespacedName{}},
		}
		Expect(create(app)).To(MatchError(ContainSubstring("identityRef.name must not be empty")))
	})

	Describe("immutable fields", func() {
		var app *v1alpha1.AMApplication

		BeforeEach(func() {
			app = newApplication()
			Expect(create(app)).To(Succeed())
			Eventually(func() error { return manager.GetLatest(ctx, app.DeepCopy()) },
				constants.EventualTimeout, constants.Interval).Should(Succeed())
		})

		// update retries on conflicts: the controller keeps updating the status of this application,
		// whose domain does not exist.
		update := func(mutate func(*v1alpha1.AMApplication)) error {
			return retry.RetryOnConflict(retry.DefaultRetry, func() error {
				latest := app.DeepCopy()
				Expect(manager.GetLatest(ctx, latest)).To(Succeed())
				mutate(latest)
				return manager.Client().Update(ctx, latest)
			})
		}

		It("should reject a domainRef change", func() {
			Expect(update(func(a *v1alpha1.AMApplication) { a.Spec.DomainRef.Name = "other-domain" })).
				To(MatchError(ContainSubstring("domainRef is immutable")))
		})

		It("should admit setting a clientId, then reject changing or removing it", func() {
			Expect(update(func(a *v1alpha1.AMApplication) { a.Spec.ClientId = new("client") })).To(Succeed())

			Expect(update(func(a *v1alpha1.AMApplication) { a.Spec.ClientId = new("other-client") })).
				To(MatchError(ContainSubstring("clientId is immutable once set")))
			Expect(update(func(a *v1alpha1.AMApplication) { a.Spec.ClientId = nil })).
				To(MatchError(ContainSubstring("clientId is immutable once set")))
		})

		It("should admit a type change", func() {
			Expect(update(func(a *v1alpha1.AMApplication) { a.Spec.Type = "web" })).To(Succeed())
		})
	})
})
