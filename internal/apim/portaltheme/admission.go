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

package portaltheme

import (
	"context"
	"net/http"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/navigation"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

// DryRun validates the theme against APIM without persisting it. It returns APIM's validation errors as
// findings, and APIM's error when the call fails.
func DryRun(ctx context.Context, client *apim.SDKClient, dto sdk.PortalTheme) (*errors.AdmissionErrors, error) {
	if client == nil {
		return nil, nil
	}
	state, err := put(ctx, client, dto, true)
	if err != nil {
		return nil, err
	}
	return navigation.Findings(state.Errors), nil
}

// GetRemote fetches the theme from APIM as the payload a PUT would send, for drift detection.
func GetRemote(ctx context.Context, client *apim.SDKClient, dto sdk.PortalTheme) (sdk.PortalTheme, error) {
	state, err := GetState(ctx, client, dto)
	if err != nil {
		return sdk.PortalTheme{}, err
	}
	return state.Spec(), nil
}

// GetState fetches the theme state from APIM, including the server-managed fields.
func GetState(ctx context.Context, client *apim.SDKClient, dto sdk.PortalTheme) (sdk.PortalThemeState, error) {
	if client == nil {
		return sdk.PortalThemeState{}, navigation.ErrNoClient
	}
	resp, err := client.GetPortalThemeWithResponse(ctx, dto.Identity())
	if err = errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return sdk.PortalThemeState{}, err
	}
	if resp.JSON200 == nil {
		return sdk.PortalThemeState{}, errors.UnexpectedSDKResponse(resp.HTTPResponse)
	}
	return *resp.JSON200, nil
}
