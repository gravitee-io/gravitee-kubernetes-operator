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

package amclientsecret

import (
	"context"
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Create", func() {
	ctx := context.Background()

	It("should create the secret under its application, owned by it, with the value nowhere in the cluster", func() {
		fixtures := fixtures()
		createSecret(ctx, fixtures.AMClientSecret, withValue(fixtures.AMClientSecret, value))
		fixtures.Apply()
		app, cs := fixtures.AMApplication, fixtures.AMClientSecret
		expectAccepted(ctx, cs)

		Expect(cs.Status.Key).To(Equal(refs.NewNamespacedNameFromObject(cs).HRID()))
		Expect(cs.Status.ApplicationKey).To(Equal(app.Status.Key))
		Expect(cs.Status.DomainKey).To(Equal(app.Status.DomainKey))
		Expect(remoteSecret(ctx, app, cs)).To(Equal(http.StatusOK))

		By("expecting the application to own the client secret, not as its controller")

		Expect(cs.OwnerReferences).To(HaveLen(1))
		owner := cs.OwnerReferences[0]
		Expect(owner.Kind).To(Equal("AMApplication"))
		Expect(owner.UID).To(Equal(app.UID))
		Expect(owner.Controller).To(Or(BeNil(), HaveValue(BeFalse())))

		By("expecting the value neither in the CR nor in its events (AC3)")

		stored, err := json.Marshal(cs)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(stored)).ToNot(ContainSubstring(value))

		events := &corev1.EventList{}
		Expect(manager.Client().List(ctx, events, client.InNamespace(cs.Namespace))).To(Succeed())
		for _, event := range events.Items {
			Expect(event.Message).ToNot(ContainSubstring(value))
		}
	})

	It("should leave a single secret in AM when reconciled again (idempotence)", func() {
		fixtures := fixtures()
		secret := createSecret(ctx, fixtures.AMClientSecret, withValue(fixtures.AMClientSecret, value))
		fixtures.Apply()
		app, cs := fixtures.AMApplication, fixtures.AMClientSecret
		expectAccepted(ctx, cs)

		By("touching the Secret without changing the value: the watch reconciles the client secret")

		Expect(manager.GetLatest(ctx, secret)).To(Succeed())
		secret.Labels = map[string]string{"test.gravitee.io/touch": "true"}
		Expect(manager.Client().Update(ctx, secret)).To(Succeed())

		Consistently(ctx, func() error {
			if err := manager.GetLatest(ctx, cs); err != nil {
				return err
			}
			return assert.AMClientSecretAccepted(cs)
		}, "2s", "200ms").Should(Succeed())
		Expect(remoteSecrets(ctx, app)).To(HaveLen(1))
	})

	It("should report AM's immutability error when the value changes in the Secret (AC6)", func() {
		fixtures := fixtures()
		secret := createSecret(ctx, fixtures.AMClientSecret, withValue(fixtures.AMClientSecret, value))
		fixtures.Apply()
		app, cs := fixtures.AMApplication, fixtures.AMClientSecret
		expectAccepted(ctx, cs)

		Expect(manager.GetLatest(ctx, secret)).To(Succeed())
		secret.Data = nil
		secret.StringData = withValue(cs, otherValue)
		Expect(manager.Client().Update(ctx, secret)).To(Succeed())

		expectCondition(ctx, cs, k8s.ConditionAccepted, "rotate it with a new AMClientSecret")
		Expect(cs.Status.Conditions).ToNot(ContainElement(HaveField("Message", ContainSubstring(otherValue))))
		Expect(remoteSecret(ctx, app, cs)).To(Equal(http.StatusOK))
	})

	It("should rotate with a second client secret, then the delete of the first (AC5)", func() {
		fixtures := fixtures()
		firstSecret := createSecret(ctx, fixtures.AMClientSecret, withValue(fixtures.AMClientSecret, value))
		fixtures.Apply()
		app, first := fixtures.AMApplication, fixtures.AMClientSecret
		expectAccepted(ctx, first)

		By("creating a second client secret with a new Secret")

		second := &v1alpha1.AMClientSecret{
			ObjectMeta: metaV1.ObjectMeta{Name: first.Name + "-next", Namespace: first.Namespace},
			Spec:       *first.Spec.DeepCopy(),
		}
		second.Spec.SecretRef.Name += "-next"
		createSecret(ctx, second, withValue(second, otherValue))
		Expect(manager.Client().Create(ctx, second)).To(Succeed())
		expectAccepted(ctx, second)
		Expect(remoteSecrets(ctx, app)).To(HaveLen(2))

		By("deleting the first client secret: its Secret is kept")

		deleteAndWait(ctx, "AMClientSecret", first)
		Expect(remoteSecrets(ctx, app)).To(ConsistOf(HaveField("Key", second.Status.Key)))
		Expect(manager.GetLatest(ctx, firstSecret)).To(Succeed())
	})

	It("should wait for its application, its Secret and the entry when applied before them", func() {
		fixtures := fixtures()
		app, cs := fixtures.AMApplication, fixtures.AMClientSecret
		fixtures.AMApplication, fixtures.AMClientSecret = nil, nil
		fixtures.Apply()

		By("creating the client secret alone")

		Expect(manager.Client().Create(ctx, cs)).To(Succeed())
		expectCondition(ctx, cs, k8s.ConditionResolvedRefs, app.Name)

		By("creating the application: the client secret now waits for its Secret")

		Expect(manager.Client().Create(ctx, app)).To(Succeed())
		expectCondition(ctx, cs, k8s.ConditionResolvedRefs, "Secret ["+cs.Spec.SecretRef.Name+"] not found")
		Expect(manager.GetLatest(ctx, app)).To(Succeed())
		Expect(remoteSecrets(ctx, app)).To(BeEmpty())

		By("creating the Secret without the entry")

		secret := createSecret(ctx, cs, map[string]string{"other": value})
		expectCondition(ctx, cs, k8s.ConditionResolvedRefs, "key ["+cs.Spec.SecretRef.Key+"]")
		Expect(remoteSecrets(ctx, app)).To(BeEmpty())

		By("adding the entry: the watch picks the client secret up")

		Expect(manager.GetLatest(ctx, secret)).To(Succeed())
		secret.StringData = withValue(cs, value)
		Expect(manager.Client().Update(ctx, secret)).To(Succeed())
		expectAccepted(ctx, cs)
		Expect(remoteSecrets(ctx, app)).To(HaveLen(1))
	})
})
