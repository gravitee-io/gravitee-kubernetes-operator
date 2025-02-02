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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/indexer"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/types"
)

func getSecretRef(instance *v1alpha1.ManagementContext) types.NamespacedName {
	nsn := instance.GetSecretRef().NamespacedName()

	if nsn.Namespace == "" {
		nsn.Namespace = instance.Namespace
	}

	return nsn
}

func isReferenced(ctx context.Context, ref refs.NamespacedName) (bool, error) {
	list := new(v1alpha1.ManagementContextList)
	if err := search.FindByFieldReferencing(
		ctx,
		indexer.SecretRefField,
		ref,
		list,
	); err != nil {
		return false, err
	}

	items, err := meta.ExtractList(list)
	if err != nil {
		return false, err
	}

	ref.Namespace = ""

	if err = search.FindByFieldReferencing(
		ctx,
		indexer.SecretRefField,
		ref,
		list,
	); err != nil {
		return false, err
	}

	currentNSItems, err := meta.ExtractList(list)
	if err != nil {
		return false, err
	}

	items = append(items, currentNSItems...)

	return len(items) > 0, nil
}
