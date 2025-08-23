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

package impl

import (
	_ "embed"

	"gopkg.in/yaml.v2"
	v1 "sigs.k8s.io/gateway-api/conformance/apis/v1"
)

//go:embed impl.yaml
var manifest []byte
var Manifest v1.Implementation

func GetReportOutputPath() string {
	return "standard-" + Manifest.Version + "-default-report.yaml"
}

func init() {
	impl := new(v1.Implementation)
	if err := yaml.Unmarshal(manifest, impl); err != nil {
		panic(err)
	}
	Manifest = *impl
}
