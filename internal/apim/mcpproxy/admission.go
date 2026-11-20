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

package mcpproxy

import (
	"context"
	"net/http"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

// DryRun validates the proxy against the platform without persisting it. A dry run answers 200 with its
// findings; a refusal body is only expected from a real apply, but is read the same way should the platform
// answer one. Any other failure returns the error.
func DryRun(ctx context.Context, client *apim.SDKClient, dto sdk.AimMcpProxy) (*errors.AdmissionErrors, error) {
	state, err := upsert(ctx, client, dto, true)
	if err != nil {
		if findings, refused := model.AutomationRefusal(err); refused {
			return errors.NewAdmissionErrorsFromStatus(findings), nil
		}
		return nil, err
	}
	return errors.NewAdmissionErrorsFromStatus(toStatusErrors(state.Errors)), nil
}

// GetRemote reads the proxy as the payload a PUT would send, in the order the mapper produces, for drift
// detection. The tools status lives in the state only and is never compared.
func GetRemote(ctx context.Context, client *apim.SDKClient, dto sdk.AimMcpProxy) (sdk.AimMcpProxy, error) {
	resp, err := client.GetAimMcpProxyWithResponse(ctx, dto.Identity())
	if err = errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return sdk.AimMcpProxy{}, err
	}
	if resp.JSON200 == nil {
		return sdk.AimMcpProxy{}, errors.UnexpectedSDKResponse(resp.HTTPResponse)
	}
	return Canonical(resp.JSON200.Spec()), nil
}
