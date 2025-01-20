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

package base

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/ctxref"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

func ValidateCreate(ctx context.Context, obj runtime.Object) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	// Should be the first validation, it will also compile the templates internally
	errs.Add(admission.CompileAndValidateTemplate(ctx, obj))

	if errs.IsSevere() {
		return errs
	}

	errs.Add(ctxref.Validate(ctx, obj))

	if api, ok := obj.(core.ApiDefinitionObject); ok {
		errs.Add(validatePlans(api))
		errs.Add(validateNoConflictingPath(ctx, api))
		errs.MergeWith(validateResourceOrRefs(ctx, api))
		errs.MergeWith(validatePages(api))
	}

	return errs
}

func ValidateUpdate(ctx context.Context, obj runtime.Object) *errors.AdmissionErrors {
	return ValidateCreate(ctx, obj)
}
