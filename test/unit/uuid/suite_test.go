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

package uuid

import (
	"testing"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUUID(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "package uuid")
}

var _ = Describe("FromStrings", func() {
	It("generates a predictable UUID", func() {
		firstRun := uuid.FromStrings("foo", "bar")
		secondRun := uuid.FromStrings("foo", "bar")
		Expect(secondRun).To(Equal(firstRun))
	})
})

var _ = Describe("JavaUUIDFromBytes", func() {
	It("generates a predictable UUID same as Java", func() {
		firstRun := uuid.JavaUUIDFromBytes("foo")
		secondRun := uuid.JavaUUIDFromBytes("foo")
		Expect(secondRun).To(Equal(firstRun))
	})
})
