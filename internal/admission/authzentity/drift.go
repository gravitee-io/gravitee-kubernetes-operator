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

package authzentity

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/template"
)

func validateUpdate(
	ctx context.Context,
	oldObj *v1alpha1.AuthzEntity,
	newObj *v1alpha1.AuthzEntity,
) *errors.AdmissionErrors {
	// A metadata-only update, the controller's own finalizer and annotation writes included,
	// does not need the platform: only a spec change is dry-run. Hashed before validateContext,
	// which compiles templates into newObj while oldObj keeps them: hashed after, a templated spec
	// would always look changed.
	specChanged := oldObj.Spec.Hash() != newObj.Spec.Hash()

	errs := validateContext(ctx, newObj)
	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(validateTargetKinds(newObj))
	if errs.IsSevere() {
		return errs
	}

	// Resolved on every update, metadata-only ones included: drift resolves them too, and would
	// refuse the controller's own writes while an AuthzPdp is pending.
	targetErrs, resolved := resolveTargets(ctx, newObj)
	errs.MergeWith(targetErrs)
	if !resolved {
		return errs
	}

	if specChanged {
		errs.MergeWith(validateDryRun(ctx, newObj))
		if errs.IsSevere() {
			return errs
		}
	}

	// An entity the platform never received, one waiting on an AuthzPdp for instance, has nothing
	// to drift from: comparing would apply the remote-missing policy (deny by default) and refuse the
	// very edit that unblocks it.
	if oldObj.Status.ID == "" {
		return errs
	}

	// The platform refused the last apply: the stored spec is not what it holds, and this edit is how
	// the entity gets fixed. A RESOURCE the platform deleted with its parent is one: applied again, it
	// is refused for that parent, and the edit dropping it would read the missing remote as drift. The
	// dry run above has checked the new spec. What this gives up: an out-of-band edit made while the
	// entity is refused goes undetected. Drift compares the platform's copy with the stored spec, which
	// the platform never held, so it would refuse every edit that fixes the entity rather than tell
	// such an edit apart. The skip ends with the next successful apply, which clears status.errors.
	if len(oldObj.Status.Errors.Severe) > 0 {
		return errs
	}

	// An old target that no longer resolves is usually a deleted AuthzPdp, whose scope the platform
	// has already removed from the entity: comparing the old targets against it would refuse the
	// edit that retargets the entity.
	if _, err := dynamic.ResolveAuthzTargets(ctx, oldObj, oldObj.Spec.Targets); err != nil {
		return errs
	}

	errs.MergeWith(drift.ValidateDrift(ctx, oldObj, newObj, resolveRefs, getRemoteAuthzEntity,
		drift.MapDTO(model.ToAuthzEntityDTO)))
	return errs
}

// resolveRefs compiles templates on the copy drift compares, then resolves its refs as the
// reconcile does, so both sides are compared with what the operator sends.
func resolveRefs(ctx context.Context, entity *v1alpha1.AuthzEntity) error {
	if err := template.Compile(ctx, entity, false); err != nil {
		return err
	}
	return dynamic.ResolveAuthzEntityRefs(ctx, entity)
}

// getRemoteAuthzEntity reads the platform's copy and returns the payload half.
func getRemoteAuthzEntity(apimClient *apim.APIM, entity *v1alpha1.AuthzEntity) (any, error) {
	hrid := refs.NewNamespacedNameFromObject(entity).HRID()
	remote, err := apimClient.AuthzEntities.GetByHRID(hrid)
	if err != nil {
		return nil, err
	}
	return remote.AuthzEntityDTO, nil
}
