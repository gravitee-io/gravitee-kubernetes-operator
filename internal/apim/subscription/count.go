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

package subscription

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// RecountParents sets the subscription count of the API and of the application the subscription references to
// the number of Subscriptions referencing each of them and not being deleted. A parent already gone is skipped.
// The counts are informative (deletes are guarded by the Subscriptions themselves), and a retried reconcile
// sets the same numbers again.
func RecountParents(ctx context.Context, sub *v1alpha1.Subscription) error {
	ns := sub.GetNamespace()

	api, err := dynamic.ResolveAPI(ctx, APIRef(sub), ns)
	if client.IgnoreNotFound(err) != nil {
		return err
	}
	if err == nil {
		if err := recount(ctx, api, apiSubscriptionsField(api)); err != nil {
			return err
		}
	}

	app, err := dynamic.ResolveApplication(ctx, AppRef(sub), ns)
	if apierrors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return recount(ctx, app, search.AppSubsField)
}

func apiSubscriptionsField(api core.ApiDefinitionObject) search.IndexField {
	if _, ok := api.(*v1alpha1.ApiDefinition); ok {
		return search.ApiV2SubsField
	}
	return search.ApiV4SubsField
}

// recount updates the status inside a retry on conflict: the API and application controllers write the same
// status. Each attempt lists the Subscriptions again, so it never writes a count older than the object it read.
func recount(ctx context.Context, parent core.Object, field search.IndexField) error {
	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		count, err := countActiveSubscriptions(ctx, parent, field)
		if err != nil {
			return err
		}
		if err := k8s.GetClient().Get(ctx, client.ObjectKeyFromObject(parent), parent); err != nil {
			return client.IgnoreNotFound(err)
		}
		status, ok := parent.GetStatus().(core.SubscribableStatus)
		if !ok || status.GetSubscriptionCount() == count {
			return nil
		}
		status.SetSubscriptionCount(count)
		return k8s.GetClient().Status().Update(ctx, parent)
	})
}

func countActiveSubscriptions(ctx context.Context, parent core.Object, field search.IndexField) (uint, error) {
	subs, err := search.FindSubscriptionsReferencing(ctx, field, parent)
	if err != nil {
		return 0, err
	}
	var count uint
	for i := range subs.Items {
		if !subs.Items[i].IsBeingDeleted() {
			count++
		}
	}
	return count, nil
}
