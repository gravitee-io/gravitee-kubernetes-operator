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
	"fmt"
	"net/http"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/catalogmcpserver"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/status"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

// Response is what Upsert returns: the state the platform answers.
type Response struct {
	State sdk.AimCatalogMcpServerState
}

// GetOrgID is the organization the state reports, as the status has always shown it.
func (r Response) GetOrgID() string {
	return valueOf(r.State.OrganizationId)
}

// GetEnvID is the environment the state reports, as the status has always shown it.
func (r Response) GetEnvID() string {
	return valueOf(r.State.EnvironmentId)
}

// Upsert creates or updates the server in the catalog. The platform discovers the upstream before it
// persists anything and answers the state: ids, findings and what it discovered. A refusal is a 400
// whose body carries the findings, see model.AutomationRefusal.
func Upsert(ctx context.Context, client *apim.SDKClient, dto sdk.AimCatalogMcpServer) (Response, error) {
	state, err := upsert(ctx, client, dto, false)
	if err != nil {
		return Response{}, err
	}
	return Response{State: state}, nil
}

func upsert(ctx context.Context, client *apim.SDKClient, dto sdk.AimCatalogMcpServer, dryRun bool) (
	sdk.AimCatalogMcpServerState, error,
) {
	resp, err := client.UpsertAimCatalogMcpServerWithResponse(ctx,
		&sdk.UpsertAimCatalogMcpServerParams{DryRun: &dryRun}, dto)
	if err = errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return sdk.AimCatalogMcpServerState{}, err
	}
	if resp.JSON200 == nil {
		return sdk.AimCatalogMcpServerState{}, errors.UnexpectedSDKResponse(resp.HTTPResponse)
	}
	return *resp.JSON200, nil
}

// OnUpsertError records the findings of a refused apply in the status. A refusal, a failed discovery of the
// upstream included, is often transient: it stays recoverable (backoff) rather than terminal until the next
// spec change. Any other failure keeps the lifecycle's classification.
func OnUpsertError(_ context.Context, obj *v1alpha1.CatalogMcpServer, err error) error {
	if !errors.IsBadRequest(err) {
		return nil
	}
	if findings, refused := model.AutomationRefusal(err); refused {
		obj.Status.Errors = findings
	}
	return fmt.Errorf("catalog mcp server [%s] was refused by the automation api: %w", obj.GetName(), err)
}

// Delete removes the server from the catalog. A server already gone counts as deleted.
func Delete(ctx context.Context, client *apim.SDKClient, dto sdk.AimCatalogMcpServer) error {
	resp, err := client.DeleteAimCatalogMcpServerWithResponse(ctx, dto.Identity())
	return errors.IgnoreNotFound(errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}))
}

// UpdateStatus copies the upsert state into the CR status, field by field so that the conditions are
// kept, and marks the server Automation API managed.
func UpdateStatus(_ context.Context, obj *v1alpha1.CatalogMcpServer, resp Response) error {
	obj.Status.ID = valueOf(resp.State.Id)
	obj.Status.OrgID = resp.GetOrgID()
	obj.Status.EnvID = resp.GetEnvID()
	obj.Status.HRID = refs.NewNamespacedNameFromObject(obj).HRID()
	obj.Status.Errors = ToStatusErrors(resp.State.Errors)
	obj.Status.AuthPolicyRef = valueOf(resp.State.AuthPolicyRef)
	obj.Status.Discovered = toDiscovered(resp.State)

	k8s.AddAutomationAPIManagedCondition(obj)
	return nil
}

// ToStatusErrors maps the findings of a state onto the CR status.
func ToStatusErrors(errs *sdk.AimErrors) status.Errors {
	if errs == nil {
		return status.Errors{}
	}
	return status.Errors{Severe: errs.Severe, Warning: errs.Warning}
}

// toDiscovered maps what the platform discovered on the upstream. A list the state leaves out stays nil.
func toDiscovered(state sdk.AimCatalogMcpServerState) catalogmcpserver.Discovered {
	discovered := catalogmcpserver.Discovered{
		LastSyncedAt:    valueOf(state.LastSyncedAt),
		ProtocolVersion: valueOf(state.ProtocolVersion),
	}
	if state.ServerInfo != nil {
		discovered.ServerInfo = &catalogmcpserver.ServerInfo{
			Name:    valueOf(state.ServerInfo.Name),
			Version: valueOf(state.ServerInfo.Version),
		}
	}
	if state.Tools != nil {
		discovered.Tools = make([]catalogmcpserver.Tool, 0, len(state.Tools))
		for _, tool := range state.Tools {
			discovered.Tools = append(discovered.Tools, catalogmcpserver.Tool{
				Name:          valueOf(tool.Name),
				EntityID:      valueOf(tool.EntityId),
				AuthPolicyRef: valueOf(tool.AuthPolicyRef),
			})
		}
	}
	if state.Prompts != nil {
		discovered.Prompts = make([]catalogmcpserver.Prompt, 0, len(state.Prompts))
		for _, prompt := range state.Prompts {
			discovered.Prompts = append(discovered.Prompts, catalogmcpserver.Prompt{
				Name:     valueOf(prompt.Name),
				EntityID: valueOf(prompt.EntityId),
			})
		}
	}
	if state.Resources != nil {
		discovered.Resources = make([]catalogmcpserver.Resource, 0, len(state.Resources))
		for _, resource := range state.Resources {
			discovered.Resources = append(discovered.Resources, catalogmcpserver.Resource{
				URI:      valueOf(resource.Uri),
				Name:     valueOf(resource.Name),
				EntityID: valueOf(resource.EntityId),
			})
		}
	}
	return discovered
}

func valueOf(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
