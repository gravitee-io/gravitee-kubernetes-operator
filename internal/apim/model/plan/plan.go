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

package plan

import v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"

type Status string

type Security struct {
	Type string `json:"type"`
}

type Entity struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Security Security `json:"security"`
	Status   Status   `json:"status"`
	Api      string   `json:"apiId"`
}

type List struct {
	Data []Entity `json:"data"`
}

type Updates struct {
	ToCreate []*v4.Plan
	ToUpdate map[string]*v4.Plan
	ToDelete []string
}

func PrepareUpdates() *Updates {
	return &Updates{
		ToCreate: []*v4.Plan{},
		ToUpdate: map[string]*v4.Plan{},
		ToDelete: []string{},
	}
}

type ApiKey struct {
	Id  string `json:"id"`
	Key string `json:"key"`
}

type ApiKeyList struct {
	Data []ApiKey `json:"data"`
}

type Subscription struct {
	Id string `json:"id"`
}

type SubscriptionRequest struct {
	PlanID        string `json:"planId"`
	ApplicationID string `json:"applicationId"`
}
