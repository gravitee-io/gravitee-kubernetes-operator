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
	RegisterEquivalenceFunc("ignore-remote-prefix", reflect.String, IgnoreRemotePrefix)
	RegisterEquivalenceFunc(ignoreName, reflect.Bool, Ignore)
	RegisterEquivalenceFunc(emptyIsNilName, reflect.Bool, EmptyIsNilBool)
	RegisterEquivalenceFunc(emptyIsNilName, reflect.Int, EmptyIsNilInt)
	RegisterEquivalenceFunc(emptyIsNilName, reflect.Uint, EmptyIsNilUint)
	RegisterEquivalenceFunc(emptyIsNilName, reflect.Slice, EmptyIsNilLen)
	RegisterEquivalenceFunc(emptyIsNilName, reflect.Map, EmptyIsNilLen)
	RegisterEquivalenceFunc(emptyIsNilName, reflect.Struct, EmptyIsNilStruct)
	RegisterEquivalenceFunc("empty-is-true", reflect.Bool, EmptyIsTrue)
	RegisterEquivalenceFunc(ignoreName, reflect.Struct, IgnoreSkip)
	RegisterEquivalenceFunc("unstructured", reflect.Struct, DefaultEquivalencePostPullUpObjectChildren)
}

func Ignore(_ any, _ any) Equivalence {
	return Equivalence{Equivalent: CannotCompare}
}

// EmptyIsNilString checks if the value is nil or empty string.
func EmptyIsNilString(crd any, remote any) Equivalence {
	if crd == nil && remote != nil && remote == "" {
		return Equivalence{Equivalent: Equivalent}
	}
	if remote == nil && crd != nil && crd == "" {
		return Equivalence{Equivalent: Equivalent}
	}
	return FromDeepEqual(crd, remote)
}

// Trimmed trims the value before comparing.
func Trimmed(crd any, remote any) Equivalence {
	// the registry protects us from casting panics
	crdString, _ := crd.(string)
	remoteString, _ := remote.(string)
	if strings.TrimSpace(crdString) == strings.TrimSpace(remoteString) {
		return Equivalence{Equivalent: Equivalent}
	}
	return Equivalence{Equivalent: Inequivalent}
}

// IgnoreRemotePrefix ignores the remote difference if the remote string ends with the crd string.
func IgnoreRemotePrefix(crd any, remote any) Equivalence {
	// the registry protects us from casting panics
	crdString, _ := crd.(string)
	remoteString, _ := remote.(string)
	if len(crdString) <= len(remoteString) {
		if strings.HasSuffix(remoteString, crdString) {
			return Equivalence{Equivalent: Equivalent}
		}
	}
	return FromDeepEqual(crd, remote)
}

// RFC3339 checks if the value is a valid RFC3339 time and if they represent the same time.
func RFC3339(crd any, remote any) Equivalence {
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

// EmptyIsNilInt checks if the value is nil or equal to 0.
func EmptyIsNilInt(crd any, remote any) Equivalence {
	if crd == nil && remote != nil && reflect.DeepEqual(remote, 0) {
		return Equivalence{Equivalent: Equivalent}
	}
	if remote == nil && crd != nil && reflect.DeepEqual(crd, 0) {
		return Equivalence{Equivalent: Equivalent}
	}
	return FromDeepEqual(crd, remote)
}

// EmptyIsNilUint checks if the value is nil or equal to 0.
func EmptyIsNilUint(crd any, remote any) Equivalence {
	if crd == nil && remote != nil && reflect.DeepEqual(remote, uint(0)) {
		return Equivalence{Equivalent: Equivalent}
	}
	if remote == nil && crd != nil && reflect.DeepEqual(crd, uint(0)) {
		return Equivalence{Equivalent: Equivalent}
	}
	return FromDeepEqual(crd, remote)
}

// EmptyIsNilBool checks if the value is nil or equal to false.
func EmptyIsNilBool(crd any, remote any) Equivalence {
	if crd == nil && remote != nil && reflect.DeepEqual(remote, false) {
		return Equivalence{Equivalent: Equivalent}
	}
	if remote == nil && crd != nil && reflect.DeepEqual(crd, false) {
		return Equivalence{Equivalent: Equivalent}
	}
	return FromDeepEqual(crd, remote)
}

// EmptyIsNilLen checks if the slice or map value is nil or len is equal to 0.
func EmptyIsNilLen(crd any, remote any) Equivalence {
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
func EmptyIsNilStruct(crd any, remote any) Equivalence {
	if crd == nil && remote != nil {
		crd = toZero(remote)
		e := FromDeepEqual(crd, remote)
		if e.Equivalent == Equivalent {
			// don't need to go further
			e.Skip = true
			return e
		}
	}
	if crd != nil && remote == nil {
		remote = toZero(crd)
		e := FromDeepEqual(crd, remote)
		if e.Equivalent == Equivalent {
			// don't need to go further
			e.Skip = true
			return e
		}
	}
	return Equivalence{Equivalent: CannotCompare}
}

// IgnoreSkip ignores the comparison and skips the children.
func IgnoreSkip(crd any, remote any) Equivalence {
	r := Ignore(crd, remote)
	r.Skip = true
	return r
}

// DefaultEquivalencePostPullUpObjectChildren perform s a default struct equivalence and adds a post-function moves the children of the "object" property to the root.
func DefaultEquivalencePostPullUpObjectChildren(crd any, remote any) Equivalence {
	e := defaultStructEquivalence(crd, remote)
	e.PostFunc = func(r *Result) {
		var objectChild *Result
		r.Children = slices.DeleteFunc(r.Children, func(e *Result) bool {
			if e.Property == "object" {
				if len(e.Children) > 0 {
					objectChild = e
				}
				return true
			}
			return false
		})

		if objectChild != nil {
			for _, c := range objectChild.Children {
				r.AppendChild(c, true)
			}
		}
	}
	return e
}

func EmptyIsTrue(crd any, remote any) Equivalence {
	if crd == nil && remote != nil && reflect.DeepEqual(remote, true) {
		return Equivalence{Equivalent: Equivalent}
	}
	return FromDeepEqual(crd, remote)
}

func FromDeepEqual(crd any, remote any) Equivalence {
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
