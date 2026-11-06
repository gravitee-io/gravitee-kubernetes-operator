// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Delete deletes the entity by hrid: the targets do not need resolving. The platform deletes its
// RESOURCE children and detaches its PRINCIPAL children. Over the cascade limit it answers 413 and
// deletes nothing: the control plane error is retried with backoff, so the finalizer stays and
// Accepted names the reason until the cascade fits.
func Delete(ctx context.Context, entity *v1alpha1.AuthzEntity) error {
	if !util.ContainsFinalizer(entity, core.AuthzEntityFinalizer) {
		return nil
	}

	apimClient, err := apim.FromContextRef(ctx, entity.ContextRef(), entity.GetNamespace())
	if err != nil {
		// ManagementContext already gone: APIM is unreachable, let the finalizer be removed.
		if apierrors.IsNotFound(err) {
			return nil
		}
		return err
	}

	hrid := refs.NewNamespacedNameFromObject(entity).HRID()
	if err := gerrors.IgnoreNotFound(apimClient.AuthzEntities.Delete(hrid)); err != nil {
		return gerrors.NewControlPlaneError(err)
	}

	return nil
}
