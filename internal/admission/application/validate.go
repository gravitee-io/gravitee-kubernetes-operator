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
	"encoding/pem"
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/ctxref"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	appResolve "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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
		errs.MergeWith(validateSettings(ctx, app))
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
	oldApp, ook := oldObj.(core.ApplicationObject)
	newApp, nok := newObj.(core.ApplicationObject)
	if !ook || !nok {
		return errs
	}

	if newApp.IsBeingDeleted() {
		return errs
	}

	// Should be the first validation, it will also compile the templates internally
	errs.Add(admission.CompileAndValidateTemplate(ctx, newApp))
	if errs.IsSevere() {
		return errs
	}
	errs.Add(ctxref.Validate(ctx, newApp))
	if errs.IsSevere() {
		return errs
	}
	errs.MergeWith(validateSettings(ctx, newApp))
	if errs.IsSevere() {
		return errs
	}
	errs.MergeWith(validateSettingsUpdate(oldApp, newApp))
	if errs.IsSevere() {
		return errs
	}
	errs.MergeWith(validateCertEndDatesVsSubscriptionEndDates(ctx, newApp))
	if errs.IsSevere() {
		return errs
	}
	errs.MergeWith(validateDryRun(ctx, newApp))
	return errs
}

func validateSettings(ctx context.Context, app core.ApplicationObject) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	model := app.GetModel()

	settings := model.GetSettings()
	if settings.IsOAuth() && settings.IsSimple() {
		errs.AddSevere("configuring both OAuth and simple settings is not allowed")
	}

	hasSingleCert := settings.HasTLS() && settings.GetClientCertificate() != ""
	hasMultipleCerts := settings.HasClientCertificates()

	if hasSingleCert && hasMultipleCerts {
		errs.AddSevere("clientCertificate and clientCertificates cannot be used at the same time")
		return errs
	}

	if hasSingleCert {
		errs.Add(validateSingleClientCertificate(settings.GetClientCertificate()))
	}

	if hasMultipleCerts {
		appSettings, ok := settings.(*application.Setting)
		if !ok {
			return errs
		}
		errs.MergeWith(validateClientCertificates(appSettings.GetClientCertificates()))
		if errs.IsSevere() {
			return errs
		}
		if err := appResolve.ResolveClientCertificates(ctx, appSettings, app.GetNamespace(), app.GetName()); err != nil {
			errs.AddSevere(err.Error())
		}
	}

	return errs
}

func validateSingleClientCertificate(cert string) *errors.AdmissionError {
	if b, _ := pem.Decode([]byte(cert)); b == nil {
		return errors.NewSevere("failed to parse TLS client certificate")
	}
	return errors.NewWarning("clientCertificate is deprecated, use clientCertificates instead")
}

func validateClientCertificates(certs []application.ClientCertificate) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	for i, cert := range certs {
		hasContent := cert.Content != ""
		hasRef := cert.Ref != nil

		if hasContent && hasRef {
			errs.AddSeveref("clientCertificates[%d]: content and ref cannot both be set", i)
		}
		if !hasContent && !hasRef {
			errs.AddSeveref("clientCertificates[%d]: either content or ref must be set", i)
		}
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

	cp, _ := app.DeepCopyObject().(*v1alpha1.Application)

	apim, err := apim.FromContextRef(ctx, cp.ContextRef(), cp.GetNamespace())
	if err != nil {
		errs.AddSevere(err.Error())
	}

	cp.PopulateIDs(apim.Context, k8s.IsAutomationAPIManaged(app))

	impl, ok := cp.GetModel().(*application.Application)
	if !ok {
		errs.AddSeveref("unable to call dry run (unknown type %T)", impl)
	}

	status, err := apim.Applications.DryRunCreateOrUpdate(cp)
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

func validateCertEndDatesVsSubscriptionEndDates(
	ctx context.Context,
	app core.ApplicationObject,
) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	settings, ok := app.GetModel().GetSettings().(*application.Setting)
	if !ok || !settings.HasClientCertificates() {
		return errs
	}

	maxCertEnd := maxCertificateEndDate(settings.GetClientCertificates())
	if maxCertEnd == nil {
		return errs
	}

	appRef := refs.NewNamespacedName(app.GetNamespace(), app.GetName())
	subList := &v1alpha1.SubscriptionList{}
	if err := search.FindByFieldReferencing(ctx, search.AppSubsField, appRef, subList); err != nil {
		errs.AddSevere(err.Error())
		return errs
	}

	for i := range subList.Items {
		sub := &subList.Items[i]
		endingAt := sub.Spec.GetEndingAt()
		if endingAt == nil {
			continue
		}
		subEnd, err := time.Parse(time.RFC3339, *endingAt)
		if err != nil {
			continue
		}
		if !subEnd.After(*maxCertEnd) {
			continue
		}
		api := fetchAPI(ctx, sub)
		if api == nil {
			continue
		}
		plan := api.GetPlan(sub.Spec.GetPlan())
		if plan == nil || plan.GetSecurityType() != "MTLS" {
			continue
		}
		errs.AddSeveref(
			"subscription [%s/%s] ending date [%s] is after all client certificate end dates in application [%s]",
			sub.GetNamespace(), sub.GetName(), *endingAt, app.GetRef(),
		)
	}

	return errs
}

func maxCertificateEndDate(certs []application.ClientCertificate) *time.Time {
	var maxEnd *time.Time
	for _, cert := range certs {
		if cert.EndsAt == "" {
			return nil // at least one cert has no end date, no constraint
		}
		certEnd, err := time.Parse(time.RFC3339, cert.EndsAt)
		if err != nil {
			continue
		}
		if maxEnd == nil || certEnd.After(*maxEnd) {
			maxEnd = &certEnd
		}
	}
	return maxEnd
}

func fetchAPI(ctx context.Context, sub *v1alpha1.Subscription) core.ApiDefinitionModel {
	apiRef := sub.Spec.GetApiRef()
	kind := apiRef.GetKind()
	if kind == "" {
		kind = core.CRDApiV4DefinitionResource
	}
	kind = dynamic.ResourceFromKind(kind)

	ns := apiRef.GetNamespace()
	if ns == "" {
		ns = sub.GetNamespace()
	}
	key := types.NamespacedName{Name: apiRef.GetName(), Namespace: ns}

	switch kind {
	case core.CRDApiV4DefinitionResource:
		api := &v1alpha1.ApiV4Definition{}
		if err := k8s.GetClient().Get(ctx, key, api); err != nil {
			return nil
		}
		return api
	case core.CRDApiDefinitionResource:
		api := &v1alpha1.ApiDefinition{}
		if err := k8s.GetClient().Get(ctx, key, api); err != nil {
			return nil
		}
		return api
	default:
		return nil
	}
}

func validateDelete(_ context.Context, obj runtime.Object) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	if app, ok := obj.(core.ApplicationObject); ok {
		errs.Add(validateSubscriptionCount(app))
	}
	return errs
}

func validateSubscriptionCount(app core.ApplicationObject) *errors.AdmissionError {
	st, _ := app.GetStatus().(core.SubscribableStatus)
	sc := st.GetSubscriptionCount()
	if sc > 0 {
		return errors.NewSeveref(
			"cannot delete [%s] because it is referenced in %d subscriptions. "+
				"Subscriptions must be deleted before the application. "+
				"You can review the subscriptions using the following command: "+
				"kubectl get subscriptions.gravitee.io -A "+
				"-o jsonpath='{.items[?(@.spec.application.name==\"%s\")].metadata.name}'",
			app.GetRef(), sc, app.GetName(),
		)
	}
	return nil
}
