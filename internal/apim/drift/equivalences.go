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
	"reflect"
	"slices"
	"strings"

	driftlib "github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
)

func InitRegistry() {
	driftlib.RegisterEquivalenceFunc("ignore-unknown-crd-groups", reflect.Slice, IgnoreUnknownCRDGroups)
}

// IgnoreUnknownCRDGroups ignores CRD-only group names (not present on remote after
// stripping the namespace prefix), then compares remaining membership as a set.
func IgnoreUnknownCRDGroups(crd any, remote any, context driftlib.DriftContext) driftlib.Equivalence {
	eq := driftlib.EmptyIsNilLen(crd, remote, driftlib.DriftContext{})
	if eq.Equivalent == driftlib.Equivalent {
		return eq
	}
	crdItems := toAnySlice(crd)
	remoteItems := toAnySlice(remote)

	remoteNorm := make([]string, 0, len(remoteItems))
	for _, item := range remoteItems {
		s, ok := item.(string)
		if !ok {
			return driftlib.Equivalence{Equivalent: driftlib.CannotCompare}
		}
		remoteNorm = append(remoteNorm, stripNamespacePrefix(s, context.Namespace))
	}

	crdKnown := make([]string, 0, len(crdItems))
	for _, item := range crdItems {
		s, ok := item.(string)
		if !ok {
			return driftlib.Equivalence{Equivalent: driftlib.CannotCompare}
		}
		norm := stripNamespacePrefix(s, context.Namespace)
		if slices.Contains(remoteNorm, norm) {
			crdKnown = append(crdKnown, norm)
		}
	}

	slices.Sort(crdKnown)
	slices.Sort(remoteNorm)
	if slices.Equal(crdKnown, remoteNorm) {
		return driftlib.Equivalence{Equivalent: driftlib.Equivalent, Skip: true}
	}
	return driftlib.Equivalence{Equivalent: driftlib.Inequivalent, Skip: true}
}

func stripNamespacePrefix(value, namespace string) string {
	if namespace == "" {
		return value
	}
	return strings.TrimPrefix(value, namespace+"-")
}

func toAnySlice(a any) []any {
	if a != nil {
		// ,pt check of required here we know it is a slice already
		v := reflect.ValueOf(a)
		result := make([]any, v.Len())
		for i := 0; i < v.Len(); i++ {
			result[i] = v.Index(i).Interface()
		}
		return result
	}
	return make([]any, 0)
}
