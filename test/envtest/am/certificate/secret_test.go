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
package amcertificate

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

var _ = Describe("Secret-backed keystore", func() {
	ctx := context.Background()

	// expectInAM waits for AM to hold the keystore file, untouched by GKO, and the store password.
	expectInAM := func(cert *v1alpha1.AMCertificate, keystore, storepass string) {
		file := fmt.Sprintf(`{"name":"keystore.p12","content":"%s"}`, keystore)
		Eventually(ctx, func() error {
			status, remote := remoteCertificate(ctx, cert)
			if status != http.StatusOK {
				return fmt.Errorf("status %d", status)
			}
			if content := remote.Configuration.Object["content"]; content != file {
				return fmt.Errorf("content %v", content)
			}
			if pass := remote.Configuration.Object["storepass"]; pass != storepass {
				return fmt.Errorf("storepass %v", pass)
			}
			return nil
		}, constants.EventualTimeout, constants.Interval).Should(Succeed(), cert.Name)
	}

	// expectNoSecretValues reads the CR back from the API server: the Secret's values appear nowhere in it.
	expectNoSecretValues := func(cert *v1alpha1.AMCertificate, secret *corev1.Secret) {
		Expect(manager.GetLatest(ctx, cert)).To(Succeed())
		stored, err := json.Marshal(cert)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(stored)).To(ContainSubstring("[[ secret"))
		for key, value := range secret.Data {
			Expect(string(stored)).ToNot(ContainSubstring(string(value)), key)
		}
	}

	It("should send the keystore and passwords from the Secret, and follow its rotation", func() {
		fixtures := secretFixtures().Apply()
		cert := fixtures.AMCertificate
		secret := &corev1.Secret{}
		Expect(manager.Client().Get(ctx, client.ObjectKeyFromObject(fixtures.Secrets[1]), secret)).To(Succeed())

		By("expecting AM to hold the Secret's keystore and password (AC3)")

		expectInAM(cert, string(secret.Data["keystore.b64"]), string(secret.Data["storepass"]))
		expectNoSecretValues(cert, secret)

		By("rotating the keystore and the store password in the Secret (AC4)")

		secret.Data["keystore.b64"] = []byte("cm90YXRlZC1rZXlzdG9yZQ==")
		secret.Data["storepass"] = []byte("rotated-storepass")
		Expect(manager.Client().Update(ctx, secret)).To(Succeed())

		expectInAM(cert, "cm90YXRlZC1rZXlzdG9yZQ==", "rotated-storepass")

		By("expecting no Secret value in the CR (AC2)")

		expectNoSecretValues(cert, secret)
	})
})
