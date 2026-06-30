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
	"log"
	"slices"
	"strconv"
	"strings"
)

// EquivalentStatus represents the status of the equivalence.
type EquivalentStatus byte

const (
	CannotCompare EquivalentStatus = iota
	Equivalent    EquivalentStatus = 1
	Inequivalent  EquivalentStatus = 2
)

func (s EquivalentStatus) String() string {
	switch s {
	case Inequivalent:
		return "Inequivalent"
	case Equivalent:
		return "Equivalent"
	case CannotCompare:
		fallthrough
	default:
		return "Cannot compare"
	}
}

// Equivalence represents the equivalence between a CRD and an API.
type Equivalence struct {
	// Equivalent is true if the CRD and the API are equivalent.
	Equivalent EquivalentStatus
	// Reason is the reason why the equivalence is not true.
	Reason any
	// Skip is true if the children should be ignored.
	Skip bool
	// PostFunc is a function called after the equivalence is evaluated on the result that has just been processed.
	PostFunc PostEquivalenceFunc
}

// EquivalenceFunc is a function that compares two values and returns an Equivalence.
type EquivalenceFunc func(crd any, remote any) Equivalence

// PostEquivalenceFunc is a function called after the equivalence is evaluated on the result that has just been processed.
type PostEquivalenceFunc func(r *Result)

// Result represents the result of the drift detection.
// It contains the equivalence status, the property name, optionally the index of the property,
// the CRD value, the API value, and the children results.
type Result struct {
	// Equivalence is the equivalence between the CRD and the API.
	Equivalence
	// Property is the name of the property that is compared or the key of the map entry.
	Property string
	// Index is the index of the item when the property is a slice.
	Index *int
	// CRDValue is actual the value of the CRD.
	CRDValue any
	// RemoteValue is actual the value of the Remote.
	RemoteValue any
	// Children is the list of children results (fields, slice items, map entries).
	Children []*Result
}

// String returns a string representation of the result as a pseudo-yaml tree.
func (r *Result) String() string {
	var builder strings.Builder
	format(r, &builder, -2)
	return strings.TrimSpace(strings.TrimRight(builder.String(), "\n"))
}

// DriftDetected returns true if the result or any of its children is found to be inequivalent, false otherwise.
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

// AppendChild adds a child to the result and returns it for convenience.
// If ordered is true, the children are sorted by property name to be able to compare maps.
func (r *Result) AppendChild(child *Result, ordered bool) *Result {
	r.Children = append(r.Children, child)
	if ordered {
		r.Children = slices.SortedFunc(slices.Values(r.Children), func(e1 *Result, e2 *Result) int {
			return cmp.Compare(e1.Property, e2.Property)
		})
	}
	return child
}

// Merge merges Result with one another.
// The first argument is the Result of detecting drift between "Old CRD" (in etcd) vs. Remote (O/R).
// The second argument is the Result of detecting drift between "New CRD" (applied) vs. Remote (N/R).
// 5 cases can happen:
// Case 1: Nothing changed => Ok
// Case 2: CRD change is equivalent to remote change (O/R diff, N/R equivalent) => Ok
// Case 3: Only remote changed (O/R diff, N/R diff) => Drift
// Case 4: Only CRD changed (O/R equivalent, N/R diff) => Ok
// Case 5: CRD change but is not equivalent to remote change (O/R diff, N/R diff)  => Drift
// As a result: when both O/R and N/R are drifted, the result is a drift.
func Merge(or Result, nr Result) Result {
	m := Result{}
	merge(or, nr, &m, nil)
	return m
}

func childKey(r *Result) string {
	if r.Index != nil {
		return r.Property + "_" + strconv.Itoa(*r.Index)
	}
	return r.Property
}

func merge(or Result, nr Result, merger *Result, parentIndex *int) {
	// evaluates if a drift occurred and return the result that will be used to display the drift if any.
	result := evaluateDrift(or, nr, parentIndex)
	// put it all in the merged result
	merger.Equivalence = result.Equivalence
	merger.Property = result.Property
	merger.Index = result.Index
	merger.CRDValue = result.CRDValue
	merger.RemoteValue = result.RemoteValue

	if childrenAlignAtIndex(or, nr) {
		mergeChildrenPositionally(or, nr, merger, merger.Index)
		return
	}
	mergeChildrenByKey(or, nr, merger, merger.Index)
}

func childrenAlignAtIndex(or Result, nr Result) bool {
	if len(or.Children) != len(nr.Children) {
		return false
	}
	for i := range or.Children {
		if childKey(or.Children[i]) != childKey(nr.Children[i]) {
			return false
		}
	}
	return true
}

func mergeChildrenPositionally(or Result, nr Result, merger *Result, parentIndex *int) {
	merger.Children = make([]*Result, len(or.Children))
	for i, orChild := range or.Children {
		merger.Children[i] = &Result{}
		nrChild := Result{}
		if i < len(nr.Children) {
			nrChild = *nr.Children[i]
		}
		merge(*orChild, nrChild, merger.Children[i], parentIndex)
	}
}

func mergeChildrenByKey(or Result, nr Result, merger *Result, parentIndex *int) {
	nrChildren := make(map[string]*Result, len(nr.Children))
	for _, child := range nr.Children {
		nrChildren[childKey(child)] = child
	}

	for _, orChild := range or.Children {
		key := childKey(orChild)
		mergedChild := &Result{}
		if nrChild, ok := nrChildren[key]; ok {
			merge(*orChild, *nrChild, mergedChild, parentIndex)
			delete(nrChildren, key)
		} else {
			merge(*orChild, Result{}, mergedChild, parentIndex)
		}
		merger.Children = append(merger.Children, mergedChild)
	}

	for _, nrChild := range nrChildren {
		mergedChild := &Result{}
		merge(Result{}, *nrChild, mergedChild, parentIndex)
		merger.Children = append(merger.Children, mergedChild)
	}
}

func evaluateDrift(or Result, nr Result, parentIndex *int) Result {
	// Case 1, both are Equivalent => Equivalent
	// Case 3 & 5, both are Inequivalent => Inequivalent
	if or.Equivalent == nr.Equivalent {
		if parentIndex == nil && or.Equivalent == Inequivalent && remoteOnlyDriftUnchanged(or, nr) {
			return Result{Equivalence: Equivalence{Equivalent: Equivalent}}
		}
		return nr
	}
	// Rest is no drift, so we return the result without drift to have a clean result (reason could be filled)
	// Case 2
	if nr.Equivalent == Equivalent || nr.Equivalent == CannotCompare {
		return nr
	}
	// Case 4
	if or.Equivalent == Equivalent || or.Equivalent == CannotCompare {
		return or
	}

	log.Panicf("Unable to evaluate drift between (Old vs. Remote) %v and %v (New vs. Remote)", or, nr)
	return Result{}
}

// remoteOnlyDriftUnchanged is true when old and new CRD are both absent for a remote-only
// map entry and the remote value did not change between comparisons.
func remoteOnlyDriftUnchanged(or Result, nr Result) bool {
	return or.CRDValue == nil && nr.CRDValue == nil && or.RemoteValue == nr.RemoteValue
}
