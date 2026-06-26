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

package drift

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func expectDrift(r drift.Result, expected string) {
	GinkgoHelper()
	Expect(r.String()).To(Equal(expected))
	Expect(r.DriftDetected()).To(BeTrue())
}

func expectNoDrift(r drift.Result) {
	GinkgoHelper()
	Expect(r.String()).To(BeEmpty())
	Expect(r.DriftDetected()).To(BeFalse())
}

func ptr[T any](v T) *T {
	GinkgoHelper()
	return &v
}
