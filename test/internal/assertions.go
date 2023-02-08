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
	"fmt"
	"net/http"
	"reflect"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
)

func AssertStatusIsSet(apiDefinition *gio.ApiDefinition) error {
	status := apiDefinition.Status

	if status.ID == "" {
		return fmt.Errorf("id should not be empty in status")
	}

	if status.EnvID == "" {
		return fmt.Errorf("envId should not be empty in status")
	}

	if status.OrgID == "" {
		return fmt.Errorf("envId should not be empty in status")
	}

	if status.Status == "" {
		return fmt.Errorf("status should not be empty in status")
	}

	if status.State == "" {
		return fmt.Errorf("state should not be empty in status")
	}

	return nil
}

func AssertApiEntityMatchesStatus(apiEntity *model.ApiEntity, apiDefinition *gio.ApiDefinition) error {
	if apiEntity.ID != apiDefinition.Status.ID {
		return NewAssertionError("Status id", apiEntity.ID, apiDefinition.Status.ID)
	}
	return nil
}

func AssertStatusMatches(
	apiDefinition *gio.ApiDefinition, expectedStatus gio.ApiDefinitionStatus,
) error {
	if !reflect.DeepEqual(apiDefinition.Status, expectedStatus) {
		return NewAssertionError("status", expectedStatus, apiDefinition.Status.Status)
	}

	return nil
}

func AssertNoErrorAndHTTPStatus(err error, res *http.Response, expectedStatus int) error {
	if err != nil {
		return err
	}
	if res.StatusCode != expectedStatus {
		return NewAssertionError("status", expectedStatus, res.StatusCode)
	}
	return nil
}

func AssertNoErrorAndObservedGenerationEquals(
	err error, apiDefinition *gio.ApiDefinition, expectedGeneration int64,
) error {
	if err != nil {
		return err
	}
	if apiDefinition.Status.ObservedGeneration != expectedGeneration {
		return NewAssertionError(
			"observedGeneration", expectedGeneration, apiDefinition.Status.ObservedGeneration,
		)
	}
	return nil
}

func AssertEquals(property string, expected, actual interface{}) error {
	if !reflect.DeepEqual(expected, actual) {
		return NewAssertionError(property, expected, actual)
	}
	return nil
}
