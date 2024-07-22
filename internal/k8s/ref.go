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

package k8s

import (
	"context"
	"encoding/json"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/management"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/types/k8s/custom"
	coreV1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func ResolveContext(ctx context.Context, ref custom.ResourceRef) (*management.Context, error) {
	return resolveRef(ctx, ref, schema.GroupVersionResource{
		Group:    keys.CrdGroup,
		Version:  keys.CrdVersion,
		Resource: "managementcontexts",
	}, new(management.Context))
}

func ResolveSecret(ctx context.Context, ref custom.ResourceRef) (*coreV1.Secret, error) {
	return resolveRef(ctx, ref, schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "secrets",
	}, new(coreV1.Secret))
}

func resolveRef[T any](
	ctx context.Context,
	ref custom.ResourceRef,
	gvr schema.GroupVersionResource,
	target T,
) (T, error) {
	dynamic, err := GetDynamicClient().
		Resource(gvr).
		Namespace(ref.GetNamespace()).
		Get(ctx, ref.GetName(), v1.GetOptions{})

	if err != nil {
		return target, err
	}

	spec := dynamic.Object["spec"]

	b, err := json.Marshal(spec)
	if err != nil {
		return target, err
	}

	if err = json.Unmarshal(b, target); err != nil {
		return target, err
	}

	return target, nil
}
