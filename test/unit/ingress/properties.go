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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/gateway"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
)

var _ = Describe("Gravitee kubernetes properties", func() {

	It("should parse path components with scheme", func() {
		prop := gateway.GraviteeKubeProperty("kubernetes://ns/secrets/name/key")
		Expect(prop.Namespace()).To(Equal("ns"))
		Expect(prop.Type()).To(Equal("secrets"))
		Expect(prop.Name()).To(Equal("name"))
		Expect(prop.Key()).To(Equal("key"))
	})

	It("should parse path components without scheme", func() {
		prop := gateway.GraviteeKubeProperty("/ns/secrets/name/key")
		Expect(prop.Namespace()).To(Equal("ns"))
		Expect(prop.Type()).To(Equal("secrets"))
		Expect(prop.Name()).To(Equal("name"))
		Expect(prop.Key()).To(Equal("key"))
	})

	It("should return a secret receiver", func() {
		prop := gateway.GraviteeKubeProperty("/ns/secrets/name/key")
		Expect(prop.NewReceiver()).To(BeAssignableToTypeOf(&v1.Secret{}))
	})

	It("should return a config map receiver", func() {
		prop := gateway.GraviteeKubeProperty("/ns/configmaps/name/key")
		Expect(prop.NewReceiver()).To(BeAssignableToTypeOf(&v1.ConfigMap{}))
	})

	It("should return a nil receiver", func() {
		prop := gateway.GraviteeKubeProperty("/ns/wrong/name/key")
		Expect(prop.NewReceiver()).To(BeNil())
	})

	It("should get bytes from secret", func() {
		prop := gateway.GraviteeKubeProperty("/ns/secrets/name/key")
		secret, ok := prop.NewReceiver().(*v1.Secret)
		Expect(ok).To(BeTrue())
		secret.Data = map[string][]byte{
			"key": []byte("value"),
		}
		Expect(prop.Get(secret)).To(Equal([]byte("value")))
	})

	It("should get bytes from config map", func() {
		prop := gateway.GraviteeKubeProperty("/ns/configmaps/name/key")
		cm, ok := prop.NewReceiver().(*v1.ConfigMap)
		Expect(ok).To(BeTrue())
		cm.Data = map[string]string{
			"key": "value",
		}
		Expect(prop.Get(cm)).To(Equal([]byte("value")))
	})

	It("should return nil bytes", func() {
		prop := gateway.GraviteeKubeProperty("/ns/wrong/name/key")
		obj := prop.NewReceiver()
		Expect(obj).To(BeNil())
		Expect(prop.Get(obj)).To(BeNil())
	})
})
