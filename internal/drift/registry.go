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

var equivalenceRegistry = &registry{
	registry: make(map[registryEntry]EquivalenceFunc),
}

// Get returns the EquivalenceFunc function for the given kind.
// If no function is registered for the given kind and name, it panics.
// If no name is given (no drift tag), it returns a default function for the given kind:
// - For slices and arrays, it returns a function that always returns CannotCompare
// - For structs, it returns a function that returns CannotCompare and never skip, so that children are compared.
// - For other types, it returns a function that calls reflect.DeepEqual.
func (r *registry) Get(name string, k reflect.Kind) EquivalenceFunc {
	if f, ok := r.registry[registryEntry{
		name: name,
		kind: k,
	}]; ok {
		return f
	} else if name != "" {
		log.Panicf("drift function '%s' not found for kind '%s'", name, k)
	}
	if k == reflect.Slice {
		return defaultSliceEquivalence
	}
	if k == reflect.Struct {
		return defaultStructEquivalence
	}
	return FromDeepEqual
}

// RegisterEquivalenceFunc registers a named EquivalenceFunc function 'f' to compare value of kind 'k'.
// It adds a default safeguard to make sure that both types passed to the equivalence function are of the same type.
// It also checks that neither are interface types. In both cases, it panics.
func RegisterEquivalenceFunc(name string, k reflect.Kind, f EquivalenceFunc) {
	if k == reflect.Pointer {
		panic("cannot register a pointer, use a concrete type")
	}

	equivalenceRegistry.registry[registryEntry{
		kind: k,
		name: name,
	}] = func(crd any, remote any, ctx DriftContext) Equivalence {
		assertTypes(crd, remote)
		e := f(crd, remote, ctx)
		return e
	}
}

func assertTypes(crd any, remote any) {
	crdType := reflect.TypeOf(crd)
	remoteType := reflect.TypeOf(remote)
	if crdType != nil && remoteType != nil && crdType != remoteType {
		log.Panicf("drift detection only work comparing values of same type, crd=%T remote=%T", crd, remote)
	}
	if crdType != nil && crdType.Kind() == reflect.Interface {
		log.Panicf("drift detection only compare non interface types")
	}
	if remoteType != nil && remoteType.Kind() == reflect.Interface {
		log.Panicf("drift detection only compare non interface types")
	}
}

// defaultStructEquivalence if one of the fields is exclusively nil, the whole struct is not equivalent and elements are skipped.
// Otherwise, it is marked as cannot compare, but fields are not skipped.
func defaultStructEquivalence(any, any, DriftContext) Equivalence {
	return Equivalence{
		Equivalent: CannotCompare,
		Skip:       false,
	}
}

// defaultSliceEquivalence the whole slice or array is marked as cannot compare but elements are not skipped.
func defaultSliceEquivalence(any, any, DriftContext) Equivalence {
	return Equivalence{
		Equivalent: CannotCompare,
		Skip:       false,
	}
}
