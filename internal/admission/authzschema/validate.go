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

package authzschema

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

// validateCreate runs the checks in order: the context, the derived HRID, the pdpRef kind, then
// the platform dry run with the resolved target, which refuses a target another schema holds
// (AUTHZ-6) and a document that does not parse. Schema and CEL own the spec-only rules (the
// pdpRef/scope union, the scope grammar, no *).
func validateCreate(ctx context.Context, schema *v1alpha1.AuthzSchema) *errors.AdmissionErrors {
	errs := validateContext(ctx, schema)
	if errs.IsSevere() {
		return errs
	}

	if err := hrid.Validate(refs.NewNamespacedNameFromObject(schema).HRID()); err != nil {
		errs.AddSevere(err.Error())
		return errs
	}

	errs.MergeWith(validateTargetKind(schema))
	if errs.IsSevere() {
		return errs
	}

	targetErrs, resolved := resolveTarget(ctx, schema)
	errs.MergeWith(targetErrs)
	if !resolved {
		return errs
	}

	errs.MergeWith(validateDryRun(ctx, schema))
	return errs
}

// validateContext requires a context that resolves. Resolving also compiles templates, so a dry
// run that follows carries resolved values.
func validateContext(ctx context.Context, schema *v1alpha1.AuthzSchema) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	if !schema.HasContext() {
		errs.AddSevere("a management context reference (spec.contextRef) is required")
		return errs
	}

	errs.Add(ctxref.Validate(ctx, schema))
	return errs
}

// validateTargetKind refuses a pdpRef that names a kind other than AuthzPdp.
func validateTargetKind(schema *v1alpha1.AuthzSchema) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	target := schema.Spec.Target
	if target == nil || target.PdpRef == nil {
		return errs
	}

	ref := *target.PdpRef
	if ref.Namespace == "" {
		ref.Namespace = schema.GetNamespace()
	}
	if ref.Kind != "" && dynamic.ResourceFromKind(ref.Kind) != core.CRDAuthzPdpResource {
		errs.AddSeveref("target.pdpRef [%s]: kind must be AuthzPdp, got [%s]", ref.String(), ref.Kind)
	}
	return errs
}

// resolveTarget replaces the schema's target with the scope the dry run sends, and tells whether it
// resolved. An AuthzPdp that does not exist or has not synced yet is a normal state when a
// directory is applied at once: it is a warning, and the dry run is skipped. The controller waits
// for the AuthzPdp before applying. Any other error, an AuthzPdp of another context for instance,
// is severe.
func resolveTarget(ctx context.Context, schema *v1alpha1.AuthzSchema) (*errors.AdmissionErrors, bool) {
	errs := errors.NewAdmissionErrors()

	err := dynamic.ResolveAuthzSchemaRefs(ctx, schema)
	switch {
	case err == nil:
		return errs, true
	case dynamic.IsAuthzPdpPending(err):
		errs.AddWarningf("%s: the schema is checked against the platform once it is", err)
	default:
		errs.AddSevere(err.Error())
	}
	return errs, false
}

func validateDryRun(ctx context.Context, schema *v1alpha1.AuthzSchema) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	apimClient, err := apim.FromContextRef(ctx, schema.ContextRef(), schema.GetNamespace())
	if err != nil {
		errs.AddSevere(err.Error())
		return errs
	}

	// A dry run answers 200 with its findings; a refusal body is only expected from a real apply,
	// but is read the same way should the platform answer one.
	state, err := apimClient.AuthzSchemas.DryRunCreateOrUpdate(model.ToAuthzSchemaDTO(schema))
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
