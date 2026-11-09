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
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

// CreateOrUpdate applies an entity whose refs dynamic.ResolveAuthzEntityRefs has resolved.
func CreateOrUpdate(ctx context.Context, entity *v1alpha1.AuthzEntity) error {
	if !entity.HasContext() {
		return gerrors.NewIllegalStateError(
			fmt.Errorf("authz entity [%s] has no management context", entity.GetName()),
		)
	}

	apimClient, err := apim.FromContextRef(ctx, entity.ContextRef(), entity.GetNamespace())
	if err != nil {
		return err
	}

	state, err := apimClient.AuthzEntities.CreateOrUpdate(model.ToAuthzEntityDTO(entity))
	if err != nil {
		// A 400 is how the platform refuses an apply: an entityId another source holds (AUTHZ-5), an
		// attribute off a strict schema. Admission skips the dry run while a pdpRef is pending, and
		// another source may take the entityId since: keep it recoverable (backoff) rather than
		// terminal until the next spec change.
		if gerrors.IsBadRequest(err) {
			if findings, refused := model.AutomationRefusal(err); refused {
				entity.Status.Errors = findings
			}
			return fmt.Errorf("authz entity [%s] was refused by the automation api: %w", entity.GetName(), err)
		}
		return gerrors.NewControlPlaneError(err)
	}

	k8s.AddAutomationAPIManagedCondition(entity)

	// Setting fields by fields to keep the rest (conditions) intact
	entity.Status.ID = state.ID
	entity.Status.OrgID = state.OrgID
	entity.Status.EnvID = state.EnvID
	entity.Status.HRID = refs.NewNamespacedNameFromObject(entity).HRID()
	entity.Status.Errors = state.Errors

	return nil
}
