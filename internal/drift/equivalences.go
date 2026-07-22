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
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
)

const (
	emptyIsNilName = "empty-is-nil"
	ignoreName     = "ignore"
)

// InitRegistry initializes the equivalence registry.
func InitRegistry() {
	RegisterEquivalenceFunc(emptyIsNilName, reflect.String, EmptyIsNilString)
	RegisterEquivalenceFunc(ignoreName, reflect.String, Ignore)
	RegisterEquivalenceFunc("trimmed", reflect.String, Trimmed)
	RegisterEquivalenceFunc("rfc3339", reflect.String, RFC3339)
	RegisterEquivalenceFunc("ignore-namespace-prefix", reflect.String, IgnoreNamespacePrefix)
	RegisterEquivalenceFunc("ignore-remote-only-metadata", reflect.Slice, IgnoreRemoteOnlyMetadata)
	RegisterEquivalenceFunc("ignore-crd-only-and-namespace-prefix", reflect.Slice, IgnoreCRDOnlyThenIgnoreNamespacePrefix)
	RegisterEquivalenceFunc(ignoreName, reflect.Bool, Ignore)
	RegisterEquivalenceFunc(emptyIsNilName, reflect.Bool, EmptyIsNilBool)
	RegisterEquivalenceFunc(emptyIsNilName, reflect.Int, EmptyIsNilInt)
	RegisterEquivalenceFunc(emptyIsNilName, reflect.Int32, EmptyIsNilInt)
	RegisterEquivalenceFunc(emptyIsNilName, reflect.Uint, EmptyIsNilUint)
	RegisterEquivalenceFunc(emptyIsNilName, reflect.Slice, EmptyIsNilLen)
	RegisterEquivalenceFunc(emptyIsNilName, reflect.Map, EmptyIsNilLen)
	RegisterEquivalenceFunc(emptyIsNilName, reflect.Struct, EmptyIsNilStruct)
	RegisterEquivalenceFunc("empty-is-true", reflect.Bool, EmptyIsTrue)
	RegisterEquivalenceFunc(ignoreName, reflect.Struct, IgnoreSkip)
	RegisterEquivalenceFunc("unstructured", reflect.Struct, DefaultEquivalencePostPullUpObjectChildren)
}

func Ignore(_ any, r any, c DriftContext) Equivalence {
	return Equivalence{Equivalent: CannotCompare}
}

// EmptyIsNilString checks if the value is nil or empty string.
func EmptyIsNilString(crd any, remote any, ctx DriftContext) Equivalence {
	if crd == nil && remote != nil && remote == "" {
		return Equivalence{Equivalent: Equivalent}
	}
	if remote == nil && crd != nil && crd == "" {
		return Equivalence{Equivalent: Equivalent}
	}
	return FromDeepEqual(crd, remote, ctx)
}

// Trimmed trims the value before comparing.
func Trimmed(crd any, remote any, _ DriftContext) Equivalence {
	// the registry protects us from casting panics
	crdString, _ := crd.(string)
	remoteString, _ := remote.(string)
	if strings.TrimSpace(crdString) == strings.TrimSpace(remoteString) {
		return Equivalence{Equivalent: Equivalent}
	}
	return Equivalence{Equivalent: Inequivalent}
}

// IgnoreNamespacePrefix ignores the remote difference if the remote string ends with the crd string.
func IgnoreNamespacePrefix(crd any, remote any, ctx DriftContext) Equivalence {
	crdString, _ := crd.(string)
	remoteString, _ := remote.(string)
	prefix := ctx.Namespace + "-"
	crdString = strings.TrimPrefix(crdString, prefix)
	remoteString = strings.TrimPrefix(remoteString, prefix)
	return FromDeepEqual(crdString, remoteString, ctx)
}

