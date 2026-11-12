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

package application

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"

	appmodel "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// newClientCertificate returns a self-signed, non-CA client certificate in PEM.
// APIM refuses a certificate another active application already holds, so the spec
// cannot reuse the fixture certificates other specs register.
func newClientCertificate() string {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	Expect(err).ToNot(HaveOccurred())

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 62))
	Expect(err).ToNot(HaveOccurred())

	template := x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: random.GetName()},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().AddDate(20, 0, 0),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	Expect(err).ToNot(HaveOccurred())

	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}))
}

var _ = Describe("Validate update", labels.WithContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := application.AdmissionCtrl{}

	// The certificate end-date check reads the security type of every subscribed plan.
	// A PUSH plan has none, which used to panic the webhook (GKO-3237).
	It("should not fail when a PUSH plan subscription ends after all client certificates", func() {
		fixtures := fixture.
			Builder().
			WithAPIv4(constants.ApiV4MessageWebhookPushPlanFile).
			WithApplication(constants.ApplicationWithClientIDFile).
			WithSubscription(constants.SubscriptionPushPlanFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build()

		clientID := random.GetName()
		fixtures.Application.Spec.Settings.App.ClientID = &clientID
		fixtures.Application.Spec.Settings.TLS = &appmodel.TLSSettings{
			ClientCertificates: []appmodel.ClientCertificate{
				{Name: "client", Content: newClientCertificate(), EndsAt: "2030-01-01T00:00:00Z"},
			},
		}

		subEndsAt := "2035-01-01T00:00:00Z"
		fixtures.Subscription.Spec.EndingAt = &subEndsAt

		fixtures.Apply()

		By("waiting for the subscription to be reconciled, so the application search finds it")

		Eventually(func() error {
			return assert.SubscriptionCompleted(fixtures.Subscription)
		}, constants.EventualTimeout, interval).Should(Succeed(), fixtures.Subscription.Name)

		By("checking that application update validation succeeds")

		Eventually(func() error {
			_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Application, fixtures.Application.DeepCopy())
			return assert.Equals("error", nil, err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})
})
