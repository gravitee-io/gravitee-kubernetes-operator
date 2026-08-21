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

package portal

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

func validateUpdate(ctx context.Context, oldObj *v1alpha1.Portal, newObj *v1alpha1.Portal) *errors.AdmissionErrors {
	errs := validateCreate(ctx, newObj)
	if errs.IsSevere() {
		return errs
	}
	errs.MergeWith(drift.ValidateDrift(ctx, oldObj, newObj, resolveRefs, getRemotePortal,
		drift.MapDTO(toPortalDTO)))
	return errs
}

func resolveRefs(context.Context, *v1alpha1.Portal) error {
	return nil
}

func toPortalDTO(prtl *v1alpha1.Portal) model.PortalDTO {
	return model.ToPortalDTO(prtl.Spec.Type, refs.NewNamespacedNameFromObject(prtl).HRID())
}

func getRemotePortal(apimClient *apim.APIM, prtl *v1alpha1.Portal) (any, error) {
	hrid := refs.NewNamespacedNameFromObject(prtl).HRID()
	remote, err := apimClient.Portals.GetByHRID(hrid)
	if err != nil {
		return nil, err
	}
	return remote.PortalDTO, nil
}
