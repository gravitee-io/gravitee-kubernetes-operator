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

<<<<<<< HEAD:controllers/apim/secrets/internal/update.go
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
=======
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
>>>>>>> 539e666 (fix: remove secret controller):controllers/apim/managementcontext/internal/update.go
	v1 "k8s.io/api/core/v1"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func CreateOrUpdate(
	ctx context.Context,
	instance *v1alpha1.ManagementContext,
) error {
	if instance.HasSecretRef() {
		secret := &v1.Secret{}

		nsn := getSecretRef(instance)
		if err := k8s.GetClient().Get(ctx, nsn, secret); err != nil {
			return err
		}

		if !util.ContainsFinalizer(secret, core.ManagementContextSecretFinalizer) {
			util.AddFinalizer(secret, core.ManagementContextSecretFinalizer)
			return k8s.GetClient().Update(ctx, secret)
		}
	}

<<<<<<< HEAD:controllers/apim/secrets/internal/update.go
	if len(contextRefs) == 0 {
		return nil
	}

	if !util.ContainsFinalizer(secret, keys.ManagementContextSecretFinalizer) {
		log.FromContext(ctx).Info("secret is used by some management context, adding finalizer")
		util.AddFinalizer(secret, keys.ManagementContextSecretFinalizer)
	}
	k8s.AddAnnotation(secret, keys.LastSpecHash, hash.Calculate(&secret.Data))

=======
>>>>>>> 539e666 (fix: remove secret controller):controllers/apim/managementcontext/internal/update.go
	return nil
}
