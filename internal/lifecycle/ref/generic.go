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
)

// RefSpec is one `ref:"kind[,refField[,key]]"` tag.
// Empty RefField means default <fieldName>Ref. Empty Key means use the whole fetched object.
type RefSpec struct {
	Kind     string
	RefField string
	Key      string
}

func ParseRefTag(tag, fieldName string) (RefSpec, error) {
	return RefSpec{}, ErrNotImplemented
}

func Walk(ctx context.Context, obj any, parentNs string) error {
	if reflect.ValueOf(obj).Kind() != reflect.Struct {
		return ErrUnknownKind
	}
	for f, val := range reflect.ValueOf(obj).Fields() {
		tag := f.Tag.Get("ref")
		if tag == "" {
			continue
		}
		spec, err := ParseRefTag(tag, f.Name)
		if err != nil {
			return err
		}
		resolved, err := Resolve(ctx, spec, parentNs)

	}
	return ErrNotImplemented
}

// GenericRefResolver is the default ResolveRefs. Walk obj, resolve each `ref` tag, write into the value field.
func GenericRefResolver[T any](ctx context.Context, obj T) error {

	return ErrNotImplemented
}
