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

package authzpolicy

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/template"
)

func validateUpdate(
	ctx context.Context,
	oldObj *v1alpha1.AuthzPolicy,
	newObj *v1alpha1.AuthzPolicy,
) *errors.AdmissionErrors {
	// A metadata-only update, the controller's own finalizer and annotation writes included,
	// does not need the platform: only a spec change is dry-run. Hashed before validateContext,
	// which compiles templates into newObj while oldObj keeps them: hashed after, a templated spec
	// would always look changed.
	specChanged := oldObj.Spec.Hash() != newObj.Spec.Hash()

	errs := validateContext(ctx, newObj)
	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(validateTargetKinds(newObj))
	if errs.IsSevere() {
		return errs
	}

	// Resolved on every update, metadata-only ones included: drift resolves them too, and would
	// refuse the controller's own writes while an AuthzPdp is pending.
	targetErrs, resolved := resolveTargets(ctx, newObj)
	errs.MergeWith(targetErrs)
	if !resolved {
		return errs
	}

	if specChanged {
		errs.MergeWith(validateDryRun(ctx, newObj))
		if errs.IsSevere() {
			return errs
		}
	}

	// A policy the platform never received, one waiting on an AuthzPdp for instance, has nothing
	// to drift from: comparing would apply the remote-missing policy (deny by default) and refuse the
	// very edit that unblocks it.
	if oldObj.Status.ID == "" {
		return errs
	}

	// An old target that no longer resolves is usually a deleted AuthzPdp, whose scope the platform
	// has already removed from the policy: comparing the old targets against it would refuse the
	// edit that retargets the policy.
	if _, err := dynamic.ResolveAuthzTargets(ctx, oldObj, oldObj.Spec.Targets); err != nil {
		return errs
	}

	errs.MergeWith(drift.ValidateDrift(ctx, oldObj, newObj, resolveRefs, getRemoteAuthzPolicy,
		drift.MapDTO(model.ToAuthzPolicyDTO)))
	return errs
}

// resolveRefs compiles templates on the copy drift compares, then resolves its refs as the
// reconcile does, so both sides are compared with what the operator sends.
func resolveRefs(ctx context.Context, policy *v1alpha1.AuthzPolicy) error {
	if err := template.Compile(ctx, policy, false); err != nil {
		return err
	}
	return dynamic.ResolveAuthzPolicyRefs(ctx, policy)
}

// getRemoteAuthzPolicy reads the platform's copy and returns the payload half. policyStatus is left
// out: a policy disabled in the Console is not drift, and the next apply deploys it again with the
// platform's warning.
func getRemoteAuthzPolicy(apimClient *apim.APIM, policy *v1alpha1.AuthzPolicy) (any, error) {
	hrid := refs.NewNamespacedNameFromObject(policy).HRID()
	remote, err := apimClient.AuthzPolicies.GetByHRID(hrid)
	if err != nil {
		return nil, err
	}
	return remote.AuthzPolicyDTO, nil
}
