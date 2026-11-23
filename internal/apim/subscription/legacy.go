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

// This file holds the addressing of resources created before APIM 4.12, by their UUID. Removing it leaves
// these to fix: the applyLegacy call in ToDTO, the params of the SDK calls (nil), UpdateStatus (always
// Automation API managed), and the drift mapper (ToDTO in place of ToDriftDTO).

import (
	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

// Legacy flags the identities of a DTO that are UUIDs rather than hrids. Each identity (subscription, API,
// application) is independently one or the other: a new subscription to an API created before 4.12 is a
// valid mix. The fields are the SDK query parameters they set.
type Legacy struct {
	// HridContainsUUID is set when the subscription Hrid is its UUID.
	HridContainsUUID bool
	// HridContainsApiUUID is set when ApiHrid and PlanHrid are UUIDs.
	HridContainsApiUUID bool
	// HridContainsAppUUID is set when ApplicationHrid is a UUID.
	HridContainsAppUUID bool
}

func (l Legacy) upsertParams() *sdk.UpsertApiSubscriptionParams {
	return &sdk.UpsertApiSubscriptionParams{
		HridContainsUUID:    &l.HridContainsUUID,
		HridContainsApiUUID: &l.HridContainsApiUUID,
		HridContainsAppUUID: &l.HridContainsAppUUID,
	}
}

func (l Legacy) getParams() *sdk.GetApiSubscriptionParams {
	return &sdk.GetApiSubscriptionParams{
		HridContainsUUID:    &l.HridContainsUUID,
		HridContainsApiUUID: &l.HridContainsApiUUID,
	}
}

func (l Legacy) deleteParams() *sdk.DeleteApiSubscriptionParams {
	return &sdk.DeleteApiSubscriptionParams{
		HridContainsUUID:    &l.HridContainsUUID,
		HridContainsApiUUID: &l.HridContainsApiUUID,
	}
}

// applyLegacy swaps in the UUID of each identity created before the Automation API, and flags it.
// The plan of a UUID-addressed API is its UUID too, or its hrid when the API has no such plan.
func applyLegacy(dto *DTO, sub *v1alpha1.Subscription, api core.ApiDefinitionObject, app core.ApplicationObject) {
	if SubscriptionUsesUUID(sub) {
		dto.Hrid = sub.Status.ID
		dto.Legacy.HridContainsUUID = true
	}
	if model.APIUsesUUID(api) {
		dto.ApiHrid = api.GetID()
		dto.Legacy.HridContainsApiUUID = true
		if plan := api.GetPlan(sub.Spec.Plan); plan != nil {
			dto.PlanHrid = plan.GetID()
		}
	}
	if app != nil && ApplicationUsesUUID(app) {
		dto.ApplicationHrid = app.GetID()
		dto.Legacy.HridContainsAppUUID = true
	}
}

// ToDriftDTO maps a Subscription CRD to what GetRemote returns. A subscription created before the Automation
// API is read back with UUIDs only; any other is read back as ToDTO maps it.
func ToDriftDTO(sub *v1alpha1.Subscription, api core.ApiDefinitionObject, app core.ApplicationObject) (DTO, error) {
	dto, err := ToDTO(sub, api, app)
	if err != nil || !SubscriptionUsesUUID(sub) {
		return dto, err
	}
	dto.ApiHrid = api.GetID()
	dto.ApplicationHrid = app.GetID()
	dto.Legacy = Legacy{HridContainsUUID: true, HridContainsApiUUID: true, HridContainsAppUUID: true}
	dto.PlanHrid = ""
	if plan := api.GetPlan(sub.Spec.Plan); plan != nil {
		dto.PlanHrid = plan.GetID()
	}
	return dto, nil
}

// SubscriptionUsesUUID is true only for a subscription created before the Automation API (UUID in status, no
// AutomationAPIManaged condition). A create that has not persisted the condition yet is still hrid-addressed.
func SubscriptionUsesUUID(sub *v1alpha1.Subscription) bool {
	return sub.Status.ID != "" && !k8s.IsAutomationAPIManaged(sub)
}

// ApplicationUsesUUID is true for an application created before the Automation API.
func ApplicationUsesUUID(app core.ApplicationObject) bool {
	return app.GetID() != "" && !k8s.IsAutomationAPIManaged(app)
}
