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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/notification"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

func ResolveConsoleNotificationRefs(ctx context.Context, api core.ApiDefinitionObject) error {
	notificationRefs := api.GetNotificationRefs()

	if len(notificationRefs) == 0 {
		api.SetConsoleNotification(nil)
		return nil
	}

	for _, ref := range notificationRefs {
		notif := new(v1alpha1.Notification)
		nsn := getNamespacedName(ref, api.GetNamespace())
		err := k8s.GetClient().Get(ctx, nsn, notif)
		if err != nil {
			return err
		}
		// take the first one ignore the rest
		// validation webhook takes care of checking that beforehand
		if notif.Spec.EventType == notification.EventTypeAPI &&
			notif.Spec.Target == notification.TargetConsole {
			api.SetConsoleNotification(base.ToAPIConsoleNotificationSettings(notif.Spec.Console))
			break
		}
	}

	return nil
}
