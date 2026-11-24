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

package authzpolicy

import (
	"context"
	"net/http"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/authz"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

// DryRun validates the policy with its resolved targets without persisting it. The platform refuses an entityId
// another policy holds, a text that does not parse, a template and a violation of a strict gateway's schema, and
// warns about an MCP policy the default engine never sees and a tool named by its catalog id. A dry run answers 200
// with its findings; a refusal body is read the same way should the platform answer one. Any other failure returns
// the error.
func DryRun(ctx context.Context, client *apim.SDKClient, dto sdk.AuthzPolicy) (*errors.AdmissionErrors, error) {
	state, err := upsert(ctx, client, dto, true)
	if err != nil {
		if findings, refused := model.AutomationRefusal(err); refused {
			return errors.NewAdmissionErrorsFromStatus(findings), nil
		}
		return nil, err
	}
	return errors.NewAdmissionErrorsFromStatus(authz.ToStatusErrors(state.Errors)), nil
}

// GetRemote reads the policy as the payload a PUT would send, for drift detection. Spec() leaves policyStatus out:
// a policy disabled in the Console is not drift, and the next apply deploys it again with the platform's warning.
func GetRemote(ctx context.Context, client *apim.SDKClient, dto sdk.AuthzPolicy) (sdk.AuthzPolicy, error) {
	resp, err := client.GetAuthzPolicyWithResponse(ctx, dto.Hrid)
	if err = errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return sdk.AuthzPolicy{}, err
	}
	if resp.JSON200 == nil {
		return sdk.AuthzPolicy{}, errors.UnexpectedSDKResponse(resp.HTTPResponse)
	}
	return resp.JSON200.Spec(), nil
}
