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
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	wk "github.com/gravitee-io/gravitee-kubernetes-operator/internal/webhook"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Webhook", labels.WithContext, func() {
	timeout := constants.EventualTimeout

	interval := constants.Interval
	It("should create Key, Cert and CA", func() {
		ca, cert, key := wk.GenerateCerts(context.TODO(), "localhost")

		var c tls.Certificate
		var err error
		Eventually(func() error {
			c, err = tls.X509KeyPair(cert, key)
			if err != nil {
				return err
			}

			return nil
		}, timeout, interval).Should(Succeed())

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(ca)

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:    caCertPool,
				ServerName: "localhost",
				MinVersion: tls.VersionTLS12,
			},
		}

		ts := httptest.NewUnstartedServer(http.HandlerFunc(handler))
		ts.TLS = &tls.Config{Certificates: []tls.Certificate{c}, MinVersion: tls.VersionTLS12}
		ts.StartTLS()
		defer ts.Close()

		var res *http.Response
		Eventually(func() error {
			client := &http.Client{Transport: tr}
			res, err = client.Get(ts.URL)
			if err != nil {
				return err
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				return fmt.Errorf("response code was %v; want 200", res.StatusCode)
			}
			var body []byte
			body, err = io.ReadAll(res.Body)
			if err != nil {
				return err
			}

			expected := []byte("Hello World")

			if !bytes.Equal(expected, body) {
				return fmt.Errorf("response body was '%v'; want '%v'", expected, body)
			}

			return nil
		}, timeout, interval).Should(Succeed())
	})
})

func handler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello World")
}
