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
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1beta1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model/api"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
)

// APIs brings support for managing gravitee.io APIM APIs.
type APIs struct {
	*client.Client
}

func NewAPIs(client *client.Client) *APIs {
	return &APIs{Client: client}
}

func (svc *APIs) baseURL() *http.URL {
	return svc.EnvV2("apis")
}

func (svc *APIs) GetByID(apiID string) (*api.Entity, error) {
	url := svc.baseURL().WithPath(apiID)
	api := new(api.Entity)

	if err := svc.HTTP.Get(url, api); err != nil {
		return nil, err
	}

	return api, nil
}

func (svc *APIs) Import(spec *v4.Api) (*v1beta1.ApiDefinitionStatus, error) {
	url := svc.baseURL().WithPath("_import/crd")

	status := new(v1beta1.ApiDefinitionStatus)
	if err := svc.HTTP.Put(url, spec, status); err != nil {
		return nil, err
	}

	return status, nil
}

func (svc *APIs) Delete(apiID string) error {
	url := svc.baseURL().WithPath(apiID).WithQueryParam("closePlans", "true")
	return svc.HTTP.Delete(url, nil)
}

func (svc *APIs) Stop(apiID string) error {
	url := svc.baseURL().WithPath(apiID).WithPath("_stop")
	return svc.HTTP.Post(url, nil, nil)
}
