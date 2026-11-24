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

package authzschema

import (
	"context"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hrid"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

// preCheck validates the HRID derived from the namespace and name.
func preCheck(_ context.Context, schema *v1alpha1.AuthzSchema) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	if err := hrid.Validate(refs.NewNamespacedNameFromObject(schema).HRID()); err != nil {
		errs.AddSevere(err.Error())
	}
	return errs
}

// resolveRefs checks the pdpRef kind, resolves the management context and, when its AuthzPdp is synced, the
// target onto its routing scope. A missing, deleting or unsynced AuthzPdp is a normal state when a directory is
// applied at once: the target is left as written, and deferDryRun and deferDrift skip the platform checks with a
// warning. The controller waits for the AuthzPdp before applying. Any other error, an AuthzPdp of another
// context for instance, refuses the schema.
func resolveRefs(ctx context.Context, schema *v1alpha1.AuthzSchema, ns string) error {
	if err := dynamic.ValidateAuthzTargetKind(schema, schema.Spec.Target); err != nil {
		return err
	}
	if err := apim.ResolveContextRef(ctx, schema, ns); err != nil {
		return err
	}
	if err := dynamic.ResolveAuthzSchemaRefs(ctx, schema); err != nil && !dynamic.IsAuthzPdpPending(err) {
		return err
	}
	return nil
}

// syncedDuringAdmission is why deferDryRun skips the dry run of a target whose AuthzPdp synced after resolveRefs
// read it.
const syncedDuringAdmission = "platform checks deferred (an AuthzPdp synced during admission)"

// deferDryRun skips the dry run while resolveRefs has left the pdpRef as written, which it does for a pending
// AuthzPdp only: the platform has no routing scope to check the schema against yet. The AuthzPdp is read again to
// name it. One that synced since resolveRefs read it still defers: ToDTO cannot map an unresolved pdpRef.
func deferDryRun(ctx context.Context, schema *v1alpha1.AuthzSchema) string {
	target := schema.Spec.Target
	if target == nil || target.PdpRef == nil {
		return ""
	}
	if _, err := dynamic.ResolveAuthzTarget(ctx, schema, target); err != nil {
		return fmt.Sprintf("platform checks deferred (%s)", err)
	}
	return syncedDuringAdmission
}

// deferDrift skips drift while either version cannot be mapped as the platform holds it. A pending target gives
// the dry run's reason, which the API server shows once. An old target that no longer resolves is usually a
// deleted AuthzPdp, whose scope the platform has already removed from the document: comparing the old target
// against it would refuse the edit that retargets the schema.
func deferDrift(ctx context.Context, oldObj, newObj *v1alpha1.AuthzSchema) string {
	if reason := deferDryRun(ctx, newObj); reason != "" {
		return reason
	}
	if _, err := dynamic.ResolveAuthzTarget(ctx, oldObj, oldObj.Spec.Target); err != nil {
		return "drift detection skipped, the previous target no longer resolves: " + err.Error()
	}
	return ""
}