// RFC3339 checks if the value is a valid RFC3339 time and if they represent the same time.
func RFC3339(crd any, remote any, _ DriftContext) Equivalence {
	// the registry protects us from casting panics
	crdString, _ := crd.(string)
	remoteString, _ := remote.(string)
	// avoid parsing error
	if (crdString != "") != (remoteString != "") {
		return Equivalence{Equivalent: Inequivalent}
	}
	if crdString == "" && remoteString == "" {
		return Equivalence{Equivalent: Equivalent}
	}
	crdRFC3339time, err := parseRFC3339(crdString)
	if err != nil {
		return Equivalence{Equivalent: Inequivalent, Reason: err}
	}
	remoteRFC3339time, err := parseRFC3339(remoteString)
	if err != nil {
		return Equivalence{Equivalent: Inequivalent, Reason: err}
	}
	if crdRFC3339time.Equal(remoteRFC3339time) {
		return Equivalence{Equivalent: Equivalent}
	}
	return Equivalence{Equivalent: Inequivalent}
}

// IgnoreRemoteOnlyMetadata ignores remote-only model.Metadata by name (unique);
// If some are found: provides a ItemsFilterFunc function to remove those before items are compared.
func IgnoreRemoteOnlyMetadata(crd any, remote any, _ DriftContext) Equivalence {
	eq := EmptyIsNilLen(crd, remote, DriftContext{})
	if eq.Equivalent == Equivalent {
		return eq
	}
	crdItems := toAnySlice(crd)
	remoteItems := toAnySlice(remote)

	// get names or return to trigger item comparison
	crdNames, ok := metadataNames(crdItems)
	if !ok {
		return Equivalence{Equivalent: CannotCompare}
	}
	remoteNames, ok := metadataNames(remoteItems)
	if !ok {
		return Equivalence{Equivalent: CannotCompare}
	}

	remoteOnlyNames := make([]string, 0)
	for _, remoteName := range remoteNames {
		if slices.Contains(crdNames, remoteName) {
			continue
		}
		remoteOnlyNames = append(remoteOnlyNames, remoteName)
	}
	if len(remoteOnlyNames) > 0 {
		return Equivalence{Equivalent: CannotCompare, RemoteItemsFilterFunc: metadataRemoteItemsOnlyFilterFunc(remoteOnlyNames)}
	}
	return Equivalence{Equivalent: CannotCompare}
}

// IgnoreCRDOnlyThenIgnoreNamespacePrefix ignores string not in remote difference then compare ignoring namespace prefix.
func IgnoreCRDOnlyThenIgnoreNamespacePrefix(crd any, remote any, context DriftContext) Equivalence {
	eq := EmptyIsNilLen(crd, remote, DriftContext{})
	if eq.Equivalent == Equivalent {
		return eq
	}
	crdItems := toAnySlice(crd)
	remoteItems := toAnySlice(remote)
	cleanCrdItems := make([]string, 0)
	for _, crdItem := range crdItems {
		s, ok := crdItem.(string)
		if !ok {
			return Equivalence{Equivalent: CannotCompare}
		}
		if !slices.Contains(remoteItems, crdItem) {
			continue
		}
		cleanCrdItems = append(cleanCrdItems, s)
	}
	e := IgnoreNamespacePrefix(cleanCrdItems, remoteItems, context)
	if e.Equivalent == Equivalent {
		return Equivalence{Equivalent: Equivalent, Skip: true}
	}
	return Equivalence{Equivalent: CannotCompare}
}

