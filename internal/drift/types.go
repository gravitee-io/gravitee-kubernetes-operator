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
	"cmp"
	"slices"
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

// Merge merges Result with one another.
// The first argument is the Result of detecting drift between "Old CRD" (in etcd) vs. Remote (O/R).
// The second argument is the Result of detecting drift between "New CRD" (applied) vs. Remote (N/R).
// O/R N/R have the same structure as they ran on the same struct
// 5 cases can happen:
// Case 1: Nothing changed  => Ok
// Case 2: CRD change is equivalent to remote change (O/R diff, N/R equivalent) => Ok
// Case 3: Only remote changed (O/R diff, N/R diff) => Drift
// Case 4: Only CRD changed (O/R equivalent, N/R diff) => Ok
// Case 5: CRD change but is not equivalent to remote change (O/R diff, N/R diff)  => Drift
// As a result: when both O/R and N/R are drifted, the result is a drift.
func Merge(or Result, nr Result) Result {
	m := Result{}
	merge(or, nr, &m)
	return m
}

func merge(or Result, nr Result, merger *Result) {
	result := evaluateDrift(or, nr)
	merger.Equivalence = result.Equivalence
	merger.Property = result.Property
	merger.Index = result.Index
	merger.CRDValue = result.CRDValue
	merger.APIValue = result.APIValue
	if len(or.Children) > 0 {
		merger.Children = make([]*Result, len(or.Children))
		for i, orChild := range or.Children {
			nrChild := nr.Children[i]
			merger.Children[i] = &Result{}
			merge(*orChild, *nrChild, merger.Children[i])
		}
	}
}

func evaluateDrift(or Result, nr Result) Result {
	// Case 1, both are Equivalent => Equivalent
	// Case 3 & 5, both are Inequivalent => Inequivalent
	if or.Equivalent == nr.Equivalent {
		return nr
	}
	// Rest is no drift, so we return the result without drift to have a clean result (reason could be filled)
	// Case 2
	if nr.Equivalent == Equivalent {
		return nr
	}
	// Case 4
	if or.Equivalent == Equivalent {
		return or
	}
	// should not happen but just in case
	return Result{}
}

// AppendChild adds a child to the result.
func (r *Result) AppendChild(child *Result, ordered bool) *Result {
	r.Children = append(r.Children, child)
	if ordered {
		r.Children = slices.SortedFunc(slices.Values(r.Children), func(e1 *Result, e2 *Result) int {
			return cmp.Compare(e1.Property, e2.Property)
		})
	}
	return child
}
