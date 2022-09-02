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

package managementapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi/model"
)

func (client *Client) SubscribeToPlan(
	apiId string,
	applicationId string,
	planId string,
) (*model.Subscription, error) {
	queryParams := "?application=" + applicationId + "&plan=" + planId

	url := client.buildUrl("/apis/" + apiId + "/subscriptions" + queryParams)

	req, err := http.NewRequestWithContext(client.ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to subscribe to the plan into the Management API, ApiId=%s, ApplicationId=%s, PlanId=%s",
			apiId, applicationId, planId)
	}

	resp, err := client.http.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf(
			"unable to subscribe to the plan into the Management API, ApiId=%s, ApplicationId=%s, PlanId=%s, HTTP Status: %d",
			apiId, applicationId, planId, resp.StatusCode)
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	var subscription model.Subscription
	err = json.Unmarshal(body, &subscription)
	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (client *Client) GetSubscriptionApiKey(
	apiId string,
	subscriptionId string,
) ([]model.ApiKeyEntity, error) {
	req, err := http.NewRequestWithContext(
		client.ctx,
		http.MethodGet,
		client.buildUrl("/apis/"+apiId+"/subscriptions/"+subscriptionId+"/apikeys"),
		nil,
	)

	if err != nil && apiId == "" {
		return nil, fmt.Errorf(
			"unable to look for apikey matching apiId=%s and subscriptionId=%s (%w)",
			apiId, subscriptionId, err)
	}
	resp, err := client.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("an error as occurred while performing GetSubscriptionApiKey request")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"an error as occurred trying to get apikey matching apiId=%s and subscriptionId=%s, HTTP Status: %d ",
			apiId, subscriptionId, resp.StatusCode)
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	var apiKeys []model.ApiKeyEntity

	err = json.Unmarshal(body, &apiKeys)
	if err != nil {
		return nil, err
	}

	return apiKeys, nil
}