func metadataRemoteItemsOnlyFilterFunc(remoteOnlyNames []string) ItemsFilterFunc {
	return func(items any) []any {
		filtered := make([]any, 0)
		for _, item := range toAnySlice(items) {
			if md, ok := item.(model.Metadata); ok {
				// skip the remote-only items
				if slices.Contains(remoteOnlyNames, md.GetName()) {
					continue
				}
			}
			// the rest needs to be compared
			filtered = append(filtered, item)
		}
		return filtered
	}
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

func metadataNames(items []any) ([]string, bool) {
	names := make([]string, len(items))
	for i, item := range items {
		if md, ok := item.(model.Metadata); ok {
			names[i] = md.GetName()
		} else {
			return nil, false
		}
	}
	return names, true
}

// EmptyIsNilInt checks if the value is nil or equal to 0.
func EmptyIsNilInt(crd any, remote any, ctx DriftContext) Equivalence {
	if crd == nil && remote != nil && reflect.DeepEqual(remote, 0) {
		return Equivalence{Equivalent: Equivalent}
	}
	if remote == nil && crd != nil && reflect.DeepEqual(crd, 0) {
		return Equivalence{Equivalent: Equivalent}
	}
	return FromDeepEqual(crd, remote, ctx)
}

// EmptyIsNilUint checks if the value is nil or equal to 0.
func EmptyIsNilUint(crd any, remote any, ctx DriftContext) Equivalence {
	if crd == nil && remote != nil && reflect.DeepEqual(remote, uint(0)) {
		return Equivalence{Equivalent: Equivalent}
	}
	if remote == nil && crd != nil && reflect.DeepEqual(crd, uint(0)) {
		return Equivalence{Equivalent: Equivalent}
	}
	return FromDeepEqual(crd, remote, ctx)
}

// EmptyIsNilBool checks if the value is nil or equal to false.
func EmptyIsNilBool(crd any, remote any, ctx DriftContext) Equivalence {
	if crd == nil && remote != nil && reflect.DeepEqual(remote, false) {
		return Equivalence{Equivalent: Equivalent}
	}
	if remote == nil && crd != nil && reflect.DeepEqual(crd, false) {
		return Equivalence{Equivalent: Equivalent}
	}
	return FromDeepEqual(crd, remote, ctx)
}

// EmptyIsNilLen checks if the slice or map value is nil or len is equal to 0.
func EmptyIsNilLen(crd any, remote any, _ DriftContext) Equivalence {
	var crdLen int
	var remoteLen int
	if crd != nil {
		crdLen = reflect.ValueOf(crd).Len()
	}
	if remote != nil {
		remoteLen = reflect.ValueOf(remote).Len()
	}
	if crdLen == 0 && remoteLen == 0 {
		return Equivalence{Equivalent: Equivalent, Skip: true}
	}
	return Equivalence{Equivalent: CannotCompare}
}

// EmptyIsNilStruct checks if one struct is nil and the other is an empty struct or vice versa, and reports equivalence.
func EmptyIsNilStruct(crd any, remote any, ctx DriftContext) Equivalence {
	if crd == nil && remote != nil {
		crd = toZero(remote)
		e := FromDeepEqual(crd, remote, ctx)
		if e.Equivalent == Equivalent {
			// don't need to go further
			e.Skip = true
			return e
		}
	}
	if crd != nil && remote == nil {
		remote = toZero(crd)
		e := FromDeepEqual(crd, remote, ctx)
		if e.Equivalent == Equivalent {
			// don't need to go further
			e.Skip = true
			return e
		}
	}
	return Equivalence{Equivalent: CannotCompare}
}

// IgnoreSkip ignores the comparison and skips the children.
func IgnoreSkip(crd any, remote any, ctx DriftContext) Equivalence {
	r := Ignore(crd, remote, ctx)
	r.Skip = true
	return r
}

// DefaultEquivalencePostPullUpObjectChildren perform s a default struct equivalence and adds a post-function moves the children of the "object" property to the root.
func DefaultEquivalencePostPullUpObjectChildren(crd any, remote any, ctx DriftContext) Equivalence {
	e := defaultStructEquivalence(crd, remote, ctx)
	e.PostFunc = func(r *Result) {
		var objectChild *Result
		r.children = slices.DeleteFunc(r.children, func(e *Result) bool {
			if e.Property == "object" {
				if len(e.children) > 0 {
					objectChild = e
				}
				return true
			}
			return false
		})

		if objectChild != nil {
			for _, c := range objectChild.children {
				r.AppendChild(c, true)
			}
		}
	}
	return e
}

func EmptyIsTrue(crd any, remote any, ctx DriftContext) Equivalence {
	if crd == nil && remote != nil && reflect.DeepEqual(remote, true) {
		return Equivalence{Equivalent: Equivalent}
	}
	return FromDeepEqual(crd, remote, ctx)
}

func FromDeepEqual(crd any, remote any, _ DriftContext) Equivalence {
	eq := reflect.DeepEqual(remote, crd)
	if eq {
		return Equivalence{Equivalent: Equivalent}
	}
	return Equivalence{Equivalent: Inequivalent}
}

func toZero(v any) any {
	return reflect.Zero(reflect.TypeOf(v)).Interface()
}

func parseRFC3339(value string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, value)
	if err == nil {
		return t, nil
	}
	return time.Parse(time.RFC3339Nano, value)
}
