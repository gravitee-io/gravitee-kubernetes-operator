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

package apiv4

import (
	"context"
	"net/http"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

// DryRun sends the payload with dryRun=true and reports APIM's findings: severe ones only when there are any,
// else the warnings. An API without a management context has nothing to dry-run.
func DryRun(ctx context.Context, client *Client, dto DTO) (*errors.AdmissionErrors, error) {
	if client == nil {
		return nil, nil
	}
	state, err := put(ctx, client, dto, true)
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

// GetRemote fetches the API from APIM as a DTO to compare with the CR's. Before the first sync (no status.id),
// or without a management context, there is nothing to fetch and the CR's DTO is returned.
func GetRemote(ctx context.Context, client *Client, dto DTO) (DTO, error) {
	if client == nil || !dto.Synced {
		return dto, nil
	}
	state, err := GetState(ctx, client, dto)
	if err != nil {
		return DTO{}, err
	}
	return RemoteDTO(dto, state), nil
}

// RemoteDTO is the DTO of the fetched state, carrying the flags of the CR's dto. Under ignore-groups it carries
// the CR's groups: APIM leaves the groups alone then, so they cannot drift.
func RemoteDTO(dto DTO, state sdk.ApiV4State) DTO {
	// Spec is a shallow copy of the state: only reassign its fields, never mutate shared slices or maps.
	remote := DTO{
		ApiV4:                  state.Spec(),
		HridContainsUUID:       dto.HridContainsUUID,
		DeleteID:               dto.DeleteID,
		DeleteHridContainsUUID: dto.DeleteHridContainsUUID,
		IgnoreGroups:           dto.IgnoreGroups,
		Synced:                 dto.Synced,
		Gateway:                dto.Gateway,
	}
	if dto.PlanOrder != nil {
		remote.PlanOrder = planRanks(remote.Plans)
	}
	if dto.IgnoreGroups {
		remote.Groups = dto.Groups
	}
	return remote
}

// GetState fetches the API from APIM, by hrid or, for a legacy API, by UUID.
func GetState(ctx context.Context, client *Client, dto DTO) (sdk.ApiV4State, error) {
	resp, err := client.GetApiWithResponse(ctx, dto.Identity(),
		&sdk.GetApiParams{HridContainsUUID: &dto.HridContainsUUID})
	if err = errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return sdk.ApiV4State{}, err
	}
	if resp.JSON200 == nil {
		return sdk.ApiV4State{}, errors.UnexpectedSDKResponse(resp.HTTPResponse)
	}
	return *resp.JSON200, nil
}
