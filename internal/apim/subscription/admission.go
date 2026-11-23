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

package subscription

import (
	"context"
	"net/http"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

// GetRemote fetches the subscription from APIM as the payload a PUT would send, for drift detection.
func GetRemote(ctx context.Context, client *apim.SDKClient, dto DTO) (DTO, error) {
	state, err := GetState(ctx, client, dto)
	if err != nil {
		return DTO{}, err
	}
	return DTO{
		Subscription: state.Spec(),
		ApiHrid:      dto.ApiHrid,
		Legacy:       dto.Legacy,
	}, nil
}

// GetState fetches the subscription state from APIM, including the server-managed fields (id, start date).
func GetState(ctx context.Context, client *apim.SDKClient, dto DTO) (sdk.SubscriptionState, error) {
	resp, err := client.GetApiSubscriptionWithResponse(ctx, dto.ApiHrid, dto.Hrid, dto.Legacy.getParams())
	if err = errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return sdk.SubscriptionState{}, err
	}
	if resp.JSON200 == nil {
		return sdk.SubscriptionState{}, errors.UnexpectedSDKResponse(resp.HTTPResponse)
	}
	return *resp.JSON200, nil
}
