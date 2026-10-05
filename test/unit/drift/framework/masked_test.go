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

package framework

import (
	"github.com/gravitee-io/gravitee-automation-sdk/common/pkg/unstructured"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
)

type withMaskedConfig struct {
	Config *unstructured.Stringified `json:"config,omitempty" drift:"unstructured:masked"`
}

type withUnmaskedConfig struct {
	Config *unstructured.Stringified `json:"config,omitempty" drift:"unstructured"`
}

func masked(obj map[string]any) withMaskedConfig {
	return withMaskedConfig{Config: unstructured.StringifiedFrom(obj)}
}

var _ = Describe("unstructured:masked", func() {
	DescribeTable("no drift",
		func(crd, remote withMaskedConfig) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
		},
		Entry("masked top-level value",
			masked(map[string]any{"password": "s3cret"}),
			masked(map[string]any{"password": "********"}),
		),
		Entry("masked value in a slice item",
			masked(map[string]any{"users": []any{map[string]any{"username": "a", "password": "s3cret"}}}),
			masked(map[string]any{"users": []any{map[string]any{"username": "a", "password": "********"}}}),
		),
		Entry("masked value in a nested map",
			masked(map[string]any{"store": map[string]any{"password": "s3cret"}}),
			masked(map[string]any{"store": map[string]any{"password": "********"}}),
		),
	)

	It("detects drift on unmasked values only", func() {
		crd := masked(map[string]any{"password": "s3cret", "host": "a"})
		remote := masked(map[string]any{"password": "********", "host": "b"})

		expectDrift(drift.DetectWithNamespace(crd, remote, ""), `config:
  host: "a" != "b"`)
	})

	It("detects drift on a masked value changed to another value", func() {
		crd := masked(map[string]any{"password": "s3cret"})
		remote := masked(map[string]any{"password": "other"})

		expectDrift(drift.DetectWithNamespace(crd, remote, ""), `config:
  password: "s3cret" != "other"`)
	})

	It("detects drift on a masked value without the masked argument", func() {
		crd := withUnmaskedConfig{Config: unstructured.StringifiedFrom(map[string]any{"password": "s3cret"})}
		remote := withUnmaskedConfig{Config: unstructured.StringifiedFrom(map[string]any{"password": "********"})}

		expectDrift(drift.DetectWithNamespace(crd, remote, ""), `config:
  password: "s3cret" != "********"`)
	})

	It("does not drift when the CRD changes a masked value", func() {
		oldCRD := masked(map[string]any{"password": "a"})
		newCRD := masked(map[string]any{"password": "b"})
		remote := masked(map[string]any{"password": "********"})

		or := drift.DetectWithNamespace(oldCRD, remote, "")
		nr := drift.DetectWithNamespace(newCRD, remote, "")
		expectNoDrift(drift.Merge(or, nr))
	})
})
