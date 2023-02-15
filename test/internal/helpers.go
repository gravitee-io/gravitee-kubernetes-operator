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
	"fmt"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewAssertionError(field string, expected, given any) error {
	return fmt.Errorf("expected %s to be %v, got %v", field, expected, given)
}

func UpdateSafely[T client.Object](client client.Client, objNew T) error {
	key := types.NamespacedName{
		Namespace: objNew.GetNamespace(),
		Name:      objNew.GetName(),
	}

	objLast, ok := objNew.DeepCopyObject().(T)
	if !ok {
		return fmt.Errorf("failed to copy object %v", objNew)
	}

	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		if err := client.Get(context.Background(), key, objLast); err != nil {
			return err
		}

		objNew.SetResourceVersion(objLast.GetResourceVersion())

		return client.Update(context.Background(), objNew)
	})
}
