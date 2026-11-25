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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ResolveRefs checks that the API, the application and the plan the subscription references exist.
// On delete, only the API is required: it carries the management context, so a missing API (NotFound) releases
// the subscription without calling APIM, while a missing application still lets the delete reach APIM.
func ResolveRefs(ctx context.Context, sub *v1alpha1.Subscription, ns string) error {
	apiRef := APIRef(sub)
	api, err := resolveAPI(ctx, apiRef, ns)
	if err != nil {
		return err
	}
	_, err = resolveApplication(ctx, AppRef(sub), ns)
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
			"unable to subscribe to API [%s] because plan [%s] cannot be found", apiRef, sub.GetPlan(),
		)
	}
	return nil
}

// CreateAPIMClient builds the SDK client of the management context the subscribed API references: the
// subscription has no context of its own.
func CreateAPIMClient(ctx context.Context, sub *v1alpha1.Subscription) (*apim.SDKClient, error) {
	api, err := resolveAPI(ctx, APIRef(sub), sub.GetNamespace())
	if err != nil {
		return nil, err
	}
	if !api.HasContext() {
		return nil, fmt.Errorf(
			"unable to subscribe to API [%s] because it does not reference a management context", api.GetRef(),
		)
	}
	client, err := apim.SDKFromContextRef(ctx, api.ContextRef(), api.GetNamespace())
	if err != nil {
		// The context may sit in another namespace than the subscription: its error stays in the operator log.
		log.FromContext(ctx).Error(err, "Unable to build the APIM client of the API", "api", api.GetRef())
		return nil, fmt.Errorf("unable to resolve management context [%s] of API [%s]", api.ContextRef(), api.GetRef())
	}
	return client, nil
}

// ToResolvedDTO resolves the API, with its IDs, and the application the subscription references, then maps them
// with ToDTO. On delete, a missing application leaves it out of the payload. A failed lookup is a reference
// error, which the reconcile retries.
func ToResolvedDTO(ctx context.Context, sub *v1alpha1.Subscription) (DTO, error) {
	api, app, err := resolveParents(ctx, sub)
	if err != nil {
		return DTO{}, gerrors.NewResolveRefError(err)
	}
	return ToDTO(sub, api, app)
}

func resolveParents(
	ctx context.Context, sub *v1alpha1.Subscription,
) (core.ApiDefinitionObject, core.ApplicationObject, error) {
	ns := sub.GetNamespace()
	api, err := resolveAPI(ctx, APIRef(sub), ns)
	if err != nil {
		return nil, nil, err
	}
	mctx, err := dynamic.ResolveContext(ctx, api.ContextRef(), api.GetNamespace())
	if err != nil {
		return nil, nil, err
	}
	api.PopulateIDs(mctx, k8s.IsAutomationAPIManaged(api))

	app, err := resolveApplication(ctx, AppRef(sub), ns)
	if sub.IsBeingDeleted() && apierrors.IsNotFound(err) {
		return api, nil, nil
	}
	return api, app, err
}

// APIRef returns a copy of the API reference of sub, in the namespace of sub when it names none. The dynamic
// resolvers write that namespace into the reference they get: resolving a copy leaves the Subscription spec as
// written, which admission compares with the previous one.
func APIRef(sub *v1alpha1.Subscription) *refs.NamespacedName {
	return inNamespace(sub.Spec.API, sub.GetNamespace())
}

// AppRef returns a copy of the application reference of sub, as APIRef does.
func AppRef(sub *v1alpha1.Subscription) *refs.NamespacedName {
	return inNamespace(sub.Spec.App, sub.GetNamespace())
}

func inNamespace(ref refs.NamespacedName, ns string) *refs.NamespacedName {
	if ref.Namespace == "" {
		ref.Namespace = ns
	}
	return &ref
}

func resolveAPI(ctx context.Context, ref core.ObjectRef, ns string) (core.ApiDefinitionObject, error) {
	api, err := dynamic.ResolveAPI(ctx, ref, ns)
	if err != nil {
		return nil, unresolvedRef{message: fmt.Sprintf("unable to resolve API [%s]", ref), cause: err}
	}
	return api, nil
}

func resolveApplication(ctx context.Context, ref core.ObjectRef, ns string) (core.ApplicationObject, error) {
	app, err := dynamic.ResolveApplication(ctx, ref, ns)
	if err != nil {
		return nil, unresolvedRef{message: fmt.Sprintf("unable to resolve application [%s]", ref), cause: err}
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
