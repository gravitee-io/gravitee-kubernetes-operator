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

package docs

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission"
	apimdocs "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/docs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

func validateCreate(ctx context.Context, doc *v1alpha1.Documentation) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	errs.Add(admission.CompileAndValidateTemplate(ctx, doc))
	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(validateParentRef(doc))
	if errs.IsSevere() {
		return errs
	}

	if !resolveTarget(ctx, doc, errs) {
		return errs
	}

	errs.MergeWith(validateDryRun(ctx, doc))
	return errs
}

func validateUpdate(ctx context.Context, oldDoc, newDoc *v1alpha1.Documentation) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	errs.Add(admission.CompileAndValidateTemplate(ctx, newDoc))
	if errs.IsSevere() {
		return errs
	}

	// validate
	errs.MergeWith(verifyImmutableFields(oldDoc, newDoc))
	if errs.IsSevere() {
		return errs
	}

	// validateCreate compiles templates and runs ref/dry-run validation.
	errs.MergeWith(validateCreate(ctx, newDoc))
	if errs.IsSevere() {
		return errs
	}

	mergeDriftValidation(ctx, oldDoc, newDoc, errs)
	return errs
}

func verifyImmutableFields(oldDoc *v1alpha1.Documentation, newDoc *v1alpha1.Documentation) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	// A documentation is parented by exactly one of a portal or an API; switching
	// between the two is a different (and clearer) error than moving it to a
	// different parent of the same kind.
	switch {
	case oldDoc.IsPortalDoc() && newDoc.IsApiDoc():
		errs.AddSevere("a documentation cannot be reassigned from a portal to an API")
		return errs
	case oldDoc.IsApiDoc() && newDoc.IsPortalDoc():
		errs.AddSevere("a documentation cannot be reassigned from an API to a portal")
		return errs
	}

	if refString(newDoc.Spec.Portal) != refString(oldDoc.Spec.Portal) {
		errs.AddSeveref(
			"portalRef is immutable; documentation cannot be moved to a different portal "+
				"(from [%s] to [%s])",
			refString(oldDoc.Spec.Portal), refString(newDoc.Spec.Portal),
		)
		return errs
	}

	if refString(newDoc.Spec.API) != refString(oldDoc.Spec.API) {
		errs.AddSeveref(
			"apiRef is immutable; documentation cannot be moved to a different API "+
				"(from [%s] to [%s])",
			refString(oldDoc.Spec.API), refString(newDoc.Spec.API),
		)
		return errs
	}
	return errs
}

func refString(ref *refs.NamespacedName) string {
	if ref == nil {
		return ""
	}
	return ref.String()
}

// validateParentRef enforces that exactly one of portalRef / apiRef is set, and
// that an apiRef points to a v4 API (the next-gen portal is v4-only).
func validateParentRef(doc *v1alpha1.Documentation) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	switch {
	case doc.IsPortalDoc() && doc.IsApiDoc():
		errs.AddSevere("exactly one of portalRef / apiRef must be set, but both were provided")
		return errs
	case !doc.IsPortalDoc() && !doc.IsApiDoc():
		errs.AddSevere("exactly one of portalRef / apiRef must be set, but neither was provided")
		return errs
	}

	if doc.IsApiDoc() {
		kind := doc.Spec.API.Kind
		if kind != "" && dynamic.ResourceFromKind(kind) != core.CRDApiV4DefinitionResource {
			errs.AddSeveref(
				"apiRef [%s] must be of kind ApiV4Definition (next-gen portal only supports those)",
				doc.Spec.API.Name,
			)
		}
	}

	return errs
}

func validateDryRun(ctx context.Context, doc *v1alpha1.Documentation) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	client, err := apimdocs.CreateAPIMClient(ctx, doc)
	if err != nil {
		errs.AddSevere(err.Error())
		return errs
	}
	dto, err := apimdocs.ToDTO(doc)
	if err != nil {
		errs.AddSevere(err.Error())
		return errs
	}
	findings, err := apimdocs.DryRun(ctx, client, dto)
	if err != nil {
		errs.AddSevere(err.Error())
		return errs
	}
	errs.MergeWith(findings)
	return errs
}

// resolveTarget checks that the page's portal or API exists and carries a management context, the one
// the page is synced with. It reports why not as a severe error and returns false.
func resolveTarget(ctx context.Context, doc *v1alpha1.Documentation, errs *errors.AdmissionErrors) bool {
	if doc.IsPortalDoc() {
		return resolvePortalTarget(ctx, doc, errs)
	}
	return resolveApiTarget(ctx, doc, errs)
}

func resolvePortalTarget(ctx context.Context, doc *v1alpha1.Documentation, errs *errors.AdmissionErrors) bool {
	prtl, err := dynamic.ResolvePortal(ctx, doc.GetPortalRef(), doc.GetNamespace())
	if err != nil {
		errs.AddSeveref(
			"documentation [%s] references portal [%v] that can't be resolved",
			doc.GetName(), doc.GetPortalRef(),
		)
		return false
	}
	if !prtl.HasContext() {
		errs.AddSeveref(
			"referenced portal [%v] has no management context (spec.contextRef)",
			doc.GetPortalRef(),
		)
		return false
	}
	return true
}

func resolveApiTarget(ctx context.Context, doc *v1alpha1.Documentation, errs *errors.AdmissionErrors) bool {
	api, err := dynamic.ResolveAPI(ctx, doc.GetApiRef(), doc.GetNamespace())
	if err != nil {
		errs.AddSeveref(
			"documentation [%s] references API [%v] that can't be resolved",
			doc.GetName(), doc.GetApiRef(),
		)
		return false
	}
	if !api.HasContext() {
		errs.AddSeveref(
			"referenced API [%s] has no management context (spec.contextRef)",
			api.GetName(),
		)
		return false
	}
	return true
}
