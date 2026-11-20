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

package catalogmcpserver

import (
	"context"
	"net/http"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

// DryRun validates the server against the catalog without persisting it. The platform discovers the
// upstream first, so an unreachable server or a taken entityId comes back as a severe finding. A refusal
// returns its findings, any other failure returns the error.
func DryRun(ctx context.Context, client *apim.SDKClient, dto sdk.AimCatalogMcpServer) (
	*errors.AdmissionErrors, error,
) {
	state, err := upsert(ctx, client, dto, true)
	if err != nil {
		if findings, refused := model.AutomationRefusal(err); refused {
			return errors.NewAdmissionErrorsFromStatus(findings), nil
		}
		return nil, err
	}
	return errors.NewAdmissionErrorsFromStatus(ToStatusErrors(state.Errors)), nil
}

// GetRemote reads the server from the catalog as the payload a PUT would send, for drift detection. What
// the platform discovered lives in the state only and is never compared.
func GetRemote(ctx context.Context, client *apim.SDKClient, dto sdk.AimCatalogMcpServer) (
	sdk.AimCatalogMcpServer, error,
) {
	resp, err := client.GetAimCatalogMcpServerWithResponse(ctx, dto.Identity())
	if err = errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return sdk.AimCatalogMcpServer{}, err
	}
	if resp.JSON200 == nil {
		return sdk.AimCatalogMcpServer{}, errors.UnexpectedSDKResponse(resp.HTTPResponse)
	}
	return resp.JSON200.Spec(), nil
}
