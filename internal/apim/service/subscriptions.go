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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
)

// Subscriptions brings support for managing gravitee.io APIM support for subscriptions.
// This service is used for testing purposes only and not initialized by the operator manager.
type Subscriptions struct {
	*client.Client
}

func NewSubscriptions(client *client.Client) *Subscriptions {
	return &Subscriptions{Client: client}
}

func (svc *Subscriptions) Import(spec *model.Subscription) (*model.SubscriptionStatus, error) {
	url := svc.EnvV2Target("apis").WithPath(spec.ApiID).
		WithPath("subscriptions").WithPath("spec").
		WithPath("_import")

	status := new(model.SubscriptionStatus)

	if err := svc.HTTP.Put(url.String(), spec, status); err != nil {
		return nil, err
	}

	return status, nil
}

// TODO: replace this import 👆
func (svc *Subscriptions) Subscribe(apiID, appID, planID string) (*model.SubscriptionResponse, error) {
	url := svc.EnvV2Target("apis").WithPath(apiID).WithPath("subscriptions")

	request := &model.SubscriptionRequest{
		AppID:  appID,
		PlanID: planID,
	}

	response := new(model.SubscriptionResponse)

	if err := svc.HTTP.Post(url.String(), request, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (svc *Subscriptions) Delete(spec *model.Subscription) error {
	url := svc.EnvV2Target("apis").WithPath(spec.ApiID).
		WithPath("subscriptions").WithPath("spec").WithPath(spec.ID)

	return svc.HTTP.Delete(url.String(), nil)
}

func (svc *Subscriptions) GetApiKeys(apiID, subscriptionID string) ([]model.ApiKeyEntity, error) {
	url := svc.EnvV1Target("apis").WithPath(apiID).
		WithPath("subscriptions").WithPath(subscriptionID).
		WithPath("apikeys")

	apiKeys := new([]model.ApiKeyEntity)

	if err := svc.HTTP.Get(url.String(), apiKeys); err != nil {
		return nil, err
	}

	return *apiKeys, nil
}
