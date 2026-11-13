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

package http_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	gkohttp "github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
)

var _ = Describe("NewStdClient", func() {
	saved := env.Config

	AfterEach(func() {
		env.Config = saved
	})

	transportOf := func(c *http.Client) *http.Transport {
		t, ok := c.Transport.(*http.Transport)
		Expect(ok).To(BeTrue())
		return t
	}

	It("applies the configured timeout", func() {
		env.Config.HTTPClientTimeoutSeconds = 7

		c, err := gkohttp.NewStdClient()

		Expect(err).ToNot(HaveOccurred())
		Expect(c.Timeout).To(Equal(7 * time.Second))
	})

	It("applies insecure skip verify", func() {
		env.Config.HTTPClientInsecureSkipVerify = true

		c, err := gkohttp.NewStdClient()

		Expect(err).ToNot(HaveOccurred())
		Expect(transportOf(c).TLSClientConfig.InsecureSkipVerify).To(BeTrue())
	})

	It("routes through the configured proxy with its credentials", func() {
		env.Config.HttpProxy = env.HttpProxy{
			Enabled:  true,
			URL:      "http://proxy.local:3128",
			Username: "user",
			Password: "secret",
		}

		c, err := gkohttp.NewStdClient()
		Expect(err).ToNot(HaveOccurred())

		req, _ := http.NewRequest(http.MethodGet, "https://am.example/automation", nil)
		proxyURL, err := transportOf(c).Proxy(req)
		Expect(err).ToNot(HaveOccurred())
		Expect(proxyURL.String()).To(Equal("http://user:secret@proxy.local:3128"))
	})

	It("fails when the trust store file cannot be read", func() {
		env.Config.HttpClientTrustStorePath = "/does/not/exist.pem"

		_, err := gkohttp.NewStdClient()

		Expect(err).To(MatchError(ContainSubstring("failed to read trust store file")))
	})

	It("shares one transport between clients, each keeping its own timeout", func() {
		env.Config.HTTPClientTimeoutSeconds = 3
		first, err := gkohttp.NewStdClient()
		Expect(err).ToNot(HaveOccurred())

		env.Config.HTTPClientTimeoutSeconds = 9
		second, err := gkohttp.NewStdClient()
		Expect(err).ToNot(HaveOccurred())

		Expect(transportOf(second)).To(BeIdenticalTo(transportOf(first)))
		Expect(first.Timeout).To(Equal(3 * time.Second))
		Expect(second.Timeout).To(Equal(9 * time.Second))
	})

	It("builds a new transport when a setting changes", func() {
		first, err := gkohttp.NewStdClient()
		Expect(err).ToNot(HaveOccurred())

		env.Config.HTTPClientInsecureSkipVerify = !env.Config.HTTPClientInsecureSkipVerify
		second, err := gkohttp.NewStdClient()
		Expect(err).ToNot(HaveOccurred())

		Expect(transportOf(second)).ToNot(BeIdenticalTo(transportOf(first)))
		Expect(transportOf(second).TLSClientConfig.InsecureSkipVerify).To(Equal(env.Config.HTTPClientInsecureSkipVerify))
	})

	It("loads the trust store again when its file changes", func() {
		path := filepath.Join(GinkgoT().TempDir(), "ca.pem")
		Expect(os.WriteFile(path, selfSignedCA("first"), 0o600)).To(Succeed())
		env.Config.HttpClientTrustStorePath = path

		first, err := gkohttp.NewStdClient()
		Expect(err).ToNot(HaveOccurred())
		unchanged, err := gkohttp.NewStdClient()
		Expect(err).ToNot(HaveOccurred())
		Expect(transportOf(unchanged)).To(BeIdenticalTo(transportOf(first)))

		By("rotating the trust store")
		Expect(os.WriteFile(path, append(selfSignedCA("first"), selfSignedCA("second")...), 0o600)).To(Succeed())
		rotated, err := gkohttp.NewStdClient()
		Expect(err).ToNot(HaveOccurred())
		Expect(transportOf(rotated)).ToNot(BeIdenticalTo(transportOf(first)))
	})
})

// selfSignedCA returns a PEM-encoded self-signed CA certificate.
func selfSignedCA(name string) []byte {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	Expect(err).ToNot(HaveOccurred())
	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: name},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	Expect(err).ToNot(HaveOccurred())
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
}
