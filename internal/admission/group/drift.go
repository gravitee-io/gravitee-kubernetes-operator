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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/service"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"k8s.io/apimachinery/pkg/runtime"
)

func validateUpdate(ctx context.Context, oldObj runtime.Object, newObj runtime.Object) *errors.AdmissionErrors {
	errs := validateCreate(ctx, newObj)
	if errs.IsSevere() {
		return errs
	}
	oldGrp, _ := oldObj.(*v1alpha1.Group)
	newGrp, _ := newObj.(*v1alpha1.Group)
	errs.MergeWith(drift.ValidateDrift(ctx, oldGrp, newGrp, resolveRefs, getRemoteGroup,
		drift.MapDTO(toGroupPayload)))
	return errs
}

func resolveRefs(context.Context, runtime.Object) error {
	return nil
}

func toGroupPayload(grp *v1alpha1.Group) model.GroupDTO {
	return model.ToGroupDTO(*grp.Spec.Type)
}

func getRemoteGroup(apimClient *apim.APIM, o runtime.Object, errs *errors.AdmissionErrors) any {
	grp, _ := o.(*v1alpha1.Group)
	grp.PopulateIDs(apimClient.Context, k8s.IsAutomationAPIManaged(grp))
	hrid, legacy := service.GetGroupID(grp)
	if legacy {
		errs.AddSeveref("drift detection is not supported for legacy Group [%s]", grp.GetRef())
		return nil
	}
	remote, err := apimClient.Env.GetGroupByHRID(hrid)
	if err != nil {
		errs.AddSeveref("cannot fetch Group during drift detection from HRID %s: %s", hrid, err.Error())
		return nil
	}
	return remote.GroupDTO
}
