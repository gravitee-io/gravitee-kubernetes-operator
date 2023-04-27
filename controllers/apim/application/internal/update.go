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
	"net/http"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimModel "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

func (d *Delegate) CreateOrUpdate(application *gio.Application) error {
	if err := d.createUpdateApplication(application); err != nil {
		return err
	}

	return d.createUpdateApplicationMetadata(application)
}

func (d *Delegate) createUpdateApplication(application *gio.Application) error {
	spec := &application.Spec
	spec.Origin = "KUBERNETES"
	app, err := d.apim.Applications.GetByID(application.Status.ID)
	if errors.IgnoreNotFound(err) != nil {
		return apim.NewContextError(err)
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

	mgmtApp, mgmtErr := d.apim.Applications.CreateUpdate(method, &spec.Application)
	if mgmtErr != nil {
		return apim.NewContextError(mgmtErr)
	}

	spec.ID = mgmtApp.Id
	application.Status.ID = mgmtApp.Id

	return nil
}

func (d *Delegate) createUpdateApplicationMetadata(application *gio.Application) error {
	spec := &application.Spec
	if spec.ApplicationMetaData == nil {
		application.Status.Status = gio.ProcessingStatusCompleted
		return nil
	}

	appMetaData, err := d.apim.Applications.GetMetadataByApplicationID(application.Status.ID)
	if err != nil {
		return apim.NewContextError(err)
	}

	for _, metaData := range *spec.ApplicationMetaData {
		method := http.MethodPost
		key := d.findMetadataKey(appMetaData, metaData.Name)
		if key != "" {
			// update
			metaData.Key = key
			method = http.MethodPut
		}

		_, mgmtErr := d.apim.Applications.CreateUpdateMetadata(method, spec.ID, metaData)
		if mgmtErr != nil {
			return apim.NewContextError(mgmtErr)
		}
	}

	// Delete removed metadata
	for _, metaData := range *appMetaData {
		if d.metadataIsRemoved(spec.ApplicationMetaData, metaData.Name) {
			err = d.apim.Applications.DeleteMetadata(application.Status.ID, metaData.Key)
			if errors.IgnoreNotFound(err) != nil {
				return err
			}
		}
	}

	application.Status.Status = gio.ProcessingStatusCompleted

	return nil
}

func (d *Delegate) findMetadataKey(appMetadata *[]apimModel.ApplicationMetaData, name string) string {
	for _, md := range *appMetadata {
		if md.Name == name {
			return md.Key
		}
	}

	return ""
}

func (d *Delegate) metadataIsRemoved(metaData *[]model.ApplicationMetaData, name string) bool {
	for _, md := range *metaData {
		if md.Name == name {
			return false
		}
	}

	return true
}
