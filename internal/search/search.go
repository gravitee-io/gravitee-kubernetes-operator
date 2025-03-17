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

package search

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/fields"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func FindByFieldReferencing(
	ctx context.Context,
	field IndexField,
	ref refs.NamespacedName,
	result client.ObjectList,
) error {
	filter := &client.ListOptions{
		FieldSelector: fields.SelectorFromSet(fields.Set{field.String(): ref.String()}),
	}

	return k8s.GetClient().List(ctx, result, filter)
}
