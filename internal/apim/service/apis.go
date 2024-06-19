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
	"net/http"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"

	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	xhttp "github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
)

const (
	crossIDParam     = "crossId"
	stateActionParam = "action"
	planParam        = "plan"
	applicationParam = "application"
)

var importParams = map[string]string{
	"definitionVersion": "2.0.0",
}

var deleteParams = map[string]string{
	"closePlans": "true",
}

// APIs brings support for managing gravitee.io APIM APIs.
type APIs struct {
	*client.Client
}

type httpImportMethod = func(xhttp.URL, any, any, ...xhttp.RequestTransformer) error

func NewAPIs(client *client.Client) *APIs {
	return &APIs{Client: client}
}

func (svc *APIs) GetByCrossID(crossID string) (*model.ApiListItem, error) {
	url := svc.EnvV1Target("apis").WithQueryParam(crossIDParam, crossID)
	apis := new([]model.ApiListItem)

	if err := svc.HTTP.Get(url, apis); err != nil {
		return nil, err
	}

	if len(*apis) == 0 {
		return nil, errors.NewNotFoundError()
	}

	return &(*apis)[0], nil
}

func (svc *APIs) GetByID(apiID string) (*model.ApiEntity, error) {
	url := svc.EnvV1Target("apis").WithPath(apiID)
	api := new(model.ApiEntity)

	if err := svc.HTTP.Get(url, api); err != nil {
		return nil, err
	}

	return api, nil
}

func (svc *APIs) GetV4ByID(apiID string) (*v1alpha1.ApiV4DefinitionSpec, error) {
	url := svc.EnvV2Target("apis").WithPath(apiID)
	resp := new(v1alpha1.ApiV4DefinitionSpec)

	if err := svc.HTTP.Get(url, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (svc *APIs) ImportV2(method string, spec *v2.Api) (*model.ApiEntity, error) {
	url := svc.EnvV1Target("apis/import").WithQueryParams(importParams)
	apiImport := model.ToApiImport(spec)
	api := new(model.ApiEntity)
	fun := svc.getImportFunc(method)

	if err := fun(url, apiImport, api); err != nil {
		return nil, err
	}

	return api, nil
}

func (svc *APIs) ImportV4(spec *v4.Api) (*v1alpha1.ApiV4DefinitionStatus, error) {
	url := svc.EnvV2Target("apis/_import/crd")

	status := new(v1alpha1.ApiV4DefinitionStatus)
	if err := svc.HTTP.Put(url, spec, status); err != nil {
		return nil, err
	}

	return status, nil
}

func (svc *APIs) getImportFunc(method string) httpImportMethod {
	if method == http.MethodPost {
		return svc.HTTP.Post
	}
	return svc.HTTP.Put
}

func (svc *APIs) UpdateState(apiID string, action model.Action) error {
	url := svc.EnvV1Target("apis").WithPath(apiID).WithQueryParam(stateActionParam, string(action))
	return svc.HTTP.Post(url, nil, nil)
}

func (svc *APIs) DeleteV2(apiID string) error {
	url := svc.EnvV1Target("apis").WithPath(apiID).WithQueryParams(deleteParams)
	return svc.HTTP.Delete(url, nil)
}

func (svc *APIs) DeleteV4(apiID string) error {
	url := svc.EnvV2Target("apis").WithPath(apiID).WithQueryParams(deleteParams)
	return svc.HTTP.Delete(url, nil)
}

func (svc *APIs) SetKubernetesContext(apiID string) error {
	url := svc.EnvV1Target("apis").WithPath(apiID).WithPath("definition-context")
	return svc.HTTP.Put(url, model.NewKubernetesContext(), nil)
}

func (svc *APIs) Deploy(id string) error {
	url := svc.EnvV1Target("apis").WithPath(id).WithPath("deploy")
	return svc.HTTP.Post(url, new(model.ApiDeployment), nil)
}

func (svc *APIs) ExportV2(id string) (*v1alpha1.ApiDefinition, error) {
	url := svc.EnvV1Target("apis").WithPath(id).WithPath("/crd")
	api := new(v1alpha1.ApiDefinition)
	if err := svc.HTTP.GetYAML(url, api); err != nil {
		return nil, err
	}
	return api, nil
}
