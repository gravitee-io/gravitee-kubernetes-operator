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

type Plans struct {
	*client.Client
}

func NewPlans(client *client.Client) *Plans {
	return &Plans{Client: client}
}

func (svc *Plans) apiURL(apiID string) *http.URL {
	return svc.EnvV2("apis").WithPath(apiID).WithPath("plans")
}

func (svc *Plans) ListByAPI(apiID string) (*plan.List, error) {
	url := svc.apiURL(apiID)
	list := new(plan.List)

	if err := svc.HTTP.Get(url, list); err != nil {
		return nil, err
	}

	return list, nil
}
