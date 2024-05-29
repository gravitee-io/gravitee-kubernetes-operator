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

package internal

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/indexer"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func getReferences(ctx context.Context, secret *v1.Secret, referenceType client.ObjectList) ([]runtime.Object, error) {
	ref := refs.NamespacedName{
		Namespace: secret.Namespace,
		Name:      secret.Name,
	}

	if err := search.FindByFieldReferencing(
		ctx,
		indexer.SecretRefField,
		ref,
		referenceType,
	); err != nil {
		return nil, err
	}

	items, err := meta.ExtractList(referenceType)
	if err != nil {
		return nil, err
	}

	ref.Namespace = ""

	if err = search.FindByFieldReferencing(
		ctx,
		indexer.SecretRefField,
		ref,
		referenceType,
	); err != nil {
		return nil, err
	}

	currentNSItems, err := meta.ExtractList(referenceType)
	if err != nil {
		return nil, err
	}

	items = append(items, currentNSItems...)

	return items, nil
}
