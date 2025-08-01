// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type GenericStringMap struct {
	// +kubebuilder:pruning:PreserveUnknownFields
	unstructured.Unstructured `json:",inline"`
}

func NewGenericStringMap() *GenericStringMap {
	return &GenericStringMap{
		Unstructured: unstructured.Unstructured{
			Object: make(map[string]interface{}),
		},
	}
}

func ToGenericStringMap(obj map[string]interface{}) *GenericStringMap {
	return &GenericStringMap{
		Unstructured: unstructured.Unstructured{
			Object: obj,
		},
	}
}

func (in *GenericStringMap) Put(key string, value interface{}) *GenericStringMap {
	in.Object[key] = value
	return in
}

func (in *GenericStringMap) Get(key string) interface{} {
	if v, ok := in.Object[key]; ok {
		return v
	}
	return nil
}

func (in *GenericStringMap) Remove(key string) {
	delete(in.Object, key)
}

func (in *GenericStringMap) GetString(key string) string {
	if v, found := in.Object[key]; found {
		s, ok := v.(string)
		if ok {
			return s
		}
	}
	return ""
}

func (in *GenericStringMap) GetBool(key string) bool {
	impl := in.Get(key)
	if impl == nil {
		return false
	}

	if b, ok := impl.(bool); ok {
		return b
	}

	return false
}

func (in *GenericStringMap) GetSlice(key string) []interface{} {
	if v, found := in.Object[key]; found {
		s, ok := v.([]interface{})
		if ok {
			return s
		}
	}
	return nil
}

func (in *GenericStringMap) DeepCopyInto(out *GenericStringMap) {
	// controller-gen cannot handle the interface{} type of an aliased Unstructured,
	// thus we write our own DeepCopyInto function.
	if out != nil {
		casted := in.Unstructured
		casted.Object = normalizeObject(casted.Object).(map[string]interface{}) //nolint:errcheck // map is expected

		deepCopy := casted.DeepCopy()
		out.Object = deepCopy.Object
	}
}

func normalizeObject(x interface{}) interface{} {
	switch x := x.(type) {
	case map[string]interface{}:
		if x == nil {
			// Typed nil - an interface{} that contains a type map[string]interface{} with a value of nil
			return x
		}
		clone := make(map[string]interface{}, len(x))
		for k, v := range x {
			clone[k] = normalizeObject(v)
		}
		return clone
	case []interface{}:
		if x == nil {
			// Typed nil - an interface{} that contains a type []interface{} with a value of nil
			return x
		}
		clone := make([]interface{}, len(x))
		for i, v := range x {
			clone[i] = normalizeObject(v)
		}
		return clone
	case int:
		return int64(x)
	default:
		return x
	}
}

func (in *GenericStringMap) UnmarshalYAML(data []byte) error {
	if in == nil {
		return nil
	}

	m := make(map[string]interface{})
	if err := yaml.Unmarshal(data, &m); err != nil {
		return err
	}

	in.Object = m

	return nil
}

func (in *GenericStringMap) UnmarshalJSON(data []byte) error {
	if in == nil {
		return nil
	}

	m := make(map[string]interface{})
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	in.Object = m

	return nil
}

func (in *GenericStringMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(in.Object)
}

func (in *GenericStringMap) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(in.Object)
}
