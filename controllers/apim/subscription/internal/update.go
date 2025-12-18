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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"

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

	creating := subscription.Status.ID == ""
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

	api.PopulateIDs(apimClient.Context, k8s.IsAutomationAPIManaged(subscription))

	sub := &model.Subscription{}

	setHridWithSubscriptionUUID := setSubscriptionID(subscription, sub)
	setHridWithApiUUID := setApiIDAndPlan(api, sub, spec)
	setHridWithAppUUID := setApplicationID(app, sub)
	sub.Metadata = spec.Metadata
	if spec.EndingAt != nil {
		sub.EndingAt = *spec.EndingAt
	}

	status, err := apimClient.Subscription.Import(*sub,
		subscription,
		setHridWithSubscriptionUUID,
		setHridWithApiUUID,
		setHridWithAppUUID)
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

	// Avoid increasing counter on updates
	// Counters are used to prevent deletions
	// When subscription is used by an application or an API
	if creating {
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
	}

	return nil
}

func setApiIDAndPlan(api core.ApiDefinitionObject, sub *model.Subscription, spec v1alpha1.SubscriptionSpec) bool {
	if api.GetID() == "" || k8s.IsAutomationAPIManaged(api) {
		sub.ApiID = refs.NewNamespacedNameFromObject(api).HRID()
		sub.PlanID = spec.Plan
		return false
	}
	sub.ApiID = api.GetID()
	sub.PlanID = api.GetPlan(spec.Plan).GetID()
	return true
}

func setSubscriptionID(subscription *v1alpha1.Subscription, sub *model.Subscription) bool {
	if subscription.Status.ID == "" || k8s.IsAutomationAPIManaged(subscription) {
		sub.ID = refs.NewNamespacedNameFromObject(subscription).HRID()
		return false
	}
	sub.ID = string(subscription.UID)
	return true
}

func setApplicationID(app core.ApplicationObject, sub *model.Subscription) bool {
	if app.GetID() == "" || k8s.IsAutomationAPIManaged(app) {
		sub.AppID = refs.NewNamespacedNameFromObject(app).HRID()
		return false
	}
	sub.AppID = app.GetID()
	return true
}
