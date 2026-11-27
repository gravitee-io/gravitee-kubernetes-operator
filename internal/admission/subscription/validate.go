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

package subscription

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimsubscription "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/subscription"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

var allowedPlanSecurities = []string{"API_KEY", "JWT", "OAUTH2", "MTLS"}

// validateUpdate holds the checks of an update the lifecycle does not run: templates, references and immutable
// fields are checked before it, drift after it.
func validateUpdate(ctx context.Context, sub *v1alpha1.Subscription) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	errs.Add(validateEndingAt(sub.GetEndingAt()))

	_, app, plan := resolveDependencies(ctx, sub, sub.GetNamespace(), errs)
	if errs.IsSevere() {
		return errs
	}

	if plan == nil {
		errs.AddSeveref(
			"unable to subscribe to API [%s] because plan [%s] cannot be found",
			apimsubscription.APIRef(sub), sub.GetPlan(),
		)
		return errs
	}

	errs.Add(validateApiKeys(plan, sub.GetApiKeys()))
	if errs.IsSevere() {
		return errs
	}

	errs.Add(validateConsumerConfiguration(plan, sub))
	if errs.IsSevere() {
		return errs
	}

	validateMTLS(sub, plan, app, errs)
	return errs
}

func validateImmutableProperties(oldSub, newSub *v1alpha1.Subscription) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	if newSub.GetApiRef().String() != oldSub.GetApiRef().String() {
		errs.AddSeveref(
			"API reference is immutable. Detected change from [%s] to [%s]",
			oldSub.GetApiRef(), newSub.GetApiRef(),
		)
	}

	if newSub.GetPlan() != oldSub.GetPlan() {
		errs.AddSeveref(
			"Plan is immutable. Detected change from [%s] to [%s]",
			oldSub.GetPlan(), newSub.GetPlan(),
		)
	}

	if newSub.GetAppRef().String() != oldSub.GetAppRef().String() {
		errs.AddSeveref(
			"Application reference is immutable. Detected change from [%s] to [%s]",
			oldSub.GetAppRef(), newSub.GetAppRef(),
		)
	}

	return errs
}

// validateCreate holds the checks of a create the lifecycle does not run: templates and references are checked
// before it.
func validateCreate(ctx context.Context, sub *v1alpha1.Subscription) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	ns := sub.GetNamespace()

	api, app, plan := resolveDependencies(ctx, sub, ns, errs)

	if errs.IsSevere() {
		return errs
	}

	errs.Add(validateApiSyncMode(api))
	if errs.IsSevere() {
		return errs
	}

	errs.Add(validateApiState(api))
	if errs.IsSevere() {
		return errs
	}

	errs.Add(ValidatePlan(sub, api))
	if errs.IsSevere() {
		return errs
	}

	errs.Add(validateApiKeys(plan, sub.GetApiKeys()))
	if errs.IsSevere() {
		return errs
	}

	errs.Add(validateApplicationState(app))
	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(validateApplicationSettings(plan, app))
	if errs.IsSevere() {
		return errs
	}

	errs.Add(ValidateContextRefs(ctx, api, app))
	if errs.IsSevere() {
		return errs
	}

	errs.Add(validateEndingAt(sub.GetEndingAt()))
	if errs.IsSevere() {
		return errs
	}

	validateMTLS(sub, plan, app, errs)

	return errs
}

func resolveDependencies(
	ctx context.Context,
	sub *v1alpha1.Subscription,
	ns string, errs *errors.AdmissionErrors) (core.ApiDefinitionObject, core.ApplicationObject, core.PlanModel) {
	apiRef := apimsubscription.APIRef(sub)
	api, err := dynamic.ResolveAPI(ctx, apiRef, ns)
	if err != nil {
		errs.AddSeveref(
			"unable to resolve API [%s]", apiRef,
		)
	}

	if errs.IsSevere() {
		return nil, nil, nil
	}

	appRef := apimsubscription.AppRef(sub)
	app, err := dynamic.ResolveApplication(ctx, appRef, ns)
	if err != nil {
		errs.AddSeveref(
			"unable to resolve application [%s]", appRef,
		)
	}

	plan := api.GetPlan(sub.GetPlan())
	return api, app, plan
}

// ValidatePlan checks that the subscribed plan exists and that GKO can subscribe to it.
func ValidatePlan(sub *v1alpha1.Subscription, api core.ApiDefinitionModel) *errors.AdmissionError {
	if !api.HasPlans() {
		return errors.NewSeveref(
			"unable to subscribe to API [%s] because it has no plan", apimsubscription.APIRef(sub),
		)
	}

	plan := api.GetPlan(sub.GetPlan())
	if plan == nil {
		return errors.NewSeveref(
			"unable to subscribe to API [%s] because plan [%s] cannot be found",
			apimsubscription.APIRef(sub), sub.GetPlan(),
		)
	}
	if err := validatePlanSecurityType(plan, sub.GetPlan()); err != nil {
		return err
	}
	return validateConsumerConfiguration(plan, sub)
}

