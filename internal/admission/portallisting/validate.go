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

package portallisting

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	"k8s.io/apimachinery/pkg/runtime"
)

func validateCreate(ctx context.Context, obj runtime.Object) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	errs.Add(admission.CompileAndValidateTemplate(ctx, obj))
	if errs.IsSevere() {
		return errs
	}

	listing, ok := obj.(*v1alpha1.PortalListing)
	if !ok {
		errs.AddSevere("can't cast to *v1alpha1.PortalListing")
		return errs
	}

	errs.MergeWith(validateApiKinds(listing))
	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(validateDryRun(ctx, listing))
	return errs
}

func validateUpdate(ctx context.Context, oldObj, newObj runtime.Object) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	newListing, nok := newObj.(*v1alpha1.PortalListing)
	oldListing, ook := oldObj.(*v1alpha1.PortalListing)
	if !nok || !ook {
		errs.AddSevere("can't cast to *v1alpha1.PortalListing")
		return errs
	}

	// Compare user-declared (raw) refs before any template compilation mutates newObj,
	// otherwise an unchanged templated portalRef would look like a change.
	if newListing.Spec.Portal.String() != oldListing.Spec.Portal.String() {
		errs.AddSeveref(
			"portalRef is immutable. Detected change from [%s] to [%s]",
			oldListing.Spec.Portal.String(), newListing.Spec.Portal.String(),
		)
		return errs
	}

	// validateCreate compiles templates and runs kind/dry-run validation.
	return validateCreate(ctx, newObj)
}

// sameContext reports whether two management context refs point to the same object.
// Each ref's empty namespace defaults to its own owning resource's namespace
// (portalNs / apiNs), since a contextRef without a namespace is resolved relative to
// the resource that declares it — not to the listing.
func sameContext(portalCtx core.ObjectRef, portalNs string, apiCtx core.ObjectRef, apiNs string) bool {
	pNs := portalCtx.GetNamespace()
	if pNs == "" {
		pNs = portalNs
	}
	aNs := apiCtx.GetNamespace()
	if aNs == "" {
		aNs = apiNs
	}
	return portalCtx.GetName() == apiCtx.GetName() && pNs == aNs
}

// validateApiKinds enforces that the next-gen portal only lists v4 APIs.
func validateApiKinds(listing *v1alpha1.PortalListing) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	for i := range listing.Spec.APIs {
		kind := listing.Spec.APIs[i].Ref.Kind
		if kind == "" {
			continue // defaults to ApiV4Definition
		}
		if dynamic.ResourceFromKind(kind) != core.CRDApiV4DefinitionResource {
			errs.AddSeveref(
				"API [%s] kind must be ApiV4Definition (the next-gen portal is v4-only), got [%s]",
				listing.Spec.APIs[i].Ref.Name, kind,
			)
		}
	}
	return errs
}

func validateDryRun(ctx context.Context, listing *v1alpha1.PortalListing) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	cp := listing.DeepCopy()
	ns := cp.GetNamespace()

	prtl, err := dynamic.ResolvePortal(ctx, cp.GetPortalRef(), ns)
	if err != nil {
		errs.AddSeveref(
			"portal listing [%s] references portal [%v] that can't be resolved",
			cp.GetName(), cp.GetPortalRef(),
		)
		return errs
	}

	if !prtl.HasContext() {
		errs.AddSeveref(
			"referenced portal [%v] has no management context (spec.contextRef)",
			cp.GetPortalRef(),
		)
		return errs
	}

	portalCtx := prtl.ContextRef()
	for _, apiRef := range cp.GetApiRefs() {
		api, err := dynamic.ResolveAPI(ctx, apiRef, ns)
		if err != nil {
			errs.AddSeveref(
				"portal listing [%s] references API [%v] that can't be resolved",
				cp.GetName(), apiRef,
			)
			continue
		}
		if !api.HasContext() {
			errs.AddSeveref(
				"API [%s] has no management context (spec.contextRef) and cannot be published to a portal",
				apiRef.GetName(),
			)
			continue
		}
		if !sameContext(portalCtx, prtl.GetNamespace(), api.ContextRef(), api.GetNamespace()) {
			errs.AddSeveref(
				"API [%s] management context [%v] must match the portal's management context [%v]",
				apiRef.GetName(), api.ContextRef(), portalCtx,
			)
		}
	}
	if errs.IsSevere() {
		return errs
	}

	apimClient, err := apim.FromContextRef(ctx, prtl.ContextRef(), prtl.GetNamespace())
	if err != nil {
		errs.AddSevere(err.Error())
		return errs
	}

	portalHrid := refs.NewNamespacedNameFromObject(prtl).HRID()

	status, err := apimClient.Listings.DryRunCreateOrUpdate(cp, portalHrid)
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
