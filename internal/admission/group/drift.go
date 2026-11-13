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

package group

import (
	"context"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

func validateUpdate(ctx context.Context, oldObj *v1alpha1.Group, newObj *v1alpha1.Group) *errors.AdmissionErrors {
	errs := validateCreate(ctx, newObj)
	if errs.IsSevere() {
		return errs
	}
	// The shared remote getter has no context: close over the admission one.
	getRemote := func(apimClient *apim.APIM, grp *v1alpha1.Group) (any, error) {
		return getRemoteGroup(ctx, apimClient, grp)
	}
	errs.MergeWith(drift.ValidateDrift(ctx, oldObj, newObj, resolveRefs, getRemote,
		drift.MapDTO(toGroupPayload)))
	return errs
}

func resolveRefs(context.Context, *v1alpha1.Group) error {
	return nil
}

// toGroupPayload maps the CRD to the payload compared with the remote. The same spec already went
// through the dry run in validateCreate, so a mapping error is not expected here; should one occur,
// the empty payload is reported as drift rather than silently admitted.
func toGroupPayload(grp *v1alpha1.Group) sdk.Group {
	payload, err := model.ToGroup(grp)
	if err != nil {
		return sdk.Group{}
	}
	if k8s.HasTrueAnnotation(grp, core.IgnoreMembersAnnotation) {
		payload.Members = nil
	}
	return payload
}

func getRemoteGroup(ctx context.Context, apimClient *apim.APIM, grp *v1alpha1.Group) (any, error) {
	automation := k8s.IsAutomationAPIManaged(grp)
	grp.PopulateIDs(apimClient.Context, automation)
	var (
		state sdk.GroupState
		err   error
	)
	if !automation && grp.GetID() != "" {
		state, err = apimClient.Env.GetGroupWithUUID(ctx, grp.GetID())
	} else {
		state, err = apimClient.Env.GetGroupByHRID(ctx, groupHRID(grp))
	}
	if err != nil {
		return nil, err
	}
	// Spec is the payload a PUT would send, so drift compares two values of the same type.
	// It is a shallow copy: only reassign its fields here, never mutate shared slices or maps.
	remote := state.Spec()
	if k8s.HasTrueAnnotation(grp, core.IgnoreMembersAnnotation) {
		remote.Members = nil
	}
	return remote, nil
}

func groupHRID(grp *v1alpha1.Group) string {
	if grp.Spec.HRID != "" {
		return grp.Spec.HRID
	}
	return refs.NewNamespacedNameFromObject(grp).HRID()
}
