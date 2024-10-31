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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
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

	api, err := dynamic.ResolveAPI(ctx, &subscription.Spec.API, ns)
	if err != nil {
		return err
	}

	apim, err := apim.FromContextRef(ctx, api.ContextRef(), ns)
	if err != nil {
		return err
	}

	api.PopulateIDs(apim.Context)

	err = apim.Subscription.Delete(&model.Subscription{
		ApiID: api.GetID(),
		ID:    subscription.Status.ID,
	})

	if err != nil {
		return err
	}

	util.RemoveFinalizer(subscription, core.SubscriptionFinalizer)

	return nil
}
