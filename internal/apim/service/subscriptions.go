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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
)

type Subscriptions struct {
	*client.Client
}

func NewSubscriptions(client *client.Client) *Subscriptions {
	return &Subscriptions{Client: client}
}

func (svc *Subscriptions) APITarget(apiID string) *http.URL {
	return svc.EnvTarget("apis").WithPath(apiID).WithPath("subscriptions")
}

func (svc *Subscriptions) Subscribe(apiID, applicationID, planID string) (*model.Subscription, error) {
	url := svc.APITarget(apiID).WithQueryParams(
		map[string]string{
			planParam:        planID,
			applicationParam: applicationID,
		},
	)

	subscription := new(model.Subscription)

	if err := svc.HTTP.Post(url.String(), nil, subscription); err != nil {
		return nil, err
	}

	return subscription, nil
}

func (svc *Subscriptions) GetApiKeys(apiID, subscriptionID string) ([]model.ApiKeyEntity, error) {
	url := svc.APITarget(apiID).WithPath(subscriptionID).WithPath("apikeys")
	apiKeys := new([]model.ApiKeyEntity)

	if err := svc.HTTP.Get(url.String(), apiKeys); err != nil {
		return nil, err
	}

	return *apiKeys, nil
}
