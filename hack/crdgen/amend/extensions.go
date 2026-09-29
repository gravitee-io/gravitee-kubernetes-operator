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

package amend

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

type extension func(f field, value any) error

// extensions covers what neither inference nor standard oapi-codegen extensions can express.
var extensions = map[string]extension{
	"x-kubebuilder": rawMarkers,
	"x-gko-default": keepDefault,
}

// ApplyExtensions runs the registry function of each custom key, then deletes
// the key: the amended OAS carries standard extensions only.
func ApplyExtensions(doc *openapi3.T) error {
	var errs []error
	for _, f := range fields(doc) {
		for _, key := range slices.Sorted(maps.Keys(f.extensions)) {
			if !isCustom(key) {
				continue
			}
			if err := applyExtension(f, key); err != nil {
				errs = append(errs, fmt.Errorf("%s: %s: %w", f.pointer, key, err))
			}
			delete(f.extensions, key)
		}
	}
	return errors.Join(errs...)
}

func isCustom(key string) bool {
	return key == "x-kubebuilder" || strings.HasPrefix(key, "x-gko-")
}

func applyExtension(f field, key string) error {
	ext, ok := extensions[key]
	switch {
	case !ok:
		return errors.New("unknown extension")
	case !inline(f):
		return errors.New("not supported next to a $ref")
	}
	return ext(f, f.extensions[key])
}

// rawMarkers appends kubebuilder markers the OAS cannot express, e.g. validation stated only in prose.
func rawMarkers(f field, value any) error {
	list, ok := value.([]any)
	if !ok {
		return fmt.Errorf("want a list of markers, got %v", value)
	}
	markers := make([]string, len(list))
	for i, v := range list {
		s, ok := v.(string)
		if !ok {
			return fmt.Errorf("want a marker string, got %v", v)
		}
		markers[i] = "+kubebuilder:" + s
	}
	appendDescription(f.schema, markers...)
	return nil
}

// keepDefault turns the OAS default into a CRD default.
func keepDefault(f field, value any) error {
	if value != "keep" {
		return fmt.Errorf("want keep, got %v", value)
	}
	if f.schema.Default == nil {
		return errors.New("no OAS default to keep")
	}
	data, err := json.Marshal(f.schema.Default)
	if err != nil {
		return err
	}
	appendDescription(f.schema, "+kubebuilder:default="+string(data))
	return nil
}
