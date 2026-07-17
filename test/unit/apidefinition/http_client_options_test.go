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

package apidefinition_test

import (
	"encoding/json"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HttpClientOptions pool fields", func() {
	Context("JSON serialization", func() {
		It("omits pool fields when nil (preserves APIM defaults)", func() {
			opts := &base.HttpClientOptions{
				KeepAlive: true,
			}
			data, err := json.Marshal(opts)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(data)).NotTo(ContainSubstring("maxWaitQueueSize"))
			Expect(string(data)).NotTo(ContainSubstring("maxConnectionLifetime"))
		})

		It("serializes maxWaitQueueSize when explicitly set to -1", func() {
			queueSize := -1
			opts := &base.HttpClientOptions{
				MaxWaitQueueSize: &queueSize,
			}
			data, err := json.Marshal(opts)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(data)).To(ContainSubstring(`"maxWaitQueueSize":-1`))
		})

		It("serializes maxWaitQueueSize when explicitly set to 0", func() {
			queueSize := 0
			opts := &base.HttpClientOptions{
				MaxWaitQueueSize: &queueSize,
			}
			data, err := json.Marshal(opts)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(data)).To(ContainSubstring(`"maxWaitQueueSize":0`))
		})

		It("serializes maxConnectionLifetime when explicitly set", func() {
			lifetime := int64(30000)
			opts := &base.HttpClientOptions{
				MaxConnectionLifetime: &lifetime,
			}
			data, err := json.Marshal(opts)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(data)).To(ContainSubstring(`"maxConnectionLifetime":30000`))
		})
	})

	Context("JSON deserialization", func() {
		It("leaves pool fields nil when absent from JSON", func() {
			input := `{"keepAlive":true,"keepAliveTimeout":30000}`
			opts := &base.HttpClientOptions{}
			err := json.Unmarshal([]byte(input), opts)
			Expect(err).NotTo(HaveOccurred())

			Expect(opts.MaxWaitQueueSize).To(BeNil())
			Expect(opts.MaxConnectionLifetime).To(BeNil())
		})

		It("deserializes maxWaitQueueSize from JSON", func() {
			input := `{"maxWaitQueueSize":-1,"keepAlive":true,"keepAliveTimeout":30000}`
			opts := &base.HttpClientOptions{}
			err := json.Unmarshal([]byte(input), opts)
			Expect(err).NotTo(HaveOccurred())

			Expect(opts.MaxWaitQueueSize).NotTo(BeNil())
			Expect(*opts.MaxWaitQueueSize).To(Equal(-1))
		})

		It("deserializes maxConnectionLifetime from JSON", func() {
			input := `{"maxConnectionLifetime":60000,"keepAlive":true,"keepAliveTimeout":30000}`
			opts := &base.HttpClientOptions{}
			err := json.Unmarshal([]byte(input), opts)
			Expect(err).NotTo(HaveOccurred())

			Expect(opts.MaxConnectionLifetime).NotTo(BeNil())
			Expect(*opts.MaxConnectionLifetime).To(Equal(int64(60000)))
		})
	})

	Context("round-trip", func() {
		It("preserves explicit values through marshal/unmarshal", func() {
			queueSize := 100
			lifetime := int64(120000)
			original := &base.HttpClientOptions{
				KeepAlive:             true,
				KeepAliveTimeout:      30000,
				MaxWaitQueueSize:      &queueSize,
				MaxConnectionLifetime: &lifetime,
			}

			data, err := json.Marshal(original)
			Expect(err).NotTo(HaveOccurred())

			restored := &base.HttpClientOptions{}
			err = json.Unmarshal(data, restored)
			Expect(err).NotTo(HaveOccurred())

			Expect(restored.MaxWaitQueueSize).NotTo(BeNil())
			Expect(*restored.MaxWaitQueueSize).To(Equal(100))
			Expect(restored.MaxConnectionLifetime).NotTo(BeNil())
			Expect(*restored.MaxConnectionLifetime).To(Equal(int64(120000)))
		})

		It("preserves nil through marshal/unmarshal (no zero-value leak)", func() {
			original := &base.HttpClientOptions{
				KeepAlive:        true,
				KeepAliveTimeout: 30000,
			}

			data, err := json.Marshal(original)
			Expect(err).NotTo(HaveOccurred())

			restored := &base.HttpClientOptions{}
			err = json.Unmarshal(data, restored)
			Expect(err).NotTo(HaveOccurred())

			Expect(restored.MaxWaitQueueSize).To(BeNil())
			Expect(restored.MaxConnectionLifetime).To(BeNil())
		})
	})
})
