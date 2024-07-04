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

package admissionwebhook

import (
	"context"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"

	wk "github.com/gravitee-io/gravitee-kubernetes-operator/internal/webhook"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Delete", labels.WithContext, func() {
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

	})
})
