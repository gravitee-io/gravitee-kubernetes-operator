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

package base

import (
	"context"
	"net/url"
	"slices"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

func validateNoConflictingPath(ctx context.Context, api core.ApiDefinitionObject) *errors.AdmissionError {
	apiPaths := api.GetContextPaths()
	existingPaths, err := getExistingPaths(ctx, api)
	if err != nil {
		return errors.NewSevere(err.Error())
	}
	for _, apiPath := range apiPaths {
		if _, err := url.Parse(apiPath); err != nil {
			return errors.NewSevere(
				"path [%s] is invalid",
				apiPath,
			)
		}
		if slices.Contains(existingPaths, apiPath) {
			return errors.NewSevere(
				"invalid API context path [%s]. Another API with the same path already exists",
				apiPath,
			)
		}
	}
	return nil
}

func getExistingPaths(ctx context.Context, api core.ApiDefinitionObject) ([]string, error) {
	existingPaths := make([]string, 0)
	apis, err := dynamic.GetAPIs(ctx, dynamic.ListOptions{
		Namespace: api.GetNamespace(),
		Excluded: []core.ObjectRef{
			api.GetRef(),
		},
	})
	if err != nil {
		return existingPaths, err
	}

	for _, api := range apis {
		apiPaths := api.GetContextPaths()
		existingPaths = append(existingPaths, apiPaths...)
	}
	return existingPaths, nil
}
