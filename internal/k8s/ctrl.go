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
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const defaultRequeueAfter = 5 * time.Second

func RequeueError(err error) (ctrl.Result, error) {
	return ctrl.Result{RequeueAfter: defaultRequeueAfter}, err
}

func CreateOrUpdate(ctx context.Context, obj client.Object, fns ...util.MutateFn) error {
	if _, err := util.CreateOrUpdate(ctx, GetClient(), obj, func() error {
		for _, f := range fns {
			if e := f(); e != nil {
				return e
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
