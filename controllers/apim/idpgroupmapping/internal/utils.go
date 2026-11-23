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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/idpgroupmapping"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/service"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
)

// resolveGroups returns the values to store in the APIM group mapping for the spec groups.
// Expressions are kept as-is: APIM evaluates them at login and they resolve to group names.
// Any other value is a group name, resolved to its ID. Unresolved names are skipped.
func resolveGroups(
	ctx context.Context,
	env *service.Env,
	mapping *v1alpha1.IDPGroupMapping,
) ([]string, error) {
	resolved := []string{}
	for _, group := range mapping.Spec.Groups {
		if idpgroupmapping.IsGroupExpression(group) {
			resolved = append(resolved, group)
			continue
		}

		found, err := env.FindGroup(group)
		if err != nil {
			return nil, gerrors.NewControlPlaneError(err)
		}
		if found == nil || found.ID == "" {
			log.Error(ctx, nil, "unable to resolve group, skipping", "group", group, "mapping", mapping.Name)
			continue
		}
		resolved = append(resolved, found.ID)
	}
	return resolved, nil
}

func containsGroupMapping(groupMappings []model.GroupMapping, condition string, groups []string) bool {
	for _, gm := range groupMappings {
		if gm.Condition == condition && equalStringSlices(gm.Groups, groups) {
			return true
		}
	}
	return false
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	stringMap := make(map[string]bool)
	for _, item := range a {
		stringMap[item] = true
	}

	for _, item := range b {
		if !stringMap[item] {
			return false
		}
	}

	return true
}
