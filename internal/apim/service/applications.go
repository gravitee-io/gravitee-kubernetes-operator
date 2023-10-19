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

	v1alpha1 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model/application"
	xhttp "github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
)

const (
	metadataPath = "/metadata"
)

// Applications brings support for managing gravitee.io APIM applications
// This service is used for testing purposes only and not initialized by the operator manager.
type Applications struct {
	*client.Client
}

func NewApplications(client *client.Client) *Applications {
	return &Applications{Client: client}
}

func (svc *Applications) Target() *xhttp.URL {
	return svc.EnvV1("applications")
}

func (svc *Applications) Search(query string, status string) ([]application.Entity, error) {
	url := svc.Target().WithQueryParam("query", query).WithQueryParam("status", status)
	applications := new([]application.Entity)

	if err := svc.HTTP.Get(url, applications); err != nil {
		return nil, err
	}

	return *applications, nil
}

func (svc *Applications) GetByID(appId string) (*application.Entity, error) {
	if appId == "" {
		return nil, errors.NewNotFoundError()
	}

	url := svc.Target().WithPath(appId)
	application := new(application.Entity)

	if err := svc.HTTP.Get(url, application); err != nil {
		return nil, err
	}

	return application, nil
}

func (svc *Applications) GetMetadataByApplicationID(appId string) (*[]application.MetaData, error) {
	if appId == "" {
		return nil, fmt.Errorf("can't retrieve metadata without application id")
	}

	url := svc.Target().WithPath(appId).WithPath(metadataPath)
	application := new([]application.MetaData)

	if err := svc.HTTP.Get(url, application); err != nil {
		return nil, err
	}

	return application, nil
}

func (svc *Applications) CreateUpdate(method string, spec *v1alpha1.Application) (*application.Entity, error) {
	url := svc.Target()
	if spec.ID != "" {
		url = svc.Target().WithPath(spec.ID)
	}

	fun := svc.HTTP.Post
	if method == http.MethodPut {
		fun = svc.HTTP.Put
	}

	application := new(application.Entity)
	if err := fun(url, spec, application); err != nil {
		return nil, err
	}

	return application, nil
}

func (svc *Applications) Delete(appId string) error {
	url := svc.Target().WithPath(appId)
	return svc.HTTP.Delete(url, nil)
}

func (svc *Applications) CreateUpdateMetadata(method string, appId string,
	spec v1alpha1.MetaData) (*application.MetaData, error) {
	url := svc.Target().WithPath(appId).WithPath(metadataPath)
	fun := svc.HTTP.Post
	if method == http.MethodPut {
		url = url.WithPath(spec.Key)
		fun = svc.HTTP.Put
	}

	md := new(application.MetaData)
	if err := fun(url, spec, md); err != nil {
		return nil, err
	}

	return md, nil
}

func (svc *Applications) DeleteMetadata(appId string, key string) error {
	url := svc.Target().WithPath(appId).WithPath(metadataPath).WithPath(key)
	return svc.HTTP.Delete(url, nil)
}
