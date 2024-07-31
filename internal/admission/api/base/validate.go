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

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/ctxref"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/types/k8s/custom"
	"k8s.io/apimachinery/pkg/runtime"
)

func ValidateCreate(ctx context.Context, obj runtime.Object) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	errs.Add(ctxref.Validate(ctx, obj))

	if api, ok := obj.(custom.ApiDefinitionResource); ok {
		errs.Add(validateNoConflictingPath(ctx, api))
	}

	return errs
}

func ValidateUpdate(ctx context.Context, obj runtime.Object) *errors.AdmissionErrors {
	return ValidateCreate(ctx, obj)
}
