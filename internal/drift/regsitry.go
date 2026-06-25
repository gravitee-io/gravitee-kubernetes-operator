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

// Get returns the EquivalenceFunc function for the given kind.
// If no function is registered for the given kind and name, it panics.
// If no name is given (no drift tag), it returns a default function for the given kind:
// - For slices and arrays, it returns a function that always returns CannotCompare
// - For structs, it returns a function that returns Inequivalent if one of the fields is exclusively nil or returns CannotCompare
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
	if k == reflect.Array || k == reflect.Slice {
		return defaultSliceArrayEquivalence
	}
	if k == reflect.Struct {
		return defaultStructEquivalence
	}
	return FromDeepEqual
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
	}] = func(crd any, remote any) Equivalence {
		assertTypes(crd, remote)
		e := f(crd, remote)
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

func defaultStructEquivalence(crd any, remote any) Equivalence {
	// xor
	if (crd == nil) != (remote == nil) {
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

func defaultSliceArrayEquivalence(any, any) Equivalence {
	return Equivalence{
		Equivalent: CannotCompare,
		Skip:       false,
	}
}
