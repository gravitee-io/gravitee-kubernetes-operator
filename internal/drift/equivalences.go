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
	RegisterEquivalenceFunc("ignore-remote", reflect.String, IgnoreRemoteArgs)
	RegisterEquivalenceFunc("ignore-namespace-prefix", reflect.String, IgnoreNamespacePrefix)
	RegisterEquivalenceFunc("case-insensitive", reflect.String, CaseInsensitive)
	RegisterEquivalenceFunc(ignoreName, reflect.Slice, IgnoreSkip)
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
	RegisterEquivalenceFunc("unstructured", reflect.Slice, DefaultEquivalencePostPullUpObjectChildren)
	RegisterEquivalenceFunc("ignore-only", reflect.Slice, IgnoreOnlyArgs)
}

const (
	ignoreOnlyRemote  = "remote"
	ignoreOnlyCRD     = "crd"
	ignoreOnlyStripNS = "strip-ns"
)

// IgnoreOnlyArgs ignores items present only on one side, chosen by ctx.FuncArgs[0]
// ("remote" or "crd"). Items must implement Keyed.
// If FuncArgs contain "strip-ns", keys are compared after stripping
// the namespace prefix and remaining membership is compared as a set.
func IgnoreOnlyArgs(crd any, remote any, ctx DriftContext) Equivalence {
	eq := EmptyIsNilLen(crd, remote, DriftContext{})
	if eq.Equivalent == Equivalent {
		return eq
	}
	cannotCompare := Equivalence{Equivalent: CannotCompare}
	side, ok := ignoreOnlySide(ctx.FuncArgs)
	if !ok {
		return cannotCompare
	}

	crdItems, ok := asKeyed(crd)
	if !ok {
		return cannotCompare
	}
	remoteItems, ok := asKeyed(remote)
	if !ok {
		return cannotCompare
	}

	stripNS := slices.Contains(ctx.FuncArgs[1:], ignoreOnlyStripNS)
	crdNames := keys(crdItems, ctx.Namespace, stripNS)
	remoteNames := keys(remoteItems, ctx.Namespace, stripNS)
	if stripNS {
		return ignoreOnlySetCompare(side, crdNames, remoteNames)
	}
	filter := itemsOnlyFilterFunc(onlyOnSide(side, crdNames, remoteNames))
	if filter == nil {
		return cannotCompare
	}
	if side == ignoreOnlyRemote {
		return Equivalence{Equivalent: CannotCompare, RemoteItemsFilterFunc: filter}
	}
	return Equivalence{Equivalent: CannotCompare, CRDItemsFilterFunc: filter}
}

func ignoreOnlySetCompare(side string, crdNames, remoteNames []string) Equivalence {
	only := onlyOnSide(side, crdNames, remoteNames)
	left, right := crdNames, remoteNames
	if side == ignoreOnlyCRD {
		left = difference(crdNames, only)
	} else {
		right = difference(remoteNames, only)
	}
	slices.Sort(left)
	slices.Sort(right)
	if slices.Equal(left, right) {
		return Equivalence{Equivalent: Equivalent, Skip: true}
	}
	return Equivalence{Equivalent: Inequivalent, Skip: true}
}

func ignoreOnlySide(args []string) (string, bool) {
	if len(args) == 0 {
		return "", false
	}
	switch args[0] {
	case ignoreOnlyRemote, ignoreOnlyCRD:
		return args[0], true
	default:
		return "", false
	}
}

func onlyOnSide(side string, crdNames, remoteNames []string) []string {
	if side == ignoreOnlyRemote {
		return difference(remoteNames, crdNames)
	}
	return difference(crdNames, remoteNames)
}

func difference(from, without []string) []string {
	only := make([]string, 0)
	for _, name := range from {
		if slices.Contains(without, name) {
			continue
		}
		only = append(only, name)
	}
	return only
}

func itemsOnlyFilterFunc(onlyIDs []string) ItemsFilterFunc {
	if len(onlyIDs) == 0 {
		return nil
	}
	return func(items any) []any {
		if items == nil {
			return nil
		}
		v := reflect.ValueOf(items)
		filtered := make([]any, 0, v.Len())
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).Interface()
			if keyed, ok := item.(Keyed); ok {
				if slices.Contains(onlyIDs, keyed.MatchKey()) {
					continue
				}
			}
			filtered = append(filtered, item)
		}
		if len(filtered) == 0 {
			return nil
		}
		return filtered
	}
}

func asKeyed(items any) ([]Keyed, bool) {
	if items == nil {
		return []Keyed{}, true
	}
	v := reflect.ValueOf(items)
	keyed := make([]Keyed, v.Len())
	for i := 0; i < v.Len(); i++ {
		item, ok := v.Index(i).Interface().(Keyed)
		if !ok {
			return nil, false
		}
		keyed[i] = item
	}
	return keyed, true
}

