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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/apidefinition"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/apiv4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

func mergeDriftValidation(
	ctx context.Context,
	oldApi *v1alpha1.ApiV4Definition,
	newApi *v1alpha1.ApiV4Definition,
	errs *errors.AdmissionErrors,
) {
	if !newApi.HasContext() {
		return
	}

	// The API talks to APIM through its SDK client: the resolver builds that client, only once drift is known
	// to be enabled, and hands its APIM facade to the drift framework.
	var client *apiv4.Client
	resolveContext := func(ctx context.Context) (*apim.APIM, error) {
		var err error
		client, err = apiv4.CreateAPIMClient(ctx, newApi)
		if err != nil || client == nil {
			return nil, err
		}
		return client.APIM, nil
	}
	getRemote := func(_ *apim.APIM, api *v1alpha1.ApiV4Definition) (any, error) {
		dto, err := apiv4.ToDTO(api)
		if err != nil {
			return nil, err
		}
		remote, err := apiv4.GetRemote(ctx, client, dto)
		if err != nil {
			return nil, err
		}
		return apiv4.DriftDTO(remote), nil
	}
	errs.MergeWith(drift.ValidateDriftWithContext(ctx, oldApi, newApi, resolveContext, resolveApiV4Refs, getRemote,
		drift.MapDTO(toDriftDTO)))
}

func resolveApiV4Refs(ctx context.Context, api *v1alpha1.ApiV4Definition) error {
	return apidefinition.PrepareV4SpecForAutomation(ctx, api, false)
}

func toDriftDTO(api *v1alpha1.ApiV4Definition) apiv4.DTO {
	dto, err := apiv4.ToDTO(api)
	if err != nil {
		return apiv4.DTO{}
	}
	return apiv4.DriftDTO(dto)
}
