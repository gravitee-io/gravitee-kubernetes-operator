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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/navigation"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

// Upsert creates or updates the link under its portal or API in APIM.
func Upsert(ctx context.Context, client *apim.SDKClient, dto DTO) (Response, error) {
	state, err := put(ctx, client, dto, false)
	if err != nil {
		return Response{}, err
	}
	return newResponse(state), nil
}

func put(ctx context.Context, client *apim.SDKClient, dto DTO, dryRun bool) (sdk.PortalLinkState, error) {
	if client == nil {
		return sdk.PortalLinkState{}, navigation.ErrNoClient
	}
	var (
		state    *sdk.PortalLinkState
		httpResp func() (*http.Response, []byte)
		err      error
	)
	if dto.Parent.IsAPI() {
		var resp *sdk.UpsertApiLinkResponse
		resp, err = client.UpsertApiLinkWithResponse(ctx, dto.Parent.ApiHrid,
			&sdk.UpsertApiLinkParams{DryRun: &dryRun}, dto.PortalLink)
		httpResp = func() (*http.Response, []byte) { return resp.HTTPResponse, resp.Body }
		if err == nil {
			state = resp.JSON200
		}
	} else {
		var resp *sdk.UpsertPortalLinkResponse
		resp, err = client.UpsertPortalLinkWithResponse(ctx, dto.Parent.PortalHrid,
			&sdk.UpsertPortalLinkParams{DryRun: &dryRun}, dto.PortalLink)
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

// Delete removes the link from its portal or API in APIM. A link already gone counts as deleted.
func Delete(ctx context.Context, client *apim.SDKClient, dto DTO) error {
	if client == nil {
		return nil
	}
	var (
		httpResp func() (*http.Response, []byte)
		err      error
	)
	if dto.Parent.IsAPI() {
		var resp *sdk.DeleteApiLinkResponse
		resp, err = client.DeleteApiLinkWithResponse(ctx, dto.Parent.ApiHrid, dto.Identity())
		httpResp = func() (*http.Response, []byte) { return resp.HTTPResponse, resp.Body }
	} else {
		var resp *sdk.DeletePortalLinkResponse
		resp, err = client.DeletePortalLinkWithResponse(ctx, dto.Parent.PortalHrid, dto.Identity())
		httpResp = func() (*http.Response, []byte) { return resp.HTTPResponse, resp.Body }
	}
	return errors.IgnoreNotFound(errors.FromSDKResponse(err, httpResp))
}

// UpdateStatus copies the ids of the upsert state into the status, field by field so that the rest of
// the status is kept. Every write is HRID-based, so the link is Automation API managed.
func UpdateStatus(_ context.Context, obj *v1alpha1.PortalLink, resp Response) error {
	obj.Status.ID = utils.SafeDereference(resp.State.Id)
	obj.Status.OrgID = resp.GetOrgID()
	obj.Status.EnvID = resp.GetEnvID()
	k8s.AddAutomationAPIManagedCondition(obj)
	return nil
}
