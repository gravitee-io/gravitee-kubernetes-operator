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

package crdgen_test

import (
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCRDGen(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "crdgen unit tests suite")
}

// spec loads a minimal OAS whose components.schemas are the given YAML.
func spec(schemas string) *openapi3.T {
	indented := "    " + strings.ReplaceAll(strings.TrimSpace(schemas), "\n", "\n    ") + "\n"
	data := "openapi: 3.0.1\ninfo: {title: t, version: v}\npaths: {}\ncomponents:\n  schemas:\n" + indented
	doc, err := openapi3.NewLoader().LoadFromData([]byte(data))
	Expect(err).ToNot(HaveOccurred())
	return doc
}

// thing loads a spec with one Thing schema holding one "value" property, required or not.
func thing(property string, required bool) *openapi3.T {
	req := "[]"
	if required {
		req = "[value]"
	}
	return spec("Thing:\n  type: object\n  required: " + req + "\n  properties:\n    value: " + property)
}

func value(doc *openapi3.T) *openapi3.Schema {
	return doc.Components.Schemas["Thing"].Value.Properties["value"].Value
}
