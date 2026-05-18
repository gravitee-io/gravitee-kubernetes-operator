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

package dictionary

import (
	"context"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/dictionary"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/ctxref"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

func validateCreate(ctx context.Context, obj runtime.Object) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	dict, ok := obj.(*v1alpha1.Dictionary)
	if !ok {
		errs.AddSevere("can't cast to *v1alpha1.Dictionary")
		return errs
	}

	errs.Add(ctxref.Validate(ctx, dict))
	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(validateTypeConsistency(dict))
	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(validateDryRun(ctx, dict))
	return errs
}

func validateTypeConsistency(dict *v1alpha1.Dictionary) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	switch dict.Spec.DictionaryType {
	case dictionary.ManualType:
		if dict.Spec.Dynamic != nil {
			errs.AddSevere("dictionary type is MANUAL but 'dynamic' field is set, use 'manual' instead")
		}
		if dict.Spec.Manual == nil {
			errs.AddSevere("dictionary type is MANUAL but 'manual' field is not set")
		}
	case dictionary.DynamicType:
		if dict.Spec.Manual != nil {
			errs.AddSevere("dictionary type is DYNAMIC but 'manual' field is set, use 'dynamic' instead")
		}
		if dict.Spec.Dynamic == nil {
			errs.AddSevere("dictionary type is DYNAMIC but 'dynamic' field is not set")
		}
	default:
		errs.AddSevere(fmt.Sprintf("unknown dictionary type %q", dict.Spec.DictionaryType))
	}

	return errs
}

func validateDryRun(ctx context.Context, dict *v1alpha1.Dictionary) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	cp := dict.DeepCopy()

	apim, err := apim.FromContextRef(ctx, cp.ContextRef(), cp.GetNamespace())
	if err != nil {
		errs.AddSevere(err.Error())
		return errs
	}

	cp.PopulateIDs(apim.Context, true)

	status, err := apim.Dictionaries.DryRunCreateOrUpdate(cp)
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
