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

// Delete withdraws the schema from its PDP gateway and deletes it, by hrid: the target does not
// need resolving.
func Delete(ctx context.Context, schema *v1alpha1.AuthzSchema) error {
	if !util.ContainsFinalizer(schema, core.AuthzSchemaFinalizer) {
		return nil
	}

	apimClient, err := apim.FromContextRef(ctx, schema.ContextRef(), schema.GetNamespace())
	if err != nil {
		// ManagementContext already gone: APIM is unreachable, let the finalizer be removed.
		if apierrors.IsNotFound(err) {
			return nil
		}
		return err
	}

	hrid := refs.NewNamespacedNameFromObject(schema).HRID()
	if err := gerrors.IgnoreNotFound(apimClient.AuthzSchemas.Delete(hrid)); err != nil {
		return gerrors.NewControlPlaneError(err)
	}

	return nil
}