func keys(items []Keyed, namespace string, stripNS bool) []string {
	names := make([]string, len(items))
	prefix := ""
	if stripNS && namespace != "" {
		prefix = namespace + "-"
	}
	for i, item := range items {
		name := item.MatchKey()
		if prefix != "" {
			name = strings.TrimPrefix(name, prefix)
		}
		names[i] = name
	}
	return names
}

// IgnoreRemoteArgs ignores the remote difference if the remote string is in the context.FuncArgs.
func IgnoreRemoteArgs(crd any, remote any, context DriftContext) Equivalence {
	e := DefaultEquivalence(crd, remote, context)
	if e.Equivalent == Inequivalent {
		rs := asString(remote)
		if context.FuncArgs != nil {
			if slices.Contains(context.FuncArgs, rs) {
				return Equivalence{Equivalent: Equivalent}
			}
		}
	}
	return e
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
	return DefaultEquivalence(crd, remote, ctx)
}

// Trimmed trims the value before comparing.
func Trimmed(crd any, remote any, _ DriftContext) Equivalence {
	// the registry protects us from casting panics
	crdString := asString(crd)
	remoteString := asString(remote)
	if strings.TrimSpace(crdString) == strings.TrimSpace(remoteString) {
		return Equivalence{Equivalent: Equivalent}
	}
	return Equivalence{Equivalent: Inequivalent}
}

// IgnoreNamespacePrefix ignores the remote difference if the remote string ends with the crd string.
func IgnoreNamespacePrefix(crd any, remote any, ctx DriftContext) Equivalence {
	crdString := asString(crd)
	remoteString := asString(remote)
	prefix := ctx.Namespace + "-"
	crdString = strings.TrimPrefix(crdString, prefix)
	remoteString = strings.TrimPrefix(remoteString, prefix)
	return DefaultEquivalence(crdString, remoteString, ctx)
}

// RFC3339 checks if the value is a valid RFC3339 time and if they represent the same time.
func RFC3339(crd any, remote any, _ DriftContext) Equivalence {
	// the registry protects us from casting panics
	crdString := asString(crd)
	remoteString := asString(remote)
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

// CaseInsensitive checks if the value is equal ignoring the case.
func CaseInsensitive(crd any, remote any, _ DriftContext) Equivalence {
	if crd == nil && remote == nil {
		return Equivalence{Equivalent: Equivalent}
	}
	if crd == nil || remote == nil {
		return Equivalence{Equivalent: Inequivalent}
	}
	crdString := asString(crd)
	crdString = strings.ToLower(crdString)
	remoteString := asString(remote)
	remoteString = strings.ToLower(remoteString)
	if crdString == remoteString {
		return Equivalence{Equivalent: Equivalent}
	}
	return Equivalence{Equivalent: Inequivalent}
}

// EmptyIsNilInt checks if the value is nil or equal to 0.
func EmptyIsNilInt(crd any, remote any, ctx DriftContext) Equivalence {
	if crd == nil && remote != nil && reflect.DeepEqual(remote, 0) {
		return Equivalence{Equivalent: Equivalent}
	}
	if remote == nil && crd != nil && reflect.DeepEqual(crd, 0) {
		return Equivalence{Equivalent: Equivalent}
	}
	return DefaultEquivalence(crd, remote, ctx)
}

// EmptyIsNilUint checks if the value is nil or equal to 0.
func EmptyIsNilUint(crd any, remote any, ctx DriftContext) Equivalence {
	if crd == nil && remote != nil && reflect.DeepEqual(remote, uint(0)) {
		return Equivalence{Equivalent: Equivalent}
	}
	if remote == nil && crd != nil && reflect.DeepEqual(crd, uint(0)) {
		return Equivalence{Equivalent: Equivalent}
	}
	return DefaultEquivalence(crd, remote, ctx)
}

// EmptyIsNilBool checks if the value is nil or equal to false.
func EmptyIsNilBool(crd any, remote any, ctx DriftContext) Equivalence {
	if crd == nil && remote != nil && reflect.DeepEqual(remote, false) {
		return Equivalence{Equivalent: Equivalent}
	}
	if remote == nil && crd != nil && reflect.DeepEqual(crd, false) {
		return Equivalence{Equivalent: Equivalent}
	}
	return DefaultEquivalence(crd, remote, ctx)
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
		e := DefaultEquivalence(crd, remote, ctx)
		if e.Equivalent == Equivalent {
			// don't need to go further
			e.Skip = true
			return e
		}
	}
	if crd != nil && remote == nil {
		remote = toZero(crd)
		e := DefaultEquivalence(crd, remote, ctx)
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
	var e Equivalence
	if ctx.Kind == reflect.Slice {
		e = defaultSliceEquivalence(crd, remote, ctx)
	} else {
		e = defaultStructEquivalence(crd, remote, ctx)
	}
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
	return DefaultEquivalence(crd, remote, ctx)
}

func DefaultEquivalence(crd any, remote any, _ DriftContext) Equivalence {
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

// asString returns the string representation of the value. Works with pure strings and typed-strings (e.g., enum)
func asString(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.String {
		return rv.String()
	}
	return ""
}
