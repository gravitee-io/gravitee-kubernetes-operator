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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func Update(ctx context.Context, secret *v1.Secret) error {
	return ensureContextFinalizer(ctx, secret)
}

func ensureContextFinalizer(ctx context.Context, secret *v1.Secret) error {
	k8s := k8s.GetClient()

	contextRefs, err := getReferences(ctx, secret, new(v1beta1.ManagementContextList))

	if err != nil {
		return err
	}

	if len(contextRefs) == 0 {
		return nil
	}

	if !controllerutil.ContainsFinalizer(secret, keys.ManagementContextSecretFinalizer) {
		log.FromContext(ctx).Info("secret is used by some management context, adding finalizer")
		controllerutil.AddFinalizer(secret, keys.ManagementContextSecretFinalizer)
		return k8s.Update(ctx, secret)
	}

	return nil
}
