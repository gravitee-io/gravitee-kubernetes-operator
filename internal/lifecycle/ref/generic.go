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

package ref

import (
	"context"
	"reflect"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
)

// TagSpec is the parsed version of the `ref: "..."` tag.
type TagSpec struct {
	Kind     string
	RefField string
	Key      string
}

// GenericRefResolver is the default ResolveRefs. Walk obj, resolve each `ref` tag, write into the value field.
func GenericRefResolver[T any](ctx context.Context, obj T, parentNs string) error {
	return Walk(ctx, obj, parentNs)
}

// Walk obj, resolve each `ref` tag, write value into the target field.
func Walk(ctx context.Context, obj any, parentNs string) error {
	s, ok := isStruct(reflect.ValueOf(obj))
	if !ok {
		return NewWrappedError("walk", "", "", ErrStructKind)
	}
	return walkStruct(ctx, s, parentNs)
}

func walkStruct(ctx context.Context, s reflect.Value, parentNs string) error {
	targetFields := make(map[string]TagSpec)
	for targetField, val := range s.Fields() {
		if nested, ok := isStruct(val); ok && !isNamespacedName(nested) {
			if err := walkStruct(ctx, nested, parentNs); err != nil {
				return err
			}
		}
		if sl, ok := isSliceOfStruct(val); ok {
			for i := 0; i < sl.Len(); i++ {
				elem, ok := isStruct(sl.Index(i))
				if !ok {
					continue
				}
				if err := walkStruct(ctx, elem, parentNs); err != nil {
					return err
				}
			}
		}
		if !val.IsZero() {
			continue
		}
		err := storeFields(targetFields, targetField)
		if err != nil {
			return err
		}
	}
	return resolveAndAssign(ctx, s, targetFields, parentNs)
}

func storeFields(targetFields map[string]TagSpec, targetField reflect.StructField) error {
	tag := targetField.Tag.Get("ref")
	if tag == "" {
		return nil
	}
	targetName := targetField.Name
	spec, err := ParseRefTag(tag, targetName)
	if err != nil {
		return err
	}
	targetFields[targetName] = spec
	return nil
}

// ParseRefTag is one `ref:"kind[,refField[,key]]"` tag.
// Empty RefField means default <fieldName>Ref. Empty Key means use the whole fetched object.
func ParseRefTag(tag, fieldName string) (TagSpec, error) {
	elements := strings.Split(tag, ",")
	if len(elements) < 1 || len(elements) > 3 {
		return TagSpec{}, NewWrappedError("parse", "", fieldName, ErrRefTagInvalid)
	}
	kind := strings.TrimSpace(elements[0])
	if kind == "" {
		return TagSpec{}, NewWrappedError("parse", "", fieldName, ErrRefTagInvalid)
	}
	refField := fieldName + "Ref"
	if len(elements) >= 2 {
		if named := strings.TrimSpace(elements[1]); named != "" {
			refField = named
		}
	}
	key := ""
	if len(elements) == 3 {
		key = strings.TrimSpace(elements[2])
	}
	return TagSpec{Kind: kind, RefField: refField, Key: key}, nil
}

func resolveAndAssign(ctx context.Context, s reflect.Value, targetFields map[string]TagSpec, parentNs string) error {
	for targetField, spec := range targetFields {
		namespacedName, err := getRefFieldValue(s, spec.RefField)
		if err != nil {
			return err
		}

		val, err := Resolve(ctx, spec, ObjectKey(namespacedName, parentNs))
		if err != nil {
			return err
		}

		fieldVal := s.FieldByName(targetField)
		if !fieldVal.CanSet() {
			return NewWrappedError("assign", "", targetField, ErrCannotSetField)
		}

		newVal := reflect.ValueOf(val)
		if !newVal.Type().AssignableTo(fieldVal.Type()) {
			return NewWrappedError("assign", "", targetField, ErrTypeMismatch)
		}

		fieldVal.Set(newVal)
	}
	return nil
}

func getRefFieldValue(s reflect.Value, field string) (*refs.NamespacedName, error) {
	value := s.FieldByName(field)
	if !value.IsValid() {
		return nil, NewWrappedError("walk", "", field, ErrMissingRef)
	}
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil, NewWrappedError("walk", "", field, ErrMissingRef)
		}
		value = value.Elem()
	}
	if !value.IsValid() {
		return nil, NewWrappedError("walk", "", field, ErrMissingRef)
	}
	if value.Type().AssignableTo(reflect.TypeFor[refs.NamespacedName]()) {
		return value.Addr().Interface().(*refs.NamespacedName), nil
	}
	return nil, NewWrappedError("walk", "", field, ErrRefTypeUnsupported)
}

func isNamespacedName(v reflect.Value) bool {
	return v.Type() == reflect.TypeFor[refs.NamespacedName]()
}

func isSliceOfStruct(val reflect.Value) (reflect.Value, bool) {
	s := val
	if val.Kind() == reflect.Pointer {
		s = val.Elem()
	}
	if s.Kind() != reflect.Slice {
		return reflect.Value{}, false
	}
	elem := s.Type().Elem()
	if elem.Kind() == reflect.Pointer {
		elem = elem.Elem()
	}
	if elem.Kind() != reflect.Struct {
		return reflect.Value{}, false
	}
	return s, true
}

func isStruct(val reflect.Value) (reflect.Value, bool) {
	if val.Kind() == reflect.Struct {
		return val, true
	}
	if val.Kind() == reflect.Pointer && val.Elem().Kind() == reflect.Struct {
		return val.Elem(), true
	}
	return reflect.Value{}, false
}
