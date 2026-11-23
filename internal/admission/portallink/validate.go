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

package portallink

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission"
	apimportallink "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portallink"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

func validateCreate(ctx context.Context, link *v1alpha1.PortalLink) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	errs.Add(admission.CompileAndValidateTemplate(ctx, link))
	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(validateParentRef(link))
	if errs.IsSevere() {
		return errs
	}

	if !resolveTarget(ctx, link, errs) {
		return errs
	}

	errs.MergeWith(validateDryRun(ctx, link))
	return errs
}

func validateUpdate(ctx context.Context, oldLink, newLink *v1alpha1.PortalLink) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	errs.MergeWith(verifyImmutableFields(oldLink, newLink))
	if errs.IsSevere() {
		return errs
	}

	// validateCreate compiles templates, resolves the parent, and runs the dry-run validation.
	errs.MergeWith(validateCreate(ctx, newLink))
	if errs.IsSevere() {
		return errs
	}

	mergeDriftValidation(ctx, oldLink, newLink, errs)

	return errs
}

func verifyImmutableFields(oldLink, newLink *v1alpha1.PortalLink) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	// A link is parented by exactly one of a portal or an API; switching between
	// the two is a different (and clearer) error than moving it to a different
	// parent of the same kind.
	switch {
	case oldLink.IsPortalLink() && newLink.IsApiLink():
		errs.AddSevere("a portal link cannot be reassigned from a portal to an API")
		return errs
	case oldLink.IsApiLink() && newLink.IsPortalLink():
		errs.AddSevere("a portal link cannot be reassigned from an API to a portal")
		return errs
	}

	if refString(newLink.Spec.Portal) != refString(oldLink.Spec.Portal) {
		errs.AddSeveref(
			"portalRef is immutable; the link cannot be moved to a different portal "+
				"(from [%s] to [%s])",
			refString(oldLink.Spec.Portal), refString(newLink.Spec.Portal),
		)
		return errs
	}
	if refString(newLink.Spec.API) != refString(oldLink.Spec.API) {
		errs.AddSeveref(
			"apiRef is immutable; the link cannot be moved to a different API "+
				"(from [%s] to [%s])",
			refString(oldLink.Spec.API), refString(newLink.Spec.API),
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
func validateParentRef(link *v1alpha1.PortalLink) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	switch {
	case link.IsPortalLink() && link.IsApiLink():
		errs.AddSevere("exactly one of portalRef / apiRef must be set, but both were provided")
		return errs
	case !link.IsPortalLink() && !link.IsApiLink():
		errs.AddSevere("exactly one of portalRef / apiRef must be set, but neither was provided")
		return errs
	}

	if link.IsApiLink() {
		kind := link.Spec.API.Kind
		if kind != "" && dynamic.ResourceFromKind(kind) != core.CRDApiV4DefinitionResource {
			errs.AddSeveref(
				"apiRef [%s] must be of kind ApiV4Definition (next-gen portal only supports those)",
				link.Spec.API.Name,
			)
		}
	}

	return errs
}

// resolveTarget checks that the link's portal or API exists and carries a management context, the one
// the link is synced with. It reports why not as a severe error and returns false.
func resolveTarget(ctx context.Context, link *v1alpha1.PortalLink, errs *errors.AdmissionErrors) bool {
	if link.IsPortalLink() {
		return resolvePortalTarget(ctx, link, errs)
	}
	return resolveApiTarget(ctx, link, errs)
}

func resolvePortalTarget(ctx context.Context, link *v1alpha1.PortalLink, errs *errors.AdmissionErrors) bool {
	prtl, err := dynamic.ResolvePortal(ctx, link.GetPortalRef(), link.GetNamespace())
	if err != nil {
		errs.AddSeveref(
			"portal link [%s] references portal [%v] that can't be resolved",
			link.GetName(), link.GetPortalRef(),
		)
		return false
	}
	if !prtl.HasContext() {
		errs.AddSeveref(
			"referenced portal [%v] has no management context (spec.contextRef)",
			link.GetPortalRef(),
		)
		return false
	}
	return true
}

func resolveApiTarget(ctx context.Context, link *v1alpha1.PortalLink, errs *errors.AdmissionErrors) bool {
	api, err := dynamic.ResolveAPI(ctx, link.GetApiRef(), link.GetNamespace())
	if err != nil {
		errs.AddSeveref(
			"portal link [%s] references API [%v] that can't be resolved",
			link.GetName(), link.GetApiRef(),
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

func validateDryRun(ctx context.Context, link *v1alpha1.PortalLink) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	client, err := apimportallink.CreateAPIMClient(ctx, link)
	if err != nil {
		errs.AddSevere(err.Error())
		return errs
	}
	dto, err := apimportallink.ToDTO(link)
	if err != nil {
		errs.AddSevere(err.Error())
		return errs
	}
	findings, err := apimportallink.DryRun(ctx, client, dto)
	if err != nil {
		errs.AddSevere(err.Error())
		return errs
	}
	errs.MergeWith(findings)
	return errs
}
