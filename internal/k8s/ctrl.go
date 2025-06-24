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
	"fmt"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const defaultRequeueAfter = 5 * time.Second

func RequeueError(err error) (ctrl.Result, error) {
	return ctrl.Result{RequeueAfter: defaultRequeueAfter}, err
}

func CreateOrUpdate(ctx context.Context, obj client.Object, fns ...util.MutateFn) error {
	key := client.ObjectKeyFromObject(obj)
	if err := GetClient().Get(ctx, key, obj); err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
		for _, f := range fns {
			if err := mutate(f, key, obj); err != nil {
				return err
			}
		}

		if err := GetClient().Create(ctx, obj); err != nil {
			return err
		}
		return nil
	}

	for _, f := range fns {
		if err := mutate(f, key, obj); err != nil {
			return err
		}
	}

	if err := Update(ctx, obj); err != nil {
		return err
	}
	return nil
}

func mutate(f util.MutateFn, key client.ObjectKey, obj client.Object) error {
	if err := f(); err != nil {
		return err
	}
	if newKey := client.ObjectKeyFromObject(obj); key != newKey {
		return fmt.Errorf("MutateFn cannot mutate object name and/or object namespace")
	}
	return nil
}
