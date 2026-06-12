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

import "reflect"

const (
	emptyIsNilName = "empty-is-nil"
)

func Init() {
	Register(emptyIsNilName, reflect.String, EmptyIsNilString)
	Register(emptyIsNilName, reflect.Bool, EmptyIsNilBool)
	Register(emptyIsNilName, reflect.Int, EmptyIsNilInt)
	Register(emptyIsNilName, reflect.Uint, EmptyIsNilUint)
	Register(emptyIsNilName, reflect.Slice, EmptyIsNilLen)
	Register(emptyIsNilName, reflect.Map, EmptyIsNilLen)
	Register(emptyIsNilName, reflect.Struct, EmptyIsNilStruct)
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

func toZero(v any) any {
	return reflect.Zero(reflect.TypeOf(v)).Interface()
}
