package application

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"strings"
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	admission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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
const (
	updatedDescription     = "updated description"
	remoteDescription      = "remote updated description"
	localCRDDescription    = "local CRD description"
	driftDescriptionAssert = `description: "local CRD description" != "remote updated description"`
)

var _ = Describe("Validate drift", labels.WithContext, func() {
	ctx := context.Background()
	admissionCtrl := admission.AdmissionCtrl{}

	It("should not drift on with minimal fields", func() {
		fixtures := fixture.
			Builder().
			WithApplication(constants.Application).
			WithContext(constants.ContextWithCredentialsFile).
			Build()
		fixtures.Apply()

		By("changing the application description")
		newApp := fixtures.Application.DeepCopy()

		setDescription(newApp, updatedDescription)

		_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Application, newApp)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should detect drift with minimal fields", func() {
		fixtures := fixture.
			Builder().
			WithApplication(constants.Application).
			WithContext(constants.ContextWithCredentialsFile).
			Build()
		fixtures.Apply()

		By("changing the remote application description")
		newApp := fixtures.Application.DeepCopy()

		validateDescriptionDrift(ctx, admissionCtrl, fixtures.Application, newApp, fixtures.Context)
	})

	It("should not drift on with all fields", func() {
		fixtures := fixture.
			Builder().
			WithApplication(constants.ApplicationDriftFull).
			WithGroup(constants.GroupFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build()

		fixtures.Application.Spec.Groups[0] = fixtures.Group.Name
		// so it is different every time
		fixtures.Application.Spec.Settings.TLS.ClientCertificates[0].Content = getClientTLSCert()
		fixtures.Apply()

		By("changing the application description")
		newApp := fixtures.Application.DeepCopy()

		setDescription(newApp, updatedDescription)

		_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Application, newApp)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should detect drift with all fields", func() {
		fixtures := fixture.
			Builder().
			WithApplication(constants.ApplicationDriftFull).
			WithGroup(constants.GroupFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build()

		fixtures.Application.Spec.Groups[0] = fixtures.Group.Name
		// so it is different every time
		fixtures.Application.Spec.Settings.TLS.ClientCertificates[0].Content = getClientTLSCert()
		fixtures.Apply()

		By("changing the remote application description")

		newApp := fixtures.Application.DeepCopy()

		validateDescriptionDrift(ctx, admissionCtrl, fixtures.Application, newApp, fixtures.Context)
	})

	It("should ignore remote drift when the annotation disables detection", func() {
		fixtures := fixture.
			Builder().
			WithApplication(constants.Application).
			WithContext(constants.ContextWithCredentialsFile).
			Build()
		fixtures.Apply()

		newApp := fixtures.Application.DeepCopy()
		validateDescriptionDrift(ctx, admissionCtrl, fixtures.Application, newApp, fixtures.Context)

		By("disabling drift detection on the application")
		setDriftAnnotation(newApp, env.FalseString)

		_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Application, newApp)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should detect drift when globally disabled if the annotation forces detection", func() {
		original := env.Config.DriftDetection.Enabled
		env.Config.DriftDetection.Enabled = false
		DeferCleanup(func() {
			env.Config.DriftDetection.Enabled = original
		})

		fixtures := fixture.
			Builder().
			WithApplication(constants.Application).
			WithContext(constants.ContextWithCredentialsFile).
			Build()
		fixtures.Apply()

		By("forcing drift detection with the annotation")
		newApp := fixtures.Application.DeepCopy()
		setDriftAnnotation(newApp, env.TrueString)
		validateDescriptionDrift(ctx, admissionCtrl, fixtures.Application, newApp, fixtures.Context)
	})

	It("should not drift when CRD realigns with remote", func() {
		fixtures := fixture.
			Builder().
			WithApplication(constants.Application).
			WithContext(constants.ContextWithCredentialsFile).
			Build()
		fixtures.Apply()

		By("changing the remote application description")
		newApp := fixtures.Application.DeepCopy()
		newApp.PopulateIDs(fixtures.Context, true)
		setDescription(newApp, remoteDescription)

		apimClient := apim.NewClient(ctx)
		_, err := apimClient.Applications.CreateOrUpdate(newApp)
		Expect(err).ToNot(HaveOccurred())

		By("updating the CRD to match the remote description")
		Eventually(func() error {
			_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Application, newApp)
			return err
		}, constants.EventualTimeout, constants.Interval).Should(Succeed())
	})

})

func setDescription(app *v1alpha1.Application, description string) {
	GinkgoHelper()
	app.Spec.Description = description
}

