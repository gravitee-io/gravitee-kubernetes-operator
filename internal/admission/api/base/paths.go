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
	"path/filepath"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

func ValidateNoConflictingPath(ctx context.Context, api core.ApiDefinitionObject) *errors.AdmissionError {
	if k8s.HasHTTPRouteOwner(api.GetOwnerReferences()) {
		return nil
	}

	apiPaths := api.GetContextPaths()
	existingPaths, err := getExistingPaths(ctx, api)
	if err != nil {
		return errors.NewSevere(err.Error())
	}
	for _, apiPath := range apiPaths {
		if _, err := url.Parse(apiPath); err != nil {
			return errors.NewSeveref(
				"path [%s] is invalid",
				apiPath,
			)
		}

		if isConflictingPath(existingPaths, apiPath) {
			return errors.NewSeveref(
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

func isConflictingPath(existingPaths []string, path string) bool {
	for _, existingPath := range existingPaths {
		if filepath.Clean(existingPath) == filepath.Clean(path) {
			return true
		}
	}
	return false
}
