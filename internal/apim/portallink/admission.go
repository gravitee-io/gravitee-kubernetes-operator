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

package portallink

import (
	"context"
	"net/http"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/navigation"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

// DryRun validates the link against APIM without persisting it. It returns APIM's validation errors as
// findings, and APIM's error when the call fails.
func DryRun(ctx context.Context, client *apim.SDKClient, dto DTO) (*errors.AdmissionErrors, error) {
	if client == nil {
		return nil, nil
	}
	state, err := put(ctx, client, dto, true)
	if err != nil {
		return nil, err
	}
	return navigation.Findings(state.Errors), nil
}

// GetRemote fetches the link from APIM as the payload a PUT would send, for drift detection.
func GetRemote(ctx context.Context, client *apim.SDKClient, dto DTO) (DTO, error) {
	state, err := GetState(ctx, client, dto)
	if err != nil {
		return DTO{}, err
	}
	return DTO{PortalLink: state.Spec(), Parent: dto.Parent}, nil
}

// GetState fetches the link state from APIM, including the server-managed fields.
func GetState(ctx context.Context, client *apim.SDKClient, dto DTO) (sdk.PortalLinkState, error) {
	var (
		state    *sdk.PortalLinkState
		httpResp func() (*http.Response, []byte)
		err      error
	)
	if dto.Parent.IsAPI() {
		var resp *sdk.GetApiLinkResponse
		resp, err = client.GetApiLinkWithResponse(ctx, dto.Parent.ApiHrid, dto.Identity())
		httpResp = func() (*http.Response, []byte) { return resp.HTTPResponse, resp.Body }
		if err == nil {
			state = resp.JSON200
		}
	} else {
		var resp *sdk.GetPortalLinkResponse
		resp, err = client.GetPortalLinkWithResponse(ctx, dto.Parent.PortalHrid, dto.Identity())
		httpResp = func() (*http.Response, []byte) { return resp.HTTPResponse, resp.Body }
		if err == nil {
			state = resp.JSON200
		}
	}
	if err = errors.FromSDKResponse(err, httpResp); err != nil {
		return sdk.PortalLinkState{}, err
	}
	if state == nil {
		resp, _ := httpResp()
		return sdk.PortalLinkState{}, errors.UnexpectedSDKResponse(resp)
	}
	return *state, nil
}
