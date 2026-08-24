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

package portal

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/ctxref"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

func validateCreate(ctx context.Context, prtl *v1alpha1.Portal) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	errs.MergeWith(validateNavigation(prtl))
	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(validateThemeKind(prtl))
	if errs.IsSevere() {
		return errs
	}

	if !prtl.HasContext() {
		errs.AddSevere("a management context reference (spec.contextRef) is required")
		return errs
	}

	errs.Add(ctxref.Validate(ctx, prtl))
	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(validateThemeRef(ctx, prtl))
	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(validateDryRun(ctx, prtl))
	return errs
}

// validateThemeKind rejects a themeRef pointing at anything but a PortalTheme,
// an empty kind meaning PortalTheme.
func validateThemeKind(prtl *v1alpha1.Portal) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	if !prtl.HasThemeRef() {
		return errs
	}

	kind := prtl.Spec.Theme.Kind
	if kind != "" && dynamic.ResourceFromKind(kind) != core.CRDPortalThemeResource {
		errs.AddSeveref(
			"themeRef [%s] is of kind %s, but only PortalTheme is supported",
			prtl.Spec.Theme.Name, kind,
		)
	}

	return errs
}

// validateThemeRef resolves the referenced PortalTheme, which the automation dry run
// echoes back without resolving.
func validateThemeRef(ctx context.Context, prtl *v1alpha1.Portal) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	if !prtl.HasThemeRef() {
		return errs
	}

	if _, err := dynamic.ResolvePortalTheme(ctx, prtl.GetThemeRef(), prtl.GetNamespace()); err != nil {
		errs.AddSeveref(
			"portal [%s] references theme [%v] that can't be resolved",
			prtl.GetName(), prtl.GetThemeRef(),
		)
	}

	return errs
}

func validateNavigation(prtl *v1alpha1.Portal) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	if !prtl.Spec.UsesDeprecatedNavigation() {
		return errs
	}

	if prtl.Spec.Structure != nil {
		errs.AddSevere("navigation and structure cannot be used at the same time")
		return errs
	}

	errs.AddWarning("spec.navigation is deprecated, use spec.structure.topNavbar instead")

	return errs
}

func validateDryRun(ctx context.Context, prtl *v1alpha1.Portal) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	cp := prtl.DeepCopy()

	apimClient, err := apim.FromContextRef(ctx, cp.ContextRef(), cp.GetNamespace())
	if err != nil {
		errs.AddSevere(err.Error())
		return errs
	}

	status, err := apimClient.Portals.DryRunCreateOrUpdate(cp)
	if err != nil {
		errs.AddSevere(err.Error())
		return errs
	}

	for _, severe := range status.Errors.Severe {
		errs.AddSevere(severe)
	}

	for _, warning := range status.Errors.Warning {
		errs.AddWarning(warning)
	}

	return errs
}
