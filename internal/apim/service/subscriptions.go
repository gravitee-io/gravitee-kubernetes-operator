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
	"strconv"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	httputil "github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

// Subscriptions brings support for managing gravitee.io APIM support for subscriptions.
// This service is used for testing purposes only and not initialized by the operator manager.
type Subscriptions struct {
	*client.Client
}

func NewSubscriptions(client *client.Client) *Subscriptions {
	return &Subscriptions{Client: client}
}

func (svc *Subscriptions) Import(spec model.SubscriptionDTO,
	conditionAware core.ConditionAwareObject,
	setHridWithUUID bool,
	setHridWithApiUUID bool,
	setHridWithAppUUID bool) (model.SubscriptionStatus, error) {
	url := svc.AutomationTarget("apis").WithPath(spec.ApiID).
		WithPath("subscriptions").
		WithQueryParam("hridContainsUUID", strconv.FormatBool(setHridWithUUID)).
		WithQueryParam("hridContainsApiUUID", strconv.FormatBool(setHridWithApiUUID)).
		WithQueryParam("hridContainsAppUUID", strconv.FormatBool(setHridWithAppUUID))

	status := new(model.SubscriptionStatus)

	automation := spec.ToAutomation()

	if err := svc.HTTP.Put(url.String(), automation, status); err != nil {
		return model.SubscriptionStatus{}, err
	}

	if !setHridWithUUID {
		k8s.AddAutomationAPIManagedCondition(conditionAware)
	}

	return *status, nil
}

// Subscribe For tests purposes only.
func (svc *Subscriptions) Subscribe(apiID, appID, planID string) (*model.SubscriptionResponseDTO, error) {
	url := svc.EnvV2Target("apis").WithPath(apiID).WithPath("subscriptions")

	request := &model.SubscriptionRequestDTO{
		AppID:  appID,
		PlanID: planID,
	}

	response := new(model.SubscriptionResponseDTO)

	if err := svc.HTTP.Post(url.String(), request, response); err != nil {
		return nil, err
	}

	return response, nil
}

// GetByHRID For tests purposes only.
func (svc *Subscriptions) GetByHRID(apiHRID, subscriptionHRID string) (*model.SubscriptionDTO, error) {
	return svc.getSubscription(svc.subscriptionURL(apiHRID, subscriptionHRID))
}

// GetByHRIDWithAPIUUID fetches an HRID-addressed subscription under a pre-Automation API UUID.
// For tests purposes only.
func (svc *Subscriptions) GetByHRIDWithAPIUUID(
	apiUUID, subscriptionHRID string,
) (*model.SubscriptionDTO, error) {
	url := svc.subscriptionURL(apiUUID, subscriptionHRID).
		WithQueryParam("hridContainsApiUUID", strconv.FormatBool(true))
	return svc.getSubscription(url)
}

// GetByHRIDWithSubUUID fetches a pre-Automation subscription UUID under an HRID-addressed API.
// For tests purposes only.
func (svc *Subscriptions) GetByHRIDWithSubUUID(
	apiHRID, subscriptionUUID string,
) (*model.SubscriptionDTO, error) {
	url := svc.subscriptionURL(apiHRID, subscriptionUUID).
		WithQueryParam("hridContainsUUID", strconv.FormatBool(true))
	return svc.getSubscription(url)
}

// GetWithUUID fetches a pre-Automation subscription and API by UUID from the Automation API.
// For tests purposes only.
func (svc *Subscriptions) GetWithUUID(apiUUID, subscriptionUUID string) (*model.SubscriptionDTO, error) {
	url := svc.subscriptionURL(apiUUID, subscriptionUUID).
		WithQueryParam("hridContainsUUID", strconv.FormatBool(true)).
		WithQueryParam("hridContainsApiUUID", strconv.FormatBool(true))
	return svc.getSubscription(url)
}

func (svc *Subscriptions) subscriptionURL(apiPath, subscriptionPath string) *httputil.URL {
	return svc.AutomationTarget("apis").WithPath(apiPath).
		WithPath("subscriptions").WithPath(subscriptionPath)
}

func (svc *Subscriptions) getSubscription(url *httputil.URL) (*model.SubscriptionDTO, error) {
	sub := new(model.AutomationSubscriptionDTO)

	if err := svc.HTTP.Get(url.String(), sub); err != nil {
		return nil, err
	}

	return sub.ToLegacy(), nil
}

func (svc *Subscriptions) Delete(api core.ApiDefinitionObject, subscription *v1alpha1.Subscription) error {
	subID, setHridWithUUID := getSubID(subscription)
	apiID, apisetHridWithUUID := subscriptionAPIID(api)

	url := svc.AutomationTarget("apis").
		WithPath(apiID).
		WithPath("subscriptions").
		WithPath(subID).
		WithQueryParam("hridContainsUUID", strconv.FormatBool(setHridWithUUID)).
		WithQueryParam("hridContainsApiUUID", strconv.FormatBool(apisetHridWithUUID))

	return svc.HTTP.Delete(url.String(), nil)
}

// GetApiKeys For tests purposes only.
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
	if model.SubscriptionUsesUUID(subscription) {
		return subscription.Status.ID, true
	}
	return refs.NewNamespacedNameFromObject(subscription).HRID(), false
}

func subscriptionAPIID(api core.ApiDefinitionObject) (string, bool) {
	if model.APIUsesUUID(api) {
		return api.GetID(), true
	}
	return refs.NewNamespacedNameFromObject(api).HRID(), false
}
