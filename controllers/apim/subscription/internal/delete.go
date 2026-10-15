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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func Delete(
	ctx context.Context,
	subscription *v1alpha1.Subscription,
) error {
	if !util.ContainsFinalizer(subscription, core.SubscriptionFinalizer) {
		return nil
	}

	ns := subscription.Namespace
	spec := subscription.Spec

	// Both references are resolved before anything is written, so a resolution error
	// cannot leave a counter decremented and decrement it again on the next attempt.
	api, err := dynamic.ResolveAPI(ctx, &spec.API, ns)
	apiGone := apierrors.IsNotFound(err)
	if err != nil && !apiGone {
		return err
	}

	// A deleted application closed its subscriptions in APIM: only its counter is left behind.
	app, err := dynamic.ResolveApplication(ctx, &spec.App, ns)
	appGone := apierrors.IsNotFound(err)
	if err != nil && !appGone {
		return err
	}

	if apiGone {
		// The API carries the management context needed to reach APIM, and deleting it
		// in APIM removed its subscriptions: there is nothing left to delete. Only an API
		// whose finalizer was removed by hand leaves the API and its subscriptions in APIM.
		log.Info(ctx,
			"API not found, releasing the subscription. If the API finalizer was removed by hand, "+
				"the API and its subscriptions remain in APIM",
			"api", spec.API.String(), "severity", "warning")
	} else {
		apim, err := apim.FromContextRef(ctx, api.ContextRef(), api.GetNamespace())
		if err != nil {
			return err
		}

		api.PopulateIDs(apim.Context, k8s.IsAutomationAPIManaged(api))

		if err := errors.IgnoreNotFound(apim.Subscription.Delete(api, subscription)); err != nil {
			return err
		}

		if err := updateSubscriptionCount(ctx, api, core.SubscribableStatus.RemoveSubscription); err != nil {
			return err
		}
	}

	if appGone {
		return nil
	}
	return updateSubscriptionCount(ctx, app, core.SubscribableStatus.RemoveSubscription)
}
