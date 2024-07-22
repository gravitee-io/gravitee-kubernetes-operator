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

package dynamic

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Convert[T any](source any, target T) (T, error) {
	b, err := json.Marshal(source)
	if err != nil {
		return target, err
	}

	if err = json.Unmarshal(b, target); err != nil {
		return target, err
	}

	return target, nil
}

func ConvertList[T any](list *unstructured.UnstructuredList) ([]T, error) {
	apis := make([]T, 0)
	t := new(T)
	for _, it := range list.Items {
		api, err := Convert(it.Object["spec"], *t)
		if err != nil {
			return apis, err
		}
		apis = append(apis, api)
	}
	return apis, nil
}
