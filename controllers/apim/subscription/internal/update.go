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
	"fmt"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

func CreateOrUpdate(ctx context.Context, subscription *v1alpha1.Subscription) error {
	ns := subscription.Namespace
	spec := subscription.Spec

	app, err := dynamic.ResolveApplication(ctx, &spec.App, ns)
	if err != nil {
		return err
	}

	api, err := dynamic.ResolveAPI(ctx, &spec.API, ns)
	if err != nil {
		return err
	}

	apimClient, err := apim.FromContextRef(ctx, api.ContextRef(), ns)
	if err != nil {
		return err
	}

	if plan := api.GetPlan(spec.Plan); plan == nil {
		return fmt.Errorf("plan %s not found", subscription.Spec.Plan)
	}

	api.PopulateIDs(apimClient.Context)

	var sub model.Subscription

	var legacyApiID, legacyAppID, legacySubscriptionID bool
	sub.ID = refs.NewNamespacedNameFromObject(subscription).HRID()
	sub.ApiID = api.GetHRID()
	sub.AppID = app.GetHRID()
	sub.PlanID = spec.Plan

	// handle legacy IDs
	if subscription.Status.ID != "" {
		legacySubscriptionID = true
		sub.ID = string(subscription.UID)
	}

	if api.GetID() != "" {
		legacyApiID = true
		sub.ApiID = api.GetID()
	}

	if app.GetID() != "" {
		legacyAppID = true
		sub.AppID = app.GetID()
	}

	if spec.EndingAt != nil {
		sub.EndingAt = *spec.EndingAt
	}

	status, err := apimClient.Subscription.Import(sub, legacySubscriptionID, legacyApiID, legacyAppID)
	if err != nil {
		return gerrors.NewControlPlaneError(err)
	}

	startedAt, err := time.Parse(time.RFC3339, status.StartingAt)
	if err != nil {
		return err
	}

	subscription.Status.StartedAt = startedAt.Format(time.RFC3339)
	subscription.Status.EndingAt = status.EndingAt
	subscription.Status.ID = status.ID

	appStatus, _ := app.GetStatus().(core.SubscribableStatus)
	apiStatus, _ := api.GetStatus().(core.SubscribableStatus)

	appStatus.AddSubscription()
	apiStatus.AddSubscription()

	if err := k8s.GetClient().Status().Update(ctx, app); err != nil {
		return err
	}

	if err := k8s.GetClient().Status().Update(ctx, api); err != nil {
		return err
	}

	return nil
}
