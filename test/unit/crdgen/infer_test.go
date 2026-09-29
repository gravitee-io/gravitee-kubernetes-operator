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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/hack/crdgen/amend"
)

var _ = Describe("Infer", func() {
	DescribeTable("appends the markers the OAS implies to the description",
		func(property string, required bool, want string) {
			doc := thing(property, required)
			amend.Infer(doc)
			Expect(value(doc).Description).To(Equal(want))
		},
		Entry("required", "{type: string, description: Doc.}", true, "Doc.\n+kubebuilder:validation:Required"),
		Entry("optional", "{type: string, description: Doc.}", false, "Doc.\n+kubebuilder:validation:Optional"),
		Entry("minLength", "{type: string, minLength: 1}", false,
			"+kubebuilder:validation:Optional\n+kubebuilder:validation:MinLength=1"),
		Entry("minLength 0", "{type: string, minLength: 0}", false, "+kubebuilder:validation:Optional"),
		Entry("maxLength", "{type: string, maxLength: 255}", false,
			"+kubebuilder:validation:Optional\n+kubebuilder:validation:MaxLength=255"),
		Entry("maxLength unbounded", "{type: string, maxLength: 2147483647}", false, "+kubebuilder:validation:Optional"),
		Entry("pattern", "{type: string, pattern: '^/.*'}", false,
			"+kubebuilder:validation:Optional\n+kubebuilder:validation:Pattern=`^/.*`"),
		Entry("minimum", "{type: integer, minimum: 1}", false,
			"+kubebuilder:validation:Optional\n+kubebuilder:validation:Minimum=1"),
		Entry("maximum", "{type: integer, maximum: 100}", false,
			"+kubebuilder:validation:Optional\n+kubebuilder:validation:Maximum=100"),
		Entry("enum", "{type: string, enum: [A, B]}", false,
			"+kubebuilder:validation:Optional\n+kubebuilder:validation:Enum=A;B"),
		Entry("uniqueItems", "{type: array, uniqueItems: true, items: {type: string}}", false,
			"+kubebuilder:validation:Optional\n+listType=set"),
		Entry("minItems", "{type: array, minItems: 1, items: {type: string}}", false,
			"+kubebuilder:validation:Optional\n+kubebuilder:validation:MinItems=1"),
		Entry("maxItems", "{type: array, maxItems: 3, items: {type: string}}", false,
			"+kubebuilder:validation:Optional\n+kubebuilder:validation:MaxItems=3"),
		Entry("date-time", "{type: string, format: date-time}", false,
			"+kubebuilder:validation:Optional\n+kubebuilder:validation:Format=date-time"),
		Entry("default bool", "{type: boolean, default: false}", false,
			"+kubebuilder:validation:Optional\nDefaults to false."),
		Entry("default string", "{type: string, default: AUTO}", false,
			"+kubebuilder:validation:Optional\nDefaults to AUTO."),
		Entry("x-go-type skips type-based rules",
			"{type: string, x-go-type: utils.GenericStringMap, maxLength: 10, pattern: a, enum: [a]}", false,
			"+kubebuilder:validation:Optional"),
		Entry("markers in rule order", "{type: string, minLength: 1, maxLength: 255, default: a}", true,
			"+kubebuilder:validation:Required\n+kubebuilder:validation:MinLength=1\n"+
				"+kubebuilder:validation:MaxLength=255\nDefaults to a."),
	)

	It("leaves the target of a $ref property unchanged", func() {
		doc := spec(`
Other: {type: object, description: Other.}
Thing:
  type: object
  properties:
    other: {$ref: '#/components/schemas/Other'}
`)
		amend.Infer(doc)

		Expect(doc.Components.Schemas["Other"].Value.Description).To(Equal("Other."))
	})

	It("writes the enum marker on a named enum schema", func() {
		doc := spec("Severity: {type: string, description: Level., enum: [ERROR, WARNING]}")
		amend.Infer(doc)

		Expect(doc.Components.Schemas["Severity"].Value.Description).
			To(Equal("Level.\n+kubebuilder:validation:Enum=ERROR;WARNING"))
	})
})
