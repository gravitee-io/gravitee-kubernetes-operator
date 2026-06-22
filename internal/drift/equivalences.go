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
	"strings"
	"time"
)

const (
	emptyIsNilName = "empty-is-nil"
	ignoreName     = "ignore"
)

func Init() {
	Register(emptyIsNilName, reflect.String, EmptyIsNilString)
	Register(ignoreName, reflect.String, Ignore)
	Register("trimmed", reflect.String, Trimmed)
	Register("rfc3339", reflect.String, RFC3339)
	Register(emptyIsNilName, reflect.Bool, EmptyIsNilBool)
	Register(emptyIsNilName, reflect.Int, EmptyIsNilInt)
	Register(emptyIsNilName, reflect.Uint, EmptyIsNilUint)
	Register(emptyIsNilName, reflect.Slice, EmptyIsNilLen)
	Register(emptyIsNilName, reflect.Map, EmptyIsNilLen)
	Register(emptyIsNilName, reflect.Struct, EmptyIsNilStruct)
	Register(ignoreName, reflect.Struct, IgnoreSkip)
}

func Ignore(_ any, _ any) Equivalence {
	return Equivalence{Equivalent: Equivalent}
}

func EmptyIsNilString(crd any, api any) Equivalence {
	if crd == nil && api != nil && api == "" {
		return Equivalence{Equivalent: Equivalent}
	}
	if api == nil && crd != nil && crd == "" {
		return Equivalence{Equivalent: Equivalent}
	}
	return FromDeepEqual(crd, api)
}

func Trimmed(crd any, api any) Equivalence {
	// the registry protects us from casting panics
	crdString, _ := crd.(string)
	apiString, _ := api.(string)
	if strings.TrimSpace(crdString) == strings.TrimSpace(apiString) {
		return Equivalence{Equivalent: Equivalent}
	}
	return Equivalence{Equivalent: Inequivalent}
}

func RFC3339(crd any, api any) Equivalence {
	// the registry protects us from casting panics
	crdString, _ := crd.(string)
	apiString, _ := api.(string)
	// avoid parsing error
	if (crdString != "") != (apiString != "") {
		return Equivalence{Equivalent: Inequivalent}
	}
	if crdString == "" && apiString == "" {
		return Equivalence{Equivalent: Equivalent}
	}
	crdRFC3339time, err := time.Parse(time.RFC3339, crdString)
	if err != nil {
		return Equivalence{Equivalent: Inequivalent, Reason: err}
	}
	apiRFC3339time, err := time.Parse(time.RFC3339, apiString)
	if err != nil {
		return Equivalence{Equivalent: Inequivalent, Reason: err}
	}
	if crdRFC3339time.Equal(apiRFC3339time) {
		return Equivalence{Equivalent: Equivalent}
	}
	return Equivalence{Equivalent: Inequivalent}
}

func EmptyIsNilInt(crd any, api any) Equivalence {
	if crd == nil && api != nil && reflect.DeepEqual(api, 0) {
		return Equivalence{Equivalent: Equivalent}
	}
	if api == nil && crd != nil && reflect.DeepEqual(crd, 0) {
		return Equivalence{Equivalent: Equivalent}
	}
	return FromDeepEqual(crd, api)
}

func EmptyIsNilUint(crd any, api any) Equivalence {
	if crd == nil && api != nil && reflect.DeepEqual(api, uint(0)) {
		return Equivalence{Equivalent: Equivalent}
	}
	if api == nil && crd != nil && reflect.DeepEqual(crd, uint(0)) {
		return Equivalence{Equivalent: Equivalent}
	}
	return FromDeepEqual(crd, api)
}
func EmptyIsNilBool(crd any, api any) Equivalence {
	if crd == nil && api != nil && reflect.DeepEqual(api, false) {
		return Equivalence{Equivalent: Equivalent}
	}
	if api == nil && crd != nil && reflect.DeepEqual(crd, false) {
		return Equivalence{Equivalent: Equivalent}
	}
	return FromDeepEqual(crd, api)
}

func EmptyIsNilLen(crd any, api any) Equivalence {
	crdLen := reflect.ValueOf(crd).Len()
	apiLen := reflect.ValueOf(api).Len()
	if crdLen == apiLen {
		return Equivalence{Equivalent: Equivalent}
	}
	return Equivalence{Equivalent: CannotCompare}
}

func EmptyIsNilStruct(crd any, api any) Equivalence {
	if crd == nil && api != nil {
		crd = toZero(api)
		e := FromDeepEqual(crd, api)
		if e.Equivalent == Equivalent {
			// don't need to go further
			e.Skip = true
			return e
		}
	}
	if crd != nil && api == nil {
		api = toZero(crd)
		e := FromDeepEqual(crd, api)
		if e.Equivalent == Equivalent {
			// don't need to go further
			e.Skip = true
			return e
		}
	}
	return Equivalence{Equivalent: CannotCompare}
}

func IgnoreSkip(_ any, _ any) Equivalence {
	return Equivalence{Equivalent: Equivalent, Skip: true}
}

func toZero(v any) any {
	return reflect.Zero(reflect.TypeOf(v)).Interface()
}
