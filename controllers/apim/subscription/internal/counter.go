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

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Subscription counters are informative: deletes are guarded by the Subscriptions themselves.
// The update is retried on conflict, as the application and API controllers write the same
// status: returning the error would requeue the reconcile and lose or repeat the change.
func updateSubscriptionCount(
	ctx context.Context,
	obj core.Object,
	change func(core.SubscribableStatus),
) error {
	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		if err := k8s.GetClient().Get(ctx, client.ObjectKeyFromObject(obj), obj); err != nil {
			return err
		}
		status, ok := obj.GetStatus().(core.SubscribableStatus)
		if !ok {
			return nil
		}
		change(status)
		return k8s.GetClient().Status().Update(ctx, obj)
	})
}
