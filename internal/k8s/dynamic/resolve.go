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

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func resolveRefSpec[T any](
	ctx context.Context,
	ref core.ObjectRef,
	parentNs string,
	gvr schema.GroupVersionResource,
	target T,
) (T, error) {
	dynamic, err := resolveDynamic(ctx, ref, parentNs, gvr)
	if err != nil {
		return target, err
	}
	return convert(dynamic.Object["spec"], target)
}

func resolveRef[T any](
	ctx context.Context,
	ref core.ObjectRef,
	parentNs string,
	gvr schema.GroupVersionResource,
	target T,
) (T, error) {
	dynamic, err := resolveDynamic(ctx, ref, parentNs, gvr)
	if err != nil {
		return target, err
	}
	return convert(dynamic.Object, target)
}

func resolveDynamic(
	ctx context.Context,
	ref core.ObjectRef,
	parentNs string,
	gvr schema.GroupVersionResource,
) (*unstructured.Unstructured, error) {
	if ref.GetNamespace() == "" {
		ref.SetNamespace(parentNs)
	}

	return GetClient().
		Resource(gvr).
		Namespace(ref.GetNamespace()).
		Get(ctx, ref.GetName(), v1.GetOptions{})
}
