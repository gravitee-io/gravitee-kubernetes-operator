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

package template

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type valMapper func(interface{}) (interface{}, error)

func traverse(ctx context.Context, obj runtime.Object) (interface{}, error) {
	inner, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return obj, err
	}

	cp, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj.DeepCopyObject())
	if err != nil {
		return obj, err
	}

	wrapper := unstructured.Unstructured{Object: inner}
	ns := wrapper.GetNamespace()

	// remove everything we don't want to compile from the object
	cp["status"] = nil
	cp["metadata"] = nil

	result, err := doTraverse(cp, func(val interface{}) (interface{}, error) {
		if v, ok := val.(string); ok {
			return exec(ctx, v, ns)
		}

		return val, nil
	})

	if err != nil {
		return nil, err
	}

	resultMap, _ := result.(map[string]interface{})

	inner["spec"] = resultMap["spec"]

	return inner, nil
}

func doTraverse(obj interface{}, mapper valMapper) (interface{}, error) {
	switch val := obj.(type) {
	case map[string]interface{}:
		for k, v := range val {
			mapped, err := doTraverse(v, mapper)
			if err != nil {
				return nil, err
			}
			val[k] = mapped
		}
	case []interface{}:
		for i, v := range val {
			mapped, err := doTraverse(v, mapper)
			if err != nil {
				return nil, err
			}
			val[i] = mapped
		}
	default:
		mapped, err := mapper(val)
		if err != nil {
			return nil, err
		}
		return mapped, nil
	}
	return obj, nil
}
