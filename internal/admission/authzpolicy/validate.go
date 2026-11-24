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
	"fmt"
	"slices"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hrid"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

// preCheck validates the HRID derived from the namespace and name.
func preCheck(_ context.Context, policy *v1alpha1.AuthzPolicy) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	if err := hrid.Validate(refs.NewNamespacedNameFromObject(policy).HRID()); err != nil {
		errs.AddSevere(err.Error())
	}
	return errs
}

// resolveRefs checks the pdpRef kinds, resolves the management context and, when every AuthzPdp they name is
// synced, the targets onto their routing scopes. A missing, deleting or unsynced AuthzPdp is a normal state when a
// directory is applied at once: the targets are left as written, and deferDryRun and deferDrift skip the platform
// checks with a warning naming every pending one. The controller waits for them before applying. Any other
// error, an AuthzPdp of another context for instance, refuses the policy.
func resolveRefs(ctx context.Context, policy *v1alpha1.AuthzPolicy, ns string) error {
	if err := validateTargetKinds(policy); err != nil {
		return err
	}
	if err := apim.ResolveContextRef(ctx, policy, ns); err != nil {
		return err
	}
	if err := dynamic.ResolveAuthzPolicyRefs(ctx, policy); err != nil && !dynamic.IsAuthzPdpPending(err) {
		return err
	}
	return nil
}

// validateTargetKinds refuses each pdpRef that names a kind other than AuthzPdp, naming its index. It runs before
// the targets resolve, which replaces the pdpRefs with scopes.
func validateTargetKinds(policy *v1alpha1.AuthzPolicy) error {
	var wrong []string
	for i, target := range policy.Spec.Targets {
		if target.PdpRef == nil {
			continue
		}
		ref := *target.PdpRef
		if ref.Namespace == "" {
			ref.Namespace = policy.GetNamespace()
		}
		if ref.Kind != "" && dynamic.ResourceFromKind(ref.Kind) != core.CRDAuthzPdpResource {
			wrong = append(wrong, fmt.Sprintf("targets[%d].pdpRef [%s]: kind must be AuthzPdp, got [%s]",
				i, ref.String(), ref.Kind))
		}
	}
	if len(wrong) > 0 {
		return fmt.Errorf("%s", strings.Join(wrong, "; "))
	}
	return nil
}

// syncedDuringAdmission is why deferDryRun skips the dry run of targets whose AuthzPdps synced after resolveRefs
// read them.
const syncedDuringAdmission = "platform checks deferred (an AuthzPdp synced during admission)"

// deferDryRun skips the dry run while resolveRefs has left the pdpRefs as written, which it does for a pending
// AuthzPdp only: the platform has no routing scope to check the policy against yet. The AuthzPdps are read again
// to name every pending one. Targets whose AuthzPdps synced since resolveRefs read them still defer: ToDTO cannot
// map an unresolved pdpRef.
func deferDryRun(ctx context.Context, policy *v1alpha1.AuthzPolicy) string {
	if !slices.ContainsFunc(policy.Spec.Targets, func(t authzpdp.Target) bool { return t.PdpRef != nil }) {
		return ""
	}
	if _, err := dynamic.ResolveAuthzTargets(ctx, policy, policy.Spec.Targets); err != nil {
		return fmt.Sprintf("platform checks deferred (%s)", err)
	}
	return syncedDuringAdmission
}

// deferDrift skips drift while the stored policy is not what the platform holds, or cannot be mapped as it:
//   - A pending target gives the dry run's reason, which the API server shows once.
//   - The platform refused the last apply: the stored spec is not what it holds, and this edit is how the policy
//     gets fixed. Drift would compare the platform's copy with both the stored spec and the edit, find both
//     different, and refuse the fix. The dry run has checked the new spec. An out-of-band edit made while the
//     policy is refused goes undetected. The skip ends with the next successful apply, which clears status.errors.
//   - An old target that no longer resolves is usually a deleted AuthzPdp, whose scope the platform has already
//     removed from the policy: comparing the old targets against it would refuse the edit that retargets the
//     policy.
func deferDrift(ctx context.Context, oldObj, newObj *v1alpha1.AuthzPolicy) string {
	if reason := deferDryRun(ctx, newObj); reason != "" {
		return reason
	}
	if len(oldObj.Status.Errors.Severe) > 0 {
		return "drift detection skipped, the platform refused the last apply"
	}
	if _, err := dynamic.ResolveAuthzTargets(ctx, oldObj, oldObj.Spec.Targets); err != nil {
		return "drift detection skipped, the previous targets no longer resolve: " + err.Error()
	}
	return ""
}
