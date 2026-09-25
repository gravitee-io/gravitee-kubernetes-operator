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
	"encoding/json"
	"reflect"
	"slices"
	"strings"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Spec fields the Automation API never receives under their CRD name.
var resolvedBeforeSending = []string{
	// Replaced by consoleNotification, groups and inlined resources before sending.
	"consoleNotification",
	"groupRefs",
	"notificationsRefs",
	"resources[].ref",
	// Resolved into the step's shared-policy-group-policy configuration.
	"flows[].connect[].sharedPolicyGroupRef",
	"flows[].interact[].sharedPolicyGroupRef",
	"flows[].publish[].sharedPolicyGroupRef",
	"flows[].request[].sharedPolicyGroupRef",
	"flows[].response[].sharedPolicyGroupRef",
	"flows[].subscribe[].sharedPolicyGroupRef",
	"plans[].flows[].connect[].sharedPolicyGroupRef",
	"plans[].flows[].interact[].sharedPolicyGroupRef",
	"plans[].flows[].publish[].sharedPolicyGroupRef",
	"plans[].flows[].request[].sharedPolicyGroupRef",
	"plans[].flows[].response[].sharedPolicyGroupRef",
	"plans[].flows[].subscribe[].sharedPolicyGroupRef",
	// Replaced by parentHrid and generalConditionsHrid.
	"pages[].parent",
	"plans[].generalConditions",
	// Not an Automation API field.
	"pages[].parentId",
}

var _ = Describe("API v4 payload completeness", func() {
	It("sends every API v4 spec field to APIM", func() {
		api := &v4.Api{}
		fillAllFields(reflect.ValueOf(api).Elem(), 0)

		specJSON, err := json.Marshal(api)
		Expect(err).ToNot(HaveOccurred())
		payloadJSON, err := json.Marshal(model.ToAPIV4DTO(api))
		Expect(err).ToNot(HaveOccurred())

		sent := jsonPaths(payloadJSON)
		var dropped []string
		for _, path := range jsonPaths(specJSON) {
			if !slices.Contains(sent, path) && !isResolvedBeforeSending(path) {
				dropped = append(dropped, path)
			}
		}
		Expect(dropped).To(BeEmpty(), "spec fields missing from the APIM payload")
	})
})

func isResolvedBeforeSending(path string) bool {
	return slices.ContainsFunc(resolvedBeforeSending, func(prefix string) bool {
		return path == prefix || strings.HasPrefix(path, prefix+".") || strings.HasPrefix(path, prefix+"[]")
	})
}

const fillMaxDepth = 12

var jsonUnmarshaler = reflect.TypeFor[json.Unmarshaler]()

// fillAllFields gives every reachable field a non-zero value. Types with their
// own JSON decoding (listeners, selectors, configuration blobs) are opaque
// pass-through payloads and are left unset.
func fillAllFields(v reflect.Value, depth int) {
	if depth > fillMaxDepth || !v.CanSet() || reflect.PointerTo(v.Type()).Implements(jsonUnmarshaler) {
		return
	}
	switch v.Kind() {
	case reflect.Pointer:
		if v.Type().Elem().Implements(jsonUnmarshaler) || reflect.PointerTo(v.Type().Elem()).Implements(jsonUnmarshaler) {
			return
		}
		v.Set(reflect.New(v.Type().Elem()))
		fillAllFields(v.Elem(), depth+1)
	case reflect.Struct:
		for i := range v.NumField() {
			if v.Type().Field(i).IsExported() {
				fillAllFields(v.Field(i), depth+1)
			}
		}
	case reflect.Slice:
		v.Set(reflect.MakeSlice(v.Type(), 1, 1))
		fillAllFields(v.Index(0), depth+1)
	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			return
		}
		value := reflect.New(v.Type().Elem()).Elem()
		fillAllFields(value, depth+1)
		v.Set(reflect.MakeMap(v.Type()))
		v.SetMapIndex(reflect.ValueOf("key").Convert(v.Type().Key()), value)
	case reflect.String:
		v.SetString("value")
	case reflect.Bool:
		v.SetBool(true)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(1)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.SetUint(1)
	case reflect.Float32, reflect.Float64:
		v.SetFloat(1)
	default:
	}
}

// jsonPaths lists the leaf paths of a JSON document. Array items collapse to
// `[]`, and the spec's plan and page maps are read as the lists the payload sends.
func jsonPaths(document []byte) []string {
	var root any
	Expect(json.Unmarshal(document, &root)).To(Succeed())
	var paths []string
	var walk func(node any, path string)
	walk = func(node any, path string) {
		switch typed := node.(type) {
		case map[string]any:
			for key, child := range typed {
				childPath := key
				if path != "" {
					childPath = path + "." + key
				}
				if (path == "plans" || path == "pages") && key == "key" {
					childPath = path + "[]"
				}
				walk(child, childPath)
			}
		case []any:
			for _, child := range typed {
				walk(child, path+"[]")
			}
		default:
			paths = append(paths, path)
		}
	}
	walk(root, "")
	slices.Sort(paths)
	return slices.Compact(paths)
}
