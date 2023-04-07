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
	"fmt"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	"golang.org/x/net/context"
	coreV1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func CreateOrUpdate(
	ctx context.Context,
	k8s client.Client,
	instance *gio.ManagementContext,
) error {
	// We only add a finalizer to our ManagementContexts to keep track of their deletion
	if !util.ContainsFinalizer(instance, keys.ManagementContextFinalizer) {
		util.AddFinalizer(instance, keys.ManagementContextFinalizer)

		if err := k8s.Update(ctx, instance); err != nil {
			err = fmt.Errorf("an error occurred while adding finalizer to the management context: %w", err)
			return err
		}
	}

	spec := instance.Spec
	if spec.HasSecretRef() {
		secret := new(coreV1.Secret)
		key := spec.SecretRef()
		key.Namespace = getSecretNamespace(instance)
		ns := key.ToK8sType()
		if err := k8s.Get(ctx, ns, secret); err != nil {
			return err
		}
		if !util.ContainsFinalizer(secret, keys.ManagementContextSecretFinalizer) {
			util.AddFinalizer(secret, keys.ManagementContextSecretFinalizer)
			return k8s.Update(ctx, secret)
		}
	}
	return nil
}

func getSecretNamespace(context *gio.ManagementContext) string {
	secretRef := context.Spec.SecretRef()
	if secretRef.Namespace != "" {
		return secretRef.Namespace
	}
	return context.Namespace
}
