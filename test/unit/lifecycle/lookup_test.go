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

package lifecycle_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/lifecycle/ref"
)

var _ = Describe("Lookup", func() {
	It("finds an exact registered kind", func() {
		_, ok := ref.Lookup("secret")
		Expect(ok).To(BeTrue())
	})

	It("finds a plural registry name from a singular tag", func() {
		ref.Register("gadgets", func() client.Object { return &corev1.Secret{} }, passthroughExtract)
		_, ok := ref.Lookup("gadget")
		Expect(ok).To(BeTrue())
		_, ok = ref.Lookup("gadgets")
		Expect(ok).To(BeTrue())
	})

	It("finds a singular registry name from a plural tag", func() {
		_, ok := ref.Lookup("secrets")
		Expect(ok).To(BeTrue())
		_, ok = ref.Lookup("amcontexts")
		Expect(ok).To(BeTrue())
	})

	It("does not invent a kind", func() {
		_, ok := ref.Lookup("nope")
		Expect(ok).To(BeFalse())
		_, ok = ref.Lookup("nopes")
		Expect(ok).To(BeFalse())
	})
})

func passthroughExtract(obj client.Object, key string) (any, error) {
	sec := obj.(*corev1.Secret)
	return sec.Data[key], nil
}
