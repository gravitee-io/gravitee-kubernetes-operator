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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1beta1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model/api"
)

func AssertApplicationStatusIsSet(application *v1beta1.Application) error {
	status := application.Status

	if status.ID == "" {
		return fmt.Errorf("id should not be empty in status")
	}

	if status.EnvID == "" {
		return fmt.Errorf("envId should not be empty in status")
	}

	if status.OrgID == "" {
		return fmt.Errorf("orgId should not be empty in status")
	}

	if status.Status == "" {
		return fmt.Errorf("status should not be empty in status")
	}

	return nil
}

func AssertApiEntityMatchesStatus(apiEntity *api.Entity, apiDefinition *v1alpha1.ApiDefinition) error {
	if apiEntity.ID != apiDefinition.Status.ID {
		return NewAssertionError("Status id", apiEntity.ID, apiDefinition.Status.ID)
	}
	return nil
}

func AssertStatusMatches(
	apiDefinition *v1alpha1.ApiDefinition, expectedStatus v1alpha1.ApiDefinitionStatus,
) error {
	if apiDefinition.Status.ID != expectedStatus.ID {
		return NewAssertionError("id", expectedStatus.ID, apiDefinition.Status.ID)
	}
	if apiDefinition.Status.CrossID != expectedStatus.CrossID {
		return NewAssertionError("crossId", expectedStatus.CrossID, apiDefinition.Status.CrossID)
	}
	if apiDefinition.Status.EnvID != expectedStatus.EnvID {
		return NewAssertionError("envId", expectedStatus.EnvID, apiDefinition.Status.EnvID)
	}
	if apiDefinition.Status.OrgID != expectedStatus.OrgID {
		return NewAssertionError("orgId", expectedStatus.OrgID, apiDefinition.Status.OrgID)
	}
	if apiDefinition.Status.Status != expectedStatus.Status {
		return NewAssertionError("status", expectedStatus.Status, apiDefinition.Status.Status)
	}
	if apiDefinition.Status.State != expectedStatus.State {
		return NewAssertionError("state", expectedStatus.State, apiDefinition.Status.State)
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

func AssertStatusCompleted(
	apiDefinition *v1alpha1.ApiDefinition,
) error {
	return AssertEquals("status", v1alpha1.ProcessingStatusCompleted, apiDefinition.Status.Status)
}

func AssertNoErrorAndStatusCompleted(
	err error, apiDefinition *v1alpha1.ApiDefinition,
) error {
	if err != nil {
		return err
	}
	return AssertEquals("status", v1alpha1.ProcessingStatusCompleted, apiDefinition.Status.Status)
}

func AssertEquals(property string, expected, actual interface{}) error {
	if !reflect.DeepEqual(expected, actual) {
		return NewAssertionError(property, expected, actual)
	}
	return nil
}

func AssertHostPrefix(hostname *Host, prefix string) error {
	if !hostname.StartsWith(prefix) {
		return fmt.Errorf("hostname %s does not start with %s", hostname, prefix)
	}
	return nil
}
