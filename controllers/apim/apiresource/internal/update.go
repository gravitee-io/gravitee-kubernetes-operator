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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1beta1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func CreateOrUpdate(
	ctx context.Context,
	k8s client.Client,
	instance *v1beta1.ApiResource,
) error {
	if !util.ContainsFinalizer(instance, keys.ApiResourceFinalizer) {
		util.AddFinalizer(instance, keys.ApiResourceFinalizer)

		if err := k8s.Update(ctx, instance); err != nil {
			return err
		}
	}

	return nil
}
