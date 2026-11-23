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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

// Upsert creates or updates the theme in APIM.
func Upsert(ctx context.Context, client *apim.SDKClient, dto sdk.PortalTheme) (Response, error) {
	state, err := put(ctx, client, dto, false)
	if err != nil {
		return Response{}, err
	}
	return newResponse(state), nil
}

func put(ctx context.Context, client *apim.SDKClient, dto sdk.PortalTheme, dryRun bool) (sdk.PortalThemeState, error) {
	resp, err := client.UpsertPortalThemeWithResponse(ctx, &sdk.UpsertPortalThemeParams{DryRun: &dryRun}, dto)
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

// Delete removes the theme from APIM. A theme already gone counts as deleted.
func Delete(ctx context.Context, client *apim.SDKClient, dto sdk.PortalTheme) error {
	if client == nil {
		return nil
	}
	resp, err := client.DeletePortalThemeWithResponse(ctx, dto.Identity())
	return errors.IgnoreNotFound(errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}))
}

// UpdateStatus copies the ids of the upsert state into the status, field by field so that the rest of
// the status is kept. Every write is HRID-based, so the theme is Automation API managed.
func UpdateStatus(_ context.Context, obj *v1alpha1.PortalTheme, resp Response) error {
	obj.Status.ID = utils.SafeDereference(resp.State.Id)
	obj.Status.OrgID = resp.GetOrgID()
	obj.Status.EnvID = resp.GetEnvID()
	k8s.AddAutomationAPIManagedCondition(obj)
	return nil
}
