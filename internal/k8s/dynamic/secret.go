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

	coreV1 "k8s.io/api/core/v1"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func ExpectResolvedSecret(ctx context.Context, ref core.ResourceRef, parentNs string) error {
	if _, err := ResolveSecret(ctx, ref, parentNs); err != nil {
		return err
	}
	return nil
}

func ResolveSecret(ctx context.Context, ref core.ResourceRef, parentNs string) (*coreV1.Secret, error) {
	return resolveRef(ctx, ref, parentNs, schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "secrets",
	}, new(coreV1.Secret))
}
