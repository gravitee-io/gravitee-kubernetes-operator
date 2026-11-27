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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	subscriptionsdk "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/subscription"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

func CreateOrUpdate(ctx context.Context, subscription *v1alpha1.Subscription) error {
	ns := subscription.Namespace
	spec := subscription.Spec

	creating := subscription.Status.ID == ""
	app, err := dynamic.ResolveApplication(ctx, &spec.App, ns)
	if err != nil {
		return err
	}

	api, err := dynamic.ResolveAPI(ctx, &spec.API, ns)
	if err != nil {
		return err
	}

	apimClient, err := subscriptionsdk.CreateAPIMClient(ctx, api)
	if err != nil {
		return err
	}

	if plan := api.GetPlan(spec.Plan); plan == nil {
		return fmt.Errorf("plan %s not found", subscription.Spec.Plan)
	}

	api.PopulateIDs(apimClient.Context, k8s.IsAutomationAPIManaged(subscription))

	dto, err := subscriptionsdk.ToDTO(subscription, api, app)
	if err != nil {
		return err
	}
	resp, err := subscriptionsdk.Upsert(ctx, apimClient, dto)
	if err != nil {
		return gerrors.NewControlPlaneError(err)
	}
	if err := subscriptionsdk.UpdateStatus(subscription, resp); err != nil {
		return err
	}

	// Counted once, on creation
	if creating {
		if err := updateSubscriptionCount(ctx, app, core.SubscribableStatus.AddSubscription); err != nil {
			return err
		}
		return updateSubscriptionCount(ctx, api, core.SubscribableStatus.AddSubscription)
	}

	return nil
}
