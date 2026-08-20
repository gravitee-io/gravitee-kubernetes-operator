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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/subscription"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hrid"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

func mergeDriftValidation(
	ctx context.Context,
	oldSub *v1alpha1.Subscription,
	newSub *v1alpha1.Subscription,
	api core.ApiDefinitionObject,
	app core.ApplicationObject,
	plan core.PlanModel,
	errs *errors.AdmissionErrors,
) {
	errs.MergeWith(
		drift.ValidateDriftWithContext(ctx, oldSub, newSub,
			resolveContext(app),
			resolveRefs,
			remoteSubscriptionGetter(api),
			dtoMapper(api, app, plan)),
	)
}

func remoteSubscriptionGetter(api core.ApiDefinitionObject) drift.RemoteObjectGetter[*v1alpha1.Subscription] {
	return func(apimClient *apim.APIM, sub *v1alpha1.Subscription, admissionErrors *errors.AdmissionErrors) any {
		if k8s.IsAutomationAPIManaged(sub) {
			apiNSName := refs.NewNamespacedNameFromObject(api)
			subNSName := refs.NewNamespacedNameFromObject(sub)
			remoteSub, err := apimClient.Subscription.GetByHRID(apiNSName.HRID(), subNSName.HRID())
			if err != nil {
				admissionErrors.AddSeveref("cannot fetch Subscription during drift detection from Api HRID %s and Subscription HRID %s: %s",
					apiNSName.HRID(),
					subNSName.HRID(),
					err.Error())
				return nil
			}
			return *remoteSub
		}
		remoteSub, err := apimClient.Subscription.GetByID(api.GetID(), sub.Status.ID)
		if err != nil {
			admissionErrors.AddSeveref("cannot fetch Subscription during drift detection from Api ID %s and Subscription ID %s: %s",
				api.GetID(),
				sub.Status.ID,
				err.Error())
			return nil
		}

		return *remoteSub
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

func dtoMapper(api core.ApiDefinitionObject, app core.ApplicationObject, plan core.PlanModel) drift.DTOMapper[*v1alpha1.Subscription] {
	return func(sub *v1alpha1.Subscription) any {
		if k8s.IsAutomationAPIManaged(sub) {
			apiNSName := refs.NewNamespacedNameFromObject(api)
			appNSName := refs.NewNamespacedNameFromObject(app)
			subNSName := refs.NewNamespacedNameFromObject(sub)
			return model.SubscriptionDTO{
				ID:                    subNSName.HRID(),
				ApiID:                 apiNSName.HRID(),
				AppID:                 appNSName.HRID(),
				PlanID:                hrid.NameToValidHRID(sub.Spec.Plan),
				StartingAt:            sub.Status.StartedAt,
				EndingAt:              utils.SafeDereference(sub.Spec.EndingAt),
				Metadata:              sub.Spec.Metadata,
				ApiKeys:               mapApiKey(sub.Spec.ApiKeys),
				ConsumerConfiguration: sub.Spec.ConsumerConfiguration.DeepCopy(),
			}
		}
		// Legacy
		return model.SubscriptionDTO{
			ID:                    sub.Status.ID,
			ApiID:                 api.GetID(),
			AppID:                 app.GetID(),
			PlanID:                plan.GetID(),
			StartingAt:            sub.Status.StartedAt,
			EndingAt:              utils.SafeDereference(sub.Spec.EndingAt),
			Metadata:              sub.Spec.Metadata,
			ApiKeys:               mapApiKey(sub.Spec.ApiKeys),
			ConsumerConfiguration: sub.Spec.ConsumerConfiguration.DeepCopy(),
		}
	}
}

func mapApiKey(keys []subscription.ApiKeySpec) []model.ApiKeySpec {
	apiKeys := make([]model.ApiKeySpec, len(keys))
	for i, key := range keys {
		apiKeys[i] = model.ApiKeySpec{
			Key:      key.Key,
			ExpireAt: key.ExpireAt,
		}
	}
	return apiKeys
}
