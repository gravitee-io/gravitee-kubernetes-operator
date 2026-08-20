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

package policygroups

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

func mergeDriftValidation(
	ctx context.Context,
	oldSpg *v1alpha1.SharedPolicyGroup,
	newSpg *v1alpha1.SharedPolicyGroup,
	errs *errors.AdmissionErrors,
) {
	errs.MergeWith(drift.ValidateDrift(ctx, oldSpg, newSpg, resolveRefs, getRemoteSharedPolicyGroup,
		drift.MapDTO(toSharedPolicyGroupPayload)))
}

func resolveRefs(context.Context, *v1alpha1.SharedPolicyGroup) error {
	return nil
}

func toSharedPolicyGroupPayload(spg *v1alpha1.SharedPolicyGroup) model.SharedPolicyGroupDTO {
	if spg.Spec.SharedPolicyGroup == nil {
		return model.SharedPolicyGroupDTO{}
	}
	return model.ToSharePolicyGroupDTO(*spg.Spec.SharedPolicyGroup)
}

func getRemoteSharedPolicyGroup(
	apimClient *apim.APIM,
	spg *v1alpha1.SharedPolicyGroup,
	errs *errors.AdmissionErrors,
) any {
	spg.PopulateIDs(apimClient.Context, k8s.IsAutomationAPIManaged(spg))
	id, legacy := sharedPolicyGroupID(spg)
	var (
		remote *model.SharedPolicyGroupDTO
		err    error
	)
	if legacy {
		remote, err = apimClient.SharedPolicyGroup.GetByID(id)
	} else {
		remote, err = apimClient.SharedPolicyGroup.GetByHRID(id)
	}
	if err != nil {
		kind := "HRID"
		if legacy {
			kind = "ID"
		}
		errs.AddSeveref(
			"cannot fetch SharedPolicyGroup during drift detection from %s %s: %s",
			kind, id, err.Error(),
		)
		return nil
	}
	return *remote
}

func sharedPolicyGroupID(spg *v1alpha1.SharedPolicyGroup) (string, bool) {
	if k8s.IsAutomationAPIManaged(spg) {
		return refs.NewNamespacedNameFromObject(spg).HRID(), false
	}
	return spg.GetID(), true
}
