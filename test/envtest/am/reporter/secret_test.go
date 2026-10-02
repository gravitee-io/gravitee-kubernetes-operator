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
package amreporter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Secret-backed password", func() {
	ctx := context.Background()

	// expectInAM waits for AM to hold the Kafka password.
	expectInAM := func(reporter *v1alpha1.AMReporter, password string) {
		Eventually(ctx, func() error {
			status, remote := remoteReporter(ctx, reporter)
			if status != http.StatusOK {
				return fmt.Errorf("status %d", status)
			}
			if sent := remote.Configuration.Object["password"]; sent != password {
				return fmt.Errorf("password %v", sent)
			}
			return nil
		}, constants.EventualTimeout, constants.Interval).Should(Succeed(), reporter.Name)
	}

	// expectNoSecretValues reads the CR back from the API server: the Secret's values appear nowhere in it,
	// spec, status and condition messages included.
	expectNoSecretValues := func(reporter *v1alpha1.AMReporter, secret *corev1.Secret) {
		Expect(manager.GetLatest(ctx, reporter)).To(Succeed())
		stored, err := json.Marshal(reporter)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(stored)).To(ContainSubstring("[[ secret"))
		for key, value := range secret.Data {
			Expect(string(stored)).ToNot(ContainSubstring(string(value)), key)
		}
	}

	It("should send the password from the Secret, and follow its rotation", func() {
		fixtures := secretFixtures().Apply()
		reporter := fixtures.AMReporter
		secret := &corev1.Secret{}
		Expect(manager.Client().Get(ctx, client.ObjectKeyFromObject(fixtures.Secrets[1]), secret)).To(Succeed())

		By("expecting AM to hold the Secret's password (AC3)")

		expectInAM(reporter, string(secret.Data["password"]))
		expectNoSecretValues(reporter, secret)

		By("rotating the password in the Secret (AC3)")

		secret.Data["password"] = []byte("rotated-password")
		Expect(manager.Client().Update(ctx, secret)).To(Succeed())

		expectInAM(reporter, "rotated-password")

		By("expecting no Secret value in the CR (AC2)")

		expectNoSecretValues(reporter, secret)
	})
})
