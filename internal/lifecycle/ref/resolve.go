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

package ref

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

func ObjectKey(ref *refs.NamespacedName, parentNs string) types.NamespacedName {
	ns := parentNs
	if ref != nil && ref.GetNamespace() != "" {
		ns = ref.GetNamespace()
	}
	name := ""
	if ref != nil {
		name = ref.GetName()
	}
	return types.NamespacedName{Namespace: ns, Name: name}
}

func Resolve(ctx context.Context, kindName string, nsn types.NamespacedName) (client.Object, error) {
	kind, ok := Lookup(kindName)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownKind, kindName)
	}
	obj := kind.New()
	if err := k8s.GetClient().Get(ctx, nsn, obj); err != nil {
		return nil, err
	}
	return obj, nil
}
