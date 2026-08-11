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
	"k8s.io/apimachinery/pkg/runtime"
)

func validateCreate(ctx context.Context, obj runtime.Object) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	errs.Add(admission.CompileAndValidateTemplate(ctx, obj))
	if errs.IsSevere() {
		return errs
	}

	link, ok := obj.(*v1alpha1.PortalLink)
	if !ok {
		errs.AddSevere("can't cast to *v1alpha1.PortalLink")
		return errs
	}

	errs.MergeWith(validateDryRun(ctx, link))
	return errs
}

func validateUpdate(ctx context.Context, oldObj, newObj runtime.Object) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	newLink, nok := newObj.(*v1alpha1.PortalLink)
	oldLink, ook := oldObj.(*v1alpha1.PortalLink)
	if !nok || !ook {
		errs.AddSevere("can't cast to *v1alpha1.PortalLink")
		return errs
	}

	// Compare user-declared (raw) refs before any template compilation mutates newObj,
	// otherwise an unchanged templated portalRef would look like a change.
	if newLink.Spec.Portal.String() != oldLink.Spec.Portal.String() {
		errs.AddSeveref(
			"portalRef is immutable. Detected change from [%s] to [%s]",
			oldLink.Spec.Portal.String(), newLink.Spec.Portal.String(),
		)
		return errs
	}

	// validateCreate compiles templates and runs the dry-run validation.
	return validateCreate(ctx, newObj)
}

func validateDryRun(ctx context.Context, link *v1alpha1.PortalLink) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	cp := link.DeepCopy()
	ns := cp.GetNamespace()

	prtl, err := dynamic.ResolvePortal(ctx, cp.GetPortalRef(), ns)
	if err != nil {
		errs.AddSeveref(
			"portal link [%s] references portal [%v] that can't be resolved",
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

	apimClient, err := apim.FromContextRef(ctx, prtl.ContextRef(), prtl.GetNamespace())
	if err != nil {
		errs.AddSevere(err.Error())
		return errs
	}

	status, err := apimClient.Links.DryRunCreateOrUpdate(cp, prtl)
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
