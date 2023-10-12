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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	kErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var _ = Describe("Deleting a management context", func() {
	Context("Not linked to an api definition", func() {
		var contextFixture *gio.ManagementContext
		var contextLookupKey types.NamespacedName

		BeforeEach(func() {
			By("Initializing the management context fixture")
			fixtureGenerator := internal.NewFixtureGenerator()

			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Context: internal.ClusterContextFile,
			})

			Expect(err).ToNot(HaveOccurred())

			contextFixture = fixtures.Context
			contextLookupKey = types.NamespacedName{Name: contextFixture.Name, Namespace: namespace}
		})

		It("Should delete the management context", func() {
			By("Creating a new management context")

			Expect(k8sClient.Create(ctx, contextFixture)).Should(Succeed())

			By("Getting created resource and expect to find it")

			createdContext := new(gio.ManagementContext)
			Eventually(func() error {
				return k8sClient.Get(ctx, contextLookupKey, createdContext)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Deleting the management context")
			Expect(k8sClient.Delete(ctx, createdContext)).ToNot(HaveOccurred())

			By("Checking the management context has been deleted")
			context := &gio.ManagementContext{}
			Eventually(func() error {
				return k8sClient.Get(ctx, contextLookupKey, context)
			}, timeout, interval).ShouldNot(Succeed())
		})
	})

	Context("Linked to an api definition", func() {
		var contextFixture *gio.ManagementContext
		var contextLookupKey types.NamespacedName
		var apiFixture *gio.ApiDefinition
		var apiLookupKey types.NamespacedName

		BeforeEach(func() {
			By("Initializing the management context fixture and api definition")
			fixtureGenerator := internal.NewFixtureGenerator()

			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Api:     internal.ApiWithContextFile,
				Context: internal.ClusterContextFile,
			})

			Expect(err).ToNot(HaveOccurred())

			contextFixture = fixtures.Context
			contextLookupKey = types.NamespacedName{Name: contextFixture.Name, Namespace: namespace}

			apiFixture = fixtures.Api
			apiLookupKey = types.NamespacedName{Name: apiFixture.Name, Namespace: namespace}
		})

		AfterEach(func() {
			Expect(k8sClient.Delete(ctx, apiFixture)).Should(Succeed())
		})

		It("Should not delete the management context", func() {
			By("Creating a new management context")
			Expect(k8sClient.Create(ctx, contextFixture)).Should(Succeed())

			By("Getting created resource and expect to find it")
			createdContext := new(gio.ManagementContext)
			Eventually(func() error {
				if err := k8sClient.Get(ctx, contextLookupKey, createdContext); err != nil {
					return err
				}
				return internal.AssertEquals("finalizer.length", 1, len(createdContext.ObjectMeta.Finalizers))
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Creating a api definition")
			Expect(k8sClient.Create(ctx, apiFixture)).Should(Succeed())

			By("Getting created resource and expect to find it")
			createdApi := new(gio.ApiDefinition)
			Consistently(func() error { // just to let gko have time to configure API definition
				return k8sClient.Get(ctx, apiLookupKey, createdApi)
			}, timeout/10, interval).Should(Succeed())

			By("Trying to delete the management context")
			Expect(k8sClient.Delete(ctx, createdContext)).ToNot(HaveOccurred())

			By("Checking the management context has not been deleted")
			context := &gio.ManagementContext{}
			Eventually(func() error {
				return k8sClient.Get(ctx, contextLookupKey, context)
			}, timeout, interval).Should(Succeed())
		})
	})

	Context("With secret Reference", func() {
		var ctx1 *gio.ManagementContext
		var ctx2 *gio.ManagementContext
		var secret *v1.Secret

		var ctx1Key types.NamespacedName
		var ctx2Key types.NamespacedName
		var secretKey types.NamespacedName

		secretName := "test-context-secret"

		BeforeEach(func() {
			secret = &v1.Secret{}
			secret.Data = map[string][]byte{
				"username": []byte("admin"),
				"password": []byte("admin"),
			}

			secret.Name = secretName
			secret.Namespace = namespace

			Expect(k8sClient.Create(ctx, secret)).Should(Succeed())

			secretKey = types.NamespacedName{
				Name:      secret.Name,
				Namespace: secret.Namespace,
			}

			fix1, err := internal.NewFixtureGenerator().NewFixtures(internal.FixtureFiles{
				Context: internal.ClusterContextFile,
			}, func(fix *internal.Fixtures) {
				fix.Context.Spec.Auth.SecretRef = &refs.NamespacedName{
					Name:      secretName,
					Namespace: namespace,
				}
			})

			Expect(err).ToNot(HaveOccurred())

			ctx1 = fix1.Context
			ctx1Key = types.NamespacedName{
				Name:      ctx1.Name,
				Namespace: ctx1.Namespace,
			}

			Expect(k8sClient.Create(ctx, ctx1)).Should(Succeed())

			fix2, err := internal.NewFixtureGenerator().NewFixtures(internal.FixtureFiles{
				Context: internal.ClusterContextFile,
			}, func(fix *internal.Fixtures) {
				fix.Context.Spec.Auth.SecretRef = &refs.NamespacedName{
					Name:      "test-context-secret",
					Namespace: namespace,
				}
			})

			Expect(err).ToNot(HaveOccurred())

			ctx2 = fix2.Context
			ctx2Key = types.NamespacedName{
				Name:      ctx2.Name,
				Namespace: ctx2.Namespace,
			}

			Expect(k8sClient.Create(ctx, ctx2)).Should(Succeed())
		})

		It("Should keep secret finalizer while referenced", func() {
			Expect(k8sClient.Delete(ctx, ctx1)).ToNot(HaveOccurred())

			Eventually(func() error {
				err := k8sClient.Get(ctx, ctx1Key, ctx1)
				if err == nil {
					return fmt.Errorf("management context %s still exists", ctx1Key)
				}
				if !kErrors.IsNotFound(err) {
					return err
				}
				return nil
			}, timeout, interval).Should(Succeed())

			Eventually(func() error {
				err := k8sClient.Get(ctx, secretKey, secret)
				if err != nil {
					return err
				}
				if !controllerutil.ContainsFinalizer(secret, keys.ManagementContextSecretFinalizer) {
					return fmt.Errorf("Expected secret to contain finalizer %s", keys.ManagementContextSecretFinalizer)
				}
				return nil
			})

			Expect(k8sClient.Delete(ctx, ctx2)).ToNot(HaveOccurred())

			Eventually(func() error {
				err := k8sClient.Get(ctx, ctx2Key, ctx2)
				if err == nil {
					return fmt.Errorf("management context %s still exists", ctx2Key)
				}
				if !kErrors.IsNotFound(err) {
					return err
				}
				return nil
			}, timeout, interval).Should(Succeed())

			Eventually(func() error {
				err := k8sClient.Get(ctx, secretKey, secret)
				if err != nil {
					return err
				}
				if controllerutil.ContainsFinalizer(secret, keys.ManagementContextSecretFinalizer) {
					return fmt.Errorf("Expected finalizer %s to be removed on secret", keys.ManagementContextSecretFinalizer)
				}
				return nil
			})

			Expect(k8sClient.Delete(ctx, secret)).ToNot(HaveOccurred())

			Eventually(func() error {
				err := k8sClient.Get(ctx, secretKey, secret)
				if err == nil {
					return fmt.Errorf("secret %s still exists", secretKey)
				}
				if !kErrors.IsNotFound(err) {
					return err
				}
				return nil
			}, timeout, interval).Should(Succeed())

		})
	})

})
