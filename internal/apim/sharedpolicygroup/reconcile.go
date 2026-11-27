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

package sharedpolicygroup

import (
	"context"
	"net/http"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/status"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
)

// Upsert creates or updates the shared policy group in APIM.
func Upsert(ctx context.Context, client *apim.SDKClient, dto DTO) (Response, error) {
	state, err := upsert(ctx, client, dto, false)
	if err != nil {
		return Response{}, err
	}
	return Response{
		State:            state,
		OrgID:            client.GetOrgID(),
		EnvID:            client.GetEnvID(),
		HridContainsUUID: dto.HridContainsUUID,
	}, nil
}

func upsert(ctx context.Context, client *apim.SDKClient, dto DTO, dryRun bool) (sdk.SharedPolicyGroupState, error) {
	params := &sdk.UpsertSharedPolicyGroupParams{DryRun: &dryRun}
	if dto.HridContainsUUID {
		params.HridContainsUUID = new(true)
	}

	resp, err := client.UpsertSharedPolicyGroupWithResponse(ctx, params, dto.SharedPolicyGroup)
	if err = errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return sdk.SharedPolicyGroupState{}, err
	}
	if resp.JSON200 == nil {
		return sdk.SharedPolicyGroupState{}, errors.UnexpectedSDKResponse(resp.HTTPResponse)
	}
	return *resp.JSON200, nil
}

// Delete removes the shared policy group from APIM. A group already gone counts as deleted.
func Delete(ctx context.Context, client *apim.SDKClient, dto DTO) error {
	resp, err := client.DeleteSharedPolicyGroupWithResponse(ctx, dto.Identity(),
		&sdk.DeleteSharedPolicyGroupParams{HridContainsUUID: &dto.HridContainsUUID})
	return errors.IgnoreNotFound(errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}))
}

// DeleteGuard refuses to delete a shared policy group a v4 API still references.
func DeleteGuard(ctx context.Context, obj *v1alpha1.SharedPolicyGroup) error {
	return search.AssertNoSharedPolicyGroupRef(ctx, obj)
}

// UpdateStatus copies the upsert response into the CR status. v4 APIs reference a legacy group by its cross id.
// A group upserted by its hrid becomes Automation API managed; a legacy one, upserted by its cross id, stays as
// it is.
func UpdateStatus(_ context.Context, obj *v1alpha1.SharedPolicyGroup, resp Response) error {
	obj.Status.ID = valueOr(resp.State.Id)
	obj.Status.CrossID = valueOr(resp.State.CrossId)
	obj.Status.OrgID = resp.GetOrgID()
	obj.Status.EnvID = resp.GetEnvID()
	obj.Status.Errors = status.Errors{}
	if resp.State.Errors != nil {
		obj.Status.Errors = status.Errors{Severe: resp.State.Errors.Severe, Warning: resp.State.Errors.Warning}
	}

	if !resp.HridContainsUUID {
		k8s.AddAutomationAPIManagedCondition(obj)
	}
	return nil
}

// SetProcessingStatus keeps the deprecated processingStatus: Completed after a successful reconcile, Failed
// otherwise.
func SetProcessingStatus(obj *v1alpha1.SharedPolicyGroup, err error) {
	if err != nil {
		obj.Status.ProcessingStatus = core.ProcessingStatusFailed
		return
	}
	obj.Status.ProcessingStatus = core.ProcessingStatusCompleted
}

func valueOr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
