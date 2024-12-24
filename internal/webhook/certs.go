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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"strings"
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	CertName              = "cert"
	KeyName               = "key"
	caName                = "ca"
	ValidatingWebhookName = "gko-validating-webhook-configurations"
	MutatingWebhookName   = "gko-mutating-webhook-configurations"
)

type Patcher struct {
	client *kubernetes.Clientset
}

func NewWebhookPatcher() *Patcher {
	conf := ctrl.GetConfigOrDie()
	cli, err := kubernetes.NewForConfig(conf)
	if err != nil {
		panic(err)
	}

	return &Patcher{
		client: cli,
	}
}

func (p *Patcher) CreateSecret(ctx context.Context, secretName, namespace, host string) error {
	ca := p.getCaFromSecret(ctx, secretName, namespace)
	if ca == nil {
		log.Global.Info("Creating a new CA secret for further usage by admission webhook")
		newCa, newCert, newKey := GenerateCerts(ctx, host)
		return p.saveCertsToSecret(ctx, secretName, namespace, CertName, KeyName, newCa, newCert, newKey)
	} else {
		log.Global.Infof("Found CA secret [%s] matching admission webhook configuration", secretName)
	}

	return nil
}

func (p *Patcher) UpdateValidationCaBundle(ctx context.Context, webhookName, secretName, ns string) error {
	webhookConfig, err := p.client.AdmissionregistrationV1().
		ValidatingWebhookConfigurations().Get(ctx, webhookName, metav1.GetOptions{})

	if errors.IsNotFound(err) {
		log.Global.Error(err, "Validating webhook configuration could not be found")
		return err
	} else if err != nil {
		log.Global.Error(err, "Unable to get validating webhook")
		return err
	}

	caBundle := p.getCaFromSecret(ctx, secretName, ns)
	for i := range webhookConfig.Webhooks {
		webhookConfig.Webhooks[i].ClientConfig.CABundle = caBundle
	}
	_, err = p.client.AdmissionregistrationV1().
		ValidatingWebhookConfigurations().Update(ctx, webhookConfig, metav1.UpdateOptions{})
	if err != nil {
		log.Global.Error(err, "Cannot update validating webhook configuration")
		return err
	}

	return nil
}

func (p *Patcher) UpdateMutationCaBundle(ctx context.Context, webhookName, secretName, ns string) error {
	webhookConfig, err := p.client.AdmissionregistrationV1().
		MutatingWebhookConfigurations().Get(ctx, webhookName, metav1.GetOptions{})

	if errors.IsNotFound(err) {
		log.Global.Error(err, "Mutating webhook configuration could not be found")
		return err
	} else if err != nil {
		log.Global.Error(err, "Unable to get mutating webhook")
		return err
	}

	caBundle := p.getCaFromSecret(ctx, secretName, ns)
	for i := range webhookConfig.Webhooks {
		webhookConfig.Webhooks[i].ClientConfig.CABundle = caBundle
	}
	_, err = p.client.AdmissionregistrationV1().
		MutatingWebhookConfigurations().Update(ctx, webhookConfig, metav1.UpdateOptions{})
	if err != nil {
		log.Global.Error(err, "Cannot update mutating webhook configuration")
		return err
	}

	return nil
}

// getCaFromSecret will check for the presence of a secret. If it exists, will return the content of the
// "ca" from the secret, otherwise will return nil.
func (p *Patcher) getCaFromSecret(ctx context.Context, secretName string, namespace string) []byte {
	log.Global.Debugf("Getting secret [%s] in namespace [%s]", secretName, namespace)

	secret, err := p.client.CoreV1().Secrets(namespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		log.Global.Error(err, "Could not get admission webhook CA secret")
		panic(err)
	}

	return secret.Data["ca"]
}

// SaveCertsToSecret saves the provided ca, cert and key into a secret in the specified namespace.
func (p *Patcher) saveCertsToSecret(ctx context.Context,
	secretName, namespace, certName, keyName string, ca, cert, key []byte) error {
	log.Global.Debugf("Saving to webhook secret '%s' in namespace '%s'", secretName, namespace)

	secret, err := p.client.CoreV1().Secrets(namespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	secret.Data = map[string][]byte{caName: ca, certName: cert, keyName: key}
	_, err = p.client.CoreV1().Secrets(namespace).Update(ctx, secret, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// GenerateCerts venerates a ca with a leaf certificate and key and returns the ca, cert and key as PEM encoded slices.
func GenerateCerts(ctx context.Context, host string) ([]byte, []byte, []byte) {
	notBefore := time.Now().Add(time.Minute * -5)
	notAfter := notBefore.Add(100 * 365 * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128) //nolint:gomnd // LSH number
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Global.Error(err, "Failed to generate admission webhook serial number")
	}
	rootKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Global.Error(err, "Failed to generate admission webhook key")
	}

	rootTemplate := x509.Certificate{
		SerialNumber:          serialNumber,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		Subject:               pkix.Name{Organization: []string{"gravitee.io"}},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &rootTemplate, &rootTemplate, &rootKey.PublicKey, rootKey)
	if err != nil {
		log.Global.Error(err, "Failed to create CA certificate for admission webhook")
	}

	ca := encodeCert(derBytes)

	leafKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Global.Error(err, "Failed to create leaf key for admission webhook certificate")
	}

	key := encodeKey(leafKey)

	serialNumber, err = rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Global.Error(err, "Failed to generate serial number for admission webhook certificate")
	}
	leafTemplate := x509.Certificate{
		SerialNumber:          serialNumber,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
		Subject:               pkix.Name{Organization: []string{"gko"}},
	}
	hosts := strings.Split(host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			leafTemplate.IPAddresses = append(leafTemplate.IPAddresses, ip)
		} else {
			leafTemplate.DNSNames = append(leafTemplate.DNSNames, h)
		}
	}

	derBytes, err = x509.CreateCertificate(rand.Reader, &leafTemplate, &rootTemplate, &leafKey.PublicKey, rootKey)
	if err != nil {
		log.Global.Error(err, "Failed to create leaf certificate for admission webhook")
	}

	cert := encodeCert(derBytes)
	return ca, cert, key
}

func encodeKey(key *ecdsa.PrivateKey) []byte {
	b, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		log.Global.Error(err, "Unable to marshal ECDSA private key for admission webhook")
	}
	return pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: b})
}

func encodeCert(derBytes []byte) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
}
