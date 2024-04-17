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

type Pages struct {
	*client.Client
}

func NewPages(client *client.Client) *Pages {
	return &Pages{Client: client}
}

func (svc *Pages) FindByAPI(apiID string, queries ...*model.PagesQuery) ([]model.Page, error) {
	url := svc.EnvV1Target("apis").WithPath(apiID).WithPath("pages")

	for _, q := range queries {
		if q.Type != "" {
			url = url.WithQueryParams(q.AsMap())
		}
	}

	pages := new([]model.Page)

	if err := svc.HTTP.Get(url.String(), pages); err != nil {
		return nil, err
	}

	return *pages, nil
}
