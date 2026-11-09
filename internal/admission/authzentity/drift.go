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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"
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
	targets, targetErrs, resolved := resolveTargets(ctx, newObj)
	errs.MergeWith(targetErrs)
	if !resolved {
		return errs
	}

	if specChanged {
		errs.MergeWith(validateDryRun(ctx, newObj, targets))
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

	// An old target that no longer resolves is usually a deleted AuthzPdp, whose scope the platform
	// has already removed from the entity: comparing the old targets against it would refuse the
	// edit that retargets the entity.
	if _, err := dynamic.ResolveAuthzTargets(ctx, oldObj, oldObj.Spec.Targets); err != nil {
		return errs
	}

	errs.MergeWith(drift.ValidateDrift(ctx, oldObj, newObj, resolveRefs, getRemoteAuthzEntity,
		drift.MapDTO(toDriftDTO)))
	return errs
}

// resolveRefs compiles templates on the copy drift compares, then replaces its targets with the
// scopes they resolve to, so both sides are compared with what the operator sends.
func resolveRefs(ctx context.Context, entity *v1alpha1.AuthzEntity) error {
	if err := template.Compile(ctx, entity, false); err != nil {
		return err
	}
	targets, err := dynamic.ResolveAuthzTargets(ctx, entity, entity.Spec.Targets)
	if err != nil {
		return err
	}
	entity.Spec.Targets = nil
	for _, target := range targets {
		entity.Spec.Targets = append(entity.Spec.Targets, authzpdp.Target{Scope: target})
	}
	return nil
}

// toDriftDTO maps a copy resolveRefs has resolved: each of its targets is a scope.
func toDriftDTO(entity *v1alpha1.AuthzEntity) model.AuthzEntityDTO {
	var targets []string
	for _, target := range entity.Spec.Targets {
		targets = append(targets, target.Scope)
	}
	return model.ToAuthzEntityDTO(entity, targets)
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
