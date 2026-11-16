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

package group

import (
	"context"
	"net/http"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

// Upsert creates or updates the group in APIM.
func Upsert(ctx context.Context, client *apim.SDKClient, dto DTO) (Response, error) {
	state, err := upsert(ctx, client, dto, false)
	if err != nil {
		return Response{}, err
	}
	return Response{State: state, OrgID: client.GetOrgID(), EnvID: client.GetEnvID()}, nil
}

func upsert(ctx context.Context, client *apim.SDKClient, dto DTO, dryRun bool) (sdk.GroupState, error) {
	params := &sdk.UpsertGroupParams{DryRun: &dryRun}
	if dto.HridContainsUUID {
		params.HridContainsUUID = new(true)
	}
	if dto.IgnoreMembers {
		params.IgnoreMembers = new(true)
	}

	resp, err := client.UpsertGroupWithResponse(ctx, params, dto.Group)
	if err = errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return sdk.GroupState{}, err
	}
	if resp.JSON200 == nil {
		return sdk.GroupState{}, errors.UnexpectedSDKResponse(resp.HTTPResponse)
	}
	return *resp.JSON200, nil
}

// Delete removes the group from APIM. A group already gone is an error.
func Delete(ctx context.Context, client *apim.SDKClient, dto DTO) error {
	resp, err := client.DeleteGroupWithResponse(ctx, dto.Identity(),
		&sdk.DeleteGroupParams{HridContainsUUID: &dto.HridContainsUUID})
	return errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	})
}

// UpdateStatus copies the upsert response into the CR status. A group addressed by its hrid becomes
// Automation API managed; a legacy group, addressed by its UUID, stays as it is.
func UpdateStatus(_ context.Context, obj *v1alpha1.Group, resp Response) error {
	legacy := obj.Status.ID != "" && !k8s.IsAutomationAPIManaged(obj)

	obj.Status.ID = ""
	if resp.State.Id != nil {
		obj.Status.ID = *resp.State.Id
	}
	obj.Status.Members = 0
	if resp.State.MemberCount != nil && *resp.State.MemberCount > 0 {
		obj.Status.Members = uint(*resp.State.MemberCount)
	}
	obj.Status.OrgID = resp.GetOrgID()
	obj.Status.EnvID = resp.GetEnvID()

	if !legacy {
		k8s.AddAutomationAPIManagedCondition(obj)
	}
	return nil
}
