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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/ctxref"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

func validateCreate(ctx context.Context, obj runtime.Object) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	if spg, ok := obj.(core.SharedPolicyGroupObject); ok {
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
	_, ook := oldObj.(core.SharedPolicyGroupObject)
	newSpg, nok := newObj.(core.SharedPolicyGroupObject)
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
		errs.MergeWith(validateDryRun(ctx, newSpg))
	}
	return errs
}

func validateDryRun(ctx context.Context, spg core.SharedPolicyGroupObject) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	cp, _ := spg.DeepCopyObject().(core.SharedPolicyGroupObject)

	apim, err := apim.FromContextRef(ctx, cp.ContextRef(), cp.GetNamespace())
	if err != nil {
		errs.AddSevere(err.Error())
	}

	cp.PopulateIDs(apim.Context)

	impl, ok := cp.GetSpec().(*v1alpha1.SharedPolicyGroupSpec)
	if !ok {
		errs.AddSeveref("unable to call dry run (unknown type %T)", impl)
	}

	status, err := apim.SharedPolicyGroup.DryRunCreateOrUpdate(impl.SharedPolicyGroup)
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

func validateDelete(_ context.Context, _ runtime.Object) *errors.AdmissionErrors {
	return errors.NewAdmissionErrors()
}
