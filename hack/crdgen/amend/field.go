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

package amend

import (
	"maps"
	"slices"

	"github.com/getkin/kin-openapi/openapi3"
)

const goTypeExtension = "x-go-type"

// field is a property (parent set) or a named schema (parent nil).
type field struct {
	pointer string // JSON pointer, for error messages
	parent  *openapi3.Schema
	name    string
	ref     string // $ref of a property, empty when inline
	schema  *openapi3.Schema
	// extensions holds the x- keys written next to the field: on the schema
	// when inline, next to the $ref otherwise.
	extensions map[string]any
}

// fields lists every named schema followed by its properties, sorted by name.
func fields(doc *openapi3.T) []field {
	var all []field
	for _, name := range slices.Sorted(maps.Keys(doc.Components.Schemas)) {
		schema := doc.Components.Schemas[name].Value
		pointer := "#/components/schemas/" + name
		all = append(all, field{pointer: pointer, name: name, schema: schema, extensions: schema.Extensions})
		for _, prop := range slices.Sorted(maps.Keys(schema.Properties)) {
			all = append(all, property(pointer, schema, prop))
		}
	}
	return all
}

func property(parentPointer string, parent *openapi3.Schema, name string) field {
	ref := parent.Properties[name]
	f := field{
		pointer:    parentPointer + "/properties/" + name,
		parent:     parent,
		name:       name,
		ref:        ref.Ref,
		schema:     ref.Value,
		extensions: ref.Value.Extensions,
	}
	if ref.Ref != "" {
		f.extensions = ref.Extensions
	}
	return f
}

func (f field) isRequired() bool {
	return f.parent != nil && slices.Contains(f.parent.Required, f.name)
}

func (f field) hasGoType() bool {
	_, ok := f.extensions[goTypeExtension]
	return ok
}
