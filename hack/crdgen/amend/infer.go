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
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

const marker = "+kubebuilder:validation:"

type rule struct {
	name    string // test names and error messages
	applies func(f field) bool
	action  func(f field) []string // lines appended to the description
}

var rules = []rule{
	{"required", isRequired, requiredMarker},
	{"optional", isOptional, optionalMarker},
	{"minLength", hasMinLength, minLengthMarker},
	{"maxLength", hasMaxLength, maxLengthMarker},
	{"pattern", hasPattern, patternMarker},
	{"minimum", hasMinimum, minimumMarker},
	{"maximum", hasMaximum, maximumMarker},
	{"enum", hasEnum, enumMarker},
	{"uniqueItems", hasUniqueItems, listTypeSetMarker},
	{"minItems", hasMinItems, minItemsMarker},
	{"maxItems", hasMaxItems, maxItemsMarker},
	{"date-time", isDateTime, dateTimeFormatMarker},
	{"default", hasDefault, defaultsToProse},
}

// Infer appends the markers the OAS implies to the description of every field.
func Infer(doc *openapi3.T) {
	for _, f := range fields(doc) {
		for _, r := range rules {
			if r.applies(f) {
				appendDescription(f.schema, r.action(f)...)
			}
		}
	}
}

func appendDescription(schema *openapi3.Schema, lines ...string) {
	if schema.Description != "" {
		lines = append([]string{schema.Description}, lines...)
	}
	schema.Description = strings.Join(lines, "\n")
}

// inline is false on a $ref property: OAS 3.0 ignores the siblings of a $ref.
func inline(f field) bool {
	return f.ref == ""
}

// typed is false when x-go-type replaces the OAS type, whose constraints no longer apply.
func typed(f field) bool {
	return inline(f) && !f.hasGoType()
}

func isRequired(f field) bool   { return inline(f) && f.isRequired() }
func isOptional(f field) bool   { return inline(f) && f.parent != nil && !f.isRequired() }
func hasMinLength(f field) bool { return typed(f) && f.schema.MinLength > 0 }
func hasMaxLength(f field) bool {
	return typed(f) && f.schema.MaxLength != nil && *f.schema.MaxLength < math.MaxInt32
}
func hasPattern(f field) bool     { return typed(f) && f.schema.Pattern != "" }
func hasMinimum(f field) bool     { return typed(f) && f.schema.Min != nil }
func hasMaximum(f field) bool     { return typed(f) && f.schema.Max != nil }
func hasEnum(f field) bool        { return typed(f) && len(f.schema.Enum) > 0 }
func hasUniqueItems(f field) bool { return inline(f) && f.schema.UniqueItems }
func hasMinItems(f field) bool    { return inline(f) && f.schema.MinItems > 0 }
func hasMaxItems(f field) bool    { return inline(f) && f.schema.MaxItems != nil }
func isDateTime(f field) bool     { return typed(f) && f.schema.Format == "date-time" }
func hasDefault(f field) bool     { return inline(f) && f.schema.Default != nil }

func requiredMarker(field) []string { return []string{marker + "Required"} }

func optionalMarker(field) []string { return []string{marker + "Optional"} }

func minLengthMarker(f field) []string {
	return []string{marker + "MinLength=" + strconv.FormatUint(f.schema.MinLength, 10)}
}

func maxLengthMarker(f field) []string {
	return []string{marker + "MaxLength=" + strconv.FormatUint(*f.schema.MaxLength, 10)}
}

func patternMarker(f field) []string {
	return []string{marker + "Pattern=`" + f.schema.Pattern + "`"}
}

func minimumMarker(f field) []string {
	return []string{marker + "Minimum=" + formatNumber(*f.schema.Min)}
}

func maximumMarker(f field) []string {
	return []string{marker + "Maximum=" + formatNumber(*f.schema.Max)}
}

func enumMarker(f field) []string {
	values := make([]string, len(f.schema.Enum))
	for i, v := range f.schema.Enum {
		values[i] = fmt.Sprint(v)
	}
	return []string{marker + "Enum=" + strings.Join(values, ";")}
}

func listTypeSetMarker(field) []string { return []string{"+listType=set"} }

func minItemsMarker(f field) []string {
	return []string{marker + "MinItems=" + strconv.FormatUint(f.schema.MinItems, 10)}
}

func maxItemsMarker(f field) []string {
	return []string{marker + "MaxItems=" + strconv.FormatUint(*f.schema.MaxItems, 10)}
}

func dateTimeFormatMarker(field) []string { return []string{marker + "Format=date-time"} }

// defaultsToProse documents the OAS default without enforcing it: AM owns defaults,
// and a CRD default would freeze today's value into every stored CR.
func defaultsToProse(f field) []string {
	return []string{"Defaults to " + formatValue(f.schema.Default) + "."}
}

func formatNumber(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

// formatValue renders a string as is and anything else as JSON.
func formatValue(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprint(v)
	}
	return string(data)
}
