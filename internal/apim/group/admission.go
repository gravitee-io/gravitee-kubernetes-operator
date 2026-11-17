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
	"net/http"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

// DryRun validates the group against APIM without persisting it. It returns APIM's validation errors as
// findings, and APIM's error when the call fails.
func DryRun(ctx context.Context, client *apim.SDKClient, dto DTO) (*errors.AdmissionErrors, error) {
	state, err := upsert(ctx, client, dto, true)
	if err != nil {
		return nil, err
	}
	errs := errors.NewAdmissionErrors()
	if state.Errors != nil {
		for _, severe := range state.Errors.Severe {
			errs.AddSevere(severe)
		}
		for _, warning := range state.Errors.Warning {
			errs.AddWarning(warning)
		}
	}
	return errs, nil
}

// GetRemote fetches the group from APIM as the payload a PUT would send, for drift detection.
// With ignore-members, it carries the CR members in place of the remote ones, so members never drift.
func GetRemote(ctx context.Context, client *apim.SDKClient, dto DTO) (DTO, error) {
	state, err := GetState(ctx, client, dto)
	if err != nil {
		return DTO{}, err
	}
	// Spec is a shallow copy of the state: only reassign its fields, never mutate shared slices or maps.
	remote := DTO{Group: state.Spec(), HridContainsUUID: dto.HridContainsUUID, IgnoreMembers: dto.IgnoreMembers}
	if dto.IgnoreMembers {
		// APIM owns the members under ignore-members: report the CR's own, so they are never compared.
		remote.Members = dto.Members
	}
	return remote, nil
}

// GetState fetches the group state from APIM, including the server-managed fields (id, member count).
func GetState(ctx context.Context, client *apim.SDKClient, dto DTO) (sdk.GroupState, error) {
	resp, err := client.GetGroupWithResponse(ctx, dto.Identity(),
		&sdk.GetGroupParams{HridContainsUUID: &dto.HridContainsUUID})
	if err = errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return sdk.GroupState{}, err
	}
	if resp.JSON200 == nil {
		return sdk.GroupState{}, errors.UnexpectedSDKResponse(resp.HTTPResponse)
	}
	return *resp.JSON200, nil
}
