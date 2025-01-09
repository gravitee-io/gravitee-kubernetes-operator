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

package policygroups

import (
	"context"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/ctxref"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

func validateCreate(ctx context.Context, obj runtime.Object) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	if spg, ok := obj.(*v1alpha1.SharedPolicyGroup); ok {
		// Should be the first validation, it will also compile the templates internally
		errs.Add(admission.CompileAndValidateTemplate(ctx, spg))
		if errs.IsSevere() {
			return errs
		}
		errs.Add(ctxref.Validate(ctx, spg))
		if errs.IsSevere() {
			return errs
		}
		errs.MergeWith(validateDryRun(ctx, spg))
	}
	return errs
}

func validateUpdate(
	ctx context.Context,
	oldObj runtime.Object,
	newObj runtime.Object,
) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	oldSpg, ook := oldObj.(*v1alpha1.SharedPolicyGroup)
	newSpg, nok := newObj.(*v1alpha1.SharedPolicyGroup)
	if ook && nok {
		// Should be the first validation, it will also compile the templates internally
		errs.Add(admission.CompileAndValidateTemplate(ctx, newSpg))
		if errs.IsSevere() {
			return errs
		}

		errs.Add(ctxref.Validate(ctx, newSpg))
		if errs.IsSevere() {
			return errs
		}

		errs.Add(validateImmutableFields(ctx, oldSpg, newSpg))
		if errs.IsSevere() {
			return errs
		}

		errs.MergeWith(validateDryRun(ctx, newSpg))
	}
	return errs
}

func validateImmutableFields(_ context.Context, oldSpg, newSpg *v1alpha1.SharedPolicyGroup) *errors.AdmissionError {
	if utils.ToStringValue(oldSpg.Spec.CrossID) != utils.ToStringValue(newSpg.Spec.CrossID) {
		return &errors.AdmissionError{
			Severity: errors.Severe,
			Message: fmt.Sprintf("can not change Shared Policy Group [%s] CrossID, "+
				"once it is created", oldSpg.Status.CrossID),
		}
	}

	if oldSpg.Spec.ApiType != newSpg.Spec.ApiType {
		return &errors.AdmissionError{
			Severity: errors.Severe,
			Message: fmt.Sprintf("can not change Shared Policy Group [%s] ApiType [%s], "+
				"once it is created", oldSpg.Status.CrossID, oldSpg.Spec.ApiType),
		}
	}

	if *oldSpg.Spec.Phase != *newSpg.Spec.Phase {
		return &errors.AdmissionError{
			Severity: errors.Severe,
			Message: fmt.Sprintf("can not change Shared Policy Group [%s] Phase [%s], "+
				"once it is created", oldSpg.Status.CrossID, *oldSpg.Spec.Phase),
		}
	}

	return nil
}

func validateDryRun(ctx context.Context, spg runtime.Object) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	cp, _ := spg.DeepCopyObject().(*v1alpha1.SharedPolicyGroup)

	apim, err := apim.FromContextRef(ctx, cp.ContextRef(), cp.GetNamespace())
	if err != nil {
		errs.AddSevere(err.Error())
	}

	cp.PopulateIDs(apim.Context)

	status, err := apim.SharedPolicyGroup.DryRunCreateOrUpdate(cp.Spec.SharedPolicyGroup)
	if err != nil {
		errs.AddSevere(err.Error())
		return errs
	}
	for _, severe := range status.Errors.Severe {
		errs.AddSevere(severe)
	}
	if errs.IsSevere() {
		return errs
	}
	for _, warning := range status.Errors.Warning {
		errs.AddWarning(warning)
	}
	return errs
}

func validateDelete(ctx context.Context, obj runtime.Object) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	spg, _ := obj.(*v1alpha1.SharedPolicyGroup)

	if err := search.AssertNoSharedPolicyGroupRef(ctx, spg); err != nil {
		errs.AddSevere(err.Error())
	}

	return errs
}
