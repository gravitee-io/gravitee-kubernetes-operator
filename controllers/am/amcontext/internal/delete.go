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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func Delete(
	ctx context.Context,
	instance *v1alpha1.AMContext,
) error {
	if !util.ContainsFinalizer(instance, core.AMContextFinalizer) {
		return nil
	}

	// No AssertNoContextRef: that helper only knows APIM kinds.
	// AM resources that reference an AMContext land in AM-7506.

	if instance.HasSecretRef() {
		secret := &v1.Secret{}

		nsn := getSecretRef(instance)
		if err := k8s.GetClient().Get(ctx, nsn, secret); err != nil {
			if apierrors.IsNotFound(err) {
				return nil
			}
			return gerrors.NewResolveRefError(err)
		}

		isRef, err := hasMoreReferences(ctx, *instance.Spec.Auth.SecretRef)
		if err != nil {
			return err
		}

		if !isRef {
			util.RemoveFinalizer(secret, core.AMContextSecretFinalizer)
		}

		if err := k8s.GetClient().Update(ctx, secret); err != nil {
			return err
		}
	}

	return nil
}
