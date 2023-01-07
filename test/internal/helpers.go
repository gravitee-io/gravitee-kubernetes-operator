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

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
)

func GetStatusContext(apiDefinition *gio.ApiDefinition, location types.NamespacedName) *gio.StatusContext {
	contexts := apiDefinition.Status.Contexts

	if contexts == nil {
		return nil
	}

	if len(contexts) == 0 {
		return nil
	}

	context, ok := contexts[location.String()]

	if !ok {
		return nil
	}

	return &context
}

func GetStatusId(apiDefinition *gio.ApiDefinition, location types.NamespacedName) string {
	context := GetStatusContext(apiDefinition, location)

	if context == nil {
		return ""
	}

	return context.ID
}

func NewAssertionError(field string, expected, given any) error {
	return fmt.Errorf("expected %s to be %v, got %v", field, expected, given)
}