func setDriftAnnotation(app *v1alpha1.Application, value string) {
	GinkgoHelper()
	annotations := app.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations[core.DriftDetectionAnnotation] = value
	app.SetAnnotations(annotations)
}

func validateDescriptionDrift(
	ctx context.Context,
	admissionCtrl admission.AdmissionCtrl,
	oldApp, newApp *v1alpha1.Application,
	mgmtContext core.ContextModel,
) {
	GinkgoHelper()
	apimClient := apim.NewClient(ctx)

	newApp.PopulateIDs(mgmtContext, true)
	setDescription(newApp, remoteDescription)

	_, err := apimClient.Applications.CreateOrUpdate(newApp)
	Expect(err).ToNot(HaveOccurred())

	setDescription(newApp, localCRDDescription)

	Eventually(func() error {
		_, err := admissionCtrl.ValidateUpdate(ctx, oldApp, newApp)
		return assert.DriftDetected(driftDescriptionAssert, err)
	}, constants.EventualTimeout, constants.Interval).Should(Succeed())
}

func getPrivateKey() *rsa.PrivateKey {
	GinkgoHelper()
	const pemPrivateKey = `
-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCxoeCUW5KJxNPxMp+KmCxKLc1Zv9Ny+4CFqcUXVUYH69L3mQ7v
IWrJ9GBfcaA7BPQqUlWxWM+OCEQZH1EZNIuqRMNQVuIGCbz5UQ8w6tS0gcgdeGX7
J7jgCQ4RK3F/PuCM38QBLaHx988qG8NMc6VKErBjctCXFHQt14lerd5KpQIDAQAB
AoGAYrf6Hbk+mT5AI33k2Jt1kcweodBP7UkExkPxeuQzRVe0KVJw0EkcFhywKpr1
V5eLMrILWcJnpyHE5slWwtFHBG6a5fLaNtsBBtcAIfqTQ0Vfj5c6SzVaJv0Z5rOd
7gQF6isy3t3w9IF3We9wXQKzT6q5ypPGdm6fciKQ8RnzREkCQQDZwppKATqQ41/R
vhSj90fFifrGE6aVKC1hgSpxGQa4oIdsYYHwMzyhBmWW9Xv/R+fPyr8ZwPxp2c12
33QwOLPLAkEA0NNUb+z4ebVVHyvSwF5jhfJxigim+s49KuzJ1+A2RaSApGyBZiwS
rWvWkB471POAKUYt5ykIWVZ83zcceQiNTwJBAMJUFQZX5GDqWFc/zwGoKkeR49Yi
MTXIvf7Wmv6E++eFcnT461FlGAUHRV+bQQXGsItR/opIG7mGogIkVXa3E1MCQARX
AAA7eoZ9AEHflUeuLn9QJI/r0hyQQLEtrpwv6rDT1GCWaLII5HJ6NUFVf4TTcqxo
6vdM4QGKTJoO+SaCyP0CQFdpcxSAuzpFcKv0IlJ8XzS/cy+mweCMwyJ1PFEc4FX6
wg/HcAJWY60xZTJDFN+Qfx8ZQvBEin6c2/h+zZi5IVY=
-----END RSA PRIVATE KEY-----
`
	block, _ := pem.Decode([]byte(pemPrivateKey))

	testPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		Fail("Failed to parse private key: %s " + err.Error())
	}

	return testPrivateKey
}

func getClientTLSCert() string {

	GinkgoHelper()
	random := rand.Reader

	ecdsaPriv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		Fail("Failed to generate ECDSA key:" + err.Error())
	}

	pub := &ecdsaPriv.PublicKey
	priv := getPrivateKey()
	sigAlgo := x509.SHA256WithRSA

	commonName := "test.example.com"
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"Gravitee"},
			Country:      []string{"FR"},
		},
		NotBefore: time.Now().Add(-24 * time.Hour),
		NotAfter:  time.Now().Add(24 * time.Hour),

		SignatureAlgorithm: sigAlgo,

		SubjectKeyId: []byte{1, 2, 3, 4},
		KeyUsage:     x509.KeyUsageDigitalSignature,

		BasicConstraintsValid: true,
		IsCA:                  false,

		OCSPServer:            []string{"http://ocsp.example.com"},
		IssuingCertificateURL: []string{"http://crt.example.com/ca1.crt"},
	}

	derBytes, err := x509.CreateCertificate(random, &template, &template, pub, priv)
	if err != nil {
		Fail("failed to create certificate: %s" + err.Error())
	}

	certData := strings.Builder{}
	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	}

	if err := pem.Encode(&certData, block); err != nil {
		Fail("failed to encode certificate: %s" + err.Error())
	}

	return certData.String()
}
