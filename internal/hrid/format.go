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
	"slices"
)

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
	hrid := NameToValidHRID(*s)
	return &hrid
}

// ToHRID convert a namespaced name to a valid HRID.
func ToHRID(ns, name string) string {
	return ns + "-" + name
}
