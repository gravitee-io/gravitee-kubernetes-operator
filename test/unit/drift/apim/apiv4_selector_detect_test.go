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

package apim

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
)

// apim-sdk v0.2.1 leaves out a selector path and operator the CR does not set, and APIM stores its defaults.
var _ = Describe("API v4 Drift detection on the selector defaults APIM stores", func() {
	It("admits an HTTP selector without path against the / and STARTS_WITH APIM stores", func() {
		crd := crdDriftV4(`{"name": "api", "flows": [{"name": "all", "enabled": true, "selectors": [{"type": "HTTP"}]}]}`)
		remote := remoteDriftV4(crd,
			`{"name": "api", "flows": [{"name": "all", "enabled": true, "selectors": [{"type": "HTTP", "path": "/", "pathOperator": "STARTS_WITH"}]}]}`)
		expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
	})
})
