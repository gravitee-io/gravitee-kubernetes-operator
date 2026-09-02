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

package portaltheme

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/ctxref"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
)

// validateDelete keeps a theme that portals still activate from being deleted, so that
// the rejection is reported to whoever runs the delete rather than to the reconcile loop.
func validateDelete(ctx context.Context, thm *v1alpha1.PortalTheme) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	if err := search.AssertNoPortalThemeRef(ctx, thm); err != nil {
		errs.AddSevere(err.Error())
	}

	return errs
}

func validateCreate(ctx context.Context, thm *v1alpha1.PortalTheme) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	if !thm.HasContext() {
		errs.AddSevere("a management context reference (spec.contextRef) is required")
		return errs
	}

	errs.Add(ctxref.Validate(ctx, thm))
	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(validateDryRun(ctx, thm))
	return errs
}

func validateDryRun(ctx context.Context, thm *v1alpha1.PortalTheme) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	cp := thm.DeepCopy()

	apimClient, err := apim.FromContextRef(ctx, cp.ContextRef(), cp.GetNamespace())
	if err != nil {
		errs.AddSevere(err.Error())
		return errs
	}

	status, err := apimClient.PortalThemes.DryRunCreateOrUpdate(cp)
	if err != nil {
		errs.AddSevere(err.Error())
		return errs
	}

	errs.MergeWith(errors.NewAdmissionErrorsFromStatus(status.Errors))

	return errs
}
