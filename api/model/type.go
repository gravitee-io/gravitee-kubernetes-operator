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

package model

import (
	"encoding/json"
	"reflect"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type GenericStringMap struct {
	// +kubebuilder:pruning:PreserveUnknownFields
	unstructured.Unstructured `json:",inline"`
}

func (in *GenericStringMap) DeepCopyInto(out *GenericStringMap) {
	// controller-gen cannot handle the interface{} type of an aliased Unstructured,
	// thus we write our own DeepCopyInto function.
	if out != nil {
		casted := in.Unstructured
		for k, v := range casted.Object {
			if reflect.TypeOf(v).Kind() == reflect.Int {
				casted.Object[k] = int64(v.(int))
			} else if reflect.TypeOf(v).Kind() == reflect.Map {
				if innerMap, ok := v.(map[string]interface{}); ok {
					nestedIn := GenericStringMap{Unstructured: unstructured.Unstructured{Object: innerMap}}
					nestedOut := GenericStringMap{}
					nestedIn.DeepCopyInto(&nestedOut)
					casted.Object[k] = nestedOut.Object
				}
			}
		}

		deepCopy := casted.DeepCopy()
		out.Object = deepCopy.Object
	}
}

func (in *GenericStringMap) UnmarshalJSON(data []byte) error {
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
