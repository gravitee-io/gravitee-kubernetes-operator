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

package hrid

import (
	"fmt"
	"regexp"
	"slices"
)

// MaxLength is the longest HRID the automation API accepts.
const MaxLength = 256

// pattern is the automation API's HRID grammar.
var pattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]+[a-zA-Z0-9]$`)

var forbiddenRunes = []rune{' ', '.'}

// NameToValidHRID remove spaces and dots only possible char for a map key that is not allowed as an HRID.
// Pages and Plan are maps, the key map is used to identify the plan.
// But it can contain spaces and dots, which are not allowed in an HRID.
func NameToValidHRID(s string) string {
	runes := []rune(s)
	for i, r := range runes {
		if slices.Contains(forbiddenRunes, r) {
			runes[i] = '-'
		}
	}
	return string(runes)
}

// NameToValidHRIDPointer same as NameToValidHRID but with pointers.
func NameToValidHRIDPointer(s *string) *string {
	if s == nil {
		return nil
	}
	return new(NameToValidHRID(*s))
}

// ToHRID convert a namespaced name to a valid HRID.
func ToHRID(ns, name string) string {
	return ns + "-" + name
}

// Validate checks a derived HRID against the automation API grammar before it is sent, so that a
// legal Kubernetes name the platform cannot address (a dotted name, a very long namespace and
// name) fails with a message naming the value.
func Validate(hrid string) error {
	if len(hrid) > MaxLength {
		return fmt.Errorf(
			"hrid [%s] derived from the resource namespace and name is %d characters long, the limit is %d",
			hrid, len(hrid), MaxLength,
		)
	}
	if !pattern.MatchString(hrid) {
		return fmt.Errorf(
			"hrid [%s] derived from the resource namespace and name must match %s: rename the resource",
			hrid, pattern.String(),
		)
	}
	return nil
}
