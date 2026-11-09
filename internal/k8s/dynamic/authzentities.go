// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dynamic

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

// ResolveAuthzEntityRefs replaces the entity's targets with what they send on the wire, in
// declaration order: each scope as written, each pdpRef as the routing scope of the AuthzPdp it
// names. ToAuthzEntityDTO then only maps. The reconcile, the admission dry run and drift detection
// all resolve with it. Its errors are those of ResolveAuthzTargets, and leave the targets as they
// were.
func ResolveAuthzEntityRefs(ctx context.Context, entity *v1alpha1.AuthzEntity) error {
	scopes, err := ResolveAuthzTargets(ctx, entity, entity.Spec.Targets)
	if err != nil {
		return err
	}
	var targets []authzpdp.Target
	for _, scope := range scopes {
		targets = append(targets, authzpdp.Target{Scope: scope})
	}
	entity.Spec.Targets = targets
	return nil
}
