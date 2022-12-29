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
	"k8s.io/apimachinery/pkg/types"
)

func AssertStatusContextIsSet(apiDefinition *gio.ApiDefinition) error {
	contexts := apiDefinition.Status.Contexts

	if contexts == nil {
		return fmt.Errorf("contexts should not be nil")
	}

	if len(contexts) == 0 {
		return fmt.Errorf("status contexts should not be empty")
	}

	for location, context := range contexts {
		if context.ID == "" {
			return fmt.Errorf("id should not be empty for context %s", location)
		}

		if context.EnvID == "" {
			return fmt.Errorf("envId should not be empty for context %s", location)
		}

		if context.OrgID == "" {
			return fmt.Errorf("envId should not be empty for context %s", location)
		}

		if context.Status == "" {
			return fmt.Errorf("status should not be empty for context %s", location)
		}

		if context.State == "" {
			return fmt.Errorf("state should not be empty for context %s", location)
		}
	}

	return nil
}

func AssertApiEntityMatchesStatusContext(apiEntity *model.ApiEntity, apiDefinition *gio.ApiDefinition) error {
	contexts := apiDefinition.Status.Contexts

	if contexts == nil {
		return fmt.Errorf("contexts should not be nil")
	}

	if len(contexts) == 0 {
		return fmt.Errorf("status contexts should not be empty")
	}

	found := false

	for _, context := range contexts {
		if context.ID == apiEntity.ID {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("api %s not found in status", apiEntity.ID)
	}

	return nil
}

func AssertStatusContextMatches(
	apiDefinition *gio.ApiDefinition, location types.NamespacedName, expectedContext *gio.StatusContext,
) error {
	context := GetStatusContext(apiDefinition, location)

	if context == nil {
		return fmt.Errorf("context %s not found in status", location)
	}

	if !reflect.DeepEqual(context, expectedContext) {
		return fmt.Errorf(
			"expected status context %s to match %v, got %v ",
			location, expectedContext, context,
		)
	}

	return nil
}

func AssertNoErrorAndHTTPStatus(err error, res *http.Response, expectedStatus int) error {
	if err != nil {
		return err
	}
	if res.StatusCode != expectedStatus {
		return fmt.Errorf("expected status %d, got %d", expectedStatus, res.StatusCode)
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
		return fmt.Errorf(
			"expected observed generation %d, got %d",
			expectedGeneration,
			apiDefinition.Status.ObservedGeneration,
		)
	}
	return nil
}
