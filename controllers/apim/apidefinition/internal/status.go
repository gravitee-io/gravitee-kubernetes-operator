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

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const ConditionConflicted = "ReconcileFailed"

func UpdateStatusSuccess(ctx context.Context, api core.Object) error {
	if !api.GetDeletionTimestamp().IsZero() {
		return nil
	}
	acb := k8s.NewAcceptedConditionBuilder(api.GetGeneration()).Accept("Successfully reconciled").Build()
	api.GetStatus().SetConditions([]metav1.Condition{*acb})

	// Deprecated
	api.GetStatus().SetProcessingStatus(core.ProcessingStatusCompleted)
	return k8s.GetClient().Status().Update(ctx, api)
}

func UpdateStatusFailure(ctx context.Context, api core.Object, err error) error {
	acb := k8s.NewAcceptedConditionBuilder(api.GetGeneration()).Message(err.Error()).Reason(ConditionConflicted).Build()
	api.GetStatus().SetConditions([]metav1.Condition{*acb})

	// Deprecated
	api.GetStatus().SetProcessingStatus(core.ProcessingStatusFailed)
	return k8s.GetClient().Status().Update(ctx, api)
}
