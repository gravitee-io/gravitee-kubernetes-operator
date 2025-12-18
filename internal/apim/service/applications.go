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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"strconv"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
)

const (
	applicationsPath = "/applications"
	metadataPath     = "/metadata"
)

// Applications brings support for managing gravitee.io APIM applications.
type Applications struct {
	*client.Client
}

func NewApplications(client *client.Client) *Applications {
	return &Applications{Client: client}
}

// Search For tests purposes only.
func (svc *Applications) Search(query string, status string) ([]model.Application, error) {
	url := svc.EnvV1Target(applicationsPath).WithQueryParam("query", query).WithQueryParam("status", status)
	applications := new([]model.Application)

	if err := svc.HTTP.Get(url.String(), applications); err != nil {
		return nil, err
	}

	return *applications, nil
}

// GetByID For tests purposes only.
func (svc *Applications) GetByID(appID string) (*model.Application, error) {
	if appID == "" {
		return nil, errors.NewNotFoundError()
	}

	url := svc.EnvV1Target(applicationsPath).WithPath(appID)
	app := new(model.Application)

	if err := svc.HTTP.Get(url.String(), app); err != nil {
		return nil, err
	}

	return app, nil
}

// GetMetadataByApplicationID For tests purposes only.
func (svc *Applications) GetMetadataByApplicationID(appID string) (*[]model.ApplicationMetaData, error) {
	if appID == "" {
		return nil, fmt.Errorf("can't retrieve metadata without application id")
	}

	url := svc.EnvV1Target(applicationsPath).WithPath(appID).WithPath(metadataPath)
	app := new([]model.ApplicationMetaData)

	if err := svc.HTTP.Get(url.String(), app); err != nil {
		return nil, err
	}

	return app, nil
}

func (svc *Applications) CreateOrUpdate(app *v1alpha1.Application) (*application.Status, error) {
	return svc.createOrUpdate(app, false)
}

func (svc *Applications) DryRunCreateOrUpdate(app *v1alpha1.Application) (*application.Status, error) {
	return svc.createOrUpdate(app, true)
}

func (svc *Applications) createOrUpdate(app *v1alpha1.Application, dryRun bool) (*application.Status, error) {

	url := svc.AutomationTarget(applicationsPath).
		WithQueryParam("dryRun", strconv.FormatBool(dryRun))

	usesHRID := app.Spec.HRID != ""

	// even if the resource was created with a previous version,
	// Automation API still need an HRID, so one is set.
	// API ignores it because an ID is set in that case see PopulateIDs().
	if !usesHRID {
		app.Spec.HRID = app.GetID()
		url = url.WithQueryParam("legacyID", strconv.FormatBool(true))
	}
	status := new(application.Status)

	if err := svc.HTTP.Put(url.String(), app.Spec, status); err != nil {
		return nil, err
	}

	status.UseHRID = usesHRID

	return status, nil
}

func (svc *Applications) Delete(app *v1alpha1.Application) error {
	appID, legacy := getAppID(app)
	url := svc.AutomationTarget(applicationsPath).WithPath(appID).WithQueryParam("legacy", strconv.FormatBool(legacy))
	return svc.HTTP.Delete(url.String(), nil)
}

func getAppID(application *v1alpha1.Application) (string, bool) {
	var id string
	var legacy bool
	if application.Status.UseHRID {
		id = refs.NewNamespacedNameFromObject(application).HRID()
		legacy = false
	} else {
		id = application.GetID()
		legacy = true
	}
	return id, legacy
}
