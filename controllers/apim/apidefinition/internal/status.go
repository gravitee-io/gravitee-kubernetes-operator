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
	"regexp"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

func UpdateStatusSuccess(ctx context.Context, api core.ConditionAwareObject) error {
	if api.IsBeingDeleted() {
		return nil
	}

	k8s.AddSuccessfulConditions(api)

	// Sometimes we may just set warnings instead of rejecting the API
	// This part of the code tries to set the condition for missing groups
	errors := api.GetStatus().GetErrors()
	if errors.Warning != nil {
		groupNotFoundErrorMessage := make([]string, 0)
		for i := 0; i < len(errors.Warning); i++ {
			w := errors.Warning[i]
			if strings.HasPrefix(w, "Group [") {
				re := regexp.MustCompile(`^Group \[.*] could not be found in environment \[.*]$`)
				if re.MatchString(w) {
					groupNotFoundErrorMessage = append(groupNotFoundErrorMessage, w)
				}
			}
		}
		if len(groupNotFoundErrorMessage) != 0 {
			k8s.SetCondition(
				api.(core.ConditionAware), //nolint:errcheck // api is ConditionAware
				k8s.
					NewResolvedRefsConditionBuilder(api.GetGeneration()).
					RejectGroupNotFound(strings.Join(groupNotFoundErrorMessage, ", ")).
					Build(),
			)
		}
	}

	// Deprecated
	api.GetStatus().SetProcessingStatus(core.ProcessingStatusCompleted)
	return k8s.GetClient().Status().Update(ctx, api)
}

func UpdateStatusFailure(ctx context.Context, api core.Object, err error) error {
	k8s.ErrorToCondition(api, err)

	// Deprecated
	api.GetStatus().SetProcessingStatus(core.ProcessingStatusFailed)
	return k8s.GetClient().Status().Update(ctx, api)
}
