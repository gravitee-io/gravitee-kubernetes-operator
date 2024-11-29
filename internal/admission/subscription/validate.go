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
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	"k8s.io/apimachinery/pkg/runtime"
)

var allowedPlanSecurities = []string{"JWT", "OAUTH2", "MTLS"}

func validateUpdate(
	_ context.Context,
	oldObj runtime.Object,
	newObj runtime.Object,
) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	oldSub, ook := oldObj.(core.SubscriptionObject)
	newSub, nok := newObj.(core.SubscriptionObject)
	if ook && nok {
		errs.MergeWith(validateImmutableProperties(oldSub, newSub))
		if errs.IsSevere() {
			return errs
		}
		errs.Add(validateEndingAt(newSub.GetEndingAt()))
	}
	return errs
}

func validateImmutableProperties(
	oldSub core.SubscriptionObject,
	newSub core.SubscriptionObject,
) *errors.AdmissionErrors {
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

func validateCreate(ctx context.Context, obj runtime.Object) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	sub, ok := obj.(core.SubscriptionObject)
	if !ok {
		return errs
	}

	ns := sub.GetNamespace()

	errs.Add(validateApiKind(sub))

	api, err := dynamic.ResolveAPI(ctx, sub.GetApiRef(), ns)
	if err != nil {
		errs.AddSeveref(
			"unable to resolve API [%s]", sub.GetApiRef(),
		)
	}

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

	errs.Add(validatePlan(sub, api))
	if errs.IsSevere() {
		return errs
	}

	app, err := dynamic.ResolveApplication(ctx, sub.GetAppRef(), ns)
	if err != nil {
		errs.AddSeveref(
			"unable to resolve application [%s]", sub.GetAppRef(),
		)
	}

	if errs.IsSevere() {
		return errs
	}

	errs.Add(validateApplicationState(app))
	if errs.IsSevere() {
		return errs
	}

	plan := api.GetPlan(sub.GetPlan())

	errs.MergeWith(validateApplicationSettings(plan, app))
	if errs.IsSevere() {
		return errs
	}

	errs.Add(validateContextRefs(api, app))
	if errs.IsSevere() {
		return errs
	}

	errs.Add(validateEndingAt(sub.GetEndingAt()))

	return errs
}

func validatePlan(sub core.SubscriptionModel, api core.ApiDefinitionModel) *errors.AdmissionError {
	if !api.HasPlans() {
		return errors.NewSeveref(
			"unable to subscribe to API [%s] because it has no plan", sub.GetApiRef(),
		)
	}

	plan := api.GetPlan(sub.GetPlan())
	if reflect.ValueOf(plan).IsNil() {
		return errors.NewSeveref(
			"unable to subscribe to API [%s] because plan [%s] cannot be found",
			sub.GetApiRef(), sub.GetPlan(),
		)
	}
	return validatePlanSecurityType(plan, sub.GetPlan())
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
	if !slices.Contains(allowedPlanSecurities, plan.GetSecurityType()) {
		return errors.NewSeveref(
			"unable to subscribe to plan [%s] because security type is not one of [%s]",
			planName, strings.Join(allowedPlanSecurities, ","),
		)
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

	if slices.Contains([]string{"JWT", "OAUTH"}, plan.GetSecurityType()) {
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

func validateApiKind(sub core.SubscriptionObject) *errors.AdmissionError {
	apiKind := sub.GetApiRef().GetKind()
	if apiKind == "" {
		return nil // will be defaulted later on
	}
	kind := dynamic.PluralizeKind(apiKind)
	if kind != core.CRDApiDefinitionResource && kind != core.CRDApiV4DefinitionResource {
		return errors.NewSeveref(
			"API kind is required and should be either ApiDefinition or ApiV4Definition, got [%s]",
			kind,
		)
	}
	return nil
}

func validateContextRefs(api core.ApiDefinitionObject, app core.ApplicationObject) *errors.AdmissionError {
	apiCtx, appCtx := api.ContextRef(), app.ContextRef()

	mismatch := appCtx.GetName() != apiCtx.GetName()
	mismatch = mismatch || appCtx.GetNamespace() != apiCtx.GetNamespace()
	if mismatch {
		return errors.NewSeveref(
			"management contexts must match between application [%s] and API [%s], got [%v] and [%v]",
			app.GetRef(),
			api.GetRef(),
			app.ContextRef(),
			api.ContextRef(),
		)
	}
	return nil
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
