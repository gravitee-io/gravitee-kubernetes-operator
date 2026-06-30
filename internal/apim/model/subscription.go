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

package model

import "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/subscription"

type SubscriptionResponse struct {
	ID string `json:"id"`
}

type SubscriptionRequest struct {
	AppID  string `json:"applicationId"`
	PlanID string `json:"planId"`
}

type ApiKeySpec struct {
	Key      string  `json:"key"`
	ExpireAt *string `json:"expireAt,omitempty" drift:"rfc3339"`
}

type SubscriptionDTO struct {
	ID                    string                             `json:"id"`
	ApiID                 string                             `json:"apiId"`
	AppID                 string                             `json:"applicationId"`
	PlanID                string                             `json:"planId"`
	Status                string                             `json:"status"`
	StartingAt            string                             `json:"startingAt" drift:"rfc3339"`
	EndingAt              string                             `json:"endingAt" drift:"rfc3339"`
	Metadata              map[string]string                  `json:"metadata,omitempty" drift:"empty-is-nil"`
	ApiKeys               []ApiKeySpec                       `json:"apiKeys,omitempty" drift:"empty-is-nil"`
	ConsumerConfiguration subscription.ConsumerConfiguration `json:"consumerConfiguration,omitempty"`
}

type AutomationSubscriptionDTO struct {
	HRID                  string                             `json:"hrid"`
	ApplicationHrid       string                             `json:"applicationHrid"`
	PlanHrid              string                             `json:"planHrid"`
	ApiHrid               string                             `json:"apiHrid"`
	Status                string                             `json:"status"`
	StartingAt            string                             `json:"startingAt"`
	EndingAt              string                             `json:"endingAt"`
	Metadata              map[string]string                  `json:"metadata,omitempty"`
	ApiKeys               []ApiKeySpec                       `json:"apiKeys,omitempty"`
	ConsumerConfiguration subscription.ConsumerConfiguration `json:"consumerConfiguration,omitempty"`
}

func (a *AutomationSubscriptionDTO) ToLegacy() *SubscriptionDTO {
	return &SubscriptionDTO{
		ID:                    a.HRID,
		ApiID:                 a.ApiHrid,
		AppID:                 a.ApplicationHrid,
		PlanID:                a.PlanHrid,
		Status:                a.Status,
		StartingAt:            a.StartingAt,
		EndingAt:              a.EndingAt,
		Metadata:              a.Metadata,
		ApiKeys:               a.ApiKeys,
		ConsumerConfiguration: a.ConsumerConfiguration,
	}
}

type SubscriptionStatus struct {
	ID         string `json:"id,omitempty"`
	StartingAt string `json:"startingAt,omitempty"`
	EndingAt   string `json:"endingAt,omitempty"`
}

func (s *SubscriptionDTO) ToAutomation() AutomationSubscriptionDTO {
	return AutomationSubscriptionDTO{
		HRID:                  s.ID,
		ApiHrid:               s.ApiID,
		ApplicationHrid:       s.AppID,
		PlanHrid:              s.PlanID,
		Status:                s.Status,
		StartingAt:            s.StartingAt,
		EndingAt:              s.EndingAt,
		Metadata:              s.Metadata,
		ApiKeys:               s.ApiKeys,
		ConsumerConfiguration: s.ConsumerConfiguration,
	}
}
