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

// Package subscription maps the Subscription CRD to the APIM Automation API through the SDK.
// Everything about resources created before APIM 4.12, addressed by their UUID, is in legacy.go.
package subscription

import (
	"fmt"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hrid"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/mapper"
)

// DTO is the subscription payload with what Upsert, GetRemote and Delete need to address it. None of the
// extra fields is sent (Upsert sends the embedded sdk.Subscription) and drift detection ignores them.
// The API is not compared: it is in the URL path, and the API ref of a Subscription is immutable.
type DTO struct {
	sdk.Subscription
	// ApiHrid addresses the API in the URL path.
	ApiHrid string `json:"-" drift:"ignore"`
	// Legacy flags the identities that are UUIDs.
	Legacy Legacy `json:"-" drift:"ignore"`
}

// ToDTO maps a Subscription CRD to the Automation API payload, addressing its API, plan and application as
// RecordRefs wrote them in sub.Resolved.
func ToDTO(sub *v1alpha1.Subscription) (DTO, error) {
	payload, err := mapper.MapViaJSON[sdk.Subscription](sub.Spec.Type)
	if err != nil {
		return DTO{}, fmt.Errorf("map subscription [%s] to the APIM payload: %w", sub.Name, err)
	}

	dto := DTO{Subscription: payload, ApiHrid: sub.Resolved.ApiHrid}
	dto.Hrid = refs.NewNamespacedNameFromObject(sub).HRID()
	dto.PlanHrid = sub.Resolved.PlanHrid
	dto.ApplicationHrid = sub.Resolved.ApplicationHrid
	applyLegacy(&dto, sub)
	return dto, nil
}

// RecordRefs writes in sub.Resolved how ToDTO addresses the API, its plan and the application: by their hrids, or by
// UUID for those created before the Automation API. A nil app (deleted application) leaves the application out.
// ResolveRefs calls it once the references are found and the plan UUIDs of a legacy API are filled.
func RecordRefs(sub *v1alpha1.Subscription, api core.ApiDefinitionObject, app core.ApplicationObject) {
	parents := v1alpha1.SubscriptionResolved{
		ApiHrid:  refs.NewNamespacedNameFromObject(api).HRID(),
		PlanHrid: hrid.NameToValidHRID(sub.Spec.Plan),
	}
	if app != nil {
		parents.ApplicationHrid = refs.NewNamespacedNameFromObject(app).HRID()
	}
	applyLegacyParents(&parents, sub, api, app)
	sub.Resolved = parents
}

// Response is what Upsert returns: the subscription state APIM answers, and the flags the upsert was sent with.
type Response struct {
	State  sdk.SubscriptionState
	Legacy Legacy
}

// GetOrgID is empty: a subscription response carries no organization.
func (Response) GetOrgID() string {
	return ""
}

// GetEnvID is empty: a subscription response carries no environment.
func (Response) GetEnvID() string {
	return ""
}
