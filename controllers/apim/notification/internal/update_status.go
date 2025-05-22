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
	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const resolveRefCondition = "ResolveRef"
const acceptedCondition = "Accepted"

func SetGroupRefsConditions(
	ctx context.Context,
	client client.Client,
	err error,
	notification *v1alpha1.Notification) error {
	if err != nil {
		changed := meta.SetStatusCondition(notification.Status.Conditions, v1.Condition{
			Type:    resolveRefCondition,
			Status:  v1.ConditionFalse,
			Reason:  "GroupsResolveRefFailed",
			Message: err.Error(),
		})
		if err := updateAndRefresh(ctx, client, changed, notification); err != nil {
			return err
		}
		changed = meta.SetStatusCondition(notification.Status.Conditions, v1.Condition{
			Type:    acceptedCondition,
			Status:  v1.ConditionFalse,
			Reason:  "GroupsRefsResolveFailed",
			Message: err.Error(),
		})
		if err := updateAndRefresh(ctx, client, changed, notification); err != nil {
			return err
		}
		return err
	}

	changed := meta.SetStatusCondition(notification.Status.Conditions, v1.Condition{
		Type:    resolveRefCondition,
		Status:  v1.ConditionTrue,
		Reason:  "GroupsRefsResolved",
		Message: "Successfully resolved groups references",
	})
	if err := updateAndRefresh(ctx, client, changed, notification); err != nil {
		return err
	}
	return nil
}

func updateAndRefresh(
	ctx context.Context,
	client client.Client,
	changed bool,
	notification *v1alpha1.Notification) error {
	if changed {
		err := client.Status().Update(ctx, notification)
		if err != nil {
			return err
		}
		// refresh the resource after update, as more status may be added
		if err := client.Get(ctx,
			types.NamespacedName{
				Name:      notification.Name,
				Namespace: notification.Namespace},
			notification); err != nil {
			return err
		}
	}
	return nil
}

func SetAcceptedCondition(ctx context.Context, client client.Client, notification *v1alpha1.Notification) error {
	changed := meta.SetStatusCondition(notification.Status.Conditions, v1.Condition{
		Type:    acceptedCondition,
		Status:  v1.ConditionTrue,
		Reason:  "Reconciled",
		Message: "Successfully reconciled",
	})

	if changed {
		if err := client.Status().Update(ctx, notification); err != nil {
			return err
		}
	}
	return nil
}
