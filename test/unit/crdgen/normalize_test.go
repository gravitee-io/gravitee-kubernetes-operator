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
	"maps"
	"slices"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/hack/crdgen/amend"
)

var _ = Describe("Normalize", func() {
	It("drops readOnly properties from properties and required", func() {
		doc := spec(`
Thing:
  type: object
  required: [key, createdAt]
  properties:
    key: {type: string}
    createdAt: {type: string, readOnly: true}
`)
		amend.DropReadOnly(doc)

		t := doc.Components.Schemas["Thing"].Value
		Expect(slices.Sorted(maps.Keys(t.Properties))).To(Equal([]string{"key"}))
		Expect(t.Required).To(Equal([]string{"key"}))
	})

	DescribeTable("RejectUnsupported",
		func(schemas, wantErr string) {
			err := amend.RejectUnsupported(spec(schemas))
			if wantErr == "" {
				Expect(err).ToNot(HaveOccurred())
			} else {
				Expect(err).To(MatchError(ContainSubstring(wantErr)))
			}
		},
		Entry("rejects a oneOf property", `
Thing:
  type: object
  properties:
    value:
      oneOf: [{type: string}, {type: integer}]
`, "#/components/schemas/Thing/properties/value: oneOf"),
		Entry("rejects an anyOf schema", `
Thing:
  anyOf: [{type: string}, {type: integer}]
`, "#/components/schemas/Thing: anyOf"),
		Entry("rejects a required $ref property", `
Other: {type: object}
Thing:
  type: object
  required: [other]
  properties:
    other: {$ref: '#/components/schemas/Other'}
`, "#/components/schemas/Thing/properties/other: a required $ref"),
		Entry("accepts a oneOf property with x-go-type", `
Thing:
  type: object
  properties:
    value:
      x-go-type: any
      oneOf: [{type: string}, {type: integer}]
`, ""),
		Entry("accepts an optional $ref property", `
Other: {type: object}
Thing:
  type: object
  properties:
    other: {$ref: '#/components/schemas/Other'}
`, ""),
	)

	It("keeps inline string enums plain strings, not named enums", func() {
		doc := spec(`
Level: {type: string, enum: [A, B]}
Thing:
  type: object
  properties:
    inline: {type: string, enum: [A, B]}
    plain: {type: string}
    level: {$ref: '#/components/schemas/Level'}
`)
		amend.FlattenInlineEnums(doc)

		props := doc.Components.Schemas["Thing"].Value.Properties
		Expect(props["inline"].Value.Extensions).To(HaveKeyWithValue("x-go-type", "string"))
		Expect(props["plain"].Value.Extensions).ToNot(HaveKey("x-go-type"))
		Expect(doc.Components.Schemas["Level"].Value.Extensions).ToNot(HaveKey("x-go-type"))
	})
})
