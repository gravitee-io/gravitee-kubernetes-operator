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
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// ResolveRefs checks that the API, the application and the plan the subscription references exist.
// On delete, only the API is required: it carries the management context, so a missing API (NotFound) releases
// the subscription without calling APIM, while a missing application still lets the delete reach APIM.
func ResolveRefs(ctx context.Context, sub *v1alpha1.Subscription, ns string) error {
	api, err := resolveAPI(ctx, sub, ns)
	if err != nil {
		return err
	}
	_, err = resolveApplication(ctx, sub, ns)
	if sub.IsBeingDeleted() {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err != nil {
		return err
	}
	if api.GetPlan(sub.GetPlan()) == nil {
		return fmt.Errorf(
			"unable to subscribe to API [%s] because plan [%s] cannot be found", sub.GetApiRef(), sub.GetPlan(),
		)
	}
	return nil
}

// CreateAPIMClient builds the SDK client of the management context the subscribed API references: the
// subscription has no context of its own.
func CreateAPIMClient(ctx context.Context, sub *v1alpha1.Subscription) (*apim.SDKClient, error) {
	api, err := resolveAPI(ctx, sub, sub.GetNamespace())
	if err != nil {
		return nil, err
	}
	if !api.HasContext() {
		return nil, fmt.Errorf(
			"unable to subscribe to API [%s] because it does not reference a management context", api.GetRef(),
		)
	}
	return apim.SDKFromContextRef(ctx, api.ContextRef(), api.GetNamespace())
}

// ToResolvedDTO resolves the API, with its IDs, and the application the subscription references, then maps them
// with ToDTO. On delete, a missing application leaves it out of the payload.
func ToResolvedDTO(ctx context.Context, sub *v1alpha1.Subscription) (DTO, error) {
	ns := sub.GetNamespace()
	api, err := resolveAPI(ctx, sub, ns)
	if err != nil {
		return DTO{}, err
	}
	mctx, err := dynamic.ResolveContext(ctx, api.ContextRef(), api.GetNamespace())
	if err != nil {
		return DTO{}, err
	}
	api.PopulateIDs(mctx, k8s.IsAutomationAPIManaged(api))

	app, err := resolveApplication(ctx, sub, ns)
	if sub.IsBeingDeleted() && apierrors.IsNotFound(err) {
		app, err = nil, nil
	}
	if err != nil {
		return DTO{}, err
	}
	return ToDTO(sub, api, app)
}

func resolveAPI(ctx context.Context, sub *v1alpha1.Subscription, ns string) (core.ApiDefinitionObject, error) {
	api, err := dynamic.ResolveAPI(ctx, sub.GetApiRef(), ns)
	if err != nil {
		return nil, unresolvedRef{message: fmt.Sprintf("unable to resolve API [%s]", sub.GetApiRef()), cause: err}
	}
	return api, nil
}

func resolveApplication(ctx context.Context, sub *v1alpha1.Subscription, ns string) (core.ApplicationObject, error) {
	app, err := dynamic.ResolveApplication(ctx, sub.GetAppRef(), ns)
	if err != nil {
		return nil, unresolvedRef{
			message: fmt.Sprintf("unable to resolve application [%s]", sub.GetAppRef()),
			cause:   err,
		}
	}
	return app, nil
}

// unresolvedRef reads as the admission message of a reference that cannot be resolved, and keeps its cause so
// apierrors.IsNotFound still tells a missing reference apart.
type unresolvedRef struct {
	message string
	cause   error
}

func (e unresolvedRef) Error() string {
	return e.message
}

func (e unresolvedRef) Unwrap() error {
	return e.cause
}
