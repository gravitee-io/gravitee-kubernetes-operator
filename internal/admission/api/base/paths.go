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
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Conflict struct {
	ID                string
	Path              string
	CreationTimestamp metaV1.Time
	Tags              []string
}

func (p Conflict) IsZero() bool {
	return p.ID == "" && p.Path == ""
}

func ValidateNoConflictingPath(ctx context.Context, api core.ApiDefinitionObject) *errors.AdmissionError {
	if k8s.HasHTTPRouteOwner(api.GetOwnerReferences()) {
		return nil
	}
	if conflicting, err := FindConflictingPath(ctx, api); err != nil {
		return errors.NewSevere(err.Error())
	} else if !conflicting.IsZero() {
		return errors.NewSeveref(
			"invalid API context path [%s]. API [%s] is already defined with the same path",
			conflicting.Path, conflicting.ID,
		)
	}
	return nil
}

func FindConflictingPath(ctx context.Context, api core.ApiDefinitionObject) (Conflict, error) {
	apiPaths := api.GetContextPaths()
	existingPaths, err := getExistingPaths(ctx, api)
	if err != nil {
		return Conflict{}, errors.NewSevere(err.Error())
	}
	for _, apiPath := range apiPaths {
		if _, err := url.Parse(apiPath); err != nil {
			return Conflict{}, errors.NewSeveref(
				"path [%s] is invalid",
				apiPath,
			)
		}

		if conflictingPath := findConflictingPathAPI(existingPaths, apiPath); !conflictingPath.IsZero() {
			return conflictingPath, nil
		}
	}
	return Conflict{}, nil
}

func getExistingPaths(ctx context.Context, api core.ApiDefinitionObject) ([]Conflict, error) {
	existingPaths := make([]Conflict, 0)
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
		paths := make([]Conflict, 0)
		apiPaths := api.GetContextPaths()
		for _, path := range apiPaths {
			paths = append(paths, Conflict{
				api.GetNamespace() + "/" + api.GetName(),
				path,
				api.GetCreationTimestamp(),
				api.GetTags(),
			})
		}
		existingPaths = append(existingPaths, paths...)
	}
	return existingPaths, nil
}

func findConflictingPathAPI(existingPaths []Conflict, path string) Conflict {
	for _, existingPath := range existingPaths {
		if filepath.Clean(existingPath.Path) == filepath.Clean(path) {
			return existingPath
		}
	}
	return Conflict{}
}
