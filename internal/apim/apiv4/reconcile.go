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

package apiv4

import (
	"context"
	goerrors "errors"
	"net/http"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/status"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/apidefinition"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

// Upsert sends the API to APIM when it has a management context, then writes or deletes its gateway ConfigMap.
func Upsert(ctx context.Context, client *Client, dto DTO) (Response, error) {
	resp := Response{HridContainsUUID: dto.HridContainsUUID}
	if client != nil {
		state, err := put(ctx, client, dto, false)
		if err != nil {
			return Response{}, err
		}
		resp.State = state
		resp.OrgID = valueOr(state.OrganizationId, client.GetOrgID())
		resp.EnvID = valueOr(state.EnvironmentId, client.GetEnvID())
	}
	if err := syncGateway(ctx, dto.Gateway, resp); err != nil {
		return Response{}, &GatewaySyncError{Response: resp, Err: err}
	}
	return resp, nil
}

// GatewaySyncError is a failed gateway ConfigMap write after APIM accepted the API. Response holds what APIM
// returned, so the caller can still copy it into the status.
type GatewaySyncError struct {
	Response Response
	Err      error
}

func (e *GatewaySyncError) Error() string {
	return e.Err.Error()
}

func (e *GatewaySyncError) Unwrap() error {
	return e.Err
}

// AcceptedResponse returns what APIM returned when err is a GatewaySyncError: APIM accepted the API, only its
// gateway ConfigMap failed.
func AcceptedResponse(err error) (Response, bool) {
	var syncErr *GatewaySyncError
	if goerrors.As(err, &syncErr) {
		return syncErr.Response, true
	}
	return Response{}, false
}

// put sends dto.ApiV4 through the SDK, or through GKO's HTTP client for a legacy API (the SDK cannot send
// hridContainsUUID on PUT /apis yet, see APIs.UpsertV4Legacy). Both send the same body and decode the same state.
func put(ctx context.Context, client *Client, dto DTO, dryRun bool) (*sdk.ApiV4State, error) {
	if dto.HridContainsUUID {
		return client.APIM.APIs.UpsertV4Legacy(dto.ApiV4, dryRun, dto.IgnoreGroups)
	}
	params := &sdk.UpsertApiParams{DryRun: &dryRun}
	if dto.IgnoreGroups {
		params.IgnoreGroups = new(true)
	}
	resp, err := client.UpsertApiWithResponse(ctx, params, dto.ApiV4)
	if err = errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, errors.UnexpectedSDKResponse(resp.HTTPResponse)
	}
	return resp.JSON200, nil
}

// syncGateway writes the gateway ConfigMap of an API that syncs from Kubernetes (deletes it while the API is
// stopped), and deletes the one an API that syncs from management still owns.
func syncGateway(ctx context.Context, gateway GatewaySync, resp Response) error {
	if !gateway.Enabled {
		return apidefinition.DeleteOwnedConfigMap(ctx, gateway.Name, gateway.Namespace, gateway.Owner.UID)
	}
	if gateway.Stopped {
		return apidefinition.DeleteConfigMap(ctx, gateway.Name, gateway.Namespace)
	}
	return apidefinition.SaveConfigMap(ctx, apidefinition.GatewayConfigMap{
		Name:              gateway.Name,
		Namespace:         gateway.Namespace,
		Owner:             gateway.Owner,
		OrgID:             resp.OrgID,
		EnvID:             resp.EnvID,
		DefinitionVersion: gatewayDefinitionVersion,
		Definition:        gateway.Definition,
	})
}

// Delete deletes the API in APIM. An API without a management context has nothing to delete.
func Delete(ctx context.Context, client *Client, dto DTO) error {
	if client == nil {
		return nil
	}
	resp, err := client.DeleteApiWithResponse(ctx, dto.DeleteID,
		&sdk.DeleteApiParams{HridContainsUUID: &dto.DeleteHridContainsUUID})
	return errors.IgnoreNotFound(errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}))
}

// UpdateStatus copies the APIM state into the status. The subscription count is not APIM's and is kept.
func UpdateStatus(_ context.Context, obj *v1alpha1.ApiV4Definition, resp Response) error {
	obj.Status.ProcessingStatus = ""
	if resp.State == nil {
		return nil
	}

	state := resp.State
	obj.Status.ID = valueOr(state.Id, "")
	obj.Status.CrossID = valueOr(state.CrossId, "")
	obj.Status.OrgID = resp.GetOrgID()
	obj.Status.EnvID = resp.GetEnvID()
	obj.Status.State = ""
	if state.State != nil {
		obj.Status.State = base.ApiState(*state.State)
	}
	obj.Status.Errors = status.Errors{}
	if state.Errors != nil {
		obj.Status.Errors = status.Errors{Severe: state.Errors.Severe, Warning: state.Errors.Warning}
	}
	obj.Status.Plans = plansAsMap(state.Plans)

	if !resp.HridContainsUUID {
		k8s.AddAutomationAPIManagedCondition(obj)
	}
	return nil
}

func plansAsMap(plans []sdk.PlanV4) map[string]string {
	result := make(map[string]string, len(plans))
	for _, plan := range plans {
		result[plan.Hrid] = valueOr(plan.Id, "")
	}
	return result
}

func valueOr(value *string, fallback string) string {
	if value == nil || *value == "" {
		return fallback
	}
	return *value
}
