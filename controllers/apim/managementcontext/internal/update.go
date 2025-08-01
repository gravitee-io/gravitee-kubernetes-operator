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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	v1 "k8s.io/api/core/v1"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func CreateOrUpdate(
	ctx context.Context,
	instance *v1alpha1.ManagementContext,
) error {
	if instance.HasSecretRef() {
		secret := &v1.Secret{}

		nsn := getSecretRef(instance)
		if err := k8s.GetClient().Get(ctx, nsn, secret); err != nil {
			return gerrors.NewResolveRefError(err)
		}

		if !util.ContainsFinalizer(secret, core.ManagementContextSecretFinalizer) {
			util.AddFinalizer(secret, core.ManagementContextSecretFinalizer)
			return k8s.Update(ctx, secret)
		}
	}

	return nil
}
