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

// CreateOrUpdate applies and deploys the policy with the targets the caller resolved, in
// declaration order: each scope as written, each pdpRef as the routing scope of the AuthzPdp it
// names.
func CreateOrUpdate(ctx context.Context, policy *v1alpha1.AuthzPolicy, targets []string) error {
	if !policy.HasContext() {
		return gerrors.NewIllegalStateError(
			fmt.Errorf("authz policy [%s] has no management context", policy.GetName()),
		)
	}

	apimClient, err := apim.FromContextRef(ctx, policy.ContextRef(), policy.GetNamespace())
	if err != nil {
		return err
	}

	state, err := apimClient.AuthzPolicies.CreateOrUpdate(model.ToAuthzPolicyDTO(policy, targets))
	if err != nil {
		// A 400 is how the platform refuses an apply: an entityId another policy holds, a text that
		// does not parse, a violation of a strict gateway's schema. Admission skips the dry run while
		// a pdpRef is pending, and another policy may take the entityId since: keep it recoverable
		// (backoff) rather than terminal until the next spec change.
		if gerrors.IsBadRequest(err) {
			if findings, refused := model.AutomationRefusal(err); refused {
				policy.Status.Errors = findings
			}
			return fmt.Errorf("authz policy [%s] was refused by the automation api: %w", policy.GetName(), err)
		}
		return gerrors.NewControlPlaneError(err)
	}

	k8s.AddAutomationAPIManagedCondition(policy)

	// Setting fields by fields to keep the rest (conditions) intact
	policy.Status.ID = state.ID
	policy.Status.OrgID = state.OrgID
	policy.Status.EnvID = state.EnvID
	policy.Status.HRID = refs.NewNamespacedNameFromObject(policy).HRID()
	policy.Status.PolicyStatus = state.PolicyStatus
	policy.Status.Errors = state.Errors

	return nil
}
