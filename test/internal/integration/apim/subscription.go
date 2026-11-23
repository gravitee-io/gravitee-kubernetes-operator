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

package apim

import (
	"context"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/service"
	apimsubscription "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/subscription"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

// Subscriptions reads and writes subscriptions in APIM outside the operator. Subscribe and GetApiKeys
// come from the Management API service.
type Subscriptions struct {
	*service.Subscriptions
	sdk *apim.SDKClient
}

// SubscriptionDTO is the subscription Import writes: identities and an end date, nothing else.
type SubscriptionDTO struct {
	ID       string
	ApiID    string
	AppID    string
	PlanID   string
	EndingAt string
}

// GetByHRID gets a subscription by its hrid, under an API addressed by its hrid.
func (s *Subscriptions) GetByHRID(apiHrid, subscriptionHrid string) (sdk.SubscriptionState, error) {
	return apimsubscription.GetState(context.Background(), s.sdk, apimsubscription.DTO{
		Subscription: sdk.Subscription{Hrid: subscriptionHrid},
		ApiHrid:      apiHrid,
	})
}

// GetByHRIDWithAPIUUID gets a subscription by its hrid, under an API created before the Automation API.
func (s *Subscriptions) GetByHRIDWithAPIUUID(apiUUID, subscriptionHrid string) (sdk.SubscriptionState, error) {
	return apimsubscription.GetState(context.Background(), s.sdk, apimsubscription.DTO{
		Subscription: sdk.Subscription{Hrid: subscriptionHrid},
		ApiHrid:      apiUUID,
		Legacy:       apimsubscription.Legacy{HridContainsApiUUID: true},
	})
}

// Import creates or updates a subscription in APIM outside the operator. As the operator does, it marks obj
// Automation API managed unless the subscription is addressed by its UUID.
func (s *Subscriptions) Import(
	spec SubscriptionDTO,
	obj core.ConditionAwareObject,
	hridContainsUUID, hridContainsApiUUID, hridContainsAppUUID bool,
) (sdk.SubscriptionState, error) {
	dto := apimsubscription.DTO{
		Subscription: sdk.Subscription{
			Hrid:            spec.ID,
			ApplicationHrid: spec.AppID,
			PlanHrid:        spec.PlanID,
		},
		ApiHrid: spec.ApiID,
		Legacy: apimsubscription.Legacy{
			HridContainsUUID:    hridContainsUUID,
			HridContainsApiUUID: hridContainsApiUUID,
			HridContainsAppUUID: hridContainsAppUUID,
		},
	}
	if spec.EndingAt != "" {
		dto.EndingAt = &spec.EndingAt
	}

	resp, err := apimsubscription.Upsert(context.Background(), s.sdk, dto)
	if err != nil {
		return sdk.SubscriptionState{}, err
	}
	if !hridContainsUUID {
		k8s.AddAutomationAPIManagedCondition(obj)
	}
	return resp.State, nil
}
