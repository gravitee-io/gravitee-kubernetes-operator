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

package drift

import (
	"strings"
)

// EquivalentStatus represents the status of the equivalence.
type EquivalentStatus byte

const (
	CannotCompare EquivalentStatus = iota
	Equivalent    EquivalentStatus = 1
	Inequivalent  EquivalentStatus = 2
)

// Equivalence represents the equivalence between a CRD and an API.
type Equivalence struct {
	Equivalent EquivalentStatus
	Reason     any
	Skip       bool
}

// EquivalenceFunc is a function that compares two values and returns an Equivalence.
type EquivalenceFunc func(crd any, api any) Equivalence

// Result represents the result of the drift detection.
// It contains the equivalence status, the property name, optionally the index of the property,
// the CRD value, the API value, and the children results.
type Result struct {
	Equivalence
	Property string
	Index    *int
	CRDValue any
	APIValue any
	Children []*Result
}

// String returns a string representation of the result as a pseudo-yaml tree.
func (r *Result) String() string {
	var builder strings.Builder
	format(r, &builder, -2)
	return strings.TrimSpace(strings.TrimRight(builder.String(), "\n"))
}


// DriftDetected returns true if the result is equivalent and all of its children are equivalent.
func (r *Result) DriftDetected() bool {
	if r.Equivalent == Inequivalent {
		return true
	}
	if len(r.Children) > 0 {
		for _, child := range r.Children {
			if child.DriftDetected() {
				return true
			}
		}
	}
	return false
}

// AppendChild adds a child to the result.
func (r *Result) AppendChild(child *Result) *Result {
	r.Children = append(r.Children, child)
	return child
}
