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

// ResolveAuthzSchemaRefs replaces the schema's target with what it sends on the wire: a scope as
// written, a pdpRef as the routing scope of the AuthzPdp it names, no target as none.
// ToAuthzSchemaDTO then only maps. The reconcile, the admission dry run and drift detection all
// resolve with it. Its errors are those of ResolveAuthzTarget, and leave the target as it was.
func ResolveAuthzSchemaRefs(ctx context.Context, schema *v1alpha1.AuthzSchema) error {
	scope, err := ResolveAuthzTarget(ctx, schema, schema.Spec.Target)
	if err != nil {
		return err
	}
	if scope == "" {
		schema.Spec.Target = nil
		return nil
	}
	schema.Spec.Target = &authzpdp.Target{Scope: scope}
	return nil
}
