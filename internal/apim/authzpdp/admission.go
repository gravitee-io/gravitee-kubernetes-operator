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

package authzpdp

import (
	"context"
	"net/http"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

// DryRun validates the PDP without persisting it. The platform refuses a routing scope another PDP holds, a
// strict value that disagrees with the engine and an endpoint path another API uses. A dry run answers 200 with
// its findings; a refusal body is read the same way should the platform answer one. Any other failure, the 403
// of a strict change without the permission included, returns the error.
func DryRun(ctx context.Context, client *apim.SDKClient, dto sdk.AuthzPdp) (*errors.AdmissionErrors, error) {
	state, err := upsert(ctx, client, dto, true)
	if err != nil {
		if findings, refused := model.AutomationRefusal(err); refused {
			return errors.NewAdmissionErrorsFromStatus(findings), nil
		}
		return nil, err
	}
	return errors.NewAdmissionErrorsFromStatus(ToStatusErrors(state.Errors)), nil
}

// GetRemote reads the PDP as the payload a PUT would send, for drift detection. The routing scope and the
// discovery path live in the state only and are never compared. A strict change on a sibling of the same
// targetPdpId shows up as drift, which is what happened: the platform changed outside this resource.
func GetRemote(ctx context.Context, client *apim.SDKClient, dto sdk.AuthzPdp) (sdk.AuthzPdp, error) {
	resp, err := client.GetAuthzPdpWithResponse(ctx, dto.Hrid)
	if err = errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return sdk.AuthzPdp{}, err
	}
	if resp.JSON200 == nil {
		return sdk.AuthzPdp{}, errors.UnexpectedSDKResponse(resp.HTTPResponse)
	}
	return resp.JSON200.Spec(), nil
}
