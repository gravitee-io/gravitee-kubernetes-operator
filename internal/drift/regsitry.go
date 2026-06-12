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
	"log"
	"reflect"
)

type registryEntry struct {
	kind reflect.Kind
	name string
}

type registry struct {
	registry map[registryEntry]EquivalenceFunc
}

func (r *registry) Get(name string, t reflect.Kind) EquivalenceFunc {
	if f, ok := r.registry[registryEntry{
		name: name,
		kind: t,
	}]; ok {
		return f
	}
	if name != "" {
		log.Panicf("drift function '%s' not found", name)
	}
	if t == reflect.Array || t == reflect.Slice {
		return defaultSliceArrayEquivalence
	}
	if t == reflect.Struct {
		return defaultStructEquivalence
	}
	return defaultEquivalence
}

func defaultStructEquivalence(crd any, api any) Equivalence {
	// xor
	if (crd == nil) != (api == nil) {
		return Equivalence{
			Equivalent: Inequivalent,
			Skip:       true,
		}
	}
	return Equivalence{
		Equivalent: CannotCompare,
		Skip:       false,
	}
}

var equivalenceRegistry = &registry{
	registry: make(map[registryEntry]EquivalenceFunc),
}

// Register registers a named EquivalenceFunc function 'f' to compare value of kind 'k'.
func Register(name string, k reflect.Kind, f EquivalenceFunc) {
	if k == reflect.Pointer {
		panic("cannot register a pointer to a struct, use a concrete type or an interface")
	}

	equivalenceRegistry.registry[registryEntry{
		kind: k,
		name: name,
	}] = func(crd any, api any) Equivalence {
		assertTypes(crd, api)
		e := f(crd, api)
		return e
	}
}

func assertTypes(crd any, api any) {
	crdType := reflect.TypeOf(crd)
	apiType := reflect.TypeOf(api)
	if crdType != nil && apiType != nil && crdType != apiType {
		log.Panicf("drift detection only work comparing values of same type, crd=%T api=%T", crd, api)
	}
	if crdType != nil && crdType.Kind() == reflect.Interface {
		log.Panicf("drift detection only compare non interface types")
	}
	if apiType != nil && apiType.Kind() == reflect.Interface {
		log.Panicf("drift detection only compare non interface types")
	}
}

func defaultEquivalence(crd any, api any) Equivalence {
	return FromDeepEqual(crd, api)
}

func defaultSliceArrayEquivalence(any, any) Equivalence {
	return Equivalence{
		Equivalent: CannotCompare,
		Skip:       false,
	}
}
