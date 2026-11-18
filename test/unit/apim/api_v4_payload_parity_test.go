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

package apim_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/yaml"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/apiv4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

// The fields GKO used to send that the Automation API does not declare on PUT /apis (APIM-16599). The typed
// SDK payload drops them; any other path missing from it is a field the SDK decode lost silently.
var undeclaredV4Paths = []string{
	"id",
	"definitionVersion",
	"definitionContext",
	"definitionContext.origin",
	"definitionContext.syncFrom",
	"plans[].crossId",
	"plans[].order",
	"plans[].comment_required",
	"plans[].definitionVersion",
	"pages[].order",
	"pages[].api",
}

var _ = Describe("ApiV4Definition payload parity", func() {
	files, err := filepath.Glob("../../../examples/apim/api_definition/v4/*.yml")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		for i, api := range apiV4Documents(file) {
			It(fmt.Sprintf("sends every declared field of %s #%d", filepath.Base(file), i), func() {
				before := payloadPaths(legacyV4Payload(api))
				after := payloadPaths(sdkV4Payload(api))

				var missing []string
				for path := range before {
					if !after[path] && !slices.Contains(undeclaredV4Paths, path) {
						missing = append(missing, path)
					}
				}
				sort.Strings(missing)
				Expect(missing).To(BeEmpty(), "fields GKO sent that the SDK payload drops")
			})
		}
	}
})

// legacyV4Payload is what GKO sent before the SDK: model.ToAPIV4DTO of the populated spec.
func legacyV4Payload(api *v1alpha1.ApiV4Definition) any {
	GinkgoHelper()
	cp := api.DeepCopy()
	cp.PopulateIDsWithContext(cp.HasContext(), k8s.IsAutomationAPIManaged(cp))
	payload := model.ToAPIV4DTO(&cp.Spec.Api)
	if cp.HasContext() && cp.Spec.ID != "" {
		payload.HRID = cp.Spec.ID
	}
	return payload
}

func sdkV4Payload(api *v1alpha1.ApiV4Definition) any {
	GinkgoHelper()
	dto, err := apiv4.ToDTO(api)
	Expect(err).ToNot(HaveOccurred())
	return dto.ApiV4
}

func apiV4Documents(file string) []*v1alpha1.ApiV4Definition {
	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	var apis []*v1alpha1.ApiV4Definition
	for _, doc := range bytes.Split(data, []byte("\n---")) {
		var meta struct {
			Kind string `json:"kind"`
		}
		if err := yaml.Unmarshal(doc, &meta); err != nil || meta.Kind != "ApiV4Definition" {
			continue
		}
		api := new(v1alpha1.ApiV4Definition)
		if err := yaml.Unmarshal(doc, api); err != nil {
			panic(fmt.Sprintf("%s: %v", file, err))
		}
		if api.Namespace == "" {
			api.Namespace = "default"
		}
		apis = append(apis, api)
	}
	return apis
}

// jsonPaths lists the paths of the non-null values of v marshalled to JSON, list indexes folded to [].
func payloadPaths(v any) map[string]bool {
	GinkgoHelper()
	data, err := json.Marshal(v)
	Expect(err).ToNot(HaveOccurred())
	var decoded any
	Expect(json.Unmarshal(data, &decoded)).To(Succeed())
	paths := map[string]bool{}
	collectPayloadPaths("", decoded, paths)
	return paths
}

func collectPayloadPaths(prefix string, v any, paths map[string]bool) {
	switch value := v.(type) {
	case map[string]any:
		for key, child := range value {
			if isEmptyPayloadValue(child) {
				continue
			}
			path := strings.TrimPrefix(prefix+"."+key, ".")
			paths[path] = true
			collectPayloadPaths(path, child, paths)
		}
	case []any:
		for _, child := range value {
			collectPayloadPaths(prefix+"[]", child, paths)
		}
	}
}

// isEmptyPayloadValue is a null, or an empty list or object. The SDK omits the empty collections GKO sent as
// [] (plan tags and characteristics, selector methods); the integration suites cover that APIM reads them alike.
func isEmptyPayloadValue(v any) bool {
	switch value := v.(type) {
	case nil:
		return true
	case []any:
		return len(value) == 0
	case map[string]any:
		return len(value) == 0
	}
	return false
}
