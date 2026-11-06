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

// CreateOrUpdate applies the schema with the target the caller resolved: the scope as written, or
// the routing scope of the AuthzPdp its pdpRef names.
func CreateOrUpdate(ctx context.Context, schema *v1alpha1.AuthzSchema, target string) error {
	if !schema.HasContext() {
		return gerrors.NewIllegalStateError(
			fmt.Errorf("authz schema [%s] has no management context", schema.GetName()),
		)
	}

	apimClient, err := apim.FromContextRef(ctx, schema.ContextRef(), schema.GetNamespace())
	if err != nil {
		return err
	}

	state, err := apimClient.AuthzSchemas.CreateOrUpdate(model.ToAuthzSchemaDTO(schema, target))
	if err != nil {
		// A 400 is how the platform refuses an apply: a target another schema holds (AUTHZ-6), a
		// document that does not parse. Admission skips the dry run while the pdpRef is pending, and
		// another schema may take the target since: keep it recoverable (backoff) rather than
		// terminal until the next spec change. The target is left as it was.
		if gerrors.IsBadRequest(err) {
			if findings, refused := model.AutomationRefusal(err); refused {
				schema.Status.Errors = findings
			}
			return fmt.Errorf("authz schema [%s] was refused by the automation api: %w", schema.GetName(), err)
		}
		return gerrors.NewControlPlaneError(err)
	}

	k8s.AddAutomationAPIManagedCondition(schema)

	// Setting fields by fields to keep the rest (conditions) intact
	schema.Status.ID = state.ID
	schema.Status.OrgID = state.OrgID
	schema.Status.EnvID = state.EnvID
	schema.Status.HRID = refs.NewNamespacedNameFromObject(schema).HRID()
	schema.Status.Errors = state.Errors
	schema.Status.Target = state.Target

	return nil
}
