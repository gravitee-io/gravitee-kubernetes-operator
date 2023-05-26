// Copyright (C) 2015 The Gravitee team (HTTP://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         HTTP://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"fmt"
	"net/http"

	errors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"

	kModel "github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
)

// Applications brings support for managing gravitee.io APIM applications
// This service is used for testing purposes only and not initialized by the operator manager.
type Applications struct {
	*client.Client
}

func NewApplications(client *client.Client) *Applications {
	return &Applications{Client: client}
}

func (svc *Applications) Search(query string, status string) ([]model.Application, error) {
	url := svc.EnvTarget("/applications").WithQueryParam("query", query).WithQueryParam("status", status)
	applications := new([]model.Application)

	if err := svc.HTTP.Get(url.String(), applications); err != nil {
		return nil, err
	}

	return *applications, nil
}

func (svc *Applications) GetByID(appId string) (*model.Application, error) {
	if appId == "" {
		return nil, errors.NewNotFoundError()
	}

	url := svc.EnvTarget(fmt.Sprintf("applications/%s", appId))
	application := new(model.Application)

	if err := svc.HTTP.Get(url.String(), application); err != nil {
		return nil, err
	}

	return application, nil
}

func (svc *Applications) GetMetadataByApplicationID(appId string) (*[]model.ApplicationMetaData, error) {
	if appId == "" {
		return nil, fmt.Errorf("can't retrieve metadata without application id")
	}

	url := svc.EnvTarget(fmt.Sprintf("applications/%s/metadata", appId))
	application := new([]model.ApplicationMetaData)

	if err := svc.HTTP.Get(url.String(), application); err != nil {
		return nil, err
	}

	return application, nil
}

func (svc *Applications) CreateUpdate(method string, spec *kModel.Application) (*model.Application, error) {
	url := svc.EnvTarget("applications")
	if spec.ID != "" {
		url = svc.EnvTarget(fmt.Sprintf("applications/%s", spec.ID))
	}

	fun := svc.HTTP.Post
	if method == http.MethodPut {
		fun = svc.HTTP.Put
	}

	application := new(model.Application)
	if err := fun(url.String(), spec, application); err != nil {
		return nil, err
	}

	return application, nil
}

func (svc *Applications) Delete(appId string) error {
	url := svc.EnvTarget(fmt.Sprintf("applications/%s", appId))
	return svc.HTTP.Delete(url.String(), nil)
}

func (svc *Applications) CreateUpdateMetadata(method string, appId string,
	spec kModel.ApplicationMetaData) (*model.ApplicationMetaData, error) {
	url := svc.EnvTarget(fmt.Sprintf("applications/%s/metadata", appId))
	fun := svc.HTTP.Post
	if method == http.MethodPut {
		url = svc.EnvTarget(fmt.Sprintf("applications/%s/metadata/%s", appId, spec.Key))
		fun = svc.HTTP.Put
	}

	md := new(model.ApplicationMetaData)
	if err := fun(url.String(), spec, md); err != nil {
		return nil, err
	}

	return md, nil
}

func (svc *Applications) DeleteMetadata(appId string, key string) error {
	url := svc.EnvTarget(fmt.Sprintf("applications/%s/metadata/%s", appId, key))
	return svc.HTTP.Delete(url.String(), nil)
}