// validateConsumerConfiguration mirrors APIM: a PUSH subscription needs the entrypoint
// that delivers its messages.
func validateConsumerConfiguration(plan core.PlanModel, sub *v1alpha1.Subscription) *errors.AdmissionError {
	if !plan.IsPush() {
		return nil
	}
	consumer := sub.Spec.ConsumerConfiguration
	if consumer == nil || consumer.EntrypointID == "" {
		return errors.NewSeveref(
			"consumerConfiguration.entrypointId is required to subscribe to PUSH plan [%s]",
			sub.GetPlan(),
		)
	}
	return nil
}

func validateMTLS(
	sub core.SubscriptionObject,
	plan core.PlanModel,
	app core.ApplicationObject,
	errs *errors.AdmissionErrors) {
	if plan.GetSecurityType() == "MTLS" {
		errs.Add(validateCertEndDatesVsSubscriptionEndDate(sub, app))
	}
}

func validateApiState(api core.ApiDefinitionObject) *errors.AdmissionError {
	if api.GetStatus() == nil || api.GetStatus().IsFailed() {
		return errors.NewSeveref(
			"unable to subscribe to API [%s] because it is in a failed state", api.GetRef(),
		)
	}
	if api.IsStopped() {
		return errors.NewSeveref(
			"unable to subscribe to API [%s] because it is not started", api.GetRef(),
		)
	}
	return nil
}

func validateApplicationState(app core.ApplicationObject) *errors.AdmissionError {
	if app.GetStatus() == nil || app.GetStatus().IsFailed() {
		return errors.NewSeveref(
			"unable to subscribe from application [%s] because it is in a failed state", app.GetRef(),
		)
	}
	return nil
}

func validatePlanSecurityType(plan core.PlanModel, planName string) *errors.AdmissionError {
	if plan.IsPush() {
		return nil
	}
	if !slices.Contains(allowedPlanSecurities, plan.GetSecurityType()) {
		return errors.NewSeveref(
			"unable to subscribe to plan [%s] because security type is not one of [%s]",
			planName, strings.Join(allowedPlanSecurities, ","),
		)
	}
	return nil
}

func validateApiKeys(plan core.PlanModel, apiKeys []core.ApiKeyModel) *errors.AdmissionError {
	if len(apiKeys) == 0 {
		return nil
	}

	if plan.GetSecurityType() != "API_KEY" {
		got := fmt.Sprintf("[%s]", plan.GetSecurityType())
		if plan.IsPush() {
			got = "PUSH plan"
		}
		return errors.NewSeveref("apiKeys can only be set when subscribing to an API_KEY plan, got %s", got)
	}

	seen := make(map[string]bool, len(apiKeys))
	for _, k := range apiKeys {
		if seen[k.GetKey()] {
			return errors.NewSeveref(
				"duplicate key [%s] in apiKeys",
				k.GetKey(),
			)
		}
		seen[k.GetKey()] = true

		if k.GetExpireAt() != nil {
			if err := validateEndingAt(k.GetExpireAt()); err != nil {
				return errors.NewSeveref(
					"invalid expireAt for key [%s]: %s",
					k.GetKey(), err.Message,
				)
			}
		}
	}

	return nil
}

func validateApiSyncMode(api core.ApiDefinitionObject) *errors.AdmissionError {
	if !api.IsSyncFromManagement() {
		return errors.NewSeveref(
			"unable to subscribe to API [%s] because its definition is not synced from the management API (%s)",
			api.GetRef(),
			"sourcing subscriptions from a Kubernetes cluster is not supported at the moment",
		)
	}
	return nil
}

func validateApplicationSettings(plan core.PlanModel, app core.ApplicationObject) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	if slices.Contains([]string{"JWT", "OAUTH2"}, plan.GetSecurityType()) {
		errs.Add(validateClientID(app))
	}

	if plan.GetSecurityType() == "MTLS" {
		errs.Add(validateClientCertificate(app))
	}
	return errs
}

func validateClientCertificate(app core.ApplicationObject) *errors.AdmissionError {
	settings := app.GetModel().GetSettings()
	if !settings.HasTLS() {
		return errors.NewSeveref(
			"unable to subscribe to MTLS plan from application [%s] because it does not have any client certificate",
			app.GetRef(),
		)
	}
	return nil
}

