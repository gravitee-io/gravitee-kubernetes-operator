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

package application

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/ctxref"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

func validateCreate(ctx context.Context, obj runtime.Object) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	if app, ok := obj.(core.ApplicationObject); ok {
		// Should be the first validation, it will also compile the templates internally
		errs.Add(admission.CompileAndValidateTemplate(ctx, app))
		if errs.IsSevere() {
			return errs
		}
		errs.Add(ctxref.Validate(ctx, app))
		if errs.IsSevere() {
			return errs
		}
		errs.MergeWith(validateSettings(app))
		if errs.IsSevere() {
			return errs
		}
		errs.MergeWith(validateDryRun(ctx, app))
	}
	return errs
}

func validateUpdate(
	ctx context.Context,
	oldObj runtime.Object,
	newObj runtime.Object,
) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	if errs.IsSevere() {
		return errs
	}
	oldApp, ook := oldObj.(core.ApplicationObject)
	newApp, nok := newObj.(core.ApplicationObject)
	if ook && nok {
		// Should be the first validation, it will also compile the templates internally
		errs.Add(admission.CompileAndValidateTemplate(ctx, newApp))
		if errs.IsSevere() {
			return errs
		}
		errs.Add(ctxref.Validate(ctx, newApp))
		if errs.IsSevere() {
			return errs
		}
		errs.MergeWith(validateSettings(newApp))
		if errs.IsSevere() {
			return errs
		}
		errs.MergeWith(validateSettingsUpdate(oldApp, newApp))
		if errs.IsSevere() {
			return errs
		}
		errs.MergeWith(validateDryRun(ctx, newApp))
	}
	return errs
}

func validateSettings(app core.ApplicationObject) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	model := app.GetModel()

	settings := model.GetSettings()
	if settings.IsOAuth() && settings.IsSimple() {
		errs.AddSevere("configuring both OAuth and simple settings is not allowed")
	}

	return errs
}

func validateSettingsUpdate(oldApp, newApp core.ApplicationObject) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	oldModel, newModel := oldApp.GetModel(), newApp.GetModel()
	oldSettings, newSettings := oldModel.GetSettings(), newModel.GetSettings()

	if oldSettings.IsOAuth() && newSettings.IsSimple() {
		errs.AddSevere("moving from OAuth to simple settings is not allowed")
	} else if oldSettings.IsSimple() && newSettings.IsOAuth() {
		errs.AddSevere("moving from simple to Oauth settings is not allowed")
	}

	if newSettings.GetOAuthType() != oldSettings.GetOAuthType() {
		errs.AddSevere("updating OAuth application type is not allowed")
	}

	return errs
}

func validateDryRun(ctx context.Context, app core.ApplicationObject) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	cp, _ := app.DeepCopyObject().(core.ApplicationObject)

	apim, err := apim.FromContextRef(ctx, cp.ContextRef(), cp.GetNamespace())
	if err != nil {
		errs.AddSevere(err.Error())
	}

	cp.PopulateIDs(apim.Context)

	impl, ok := cp.GetModel().(*application.Application)
	if !ok {
		errs.AddSeveref("unable to call dry run (unknown type %T)", impl)
	}

	status, err := apim.Applications.DryRunCreateOrUpdate(impl)
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
