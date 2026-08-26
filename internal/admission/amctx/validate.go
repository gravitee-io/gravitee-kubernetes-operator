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

package amctx

import (
	"context"
	"regexp"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

func validateCreate(ctx context.Context, obj *v1alpha1.AMContext) *gerrors.AdmissionErrors {
	errs := gerrors.NewAdmissionErrors()

	errs.Add(admission.CompileAndValidateTemplate(ctx, obj))
	if errs.IsSevere() {
		return errs
	}

	errs.Add(validateRequiredField(obj))
	if errs.IsSevere() {
		return errs
	}

	errs.Add(validateSecretRef(ctx, obj))
	if errs.IsSevere() {
		return errs
	}

	errs.Add(validateContextIsAvailable(ctx, obj))
	return errs
}

func validateRequiredField(obj *v1alpha1.AMContext) *gerrors.AdmissionError {
	if obj.Spec.Context == nil {
		return gerrors.NewSevere("[baseUrl] is mandatory")
	}

	if err := checkEmpty(obj.GetURL(), "[baseUrl]"); err != nil {
		return err
	}

	if ok, _ := regexp.Match("^http(s?)://.+$", []byte(obj.GetURL())); !ok {
		return gerrors.NewSevere("[baseUrl] is not a valid URL")
	}

	if err := checkEmpty(obj.GetOrgID(), "[orgId]"); err != nil {
		return err
	}

	if err := checkEmpty(obj.GetEnvID(), "[envId]"); err != nil {
		return err
	}

	if !obj.HasAuthentication() {
		return gerrors.NewSevere("[auth] is mandatory")
	}
	if obj.Spec.Auth.BearerToken == "" && !obj.HasSecretRef() {
		return gerrors.NewSevere("[auth] bearerToken or secretRef is mandatory")
	}
	return nil
}

func checkEmpty(s, field string) *gerrors.AdmissionError {
	if s == "" || strings.TrimSpace(s) == "" {
		return gerrors.NewSevere(field + " is mandatory")
	}
	return nil
}

func validateSecretRef(ctx context.Context, obj *v1alpha1.AMContext) *gerrors.AdmissionError {
	if !obj.HasSecretRef() {
		return nil
	}
	if err := dynamic.ExpectResolvedSecret(ctx, obj.GetSecretRef(), obj.GetNamespace()); err != nil {
		return gerrors.NewSeveref(
			"secret [%v] doesn't exist in the cluster",
			obj.GetSecretRef(),
		)
	}
	return nil
}

func validateContextIsAvailable(ctx context.Context, obj *v1alpha1.AMContext) *gerrors.AdmissionError {
	client, err := am.FromContext(ctx, obj, obj.GetNamespace())
	if err != nil {
		return gerrors.NewSevere(err.Error())
	}

	err = client.Domains.Probe()
	if gerrors.IsNetworkError(err) {
		return gerrors.NewWarningf(
			"unable to reach AM, [%s] is not available",
			obj.GetURL(),
		)
	}
	if err != nil {
		return gerrors.NewSevere(err.Error())
	}
	return nil
}
