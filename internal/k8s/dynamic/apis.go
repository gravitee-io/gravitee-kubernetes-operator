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

	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type ListOptions struct {
	Namespace string
	Excluded  []core.ResourceRef
}

func GetAPIs(ctx context.Context, opts ListOptions) ([]core.ApiDefinition, error) {
	v2Apis, err := getV2Apis(ctx, opts)
	if err != nil {
		return nil, err
	}
	v4Apis, err := getV4Apis(ctx, opts)
	if err != nil {
		return nil, err
	}
	apis := make([]core.ApiDefinition, 0)
	apis = append(apis, v2Apis...)
	apis = append(apis, v4Apis...)
	return apis, nil
}

func getV2Apis(ctx context.Context, opts ListOptions) ([]core.ApiDefinition, error) {
	resource := getResource(ApiGVR, opts.Namespace)
	list, err := resource.List(ctx, metav1.ListOptions{})
	apis := make([]core.ApiDefinition, 0)
	if err != nil {
		return nil, err
	}
	for _, item := range list.Items {
		if isExcluded(item, opts) {
			continue
		}
		api := new(v2.Api)
		api, err := convert(item.Object["spec"], api)
		if err != nil {
			return nil, err
		}
		apis = append(apis, api)
	}
	return apis, nil
}

func getV4Apis(ctx context.Context, opts ListOptions) ([]core.ApiDefinition, error) {
	resource := getResource(ApiV4GVR, opts.Namespace)
	list, err := resource.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	apis := make([]core.ApiDefinition, 0)
	for _, item := range list.Items {
		if isExcluded(item, opts) {
			continue
		}
		api := new(v4.Api)
		api, err := convert(item.Object["spec"], api)
		if err != nil {
			return nil, err
		}
		apis = append(apis, api)
	}
	return apis, nil
}

func isExcluded(item unstructured.Unstructured, opts ListOptions) bool {
	ns := item.Object["metadata"].(map[string]interface{})["namespace"]
	name := item.Object["metadata"].(map[string]interface{})["name"]
	for _, excluded := range opts.Excluded {
		if excluded.GetName() == name && excluded.GetNamespace() == ns {
			return true
		}
	}
	return false
}

func getResource(gvr schema.GroupVersionResource, ns string) dynamic.ResourceInterface {
	if env.Config.CheckApiContextPathConflictInCluster {
		return GetClient().Resource(gvr)
	}
	return GetClient().Resource(gvr).Namespace(ns)
}
