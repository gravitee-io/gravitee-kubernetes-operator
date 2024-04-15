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

package manager

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var cli client.Client

func Client() client.Client {
	return cli
}

func GetLatest[T client.Object](obj T) error {
	key := types.NamespacedName{
		Namespace: obj.GetNamespace(),
		Name:      obj.GetName(),
	}

	if err := Client().Get(context.Background(), key, obj); err != nil {
		return err
	}

	return nil
}

func Delete[T client.Object](obj T) error {
	err := GetLatest(obj)
	if err != nil {
		return err
	}
	return Client().Delete(context.Background(), obj)
}

func UpdateSafely[T client.Object](objNew T) error {
	key := types.NamespacedName{
		Namespace: objNew.GetNamespace(),
		Name:      objNew.GetName(),
	}

	objLast, ok := objNew.DeepCopyObject().(T)
	if !ok {
		return fmt.Errorf("failed to copy object %v", objNew)
	}

	ctx := context.Background()

	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		if err := Client().Get(ctx, key, objLast); err != nil {
			return err
		}

		objNew.SetResourceVersion(objLast.GetResourceVersion())
		objNew.SetGeneration(objLast.GetGeneration())

		return Client().Update(ctx, objNew)
	})
}
