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
	"fmt"
	"log"
	"reflect"
)

const typeIndex = 0

// unionShape is the struct a generated SDK union (OneOf*) is compared as. The SDK holds one pointer per variant,
// tagged `oneof:"<TYPE>"`, while the wire carries the variant that is set as one object, its type beside its fields.
// The flat struct holds the type, then every property of every variant once, so a type change reads as one type line
// and a field shared by two variants is compared once. Old and new comparisons then have the same properties
// whichever variants they set, which Merge pairs on.
type unionShape struct {
	flat reflect.Type
	// tags[v] is the oneof tag of variant v, the type the wire carries when it is set.
	tags []string
	// indexes[v][f] is the field of flat that field f of variant v is copied into.
	indexes [][]int
}

// asWireShape returns a generated SDK union flattened as the wire carries it, any other struct as it is.
func asWireShape(t reflect.Type, crd, remote any) (reflect.Type, any, any) {
	if !isUnion(t) {
		return t, crd, remote
	}
	shape := newUnionShape(t)
	return shape.flat, shape.flatten(crd), shape.flatten(remote)
}

func isUnion(t reflect.Type) bool {
	if t.Kind() != reflect.Struct || t.NumField() == 0 {
		return false
	}
	for i := range t.NumField() {
		field := t.Field(i)
		_, tagged := field.Tag.Lookup("oneof")
		if !tagged || field.Type.Kind() != reflect.Pointer || field.Type.Elem().Kind() != reflect.Struct {
			return false
		}
	}
	return true
}

// newUnionShape shares a property between variants only when they compare it alike, same type and same drift tag;
// otherwise each variant keeps its own.
func newUnionShape(union reflect.Type) unionShape {
	fields := []reflect.StructField{{Name: "Type", Type: reflect.TypeFor[*string](), Tag: `json:"type"`}}
	byProperty := map[string]int{"type": typeIndex}
	shape := unionShape{tags: make([]string, union.NumField()), indexes: make([][]int, union.NumField())}

	for v := range union.NumField() {
		shape.tags[v] = union.Field(v).Tag.Get("oneof")
		variant := union.Field(v).Type.Elem()
		shape.indexes[v] = make([]int, variant.NumField())
		for f := range variant.NumField() {
			field := variant.Field(f)
			if !field.IsExported() {
				log.Panicf("detect drift does not support unexported fields, '%s.%s' is unexported.", variant, field.Name)
			}
			property := getProperty(field)
			index, known := byProperty[property]
			if !known || (index != typeIndex && !comparesAlike(fields[index], field)) {
				index = len(fields)
				byProperty[property] = index
				fields = append(fields, reflect.StructField{
					Name: fmt.Sprintf("F%d", index),
					Type: optional(field.Type),
					Tag:  field.Tag,
				})
			}
			shape.indexes[v][f] = index
		}
	}

	shape.flat = reflect.StructOf(fields)
	return shape
}

func comparesAlike(flat reflect.StructField, field reflect.StructField) bool {
	return flat.Type == optional(field.Type) && flat.Tag.Get("drift") == field.Tag.Get("drift")
}

// optional is the type a variant field takes in the flat struct, nil when the variant that is set does not carry it.
func optional(t reflect.Type) reflect.Type {
	if kind := t.Kind(); kind == reflect.Pointer || kind == reflect.Slice || kind == reflect.Map ||
		kind == reflect.Interface {
		return t
	}
	return reflect.PointerTo(t)
}

// flatten copies a value of the union into the flat struct. The type comes from the oneof tag, as the SDK writes it
// on the wire, not from the variant's own type field.
func (s unionShape) flatten(union any) any {
	if union == nil {
		return nil
	}
	value := reflect.ValueOf(union)
	flat := reflect.New(s.flat).Elem()
	for v := range value.NumField() {
		variant := value.Field(v)
		if variant.IsNil() {
			continue
		}
		flat.Field(typeIndex).Set(reflect.ValueOf(new(s.tags[v])))
		for f, index := range s.indexes[v] {
			if index != typeIndex {
				setOptional(flat.Field(index), variant.Elem().Field(f))
			}
		}
	}
	return flat.Interface()
}

func setOptional(target reflect.Value, value reflect.Value) {
	if target.Type() == value.Type() {
		target.Set(value)
		return
	}
	pointer := reflect.New(value.Type())
	pointer.Elem().Set(value)
	target.Set(pointer)
}
