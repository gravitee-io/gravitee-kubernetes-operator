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

// Package amend holds the crdgen steps that turn an OAS into one oapi-codegen generates CRD-ready structs
// from: Normalize, Infer, ApplyExtensions, then FlattenInlineEnums.
package amend

import (
	"errors"
	"fmt"
	"slices"

	"github.com/getkin/kin-openapi/openapi3"
)

// Normalize drops readOnly properties and examples, and rejects what the CRD cannot express.
func Normalize(doc *openapi3.T) error {
	DropReadOnly(doc)
	dropExamples(doc)
	return RejectUnsupported(doc)
}

// DropReadOnly removes server-managed properties: a CRD spec never carries them.
func DropReadOnly(doc *openapi3.T) {
	for _, schema := range doc.Components.Schemas {
		for name, prop := range schema.Value.Properties {
			if prop.Value.ReadOnly {
				delete(schema.Value.Properties, name)
				schema.Value.Required = slices.DeleteFunc(schema.Value.Required, func(r string) bool { return r == name })
			}
		}
	}
}

// dropExamples removes examples from schemas so they don't make it to the generated CRD.
func dropExamples(doc *openapi3.T) {
	for _, schema := range doc.Components.Schemas {
		schema.Value.Example = nil
		for _, prop := range schema.Value.Properties {
			prop.Value.Example = nil
		}
	}
}

// RejectUnsupported fails on constructs the generated CRD cannot express.
// A required $ref property is one of them: OAS 3.0 ignores its description,
// so its Required marker cannot be written and the package default would make it optional.
func RejectUnsupported(doc *openapi3.T) error {
	var errs []error
	for _, f := range fields(doc) {
		if f.hasGoType() {
			continue
		}
		switch {
		case len(f.schema.OneOf) > 0 && f.ref == "":
			errs = append(errs, fmt.Errorf("%s: oneOf is not supported", f.pointer))
		case len(f.schema.AnyOf) > 0 && f.ref == "":
			errs = append(errs, fmt.Errorf("%s: anyOf is not supported", f.pointer))
		case f.ref != "" && f.isRequired():
			errs = append(errs, fmt.Errorf("%s: a required $ref property is not supported", f.pointer))
		}
	}
	return errors.Join(errs...)
}

// FlattenInlineEnums keeps an inline string enum a plain string. oapi-codegen would otherwise
// hoist it into a named type carrying the property's description, so its Enum marker would
// apply twice. Runs after infer: the marker is already written and x-go-type now hides the enum.
func FlattenInlineEnums(doc *openapi3.T) {
	for _, f := range fields(doc) {
		if f.parent == nil || !typed(f) || len(f.schema.Enum) == 0 || !f.schema.Type.Is(openapi3.TypeString) {
			continue
		}
		if f.schema.Extensions == nil {
			f.schema.Extensions = map[string]any{}
		}
		f.schema.Extensions[goTypeExtension] = "string"
	}
}
