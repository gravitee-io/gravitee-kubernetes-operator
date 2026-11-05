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
	"github.com/getkin/kin-openapi/openapi3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"

	"github.com/gravitee-io/gravitee-kubernetes-operator/hack/crdgen/amend"
)

const renamedSpec = `
openapi: 3.0.1
info: {title: t, version: v}
paths:
  /things:
    put:
      requestBody:
        content:
          application/json:
            schema: {$ref: '#/components/schemas/AutomationThing'}
      responses:
        '200': {description: ok}
components:
  schemas:
    AutomationThing:
      x-gko-name: Thing
      type: object
      properties:
        settings: {$ref: '#/components/schemas/AutomationSettings'}
        tags:
          type: array
          items: {$ref: '#/components/schemas/Tag'}
    AutomationSettings:
      x-gko-name: Settings
      type: object
      properties:
        enabled: {type: boolean}
    Tag:
      type: string
`

func renamed(data string) (*openapi3.T, error) {
	GinkgoHelper()
	var root yaml.Node
	Expect(yaml.Unmarshal([]byte(data), &root)).To(Succeed())
	if err := amend.RenameSchemas(&root); err != nil {
		return nil, err
	}
	out, err := yaml.Marshal(&root)
	Expect(err).ToNot(HaveOccurred())
	doc, err := openapi3.NewLoader().LoadFromData(out)
	Expect(err).ToNot(HaveOccurred())
	return doc, nil
}

var _ = Describe("RenameSchemas", func() {
	It("renames the schemas and every $ref to them, and drops x-gko-name", func() {
		doc, err := renamed(renamedSpec)
		Expect(err).ToNot(HaveOccurred())

		Expect(doc.Components.Schemas).To(HaveKey("Thing"))
		Expect(doc.Components.Schemas).To(HaveKey("Settings"))
		Expect(doc.Components.Schemas).To(HaveKey("Tag"))
		Expect(doc.Components.Schemas).ToNot(HaveKey("AutomationThing"))
		Expect(doc.Components.Schemas).ToNot(HaveKey("AutomationSettings"))

		thing := doc.Components.Schemas["Thing"].Value
		Expect(thing.Extensions).ToNot(HaveKey("x-gko-name"))
		Expect(thing.Properties["settings"].Ref).To(Equal("#/components/schemas/Settings"))
		Expect(thing.Properties["tags"].Value.Items.Ref).To(Equal("#/components/schemas/Tag"))

		body := doc.Paths.Find("/things").Put.RequestBody.Value.Content["application/json"].Schema
		Expect(body.Ref).To(Equal("#/components/schemas/Thing"))
	})

	It("fails when a new name is already taken", func() {
		_, err := renamed(renamedSpec + "    Thing: {type: string}\n")

		Expect(err).To(MatchError(ContainSubstring("Thing: two schemas have this name")))
	})

	It("fails on a name that is not a string", func() {
		_, err := renamed(`
openapi: 3.0.1
info: {title: t, version: v}
paths: {}
components:
  schemas:
    AutomationThing: {type: object, x-gko-name: [Thing]}
`)

		Expect(err).To(MatchError(ContainSubstring("AutomationThing: x-gko-name: want a schema name")))
	})
})
