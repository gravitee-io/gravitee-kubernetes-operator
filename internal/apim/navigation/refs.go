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

package navigation

import (
	"context"
	"errors"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

// ParentRefs are the parent references of a link or documentation page, with what the messages name.
type ParentRefs struct {
	// Kind words the resource in the messages: "portal link", "documentation".
	Kind      string
	Name      string
	Namespace string
	Portal    *refs.NamespacedName
	API       *refs.NamespacedName
}

// Resolve checks that exactly one parent is set and that an API parent is a v4 API, then reads the
// parent and checks it has a management context, with the messages the webhooks always reported. A
// parent the cluster does not hold stays a NotFound. A parent without a context is returned without
// error while deleting: nothing was ever synced.
//
// While deleting, refs that fail the shape check are a NotFound too: they name no parent GKO could have
// synced the resource under, so the lifecycle releases the finalizer instead of retrying forever.
func (r ParentRefs) Resolve(ctx context.Context, deleting bool) (core.ContextAwareObject, error) {
	if err := r.checkShape(); err != nil {
		if deleting {
			return nil, gerrors.NewRefNotFoundError(
				fmt.Sprintf("%s [%s]: %v, nothing to delete in APIM", r.Kind, r.Name, err),
				apierrors.NewNotFound(schema.GroupResource{Group: core.CRDGroup}, r.Name),
			)
		}
		return nil, err
	}
	if r.API != nil {
		return r.resolveAPI(ctx, deleting)
	}
	return r.resolvePortal(ctx, deleting)
}

func (r ParentRefs) checkShape() error {
	switch {
	case r.Portal != nil && r.API != nil:
		return errors.New("exactly one of portalRef / apiRef must be set, but both were provided")
	case r.Portal == nil && r.API == nil:
		return errors.New("exactly one of portalRef / apiRef must be set, but neither was provided")
	}
	if r.API != nil && r.API.Kind != "" && dynamic.ResourceFromKind(r.API.Kind) != core.CRDApiV4DefinitionResource {
		return fmt.Errorf("apiRef [%s] must be of kind ApiV4Definition (next-gen portal only supports those)", r.API.Name)
	}
	return nil
}

func (r ParentRefs) resolvePortal(ctx context.Context, deleting bool) (core.ContextAwareObject, error) {
	prtl, err := dynamic.ResolvePortal(ctx, r.Portal, r.Namespace)
	if err != nil {
		return nil, gerrors.NewRefNotFoundError(
			fmt.Sprintf("%s [%s] references portal [%v] that can't be resolved", r.Kind, r.Name, r.Portal), err,
		)
	}
	if !prtl.HasContext() && !deleting {
		return nil, fmt.Errorf("referenced portal [%v] has no management context (spec.contextRef)", r.Portal)
	}
	return prtl, nil
}

func (r ParentRefs) resolveAPI(ctx context.Context, deleting bool) (core.ContextAwareObject, error) {
	api, err := dynamic.ResolveAPI(ctx, r.API, r.Namespace)
	if err != nil {
		return nil, gerrors.NewRefNotFoundError(
			fmt.Sprintf("%s [%s] references API [%v] that can't be resolved", r.Kind, r.Name, r.API), err,
		)
	}
	if !api.HasContext() && !deleting {
		return nil, fmt.Errorf("referenced API [%s] has no management context (spec.contextRef)", api.GetName())
	}
	return api, nil
}

// ResolveContext reads the management context obj references, so that a context deleted first is a
// NotFound: on delete, the lifecycle then releases the finalizer without calling APIM.
func ResolveContext(ctx context.Context, obj core.ContextAwareObject) error {
	if _, err := dynamic.ResolveContext(ctx, obj.ContextRef(), obj.GetNamespace()); err != nil {
		return gerrors.NewRefNotFoundError(
			fmt.Sprintf(
				"resource [%s] references management context [%v] that doesn't exist in the cluster",
				obj.GetName(), obj.ContextRef(),
			), err,
		)
	}
	return nil
}

// ErrNoClient is returned by a call that needs APIM when the resource has no management context, and so
// no client.
var ErrNoClient = errors.New("no management context: APIM cannot be called")
