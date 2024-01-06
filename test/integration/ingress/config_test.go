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
	v1 "k8s.io/api/core/v1"
)

var _ = Describe("Gravitee gateway config", func() {
	DescribeTable("validate keystore",
		func(given gateway.KeystoreConfig, expected error) {
			if expected == nil {
				Expect(given.Validate()).To(BeNil())
				return
			}
			Expect(given.Validate()).To(Equal(expected))
		},
		Entry(
			"With valid config and plain password", gateway.KeystoreConfig{
				Type:     "jks",
				Location: "/ns/secrets/name/key",
				Password: "password",
			},
			nil,
		),
		Entry(
			"With valid config and kubernetes password", gateway.KeystoreConfig{
				Type:     "jks",
				Location: "/ns/secrets/name/key",
				Password: "kubernetes://ns/secrets/name/key",
			},
			nil,
		),
		Entry(
			"With wrong type", gateway.KeystoreConfig{
				Type:     "pem",
				Location: "/ns/secrets/name/key",
				Password: "password",
			},
			fmt.Errorf("expected keystore type jks, got pem"),
		),
		Entry(
			"With wrong kubernetes location path", gateway.KeystoreConfig{
				Type:     "jks",
				Location: "/ns/secrets/name",
				Password: "password",
			},
			fmt.Errorf("expected kubernetes location format /$NS/(secrets|configmaps)/$NAME/$KEY, got /ns/secrets/name"),
		),
		Entry(
			"With wrong kubernetes location type", gateway.KeystoreConfig{
				Type:     "jks",
				Location: "/ns/wrong/name/key",
				Password: "password",
			},
			fmt.Errorf("expected kubernetes location format /$NS/(secrets|configmaps)/$NAME/$KEY, got /ns/wrong/name/key"),
		),
		Entry(
			"With  empty password", gateway.KeystoreConfig{
				Type:     "jks",
				Location: "/ns/secrets/name/key",
				Password: "",
			},
			fmt.Errorf("password is required"),
		),
		Entry(
			"With wrong kubernetes password password path", gateway.KeystoreConfig{
				Type:     "jks",
				Location: "/ns/secrets/name/key",
				Password: "kubernetes://ns",
			},
			fmt.Errorf(
				"expected password location format kubernetes://$NS/(secrets|configmaps)/$NAME/$KEY, got kubernetes://ns",
			),
		),
		Entry(
			"With wrong kubernetes password type", gateway.KeystoreConfig{
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

	Describe("Gravitee kube property", func() {
		It("Should parse path components with scheme", func() {
			prop := gateway.GraviteeKubeProperty("kubernetes://ns/secrets/name/key")
			Expect(prop.Namespace()).To(Equal("ns"))
			Expect(prop.Type()).To(Equal("secrets"))
			Expect(prop.Name()).To(Equal("name"))
			Expect(prop.Key()).To(Equal("key"))
		})

		It("Should parse path components without scheme", func() {
			prop := gateway.GraviteeKubeProperty("/ns/secrets/name/key")
			Expect(prop.Namespace()).To(Equal("ns"))
			Expect(prop.Type()).To(Equal("secrets"))
			Expect(prop.Name()).To(Equal("name"))
			Expect(prop.Key()).To(Equal("key"))
		})

		It("Should return a secret receiver", func() {
			prop := gateway.GraviteeKubeProperty("/ns/secrets/name/key")
			Expect(prop.NewReceiver()).To(BeAssignableToTypeOf(&v1.Secret{}))
		})

		It("Should return a config map receiver", func() {
			prop := gateway.GraviteeKubeProperty("/ns/configmaps/name/key")
			Expect(prop.NewReceiver()).To(BeAssignableToTypeOf(&v1.ConfigMap{}))
		})

		It("Should return a nil receiver", func() {
			prop := gateway.GraviteeKubeProperty("/ns/wrong/name/key")
			Expect(prop.NewReceiver()).To(BeNil())
		})

		It("Should get bytes from secret", func() {
			prop := gateway.GraviteeKubeProperty("/ns/secrets/name/key")
			secret, ok := prop.NewReceiver().(*v1.Secret)
			Expect(ok).To(BeTrue())
			secret.Data = map[string][]byte{
				"key": []byte("value"),
			}
			Expect(prop.Get(secret)).To(Equal([]byte("value")))
		})

		It("Should get bytes from config map", func() {
			prop := gateway.GraviteeKubeProperty("/ns/configmaps/name/key")
			cm, ok := prop.NewReceiver().(*v1.ConfigMap)
			Expect(ok).To(BeTrue())
			cm.Data = map[string]string{
				"key": "value",
			}
			Expect(prop.Get(cm)).To(Equal([]byte("value")))
		})

		It("Should return nil bytes", func() {
			prop := gateway.GraviteeKubeProperty("/ns/wrong/name/key")
			obj := prop.NewReceiver()
			Expect(obj).To(BeNil())
			Expect(prop.Get(obj)).To(BeNil())
		})
	})
})
