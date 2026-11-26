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

package docs

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/navigation"
)

// ResolveRefsForAdmission checks the page's parent references and reads the parent, with the messages the
// webhook always reported.
func ResolveRefsForAdmission(ctx context.Context, obj *v1alpha1.Documentation, _ string) error {
	_, err := parentRefs(obj).Resolve(ctx, false)
	return err
}

// ResolveRefs resolves the parent and its management context for the reconcile. A parent or a context
// deleted first is a NotFound; a parent without a context, while deleting, means nothing was synced.
// An upsert waits for a portal GKO has not synced yet: APIM refuses a page of a portal it does not hold
// while the environment holds another, its default portal included.
func ResolveRefs(ctx context.Context, obj *v1alpha1.Documentation, _ string) error {
	parent, err := parentRefs(obj).Resolve(ctx, obj.IsBeingDeleted())
	if err != nil || !parent.HasContext() {
		return err
	}
	if err = navigation.ResolveContext(ctx, parent); err != nil || obj.IsBeingDeleted() {
		return err
	}
	return navigation.CheckPortalSynced(parent)
}

func parentRefs(obj *v1alpha1.Documentation) navigation.ParentRefs {
	return navigation.ParentRefs{
		Kind:      "documentation",
		Name:      obj.GetName(),
		Namespace: obj.GetNamespace(),
		Portal:    obj.Spec.Portal,
		API:       obj.Spec.API,
	}
}
