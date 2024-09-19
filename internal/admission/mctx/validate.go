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
	"regexp"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	"k8s.io/apimachinery/pkg/runtime"
)

func validateCreate(ctx context.Context, obj runtime.Object) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	// Should be the first validation, it will also compile the templates internally
	tmpErr := admission.CompileAndValidateTemplate(ctx, obj)
	if tmpErr != nil {
		errs.Add(tmpErr)
	}

	if mCtx, ok := obj.(core.ContextObject); ok {
		if !mCtx.HasCloud() || !mCtx.GetCloud().IsEnabled() {
			errs.Add(validateRequiredField(mCtx))
			errs.Add(validateSecretRef(ctx, mCtx))
		}
		errs.Add(validateContextIsAvailable(ctx, mCtx))
	}

	return errs
}

func validateRequiredField(context core.ContextObject) *errors.AdmissionError {
	err := checkEmpty(context.GetURL(), "[baseUrl]")
	if err != nil {
		return err
	}

	if ok, _ := regexp.Match("^http(s?)://.+$", []byte(context.GetURL())); !ok {
		return errors.NewSevere("[baseUrl] is not a valid URL")
	}

	err = checkEmpty(context.GetOrgID(), "[orgId]")
	if err != nil {
		return err
	}

	err = checkEmpty(context.GetEnvID(), "[envId]")
	if err != nil {
		return err
	}

	if !context.HasAuthentication() {
		return errors.NewSevere("[auth] is mandatory when cloud is not enabled")
	}
	return nil
}

func checkEmpty(s string, field string) *errors.AdmissionError {
	if s == "" || strings.TrimSpace(s) == "" {
		return errors.NewSevere(field + " is mandatory when cloud is not enabled")
	}
	return nil
}

func validateSecretRef(ctx context.Context, context core.ContextObject) *errors.AdmissionError {
	if context.HasSecretRef() {
		if err := dynamic.ExpectResolvedSecret(ctx, context.GetSecretRef(), context.GetNamespace()); err != nil {
			return errors.NewSeveref(
				"secret [%v] doesn't exist in the cluster",
				context.GetSecretRef(),
			)
		}
	}
	return nil
}

func validateContextIsAvailable(ctx context.Context, context core.ContextObject) *errors.AdmissionError {
	apim, err := apim.FromContext(ctx, context, context.GetNamespace())
	if err != nil {
		return errors.NewSevere(err.Error())
	}
	_, err = apim.Env.Get()

	if errors.IsNetworkError(err) {
		return errors.NewWarningf(
			"unable to reach APIM, [%s] is not available",
			apim.Context.GetURL(),
		)
	}

	if errors.IsUnauthorized(err) {
		return errors.NewSeveref(
			"bad credentials for context [%s]",
			context.GetName(),
		)
	}

	if errors.IsBadRequest(err) {
		return errors.NewSevere(err.Error())
	}

	return nil
}
