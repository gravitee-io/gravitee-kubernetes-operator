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

package apim

import (
	"context"
	"fmt"
	"net/http"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/apiv4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/mapper"
)

// UpsertAPIV4 PUTs the spec of api to APIM as it is (IDs are not populated again), to change the API behind the
// operator's back. An API with a spec.id is sent with its UUID as the hrid.
func (a *APIM) UpsertAPIV4(ctx context.Context, api *v1alpha1.ApiV4Definition) error {
	payload, err := mapper.MapViaJSON[sdk.ApiV4](model.ToAPIV4DTO(&api.Spec.Api))
	if err != nil {
		return err
	}
	ignoreGroups := k8s.HasTrueAnnotation(api, core.IgnoreGroupsAnnotation)
	if api.Spec.ID != "" {
		payload.Hrid = api.Spec.ID
		_, err := a.APIs.UpsertV4Legacy(payload, false, ignoreGroups)
		return err
	}
	resp, err := a.SDK.UpsertApiWithResponse(ctx, &sdk.UpsertApiParams{DryRun: new(false), IgnoreGroups: &ignoreGroups}, payload)
	return errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	})
}

// DeleteAPIV4 deletes api in APIM, addressed as the operator addresses it on delete, then checks that APIM no
// longer has it: apiv4.Delete counts a 404 as deleted, which would hide a wrong identity.
func (a *APIM) DeleteAPIV4(ctx context.Context, api *v1alpha1.ApiV4Definition) error {
	dto, err := apiv4.ToDTO(api)
	if err != nil {
		return err
	}
	client := &apiv4.Client{SDKClient: a.SDK, APIM: a.APIM}
	if err = apiv4.Delete(ctx, client, dto); err != nil {
		return err
	}
	_, err = apiv4.GetState(ctx, client, dto)
	if errors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return fmt.Errorf("API %s is still in APIM after its delete", dto.DeleteID)
}
