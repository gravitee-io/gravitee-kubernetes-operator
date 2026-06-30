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
	"k8s.io/apimachinery/pkg/runtime"
)

func mergeDriftValidation(
	ctx context.Context,
	oldObj runtime.Object,
	newObj runtime.Object,
	errs *errors.AdmissionErrors,
) {
	oldSpg, _ := oldObj.(*v1alpha1.SharedPolicyGroup)
	newSpg, _ := newObj.(*v1alpha1.SharedPolicyGroup)
	errs.MergeWith(drift.ValidateDrift(ctx, oldSpg, newSpg, resolveRefs, getRemoteSharedPolicyGroup,
		drift.MapDTO(toSharedPolicyGroupPayload)))
}

func resolveRefs(context.Context, runtime.Object) error {
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
	o runtime.Object,
	errs *errors.AdmissionErrors,
) any {
	spg, _ := o.(*v1alpha1.SharedPolicyGroup)
	spg.PopulateIDs(apimClient.Context, k8s.IsAutomationAPIManaged(spg))
	hrid, legacy := sharedPolicyGroupID(spg)
	if legacy {
		errs.AddSeveref("drift detection is not supported for legacy SharedPolicyGroup [%s]", spg.GetRef())
		return nil
	}
	remote, err := apimClient.SharedPolicyGroup.GetByHRID(hrid)
	if err != nil {
		errs.AddSeveref(
			"cannot fetch SharedPolicyGroup during drift detection from HRID %s: %s",
			hrid, err.Error(),
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
