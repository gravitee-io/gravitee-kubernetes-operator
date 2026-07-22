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

package v4

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/apidefinition"
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
	oldApi, _ := oldObj.(*v1alpha1.ApiV4Definition)
	newApi, _ := newObj.(*v1alpha1.ApiV4Definition)
	if !newApi.HasContext() {
		return
	}

	errs.MergeWith(drift.ValidateDrift(ctx, oldApi, newApi, resolveApiV4Refs, getRemoteApiV4,
		drift.MapDTO(toAPIV4DTO)))
}

func resolveApiV4Refs(ctx context.Context, o runtime.Object) error {
	api, ok := o.(*v1alpha1.ApiV4Definition)
	if !ok {
		return nil
	}
	return apidefinition.PrepareV4SpecForAutomation(ctx, api)
}

func toAPIV4DTO(api *v1alpha1.ApiV4Definition) model.APIV4DTO {
	return model.ToAPIV4DTO(&api.Spec.Api)
}

func getRemoteApiV4(apimClient *apim.APIM, o runtime.Object, errs *errors.AdmissionErrors) any {
	api, _ := o.(*v1alpha1.ApiV4Definition)
	api.PopulateIDs(apimClient.Context, k8s.IsAutomationAPIManaged(api))
	hrid := apiHRID(api)
	if !k8s.IsAutomationAPIManaged(api) && api.Spec.ID != "" {
		remote, err := apimClient.APIs.GetV4ByID(hrid)
		if err != nil {
			errs.AddSeveref("cannot fetch API v4 during drift detection from ID %s: %s", hrid, err.Error())
			return nil
		}
		return model.ToAPIV4DTO(remote)
	}
	remote, err := apimClient.APIs.GetV4ByHRID(hrid)
	if err != nil {
		errs.AddSeveref("cannot fetch API v4 during drift detection from HRID %s: %s", hrid, err.Error())
		return nil
	}
	return *remote
}

func apiHRID(api *v1alpha1.ApiV4Definition) string {
	if api.Spec.HRID != "" {
		return api.Spec.HRID
	}
	return refs.NewNamespacedNameFromObject(api).HRID()
}