func validateClientID(app core.ApplicationObject) *errors.AdmissionError {
	settings := app.GetModel().GetSettings()
	if settings.IsSimple() && settings.GetClientID() == "" {
		return errors.NewSeveref(
			"unable to subscribe from application [%s] because it does not have any client id",
			app.GetRef(),
		)
	}
	return nil
}

func ValidateContextRefs(
	ctx context.Context,
	api core.ApiDefinitionObject,
	app core.ApplicationObject,
) *errors.AdmissionError {
	if !api.HasContext() {
		return errors.NewSeveref(
			"unable to subscribe to API [%s] because it does not reference a management context",
			api.GetRef(),
		)
	}

	if !app.HasContext() {
		return errors.NewSeveref(
			"unable to subscribe from application [%s] because it does not reference a management context",
			app.GetRef(),
		)
	}

	apiCtx, err := dynamic.ResolveContext(ctx, api.ContextRef(), api.GetNamespace())
	if err != nil {
		return errors.NewSeveref(
			"unable to resolve management context [%s] of API [%s]",
			api.ContextRef(), api.GetRef(),
		)
	}

	appCtx, err := dynamic.ResolveContext(ctx, app.ContextRef(), app.GetNamespace())
	if err != nil {
		return errors.NewSeveref(
			"unable to resolve management context [%s] of application [%s]",
			app.ContextRef(), app.GetRef(),
		)
	}

	same, err := SameTarget(apiCtx, appCtx)
	if err != nil {
		return errors.NewSevere(err.Error())
	}

	if !same {
		return errors.NewSeveref(
			"management contexts of application [%s] and API [%s] must target the same environment, "+
				"got [%s] (%s) and [%s] (%s)",
			app.GetRef(), api.GetRef(),
			appCtx.GetRef(), envTargetOrEmpty(appCtx),
			apiCtx.GetRef(), envTargetOrEmpty(apiCtx),
		)
	}

	return nil
}

// SameTarget reports whether two management contexts address the same APIM environment.
// Distinct contexts may do so, e.g. one context per team with its own credentials.
func SameTarget(left core.ContextObject, right core.ContextObject) (bool, error) {
	if left.GetNamespace() == right.GetNamespace() && left.GetName() == right.GetName() {
		return true, nil
	}

	leftTarget, err := apim.EnvTarget(left)
	if err != nil {
		return false, fmt.Errorf("invalid base URL in management context [%s]: %w", left.GetRef(), err)
	}

	rightTarget, err := apim.EnvTarget(right)
	if err != nil {
		return false, fmt.Errorf("invalid base URL in management context [%s]: %w", right.GetRef(), err)
	}

	return leftTarget == rightTarget, nil
}

func envTargetOrEmpty(context core.ContextObject) string {
	target, _ := apim.EnvTarget(context)
	return target
}

func validateCertEndDatesVsSubscriptionEndDate(
	sub core.SubscriptionObject,
	app core.ApplicationObject,
) *errors.AdmissionError {
	endingAt := sub.GetEndingAt()
	if endingAt == nil {
		return nil
	}

	settings, ok := app.GetModel().GetSettings().(*application.Setting)
	if !ok || !settings.HasClientCertificates() {
		return nil
	}

	subEnd, err := time.Parse(time.RFC3339, *endingAt)
	if err != nil {
		return nil // endingAt format is validated elsewhere
	}

	certs := settings.GetClientCertificates()
	certsWithEndsAt := 0
	for _, cert := range certs {
		if cert.EndsAt == "" {
			continue
		}
		certsWithEndsAt++
		certEnd, err := time.Parse(time.RFC3339, cert.EndsAt)
		if err != nil {
			continue
		}
		if certEnd.After(subEnd) {
			return nil // at least one cert outlives the subscription
		}
	}

	if certsWithEndsAt == 0 {
		return nil // no cert has endsAt, nothing to validate
	}

	return errors.NewSeveref(
		"subscription ending date [%s] is after all client certificate end dates in application [%s]",
		*endingAt, app.GetRef(),
	)
}

func validateEndingAt(endingAt *string) *errors.AdmissionError {
	if endingAt != nil {
		t, err := time.Parse(time.RFC3339, *endingAt)
		if err != nil {
			return errors.NewSeveref(
				"ending date [%s] is not in RFC3339 format",
				*endingAt,
			)
		}
		tx := time.Now().Add(1 * time.Minute)
		if t.Local().Before(tx) {
			return errors.NewSeveref(
				"ending date [%s] should be at least one minute from now",
				*endingAt,
			)
		}
	}
	return nil
}
