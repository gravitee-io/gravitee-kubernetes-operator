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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/notification"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

func validateNotifications(ctx context.Context, api core.ApiDefinitionObject) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	if len(api.GetNotificationRefs()) > 0 {
		var consoleNotificationRef int
		for _, notificationRef := range api.GetNotificationRefs() {
			notif, err := dynamic.ResolveNotification(ctx, notificationRef, api.GetNamespace())
			if err != nil {
				errs.AddSeveref(
					"api references notification [%s] that does not exist in the cluster", notificationRef,
				)
				break
			}
			if notif.Spec.EventType != notification.EventTypeAPI {
				errs.AddSeveref(
					"api references notification [%s] is not configured for apis but for [%s]", notificationRef, notif.Spec.EventType,
				)
			}
			if notif.Spec.Target == notification.TargetConsole {
				consoleNotificationRef++
				checkConsoleNotification(ctx, notif, notificationRef, errs)
			}
			if consoleNotificationRef > 1 {
				errs.AddSeveref(
					"api references notification [%s] but there is already another console notification referenced", notificationRef)
			}
		}
	}

	return errs
}

func checkConsoleNotification(
	ctx context.Context,
	notif *v1alpha1.Notification,
	ref core.ObjectRef,
	errs *errors.AdmissionErrors) {
	if len(notif.Spec.Console.APIEvents) == 0 {
		errs.AddWarningf("api references notification [%s] configured withouut any API events", ref)
	}
	if len(notif.Spec.Console.GroupRefs) > 0 {
		for _, groupRef := range notif.Spec.Console.GroupRefs {
			r := groupRef
			if _, err := dynamic.ResolveGroup(ctx, &r, notif.GetNamespace()); err != nil {
				errs.AddWarningf(
					"api references notification [%s] configured with group [%s] that does not exist in the cluster",
					notif.GetRef(),
					r.String())
			}
		}
	}
}
