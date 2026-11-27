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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/ctxref"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
	"k8s.io/apimachinery/pkg/types"
)

// preCheck runs after the lifecycle compiled the templates and resolved the certificate refs, on create and
// update, in today's order: deprecation warnings, context, settings, then the ignore-groups warning.
func preCheck(ctx context.Context, app *v1alpha1.Application) *errors.AdmissionErrors {
	errs := WarnDeprecatedFields(app)
	errs.Add(ctxref.Validate(ctx, app))
	if errs.IsSevere() {
		return errs
	}
	errs.MergeWith(validateSettings(app))
	if errs.IsSevere() {
		return errs
	}
	if k8s.HasTrueAnnotation(app, core.IgnoreGroupsAnnotation) && len(app.Spec.Groups) > 0 {
		errs.AddWarningf("spec.groups is not applied while %s is true", core.IgnoreGroupsAnnotation)
	}
	return errs
}

// WarnDeprecatedFields warns about each field set that APIM ignores and that will be removed in 2027-12.
func WarnDeprecatedFields(app *v1alpha1.Application) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	if app.Spec.Picture != nil {
		errs.AddWarning("[spec.picture] is deprecated and ignored by APIM, it will be removed in 2027-12; " +
			"use [spec.pictureUrl]")
	}
	if app.Spec.Metadata == nil {
		return errs
	}
	for _, metadata := range *app.Spec.Metadata {
		if metadata.Hidden != nil {
			errs.AddWarningf("[spec.metadata[%s].hidden] is deprecated and ignored by APIM, "+
				"it will be removed in 2027-12", metadata.Name)
		}
	}
	return errs
}

func validateSettings(app *v1alpha1.Application) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	settings := app.Spec.Settings
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
		errs.MergeWith(validateClientCertificates(settings.GetClientCertificates()))
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

func validateSettingsUpdate(oldApp, newApp *v1alpha1.Application) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	oldSettings, newSettings := oldApp.Spec.Settings, newApp.Spec.Settings

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

func validateCertEndDatesVsSubscriptionEndDates(
	ctx context.Context,
	app *v1alpha1.Application,
) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	settings := app.Spec.Settings
	if settings == nil || !settings.HasClientCertificates() {
		return errs
	}

	maxCertEnd := maxCertificateEndDate(settings.GetClientCertificates())
	if maxCertEnd == nil {
		return errs
	}

	appRef := refs.NewNamespacedNameFromObject(app)
	subList := &v1alpha1.SubscriptionList{}
	if err := search.FindByFieldReferencing(ctx, search.AppSubsField, *appRef, subList); err != nil {
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

func validateDelete(ctx context.Context, app *v1alpha1.Application) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	errs.Add(validateNoSubscription(ctx, app))
	return errs
}

// The Subscriptions themselves are the source of truth: the subscription counter in the
// application status can be lost on a write conflict, and is not restored with the object.
func validateNoSubscription(ctx context.Context, app *v1alpha1.Application) *errors.AdmissionError {
	subs, err := search.FindSubscriptionsReferencing(ctx, search.AppSubsField, app)
	if err != nil {
		return errors.NewSevere(err.Error())
	}
	if sc := len(subs.Items); sc > 0 {
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
