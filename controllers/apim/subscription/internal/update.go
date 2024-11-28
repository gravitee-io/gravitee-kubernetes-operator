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
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
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

	apim, err := apim.FromContextRef(ctx, api.ContextRef(), ns)
	if err != nil {
		return err
	}

	if plan := api.GetPlan(spec.Plan); plan == nil {
		return fmt.Errorf("plan %s not found", subscription.Spec.Plan)
	}

	api.PopulateIDs(apim.Context)

	appID := app.GetID()
	apiID := api.GetID()
	planID := api.GetPlan(subscription.Spec.Plan).GetID()
	subscriptionID := string(subscription.UID)

	sub := &model.Subscription{
		ID:     subscriptionID,
		AppID:  appID,
		ApiID:  apiID,
		PlanID: planID,
	}

	if spec.EndingAt != nil {
		sub.EndingAt = *spec.EndingAt
	}

	status, err := apim.Subscription.Import(sub)
	if err != nil {
		return err
	}

	startedAt, err := time.Parse(time.RFC3339, status.StartingAt)
	if err != nil {
		return err
	}

	subscription.Status.StartedAt = startedAt.Format("2024-12-25T09:12:28Z")
	subscription.Status.EndingAt = status.EndingAt
	subscription.Status.ID = status.ID

	return nil
}
