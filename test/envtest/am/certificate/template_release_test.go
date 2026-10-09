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

package amcertificate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	kErrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	internal "github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/certificate"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Template references of a certificate refused by AM", func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval
	ctx := context.Background()

	It("should release the Secret when the certificate is deleted", func() {
		fixtures := secretFixtures()
		cert := fixtures.AMCertificate
		fixtures.AMCertificate = nil
		useOwnSecret(fixtures, cert)
		fixtures.Apply()
		secretKey := client.ObjectKeyFromObject(fixtures.Secrets[1])

		By("deleting the domain in AM only: the certificate's upsert fails after its templates are compiled")

		dto, err := internal.ToCertificateDTO(cert)
		Expect(err).ToNot(HaveOccurred())
		resp, err := am.NewSDKClient().DeleteDomainWithResponse(ctx, dto.DomainKey)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode()).To(Equal(http.StatusNoContent), string(resp.Body))

		Expect(manager.Client().Create(ctx, cert)).To(Succeed())

		By("expecting the Secret to reference the certificate")

		certID := fmt.Sprintf(`"%s/%s"`, cert.Namespace, cert.Name)
		Eventually(ctx, func() error {
			secret := &corev1.Secret{}
			if err := manager.Client().Get(ctx, secretKey, secret); err != nil {
				return err
			}
			if refs := secret.Annotations["gravitee.io/amcertificates"]; !strings.Contains(refs, certID) {
				return fmt.Errorf("references %q", refs)
			}
			return nil
		}, timeout, interval).Should(Succeed())

		By("expecting the refused certificate to carry its finalizer, so its deletion releases the Secret")

		Eventually(ctx, func() []string {
			Expect(manager.GetLatest(ctx, cert)).To(Succeed())
			return cert.GetFinalizers()
		}, timeout, interval).Should(ContainElement(core.AMCertificateFinalizer))

		By("deleting the certificate, then the Secret")

		Expect(manager.Client().Delete(ctx, cert.DeepCopy())).To(Succeed())
		Eventually(ctx, func() error {
			return assert.Deleted(ctx, "AMCertificate", cert)
		}, timeout, interval).Should(Succeed(), cert.Name)

		Expect(manager.Client().Delete(ctx, fixtures.Secrets[1].DeepCopy())).To(Succeed())
		Eventually(ctx, func() bool {
			err := manager.Client().Get(ctx, secretKey, &corev1.Secret{})
			return kErrors.IsNotFound(err)
		}, timeout, interval).Should(BeTrue(), "the Secret is still held by %v", core.TemplatingFinalizer)
	})
})

// useOwnSecret renames the keystore Secret the certificate templates read, so that no certificate
// left by another spec still references it when this one is deleted.
func useOwnSecret(fixtures *fixture.Objects, cert *v1alpha1.AMCertificate) {
	secret := fixtures.Secrets[1]
	shared := "`" + secret.Name + "/"
	secret.Name += fixtures.GetGeneratedSuffix()

	b, err := json.Marshal(cert)
	Expect(err).ToNot(HaveOccurred())
	b = bytes.ReplaceAll(b, []byte(shared), []byte("`"+secret.Name+"/"))
	*cert = v1alpha1.AMCertificate{}
	Expect(json.Unmarshal(b, cert)).To(Succeed())
}
