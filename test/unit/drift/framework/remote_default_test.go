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

type withRemoteDefaultsConfig struct {
	Config *unstructured.Stringified `json:"config,omitempty" drift:"unstructured:ignore-remote-default"`
}

func remoteDefaults(obj map[string]any) withRemoteDefaultsConfig {
	return withRemoteDefaultsConfig{Config: unstructured.StringifiedFrom(obj)}
}

var _ = Describe("unstructured:ignore-remote-default", func() {
	DescribeTable("no drift on a key only the remote has",
		func(crd, remote withRemoteDefaultsConfig) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
		},
		Entry("top-level key",
			remoteDefaults(map[string]any{"timeout": 5000}),
			remoteDefaults(map[string]any{"timeout": 5000, "async": false}),
		),
		Entry("key in a nested map",
			remoteDefaults(map[string]any{"http": map[string]any{"timeout": 5000}}),
			remoteDefaults(map[string]any{"http": map[string]any{"timeout": 5000, "keepAlive": true}}),
		),
		Entry("key in a slice item",
			remoteDefaults(map[string]any{"items": []any{map[string]any{"name": "a"}}}),
			remoteDefaults(map[string]any{"items": []any{map[string]any{"name": "a", "enabled": true}}}),
		),
		Entry("list of strings under a key",
			remoteDefaults(map[string]any{"timeout": 5000}),
			remoteDefaults(map[string]any{"timeout": 5000, "allowMethods": []any{"GET", "POST"}}),
		),
		Entry("list of objects under a key",
			remoteDefaults(map[string]any{"timeout": 5000}),
			remoteDefaults(map[string]any{"timeout": 5000, "headers": []any{map[string]any{"name": "x"}}}),
		),
		Entry("list under a nested key",
			remoteDefaults(map[string]any{"timeout": 5000}),
			remoteDefaults(map[string]any{"timeout": 5000, "http": map[string]any{"methods": []any{"GET"}}}),
		),
	)

	DescribeTable("drift on a list item only the remote has",
		func(crd, remote withRemoteDefaultsConfig, report string) {
			expectDrift(drift.DetectWithNamespace(crd, remote, ""), report)
		},
		Entry("string item",
			remoteDefaults(map[string]any{"allowOrigin": []any{"a"}}),
			remoteDefaults(map[string]any{"allowOrigin": []any{"a", "x"}}),
			`config:
  allowOrigin[1]: <nil> != "x"`,
		),
		Entry("object item",
			remoteDefaults(map[string]any{"headers": []any{map[string]any{"name": "a"}}}),
			remoteDefaults(map[string]any{"headers": []any{map[string]any{"name": "a"}, map[string]any{"name": "b"}}}),
			`config:
  headers[1]: <nil> != map[name:b]`,
		),
	)

	It("detects drift on a key the CRD sets", func() {
		crd := remoteDefaults(map[string]any{"timeout": 5000})
		remote := remoteDefaults(map[string]any{"timeout": 1000, "async": false})
		expectDrift(drift.DetectWithNamespace(crd, remote, ""), `config:
  timeout: 5000 != 1000`)
	})

	It("detects drift on a key only the remote has, without the argument", func() {
		crd := withUnmaskedConfig{Config: unstructured.StringifiedFrom(map[string]any{"timeout": 5000})}
		remote := withUnmaskedConfig{Config: unstructured.StringifiedFrom(map[string]any{"timeout": 5000, "async": false})}
		expectDrift(drift.DetectWithNamespace(crd, remote, ""), `config:
  async: <nil> != false`)
	})
})
