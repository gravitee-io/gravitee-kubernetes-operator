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

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
)

const (
	applicationsPath = "/applications"
	metadataPath     = "/metadata"
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
	url := svc.EnvV1Target(applicationsPath).WithQueryParam("query", query).WithQueryParam("status", status)
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

	url := svc.EnvV1Target(applicationsPath).WithPath(appId)
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

	url := svc.EnvV1Target(applicationsPath).WithPath(appId).WithPath(metadataPath)
	application := new([]model.ApplicationMetaData)

	if err := svc.HTTP.Get(url.String(), application); err != nil {
		return nil, err
	}

	return application, nil
}

func (svc *Applications) CreateOrUpdate(spec *application.Application) (*application.Status, error) {
	url := svc.EnvV2Target(fmt.Sprintf("%s/_import/crd", applicationsPath))

	status := new(application.Status)
	if err := svc.HTTP.Put(url.String(), spec, status); err != nil {
		return nil, err
	}

	return status, nil
}

func (svc *Applications) Delete(appId string) error {
	url := svc.EnvV1Target(applicationsPath).WithPath(appId)
	return svc.HTTP.Delete(url.String(), nil)
}

func (svc *Applications) CreateOrUpdateMetadata(method string, appId string,
	spec any, key string) (*model.ApplicationMetaData, error) {
	url := svc.EnvV1Target(applicationsPath).WithPath(appId).WithPath(metadataPath)
	fun := svc.HTTP.Post
	if method == http.MethodPut {
		url = url.WithPath(key)
		fun = svc.HTTP.Put
	}

	md := new(model.ApplicationMetaData)
	if err := fun(url.String(), spec, md); err != nil {
		return nil, err
	}

	return md, nil
}

func (svc *Applications) DeleteMetadata(appId string, key string) error {
	url := svc.EnvV1Target(applicationsPath).WithPath(appId).WithPath(metadataPath).WithPath(key)
	return svc.HTTP.Delete(url.String(), nil)
}
