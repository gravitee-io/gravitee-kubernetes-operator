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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi/model"
)

type Applications struct {
	*Client
}

func NewApplications(client *Client) *Applications {
	return &Applications{Client: client}
}

func (svc *Applications) Search(query string, status string) ([]model.Application, error) {
	url := svc.EnvTarget("/applications").WithQueryParam("query", query).WithQueryParam("status", status)
	applications := new([]model.Application)

	if err := svc.http.Get(url.String(), applications); err != nil {
		return nil, err
	}

	return *applications, nil
}
