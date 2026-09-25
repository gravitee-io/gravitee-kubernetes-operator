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
	"net/http"
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
})
