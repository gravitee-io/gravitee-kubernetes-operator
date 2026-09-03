// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hrid"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

// ToSubscriptionDTO maps a subscription CR and its API/application for Import.
// Each identity (subscription, API, application) is independently HRID or UUID —
// a new subscription against a pre-4.12 API is a valid mix.
// StartingAt is omitted: it is APIM-managed (see SubscriptionDTO drift tag).
func ToSubscriptionDTO(
	sub *v1alpha1.Subscription,
	api core.ApiDefinitionObject,
	app core.ApplicationObject,
) SubscriptionDTO {
	dto := subscriptionSpecDTO(sub)
	dto.ID = subscriptionRef(sub)
	dto.ApiID = apiRef(api)
	dto.AppID = appRef(app)
	dto.PlanID = planRef(api, sub.Spec.Plan)
	return dto
}

// ToSubscriptionDTOForDrift maps identities to match the remote GET.
// Management API (legacy subscription) returns UUIDs only. Automation returns
// the same mixed HRID/UUID shape as ToSubscriptionDTO.
func ToSubscriptionDTOForDrift(
	sub *v1alpha1.Subscription,
	api core.ApiDefinitionObject,
	app core.ApplicationObject,
) SubscriptionDTO {
	if !SubscriptionUsesUUID(sub) {
		return ToSubscriptionDTO(sub, api, app)
	}
	dto := subscriptionSpecDTO(sub)
	dto.ID = sub.Status.ID
	dto.ApiID = api.GetID()
	dto.AppID = app.GetID()
	if plan := api.GetPlan(sub.Spec.Plan); plan != nil {
		dto.PlanID = plan.GetID()
	}
	return dto
}

func subscriptionSpecDTO(sub *v1alpha1.Subscription) SubscriptionDTO {
	spec := sub.Spec
	return SubscriptionDTO{
		EndingAt:              utils.SafeDereference(spec.EndingAt),
		Metadata:              spec.Metadata,
		ApiKeys:               mapViaJSON[[]ApiKeyDTO](spec.ApiKeys),
		ConsumerConfiguration: mapViaJSON[*ConsumerConfigurationDTO](spec.ConsumerConfiguration),
	}
}

// SubscriptionUsesUUID is true only for a pre-Automation subscription
// (UUID in status, no AutomationAPIManaged condition). A create that has not
// persisted the condition yet is still HRID-addressed.
func SubscriptionUsesUUID(sub *v1alpha1.Subscription) bool {
	return sub.Status.ID != "" && !k8s.IsAutomationAPIManaged(sub)
}

func APIUsesUUID(api core.ApiDefinitionObject) bool {
	return api.GetID() != "" && !k8s.IsAutomationAPIManaged(api)
}

func ApplicationUsesUUID(app core.ApplicationObject) bool {
	return app.GetID() != "" && !k8s.IsAutomationAPIManaged(app)
}

// SubscriptionRemoteUsesLegacyAPI is true when the subscription is HRID-addressed
// but the API is a pre-4.12 UUID resource. GET must use GetByHRIDWithAPIUUID.
func SubscriptionRemoteUsesLegacyAPI(sub *v1alpha1.Subscription, api core.ApiDefinitionObject) bool {
	return !SubscriptionUsesUUID(sub) && APIUsesUUID(api)
}

func subscriptionRef(sub *v1alpha1.Subscription) string {
	if SubscriptionUsesUUID(sub) {
		return sub.Status.ID
	}
	return refs.NewNamespacedNameFromObject(sub).HRID()
}

func apiRef(api core.ApiDefinitionObject) string {
	if APIUsesUUID(api) {
		return api.GetID()
	}
	return refs.NewNamespacedNameFromObject(api).HRID()
}

func appRef(app core.ApplicationObject) string {
	if ApplicationUsesUUID(app) {
		return app.GetID()
	}
	return refs.NewNamespacedNameFromObject(app).HRID()
}

func planRef(api core.ApiDefinitionObject, planName string) string {
	if APIUsesUUID(api) {
		if plan := api.GetPlan(planName); plan != nil {
			return plan.GetID()
		}
	}
	return hrid.NameToValidHRID(planName)
}
