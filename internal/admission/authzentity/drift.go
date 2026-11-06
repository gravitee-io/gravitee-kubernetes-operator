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
	"slices"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzentity"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
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

	errs.MergeWith(drift.ValidateDrift(ctx, oldObj, newObj, resolveRefs, getRemoteAuthzEntity(ctx, oldObj),
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

// kindResource is the kind the platform deletes, rather than detaches, with its parent.
const kindResource = "RESOURCE"

// getRemoteAuthzEntity reads the platform's copy and returns the payload half, with the cascade of
// a deleted parent undone. Deleting an AuthzEntity deletes each RESOURCE that declares it and
// detaches each PRINCIPAL: they differ from their manifests until their spec next changes, and that
// change must not read as drift. A parent of the old spec that no AuthzEntity of its context declares
// any more is taken as deleted: the parents among those the platform no longer lists are added back,
// and a RESOURCE the platform no longer holds reads as its old spec. Every other field is compared.
// A parent another source created reads as deleted too: an entity detached from it, or a RESOURCE
// deleted, outside the operator is not reported.
func getRemoteAuthzEntity(
	ctx context.Context,
	oldObj *v1alpha1.AuthzEntity,
) drift.RemoteObjectGetter[*v1alpha1.AuthzEntity] {
	return func(apimClient *apim.APIM, entity *v1alpha1.AuthzEntity) (any, error) {
		hrid := refs.NewNamespacedNameFromObject(entity).HRID()
		remote, err := apimClient.AuthzEntities.GetByHRID(hrid)
		missing := errors.IsNotFound(err)
		if err != nil && !missing {
			return nil, err
		}
		if missing && oldObj.Spec.Kind != kindResource {
			return nil, err
		}
		if !missing && !lacksParent(remote.Parents, oldObj.Spec.Parents) {
			return remote.AuthzEntityDTO, nil
		}

		old := oldObj.DeepCopy()
		deleted, findErr := cascadedParents(ctx, old)
		if findErr != nil || len(deleted) == 0 {
			if missing {
				return nil, err
			}
			return remote.AuthzEntityDTO, nil
		}
		if missing {
			return toDriftDTO(old), nil
		}

		dto := remote.AuthzEntityDTO
		for _, parent := range deleted {
			if restored := model.ToAuthzEntityParentDTO(parent); !slices.Contains(dto.Parents, restored) {
				dto.Parents = append(slices.Clone(dto.Parents), restored)
			}
		}
		return dto, nil
	}
}

// cascadedParents resolves old as drift compares it, then returns the parents no AuthzEntity of its
// context declares any more.
func cascadedParents(ctx context.Context, old *v1alpha1.AuthzEntity) ([]authzentity.Parent, error) {
	if err := resolveRefs(ctx, old); err != nil {
		return nil, err
	}
	return search.FindUndeclaredAuthzEntityParents(ctx, old)
}

// lacksParent tells if the platform no longer lists one of the parents the old spec declares.
func lacksParent(remote []model.AuthzEntityParentDTO, declared []authzentity.Parent) bool {
	for _, parent := range declared {
		if !slices.Contains(remote, model.ToAuthzEntityParentDTO(parent)) {
			return true
		}
	}
	return false
}
