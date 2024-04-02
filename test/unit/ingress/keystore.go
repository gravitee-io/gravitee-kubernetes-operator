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

package ingress

import (
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/gateway"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gravitee gateway config", func() {
	DescribeTable("should validate keystore",
		func(given gateway.KeystoreConfig, expected error) {
			if expected == nil {
				Expect(given.Validate()).To(BeNil())
				return
			}
			Expect(given.Validate()).To(Equal(expected))
		},
		Entry(
			"with valid config and plain password", gateway.KeystoreConfig{
				Type:     "jks",
				Location: "/ns/secrets/name/key",
				Password: "password",
			},
			nil,
		),
		Entry(
			"with valid config and kubernetes password", gateway.KeystoreConfig{
				Type:     "jks",
				Location: "/ns/secrets/name/key",
				Password: "kubernetes://ns/secrets/name/key",
			},
			nil,
		),
		Entry(
			"with wrong type", gateway.KeystoreConfig{
				Type:     "pem",
				Location: "/ns/secrets/name/key",
				Password: "password",
			},
			fmt.Errorf("expected keystore type jks, got pem"),
		),
		Entry(
			"with wrong kubernetes location path", gateway.KeystoreConfig{
				Type:     "jks",
				Location: "/ns/secrets/name",
				Password: "password",
			},
			fmt.Errorf("expected kubernetes location format /$NS/(secrets|configmaps)/$NAME/$KEY, got /ns/secrets/name"),
		),
		Entry(
			"with wrong kubernetes location type", gateway.KeystoreConfig{
				Type:     "jks",
				Location: "/ns/wrong/name/key",
				Password: "password",
			},
			fmt.Errorf("expected kubernetes location format /$NS/(secrets|configmaps)/$NAME/$KEY, got /ns/wrong/name/key"),
		),
		Entry(
			"with  empty password", gateway.KeystoreConfig{
				Type:     "jks",
				Location: "/ns/secrets/name/key",
				Password: "",
			},
			fmt.Errorf("password is required"),
		),
		Entry(
			"with wrong kubernetes password password path", gateway.KeystoreConfig{
				Type:     "jks",
				Location: "/ns/secrets/name/key",
				Password: "kubernetes://ns",
			},
			fmt.Errorf(
				"expected password location format kubernetes://$NS/(secrets|configmaps)/$NAME/$KEY, got kubernetes://ns",
			),
		),
		Entry(
			"with wrong kubernetes password type", gateway.KeystoreConfig{
				Type:     "jks",
				Location: "kubernetes://ns/wrong/name/key",
				Password: "",
			},
			fmt.Errorf(
				"%s %s",
				"expected kubernetes location format /$NS/(secrets|configmaps)/$NAME/$KEY,",
				"got kubernetes://ns/wrong/name/key",
			),
		),
	)
})
