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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"strconv"
)

// Subscriptions brings support for managing gravitee.io APIM support for subscriptions.
// This service is used for testing purposes only and not initialized by the operator manager.
type Subscriptions struct {
	*client.Client
}

type AutomationSubscription struct {
	HRID            string `json:"hrid"`
	ApplicationHrid string `json:"applicationHrid"`
	PlanHrid        string `json:"planHrid"`
	ApiHrid         string `json:"apiHrid"`
	Status          string `json:"status"`
	StartingAt      string `json:"startingAt"`
	EndingAt        string `json:"endingAt"`
}

func NewSubscriptions(client *client.Client) *Subscriptions {
	return &Subscriptions{Client: client}
}

func (svc *Subscriptions) Import(spec model.Subscription, legacySubscriptionID bool, legacyApiID bool, legacyAppID bool) (model.SubscriptionStatus, error) {
	url := svc.AutomationTarget("apis").WithPath(spec.ApiID).
		WithPath("subscriptions").
		WithQueryParam("legacyID", strconv.FormatBool(legacySubscriptionID)).
		WithQueryParam("legacyApiID", strconv.FormatBool(legacyApiID)).
		WithQueryParam("legacyAppID", strconv.FormatBool(legacyAppID))

	sub := AutomationSubscription{
		HRID:            spec.ID,
		ApiHrid:         spec.ApiID,
		ApplicationHrid: spec.AppID,
		PlanHrid:        spec.PlanID,
		Status:          spec.Status,
		StartingAt:      spec.StartingAt,
		EndingAt:        spec.EndingAt,
	}
	status := new(model.SubscriptionStatus)

	if err := svc.HTTP.Put(url.String(), sub, status); err != nil {
		return model.SubscriptionStatus{}, err
	}

	// If managed with HRID, we don't need IDs
	if !legacySubscriptionID {
		status.UseHRID = true
	}
	return *status, nil
}

// TODO: replace this import 👆
// FOR TESTS
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

func (svc *Subscriptions) Delete(api core.ApiDefinitionObject, subscription *v1alpha1.Subscription) error {

	subID, legacyID := getSubID(subscription)
	apiID, apiLegacyID := getApiID(api)

	url := svc.AutomationTarget("apis").WithPath(apiID).WithPath("subscriptions").WithPath(subID).
		WithQueryParam("legacyID", strconv.FormatBool(legacyID)).WithQueryParam("legacyApiID", strconv.FormatBool(apiLegacyID))

	return svc.HTTP.Delete(url.String(), nil)
}

// FOR TESTS
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

func getSubID(subscription *v1alpha1.Subscription) (string, bool) {
	var id string
	var legacy bool
	if subscription.GetID() == "" {
		id = refs.NewNamespacedNameFromObject(subscription).HRID()
		legacy = false
	} else {
		id = subscription.GetID()
		legacy = true
	}
	return id, legacy
}
