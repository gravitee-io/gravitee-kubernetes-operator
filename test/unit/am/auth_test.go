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

package am_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	ammodel "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
)

var _ = Describe("AM bearer auth", func() {
	It("stores a token set from a secret", func() {
		auth := &ammodel.Auth{}
		auth.SetToken("from-secret")
		Expect(auth.GetBearerToken()).To(Equal("from-secret"))
	})

	It("keeps a secretRef", func() {
		auth := &ammodel.Auth{}
		auth.SetSecretRef(&refs.NamespacedName{Namespace: "gravitee", Name: "am-automation-token"})
		Expect(auth.GetSecretRef().GetName()).To(Equal("am-automation-token"))
		Expect(auth.GetSecretRef().GetNamespace()).To(Equal("gravitee"))
	})
})
