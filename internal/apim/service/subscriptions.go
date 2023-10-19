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

package service

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model/plan"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
)

// Subscriptions brings support for managing gravitee.io APIM support for subscriptions.
// This service is used for testing purposes only and not initialized by the operator manager.
type Subscriptions struct {
	*client.Client
}

func NewSubscriptions(client *client.Client) *Subscriptions {
	return &Subscriptions{Client: client}
}

func (svc *Subscriptions) API(apiID string) *http.URL {
	return svc.EnvV2("apis").WithPath(apiID).WithPath("subscriptions")
}

func (svc *Subscriptions) Subscribe(apiID, applicationID, planID string) (*plan.Subscription, error) {
	request := &plan.SubscriptionRequest{
		PlanID:        planID,
		ApplicationID: applicationID,
	}

	subscription := new(plan.Subscription)

	if err := svc.HTTP.Post(svc.API(apiID), request, subscription); err != nil {
		return nil, err
	}

	return subscription, nil
}

func (svc *Subscriptions) GetApiKeys(apiID, subscriptionID string) ([]plan.ApiKey, error) {
	url := svc.API(apiID).WithPath(subscriptionID).WithPath("api-keys")
	apiKeys := new(plan.ApiKeyList)

	if err := svc.HTTP.Get(url, apiKeys); err != nil {
		return nil, err
	}

	return apiKeys.Data, nil
}
