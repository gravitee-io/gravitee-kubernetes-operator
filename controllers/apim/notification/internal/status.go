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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

func UpdateStatusSuccess(ctx context.Context, notification *v1alpha1.Notification) error {
	if notification.IsBeingDeleted() {
		return nil
	}

	k8s.AddSuccessfulConditions(notification)

	return k8s.GetClient().Status().Update(ctx, notification)
}

func UpdateStatusFailure(ctx context.Context, notification *v1alpha1.Notification, err error) error {
	k8s.ErrorToCondition(notification, err)

	return k8s.GetClient().Status().Update(ctx, notification)
}
