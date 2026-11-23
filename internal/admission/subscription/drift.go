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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	subscriptionsdk "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/subscription"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

func mergeDriftValidation(
	ctx context.Context,
	oldSub *v1alpha1.Subscription,
	newSub *v1alpha1.Subscription,
	api core.ApiDefinitionObject,
	app core.ApplicationObject,
	errs *errors.AdmissionErrors,
) {
	errs.MergeWith(
		drift.ValidateDriftWithContext(ctx, oldSub, newSub,
			resolveContext(app),
			resolveRefs,
			remoteSubscriptionGetter(ctx, api, app),
			dtoMapper(api, app)),
	)
}

func remoteSubscriptionGetter(
	ctx context.Context,
	api core.ApiDefinitionObject,
	app core.ApplicationObject,
) drift.RemoteObjectGetter[*v1alpha1.Subscription] {
	return func(apimClient *apim.APIM, sub *v1alpha1.Subscription) (any, error) {
		client, err := apimClient.SDK()
		if err != nil {
			return nil, err
		}
		dto, err := subscriptionsdk.ToDTO(sub, api, app)
		if err != nil {
			return nil, err
		}
		return subscriptionsdk.GetRemote(ctx, client, dto)
	}
}

func resolveRefs(context.Context, *v1alpha1.Subscription) error {
	return nil
}

func resolveContext(app core.ContextAwareObject) drift.ContextResolver {
	return func(ctx context.Context) (*apim.APIM, error) {
		return apim.FromContextRef(ctx, app.ContextRef(), app.GetNamespace())
	}
}

func dtoMapper(
	api core.ApiDefinitionObject,
	app core.ApplicationObject,
) drift.DTOMapper[*v1alpha1.Subscription] {
	return func(sub *v1alpha1.Subscription) any {
		// The CRD schema makes the spec fit the SDK payload: mapping it cannot fail.
		dto, _ := subscriptionsdk.ToDriftDTO(sub, api, app)
		return dto
	}
}
