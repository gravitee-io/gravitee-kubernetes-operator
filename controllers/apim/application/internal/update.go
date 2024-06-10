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

package internal

import (
	"context"
	"net/http"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimModel "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/types/k8s/custom"
)

func CreateOrUpdate(ctx context.Context, application *v1alpha1.Application) error {
	if err := createUpdateApplication(ctx, application); err != nil {
		return err
	}

	return createUpdateApplicationMetadata(ctx, application)
}

func createUpdateApplication(ctx context.Context, application *v1alpha1.Application) error {
	spec := &application.Spec
	spec.Origin = "KUBERNETES"

	apim, err := apim.FromContextRef(ctx, spec.Context)
	if err != nil {
		return err
	}

	app, err := apim.Applications.GetByID(application.Status.ID)
	if errors.IgnoreNotFound(err) != nil {
		return errors.NewContextError(err)
	}

	method := http.MethodPost
	if app != nil {
		method = http.MethodPut
		spec.ID = app.Id
		// to avoid getting error from APIM because of having no settings
		if spec.Settings == nil {
			spec.Settings = app.Settings
		}
	}

	mgmtApp, mgmtErr := apim.Applications.CreateUpdate(method, &spec.Application)
	if mgmtErr != nil {
		return errors.NewContextError(mgmtErr)
	}

	spec.ID = mgmtApp.Id
	application.Status.ID = mgmtApp.Id
	application.Status.EnvID = apim.EnvID()
	application.Status.OrgID = apim.OrgID()

	return nil
}

func createUpdateApplicationMetadata(ctx context.Context, app *v1alpha1.Application) error {
	spec := &app.Spec
	if spec.MetaData == nil {
		app.Status.Status = custom.ProcessingStatusCompleted
		return nil
	}

	apimCli, err := apim.FromContextRef(ctx, spec.Context)
	if err != nil {
		return err
	}

	appMetaData, err := apimCli.Applications.GetMetadataByApplicationID(app.Status.ID)
	if err != nil {
		return errors.NewContextError(err)
	}

	// All this method Will be removed once we created a dedicated endpoint in APIM for creating Application using GKO
	for _, metaData := range *spec.MetaData {
		md := struct {
			application.MetaData
			ApplicationId string `json:"applicationId"`
			Key           string `json:"key,omitempty"`
		}{
			MetaData:      metaData,
			ApplicationId: spec.ID,
		}
		method := http.MethodPost
		key := findMetadataKey(appMetaData, md.Name)
		if key != "" {
			// update
			md.Key = key
			method = http.MethodPut
		}

		_, mgmtErr := apimCli.Applications.CreateUpdateMetadata(method, spec.ID, md, key)
		if mgmtErr != nil {
			return errors.NewContextError(mgmtErr)
		}
	}

	// Delete removed metadata
	for _, metaData := range *appMetaData {
		if metadataIsRemoved(spec.MetaData, metaData.Name) {
			err = apimCli.Applications.DeleteMetadata(app.Status.ID, metaData.Key)
			if errors.IgnoreNotFound(err) != nil {
				return err
			}
		}
	}

	app.Status.Status = custom.ProcessingStatusCompleted

	return nil
}

func findMetadataKey(appMetadata *[]apimModel.ApplicationMetaData, name string) string {
	for _, md := range *appMetadata {
		if md.Name == name {
			return md.Key
		}
	}

	return ""
}

func metadataIsRemoved(metaData *[]application.MetaData, name string) bool {
	for _, md := range *metaData {
		if md.Name == name {
			return false
		}
	}

	return true
}
