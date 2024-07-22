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

package mctx

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/types/k8s/custom"
	"k8s.io/apimachinery/pkg/runtime"
)

func validate(ctx context.Context, obj runtime.Object) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	if context, ok := obj.(custom.ContextResource); ok {
		errs.Add(validateSecretRef(ctx, context))
		errs.Add(validateContextIsAvailable(ctx, context))
	}

	return errs
}

func validateSecretRef(ctx context.Context, context custom.ContextResource) *errors.AdmissionError {
	if context.HasSecretRef() {
		if err := dynamic.ExpectResolvedSecret(ctx, context.GetSecretRef(), context.GetNamespace()); err != nil {
			return errors.NewSevere(
				"secret [%v] doesn't exist in the cluster",
				context.GetSecretRef(),
			)
		}
	}
	return nil
}

func validateContextIsAvailable(ctx context.Context, context custom.ContextResource) *errors.AdmissionError {
	apim, err := apim.FromContext(ctx, context, context.GetNamespace())
	if err != nil {
		return errors.NewSevere(err.Error())
	}
	_, err = apim.Env.Get()

	if err != nil {
		return errors.NewWarning(
			"unable to reach APIM, [%s] is not available",
			apim.Context.GetURL(),
		)
	}

	if errors.IsUnauthorized(err) {
		return errors.NewSevere(
			"bad credentials for context [%s]",
			context.GetName(),
		)
	}

	if errors.IsNotFound(err) {
		return errors.NewSevere(
			"environment [%s/%s] could not be found in APIM",
			apim.Context.GetOrg(), apim.Context.GetEnv(),
		)
	}

	return nil
}
