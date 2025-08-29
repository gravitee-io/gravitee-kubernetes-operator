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
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type ListOptions struct {
	Namespace string
	Excluded  []core.ObjectRef
}

func ResolveAPI(ctx context.Context, ref core.ObjectRef, parentNs string) (core.ApiDefinitionObject, error) {
	refKind := ref.GetKind()
	if ref.GetKind() == "" {
		refKind = ApiV4GVR.Resource
	}
	kind := ResourceFromKind(refKind)
	switch kind {
	case ApiGVR.Resource:
		return resolveRef(ctx, ref, parentNs, ApiGVR, new(v1alpha1.ApiDefinition))
	case ApiV4GVR.Resource:
		return resolveRef(ctx, ref, parentNs, ApiV4GVR, new(v1alpha1.ApiV4Definition))
	default:
		return nil, errors.NewSevere("API definition kind is mandatory")
	}
}

func GetAPIs(ctx context.Context, opts ListOptions) ([]core.ApiDefinitionObject, error) {
	v2Apis, err := getV2Apis(ctx, opts)
	if err != nil {
		return nil, err
	}
	v4Apis, err := getV4Apis(ctx, opts)
	if err != nil {
		return nil, err
	}
	apis := make([]core.ApiDefinitionObject, 0)
	apis = append(apis, v2Apis...)
	apis = append(apis, v4Apis...)
	return apis, nil
}

func getV2Apis(ctx context.Context, opts ListOptions) ([]core.ApiDefinitionObject, error) {
	resource := getAPIsResource(ApiGVR, opts.Namespace)
	list, err := resource.List(ctx, metav1.ListOptions{})
	apis := make([]core.ApiDefinitionObject, 0)
	if err != nil {
		return nil, err
	}
	for _, item := range list.Items {
		if isExcluded(item, opts) {
			continue
		}
		api := new(v1alpha1.ApiDefinition)
		api, err := convert(item.Object, api)
		if err != nil {
			return nil, err
		}
		apis = append(apis, api)
	}
	return apis, nil
}

func getV4Apis(ctx context.Context, opts ListOptions) ([]core.ApiDefinitionObject, error) {
	resource := getAPIsResource(ApiV4GVR, opts.Namespace)
	list, err := resource.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	apis := make([]core.ApiDefinitionObject, 0)
	for _, item := range list.Items {
		if isExcluded(item, opts) {
			continue
		}
		api := new(v1alpha1.ApiV4Definition)
		api, err := convert(item.Object, api)
		if err != nil {
			return nil, err
		}
		apis = append(apis, api)
	}
	return apis, nil
}

func isExcluded(item unstructured.Unstructured, opts ListOptions) bool {
	metadata, _ := item.Object["metadata"].(map[string]interface{})
	ns := metadata["namespace"]
	name := metadata["name"]
	for _, excluded := range opts.Excluded {
		if excluded.GetName() == name && excluded.GetNamespace() == ns {
			return true
		}
	}
	return false
}

func getAPIsResource(gvr schema.GroupVersionResource, ns string) dynamic.ResourceInterface {
	return getResource(gvr, ns)
}

func getResource(gvr schema.GroupVersionResource, ns string) dynamic.ResourceInterface {
	if ns == "" {
		return GetClient().Resource(gvr)
	}
	return GetClient().Resource(gvr).Namespace(ns)
}
