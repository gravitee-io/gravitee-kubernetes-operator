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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/ctxref"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hrid"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

// validateCreate runs the checks in order: the context, the derived HRID, the pdpRef kinds, then
// the platform dry run with the resolved targets, which refuses an entityId another policy holds,
// a text that does not parse, a template and a violation of a strict gateway's schema, and warns
// about an MCP policy the default engine never sees and a tool named by its catalog id. The CRD
// schema and CEL own the spec-only rules (the pdpRef/scope union, the entityId and scope grammars,
// immutability).
func validateCreate(ctx context.Context, policy *v1alpha1.AuthzPolicy) *errors.AdmissionErrors {
	errs := validateContext(ctx, policy)
	if errs.IsSevere() {
		return errs
	}

	if err := hrid.Validate(refs.NewNamespacedNameFromObject(policy).HRID()); err != nil {
		errs.AddSevere(err.Error())
		return errs
	}

	errs.MergeWith(validateTargetKinds(policy))
	if errs.IsSevere() {
		return errs
	}

	targetErrs, resolved := resolveTargets(ctx, policy)
	errs.MergeWith(targetErrs)
	if !resolved {
		return errs
	}

	errs.MergeWith(validateDryRun(ctx, policy))
	return errs
}

// validateContext requires a context that resolves. Resolving also compiles templates, so a dry
// run that follows carries resolved values.
func validateContext(ctx context.Context, policy *v1alpha1.AuthzPolicy) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	if !policy.HasContext() {
		errs.AddSevere("a management context reference (spec.contextRef) is required")
		return errs
	}

	errs.Add(ctxref.Validate(ctx, policy))
	return errs
}

// validateTargetKinds refuses a pdpRef that names a kind other than AuthzPdp, naming its index.
func validateTargetKinds(policy *v1alpha1.AuthzPolicy) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	for i, target := range policy.Spec.Targets {
		if target.PdpRef == nil {
			continue
		}
		ref := *target.PdpRef
		if ref.Namespace == "" {
			ref.Namespace = policy.GetNamespace()
		}
		if ref.Kind != "" && dynamic.ResourceFromKind(ref.Kind) != core.CRDAuthzPdpResource {
			errs.AddSeveref("targets[%d].pdpRef [%s]: kind must be AuthzPdp, got [%s]", i, ref.String(), ref.Kind)
		}
	}
	return errs
}

// resolveTargets replaces the policy's targets with the scopes the dry run sends, and tells whether
// they resolved. An AuthzPdp that does not exist or has not synced yet is a normal state when a
// directory is applied at once: it is a warning naming every pending one, and the dry run is
// skipped. The controller waits for them before applying. Any other error, an AuthzPdp of another
// context for instance, is severe.
func resolveTargets(ctx context.Context, policy *v1alpha1.AuthzPolicy) (*errors.AdmissionErrors, bool) {
	errs := errors.NewAdmissionErrors()

	err := dynamic.ResolveAuthzPolicyRefs(ctx, policy)
	switch {
	case err == nil:
		return errs, true
	case dynamic.IsAuthzPdpPending(err):
		errs.AddWarningf("%s: the policy is checked against the platform once every AuthzPdp it targets has synced", err)
	default:
		errs.AddSevere(err.Error())
	}
	return errs, false
}

func validateDryRun(ctx context.Context, policy *v1alpha1.AuthzPolicy) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	apimClient, err := apim.FromContextRef(ctx, policy.ContextRef(), policy.GetNamespace())
	if err != nil {
		errs.AddSevere(err.Error())
		return errs
	}

	// A dry run answers 200 with its findings; a refusal body is only expected from a real apply,
	// but is read the same way should the platform answer one.
	state, err := apimClient.AuthzPolicies.DryRunCreateOrUpdate(model.ToAuthzPolicyDTO(policy))
	if err != nil {
		errs.MergeWith(refusalToAdmissionErrors(err))
		return errs
	}

	errs.MergeWith(errors.NewAdmissionErrorsFromStatus(state.Errors))

	return errs
}

// refusalToAdmissionErrors reports a platform refusal finding by finding; anything else is
// reported as the request error.
func refusalToAdmissionErrors(err error) *errors.AdmissionErrors {
	if findings, refused := model.AutomationRefusal(err); refused {
		return errors.NewAdmissionErrorsFromStatus(findings)
	}

	errs := errors.NewAdmissionErrors()
	errs.AddSevere(err.Error())
	return errs
}
