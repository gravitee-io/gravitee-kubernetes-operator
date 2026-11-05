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

package authzpdp

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/ctxref"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hrid"
)

// validateCreate runs the checks in order: the context, the derived HRID, then the platform dry
// run, which refuses a routing scope another PDP holds, a strict value that disagrees with the
// engine and an endpoint path another API uses. Schema and CEL own the spec-only rules (the
// targetPdpId and tag grammar, their immutability).
func validateCreate(ctx context.Context, pdp *v1alpha1.AuthzPdp) *errors.AdmissionErrors {
	errs := validateContext(ctx, pdp)
	if errs.IsSevere() {
		return errs
	}

	if err := hrid.Validate(refs.NewNamespacedNameFromObject(pdp).HRID()); err != nil {
		errs.AddSevere(err.Error())
		return errs
	}

	errs.MergeWith(validateDryRun(ctx, pdp))
	return errs
}

// validateContext requires a context that resolves. Resolving also compiles templates, so a dry
// run that follows carries resolved values.
func validateContext(ctx context.Context, pdp *v1alpha1.AuthzPdp) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	if !pdp.HasContext() {
		errs.AddSevere("a management context reference (spec.contextRef) is required")
		return errs
	}

	errs.Add(ctxref.Validate(ctx, pdp))
	return errs
}

func validateDryRun(ctx context.Context, pdp *v1alpha1.AuthzPdp) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	apimClient, err := apim.FromContextRef(ctx, pdp.ContextRef(), pdp.GetNamespace())
	if err != nil {
		errs.AddSevere(err.Error())
		return errs
	}

	// A dry run answers 200 with its findings; a refusal body is only expected from a real apply,
	// but is read the same way should the platform answer one.
	state, err := apimClient.AuthzPdps.DryRunCreateOrUpdate(model.ToAuthzPdpDTO(pdp))
	if err != nil {
		errs.MergeWith(refusalToAdmissionErrors(err))
		return errs
	}

	errs.MergeWith(errors.NewAdmissionErrorsFromStatus(state.Errors))

	return errs
}

// refusalToAdmissionErrors reports a platform refusal finding by finding; anything else, the 403
// of a strict change without the permission included, is reported as the request error.
func refusalToAdmissionErrors(err error) *errors.AdmissionErrors {
	if findings, refused := model.AutomationRefusal(err); refused {
		return errors.NewAdmissionErrorsFromStatus(findings)
	}

	errs := errors.NewAdmissionErrors()
	errs.AddSevere(err.Error())
	return errs
}
