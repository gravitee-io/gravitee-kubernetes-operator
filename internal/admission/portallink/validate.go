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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

func validateCreate(ctx context.Context, link *v1alpha1.PortalLink) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	errs.Add(admission.CompileAndValidateTemplate(ctx, link))
	if errs.IsSevere() {
		return errs
	}

	portal := validatePortal(ctx, link, errs)
	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(validateDryRun(ctx, link, portal))
	return errs
}

func validateUpdate(ctx context.Context, oldLink, newLink *v1alpha1.PortalLink) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	// after the portal is resolved, validate that the portalRef hasn't changed.
	if newLink.Spec.Portal.String() != oldLink.Spec.Portal.String() {
		errs.AddSeveref(
			"portalRef is immutable. Detected change from [%s] to [%s]",
			oldLink.Spec.Portal.String(), newLink.Spec.Portal.String(),
		)
		return errs
	}

	// validateCreate compiles templates, resolve portal and runs the dry-run validation.
	errs.MergeWith(validateCreate(ctx, newLink))
	if errs.IsSevere() {
		return errs
	}

	mergeDriftValidation(ctx, oldLink, newLink, errs)

	return errs
}

func validateDryRun(ctx context.Context, link *v1alpha1.PortalLink, portal *v1alpha1.Portal) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	cp := link.DeepCopy()

	apimClient, err := apim.FromContextRef(ctx, portal.ContextRef(), portal.GetNamespace())
	if err != nil {
		errs.AddSevere(err.Error())
		return errs
	}

	status, err := apimClient.Links.DryRunCreateOrUpdate(cp, portal)
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

func validatePortal(ctx context.Context, link *v1alpha1.PortalLink, errs *errors.AdmissionErrors) *v1alpha1.Portal {
	ns := link.GetNamespace()
	portal, err := dynamic.ResolvePortal(ctx, link.GetPortalRef(), ns)
	if err != nil {
		errs.AddSeveref(
			"portal link [%s] references portal [%v] that can't be resolved",
			link.GetName(), link.GetPortalRef(),
		)
		return nil
	}

	if !portal.HasContext() {
		errs.AddSeveref(
			"referenced portal [%v] has no management context (spec.contextRef)",
			link.GetPortalRef(),
		)
		return nil
	}
	return portal
}
