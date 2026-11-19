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
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	structuralschema "k8s.io/apiextensions-apiserver/pkg/apiserver/schema"
	"k8s.io/apiextensions-apiserver/pkg/apiserver/schema/defaulting"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/yaml"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/apidefinition"
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
	schema := apiV4StructuralSchema()
	referenced := referencedExamples()

	for _, file := range files {
		for i, api := range apiV4Documents(file, schema) {
			name := fmt.Sprintf("%s #%d", filepath.Base(file), i)

			It("sends every declared field of "+name, func() {
				api := prepareExample(api, referenced)
				before := payloadPaths(legacyV4Payload(api))
				after := payloadPaths(sdkV4Payload(api))

				var missing []string
				for path := range before {
					if _, ok := after[path]; !ok && !slices.Contains(undeclaredV4Paths, path) {
						missing = append(missing, path)
					}
				}
				sort.Strings(missing)
				Expect(missing).To(BeEmpty(), "fields GKO sent that the SDK payload drops")
			})

			It("sends nothing GKO did not send, and the same values, for "+name, func() {
				api := prepareExample(api, referenced)
				before := payloadPaths(legacyV4Payload(api))
				after := payloadPaths(sdkV4Payload(api))

				var differing []string
				for path, values := range after {
					previous, ok := before[path]
					switch {
					case !ok:
						differing = append(differing, path+": not sent by GKO, sent as "+strings.Join(values, ", "))
					case !slices.Equal(values, previous):
						differing = append(differing,
							path+": GKO sent "+strings.Join(previous, ", ")+", the SDK sends "+strings.Join(values, ", "))
					}
				}
				sort.Strings(differing)
				Expect(differing).To(BeEmpty(), "fields the SDK payload adds or changes")
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

// referencedExamples are the example ApiResources, SharedPolicyGroups, Groups and Notifications the v4 examples
// reference.
func referencedExamples() []client.Object {
	var objects []client.Object
	for _, pattern := range []string{
		"../../../examples/apim/api_resource/*.yml",
		"../../../examples/apim/shared_policy_groups/*.yml",
		"../../../examples/apim/group/group.yml",
		"../../../examples/apim/notification/*.yml",
	} {
		files, err := filepath.Glob(pattern)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			objects = append(objects, exampleObjects(file)...)
		}
	}
	return objects
}

func exampleObjects(file string) []client.Object {
	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	var objects []client.Object
	for _, doc := range bytes.Split(data, []byte("\n---")) {
		var meta struct {
			Kind string `json:"kind"`
		}
		if err := yaml.Unmarshal(doc, &meta); err != nil {
			continue
		}
		var obj client.Object
		switch meta.Kind {
		case "ApiResource":
			obj = new(v1alpha1.ApiResource)
		case "SharedPolicyGroup":
			obj = new(v1alpha1.SharedPolicyGroup)
		case "Group":
			obj = new(v1alpha1.Group)
		case "Notification":
			obj = new(v1alpha1.Notification)
		default:
			continue
		}
		if err := yaml.Unmarshal(doc, obj); err != nil {
			panic(fmt.Sprintf("%s: %v", file, err))
		}
		if obj.GetNamespace() == "" {
			obj.SetNamespace("default")
		}
		objects = append(objects, obj)
	}
	return objects
}

// prepareExample resolves the references of an example as the controller does before it maps the API, against a
// cluster holding the referenced examples.
func prepareExample(api *v1alpha1.ApiV4Definition, referenced []client.Object) *v1alpha1.ApiV4Definition {
	GinkgoHelper()
	scheme := runtime.NewScheme()
	Expect(clientgoscheme.AddToScheme(scheme)).To(Succeed())
	Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
	objects := make([]client.Object, len(referenced))
	for i, obj := range referenced {
		objects[i] = obj.DeepCopyObject().(client.Object)
	}
	k8s.RegisterClient(fake.NewClientBuilder().WithScheme(scheme).WithObjects(objects...).Build())

	cp := api.DeepCopy()
	Expect(apidefinition.PrepareV4SpecForAutomation(context.Background(), cp, false)).To(Succeed())
	return cp
}

// apiV4StructuralSchema is the v1alpha1 schema of the ApiV4Definition CRD, used to default an example as the API
// server does before the operator reads it.
func apiV4StructuralSchema() *structuralschema.Structural {
	data, err := os.ReadFile("../../../crds/gravitee.io/gravitee.io_apiv4definitions.yaml")
	if err != nil {
		panic(err)
	}
	crd := new(apiextensionsv1.CustomResourceDefinition)
	if err := yaml.Unmarshal(data, crd); err != nil {
		panic(err)
	}
	for _, version := range crd.Spec.Versions {
		if version.Name != "v1alpha1" {
			continue
		}
		props := new(apiextensions.JSONSchemaProps)
		if err := apiextensionsv1.Convert_v1_JSONSchemaProps_To_apiextensions_JSONSchemaProps(
			version.Schema.OpenAPIV3Schema, props, nil,
		); err != nil {
			panic(err)
		}
		schema, err := structuralschema.NewStructural(props)
		if err != nil {
			panic(err)
		}
		return schema
	}
	panic("no v1alpha1 version in the ApiV4Definition CRD")
}

// apiV4Documents reads the ApiV4Definitions of an example file, with the CRD defaults applied.
func apiV4Documents(file string, schema *structuralschema.Structural) []*v1alpha1.ApiV4Definition {
	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	var apis []*v1alpha1.ApiV4Definition
	for _, doc := range bytes.Split(data, []byte("\n---")) {
		var obj map[string]any
		if err := yaml.Unmarshal(doc, &obj); err != nil || obj["kind"] != "ApiV4Definition" {
			continue
		}
		defaulting.Default(obj, schema)
		defaulted, err := json.Marshal(obj)
		if err != nil {
			panic(fmt.Sprintf("%s: %v", file, err))
		}
		api := new(v1alpha1.ApiV4Definition)
		if err := json.Unmarshal(defaulted, api); err != nil {
			panic(fmt.Sprintf("%s: %v", file, err))
		}
		if api.Namespace == "" {
			api.Namespace = "default"
		}
		apis = append(apis, api)
	}
	return apis
}

// payloadPaths lists the paths of the non-empty values of v marshalled to JSON, list indexes folded to [], each
// with the sorted JSON of the scalar values found there (none for an object or a list).
func payloadPaths(v any) map[string][]string {
	GinkgoHelper()
	data, err := json.Marshal(v)
	Expect(err).ToNot(HaveOccurred())
	var decoded any
	Expect(json.Unmarshal(data, &decoded)).To(Succeed())
	paths := map[string][]string{}
	collectPayloadPaths("", decoded, paths)
	for _, values := range paths {
		sort.Strings(values)
	}
	return paths
}

func collectPayloadPaths(prefix string, v any, paths map[string][]string) {
	switch value := v.(type) {
	case map[string]any:
		for key, child := range value {
			if isEmptyPayloadValue(child) {
				continue
			}
			path := strings.TrimPrefix(prefix+"."+key, ".")
			if _, ok := paths[path]; !ok {
				paths[path] = []string{}
			}
			collectPayloadPaths(path, child, paths)
		}
	case []any:
		for _, child := range value {
			collectPayloadPaths(prefix+"[]", child, paths)
		}
	default:
		scalar, _ := json.Marshal(value)
		paths[prefix] = append(paths[prefix], string(scalar))
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
