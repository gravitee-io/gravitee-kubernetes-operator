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

package apim

import (
	"context"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimapp "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
)

// Application is what the specs read of an APIM application.
type Application struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// Applications reads and writes APIM applications outside the operator.
type Applications struct {
	ctx  context.Context
	sdk  *apim.SDKClient
	mAPI *client.Client
}

func NewApplications(ctx context.Context, sdkClient *apim.SDKClient, mAPI *client.Client) *Applications {
	return &Applications{ctx: ctx, sdk: sdkClient, mAPI: mAPI}
}

// GetByHRID gets an application by HRID.
func (a *Applications) GetByHRID(hrid string) (*Application, error) {
	state, err := apimapp.GetState(a.ctx, a.sdk, apimapp.DTO{Application: sdk.Application{Hrid: hrid}})
	if err != nil {
		return nil, err
	}
	app := &Application{Name: state.Name}
	if state.Id != nil {
		app.ID = *state.Id
	}
	if state.Status != nil {
		app.Status = string(*state.Status)
	}
	return app, nil
}

// GetMetadataByApplicationID gets the metadata of an application by its UUID.
func (a *Applications) GetMetadataByApplicationID(id string) (*[]sdk.Metadata, error) {
	dto := apimapp.DTO{Application: sdk.Application{Hrid: id}, Legacy: apimapp.Legacy{HridContainsUUID: true}}
	state, err := apimapp.GetState(a.ctx, a.sdk, dto)
	if err != nil {
		return nil, err
	}
	return &state.Metadata, nil
}

// Search lists the applications matching query and status, through the Management API: the SDK has no list.
func (a *Applications) Search(query, status string) ([]Application, error) {
	url := a.mAPI.EnvV1Target("/applications").WithQueryParam("query", query).WithQueryParam("status", status)
	applications := new([]Application)
	if err := a.mAPI.HTTP.Get(url.String(), applications); err != nil {
		return nil, err
	}
	return *applications, nil
}

// CreateOrUpdate upserts an application in APIM. An application whose spec.hrid is set, as PopulateIDs does,
// is upserted by that hrid.
func (a *Applications) CreateOrUpdate(app *v1alpha1.Application) (any, error) {
	dto, err := apimapp.ToDTO(app)
	if err != nil {
		return nil, err
	}
	if app.Spec.HRID != "" {
		dto.Hrid = app.Spec.HRID
		dto.Legacy = apimapp.Legacy{}
	}
	return apimapp.Upsert(a.ctx, a.sdk, dto)
}
