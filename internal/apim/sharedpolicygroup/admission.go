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

package sharedpolicygroup

import (
	"context"
	"net/http"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

// PayloadChanged dry-runs every update, as before the lifecycle: an update APIM would refuse is refused, whatever
// changed.
func PayloadChanged(_, _ *v1alpha1.SharedPolicyGroup) bool {
	return true
}

// DryRun validates the group against APIM without persisting it. Warnings are reported only when APIM finds
// nothing severe.
func DryRun(ctx context.Context, client *apim.SDKClient, dto DTO) (*errors.AdmissionErrors, error) {
	state, err := upsert(ctx, client, dto, true)
	if err != nil {
		return nil, err
	}
	errs := errors.NewAdmissionErrors()
	if state.Errors == nil {
		return errs, nil
	}
	for _, severe := range state.Errors.Severe {
		errs.AddSevere(severe)
	}
	if errs.IsSevere() {
		return errs, nil
	}
	for _, warning := range state.Errors.Warning {
		errs.AddWarning(warning)
	}
	return errs, nil
}

// GetRemote fetches the group from APIM as the payload a PUT would send, for drift detection.
func GetRemote(ctx context.Context, client *apim.SDKClient, dto DTO) (DTO, error) {
	state, err := GetState(ctx, client, dto)
	if err != nil {
		return DTO{}, err
	}
	return DTO{SharedPolicyGroup: state.Spec(), HridContainsUUID: dto.HridContainsUUID}, nil
}

// GetState fetches the group state from APIM, server-managed fields (id, cross id) included.
func GetState(ctx context.Context, client *apim.SDKClient, dto DTO) (sdk.SharedPolicyGroupState, error) {
	resp, err := client.GetSharedPolicyGroupWithResponse(ctx, dto.Identity(),
		&sdk.GetSharedPolicyGroupParams{HridContainsUUID: &dto.HridContainsUUID})
	if err = errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return sdk.SharedPolicyGroupState{}, err
	}
	if resp.JSON200 == nil {
		return sdk.SharedPolicyGroupState{}, errors.UnexpectedSDKResponse(resp.HTTPResponse)
	}
	return *resp.JSON200, nil
}
