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

package application

import (
	"context"
	"net/http"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

// Upsert creates or updates the application in APIM.
func Upsert(ctx context.Context, client *apim.SDKClient, dto DTO) (Response, error) {
	state, err := upsert(ctx, client, dto, false)
	if err != nil {
		return Response{}, err
	}
	return Response{
		State:  state,
		OrgID:  client.GetOrgID(),
		EnvID:  client.GetEnvID(),
		Legacy: dto.Legacy,
	}, nil
}

func upsert(ctx context.Context, client *apim.SDKClient, dto DTO, dryRun bool) (sdk.ApplicationState, error) {
	params := &sdk.UpsertApplicationParams{DryRun: &dryRun}
	if dto.IgnoreGroups {
		params.IgnoreGroups = new(true)
	}

	resp, err := client.UpsertApplicationWithResponse(ctx, params, dto.Application, dto.Legacy.editors()...)
	if err = errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return sdk.ApplicationState{}, err
	}
	if resp.JSON200 == nil {
		return sdk.ApplicationState{}, errors.UnexpectedSDKResponse(resp.HTTPResponse)
	}
	return *resp.JSON200, nil
}

// Delete removes the application from APIM. An application already gone counts as deleted.
func Delete(ctx context.Context, client *apim.SDKClient, dto DTO) error {
	resp, err := client.DeleteApplicationWithResponse(ctx, dto.Identity(), dto.Legacy.editors()...)
	return errors.IgnoreNotFound(errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}))
}

// UpdateStatus copies the upsert response into the CR status. An application upserted by its hrid becomes
// Automation API managed; a legacy application, upserted by its UUID, stays as it is. The subscription count
// is left alone: the subscription controller owns it. A response without an id keeps the one the status has.
// The deprecated processing status is cleared: the lifecycle reports through the conditions, and a value left by
// the controller before it would never change again.
func UpdateStatus(_ context.Context, obj *v1alpha1.Application, resp Response) error {
	obj.Status.ProcessingStatus = ""
	if resp.State.Id != nil {
		obj.Status.ID = *resp.State.Id
	}
	obj.Status.OrgID = resp.GetOrgID()
	obj.Status.EnvID = resp.GetEnvID()

	obj.Status.Errors.Severe = nil
	obj.Status.Errors.Warning = nil
	if resp.State.Errors != nil {
		obj.Status.Errors.Severe = resp.State.Errors.Severe
		obj.Status.Errors.Warning = resp.State.Errors.Warning
	}

	if !resp.Legacy.HridContainsUUID {
		k8s.AddAutomationAPIManagedCondition(obj)
	}
	return nil
}
