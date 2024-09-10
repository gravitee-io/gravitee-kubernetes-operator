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

package webhook

import (
	"context"
	"fmt"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"

	wk "github.com/gravitee-io/gravitee-kubernetes-operator/internal/webhook"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Webhook", labels.WithContext, func() {
	timeout := constants.EventualTimeout

	interval := constants.Interval
	It("should create Key, Cert and CA", func() {
		cli := manager.Client()
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "webhook-secret" + random.GetSuffix(),
				Namespace: "default",
			},
		}

		Expect(cli.Create(context.Background(), secret)).To(Succeed())

		By("creating the secret and updating the webhook secret")
		webhookPatcher := wk.NewWebhookPatcher()
		Expect(webhookPatcher.CreateSecret(context.Background(), secret.Name,
			secret.Namespace, "webhook.server")).To(Succeed())

		Eventually(func() error {
			err := cli.Get(context.Background(), types.NamespacedName{
				Namespace: secret.Namespace,
				Name:      secret.Name,
			}, secret)

			if err != nil {
				return err
			}

			if secret.Data == nil {
				return fmt.Errorf("webhook secret data is nil")
			}

			if secret.Data["ca"] == nil {
				return fmt.Errorf("webhook secret CA is nil")
			}

			if secret.Data["cert"] == nil {
				return fmt.Errorf("webhook secret Cert is nil")
			}

			if secret.Data["key"] == nil {
				return fmt.Errorf("webhook secret Key is nil")
			}

			return nil
		}, timeout, interval).Should(Succeed())

		By("patching the webhook admission configuration")
		wvc := newWebhookAdmissionConfiguration()
		Expect(cli.Create(context.Background(), wvc)).To(Succeed())

		Expect(webhookPatcher.UpdateValidationCaBundle(context.Background(), wvc.Name, secret.Name, secret.Namespace)).To(Succeed())

		Eventually(func() error {
			if err := cli.Get(context.Background(), types.NamespacedName{
				Namespace: wvc.Namespace,
				Name:      wvc.Name,
			}, wvc); err != nil {
				return err
			}

			if wvc.Webhooks[0].ClientConfig.CABundle == nil {
				return fmt.Errorf("webhook caBundle is nil")
			}

			return nil
		})
	})

})

func newWebhookAdmissionConfiguration() *v1.ValidatingWebhookConfiguration {
	path := "/path"
	port := int32(3456) //nolint:gomnd // static port
	scope := v1.AllScopes
	failurePolicy := v1.Fail
	matchPolicy := v1.Equivalent
	sideEffects := v1.SideEffectClassNone
	timeoutSeconds := int32(10) //nolint:gomnd // default timeout

	return &v1.ValidatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "webhook-configuration" + random.GetSuffix(),
			Namespace: "default",
		},
		Webhooks: []v1.ValidatingWebhook{
			{
				Name: webhookServerName(),
				ClientConfig: v1.WebhookClientConfig{
					Service: &v1.ServiceReference{
						Namespace: "default",
						Name:      "webhook-service",
						Path:      &path,
						Port:      &port,
					},
				},
				Rules: []v1.RuleWithOperations{{
					Operations: []v1.OperationType{
						"CREATE",
					},
					Rule: v1.Rule{
						APIGroups:   []string{"gravitee.io"},
						APIVersions: []string{"v10"},
						Resources:   []string{"someapidefinitions"},
						Scope:       &scope,
					},
				}},
				FailurePolicy:           &failurePolicy,
				MatchPolicy:             &matchPolicy,
				SideEffects:             &sideEffects,
				TimeoutSeconds:          &timeoutSeconds,
				AdmissionReviewVersions: []string{"v1"},
			},
		},
	}
}

func webhookServerName() string {
	suffix := random.GetSuffix()
	index := strings.LastIndex(suffix, "-")
	return fmt.Sprintf("webhook.server.%s", suffix[index+1:])
}
