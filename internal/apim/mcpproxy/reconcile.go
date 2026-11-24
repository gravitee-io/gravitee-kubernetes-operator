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

package mcpproxy

import (
	"context"
	"fmt"
	"net/http"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/mcpproxy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/status"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

// Response is what Upsert returns: the state the platform answers.
type Response struct {
	State sdk.AimMcpProxyState
}

// GetOrgID is the organization the state reports, as the status has always shown it.
func (r Response) GetOrgID() string {
	return valueOf(r.State.OrganizationId)
}

// GetEnvID is the environment the state reports, as the status has always shown it.
func (r Response) GetEnvID() string {
	return valueOf(r.State.EnvironmentId)
}

// Upsert creates or updates the proxy. The platform answers the state: ids, findings, the lifecycle it
// observes and, for a STUDIO, the tools with their references. An apply refused by a finding answers a 400
// whose body carries the findings, see model.AutomationRefusal.
func Upsert(ctx context.Context, client *apim.SDKClient, dto sdk.AimMcpProxy) (Response, error) {
	state, err := upsert(ctx, client, dto, false)
	if err != nil {
		return Response{}, err
	}
	return Response{State: state}, nil
}

func upsert(ctx context.Context, client *apim.SDKClient, dto sdk.AimMcpProxy, dryRun bool) (
	sdk.AimMcpProxyState, error,
) {
	resp, err := client.UpsertAimMcpProxyWithResponse(ctx, &sdk.UpsertAimMcpProxyParams{DryRun: &dryRun}, dto)
	if err = errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}); err != nil {
		return sdk.AimMcpProxyState{}, err
	}
	if resp.JSON200 == nil {
		return sdk.AimMcpProxyState{}, errors.UnexpectedSDKResponse(resp.HTTPResponse)
	}
	return *resp.JSON200, nil
}

// ResolveRefs resolves the management context and, outside a delete, waits for the CatalogMcpServers a studio
// selects (ResolvedRefs=False, requeued): the platform refuses a proxy whose servers are not in the catalog yet.
// Admission resolves the context only, and defers its dry run instead of refusing.
func ResolveRefs(ctx context.Context, obj *v1alpha1.McpProxy, ns string) error {
	if err := apim.ResolveContextRef(ctx, obj, ns); err != nil {
		return err
	}
	if obj.IsBeingDeleted() {
		return nil
	}
	return dynamic.AssertCatalogMcpServersSynced(ctx, obj)
}

// OnUpsertError records the findings of a refused apply in the status. Admission ran the same checks as a dry
// run, so a refusal here comes from a state that changed since (a tool gone upstream, a context path taken): it
// stays recoverable (backoff) rather than terminal until the next spec change. Any other failure keeps the
// lifecycle's classification.
func OnUpsertError(_ context.Context, obj *v1alpha1.McpProxy, err error) error {
	if !errors.IsBadRequest(err) {
		return nil
	}
	if findings, refused := model.AutomationRefusal(err); refused {
		obj.Status.Errors = findings
	}
	return fmt.Errorf("mcp proxy [%s] was refused by the automation api: %w", obj.GetName(), err)
}

// Delete removes the proxy. A proxy already gone counts as deleted.
func Delete(ctx context.Context, client *apim.SDKClient, dto sdk.AimMcpProxy) error {
	resp, err := client.DeleteAimMcpProxyWithResponse(ctx, dto.Identity())
	return errors.IgnoreNotFound(errors.FromSDKResponse(err, func() (*http.Response, []byte) {
		return resp.HTTPResponse, resp.Body
	}))
}

// UpdateStatus copies the upsert state into the CR status, field by field so that the conditions are kept,
// and marks the proxy Automation API managed.
func UpdateStatus(_ context.Context, obj *v1alpha1.McpProxy, resp Response) error {
	obj.Status.ID = valueOf(resp.State.Id)
	obj.Status.OrgID = resp.GetOrgID()
	obj.Status.EnvID = resp.GetEnvID()
	obj.Status.HRID = refs.NewNamespacedNameFromObject(obj).HRID()
	obj.Status.Errors = toStatusErrors(resp.State.Errors)
	obj.Status.State = ""
	if resp.State.State != nil {
		obj.Status.State = string(*resp.State.State)
	}
	obj.Status.Tools = ToolsStatus(resp.State)

	k8s.AddAutomationAPIManagedCondition(obj)
	return nil
}

// ToolsStatus returns the tools a STUDIO exposes, for the CRD status. It reads the state's tools list, which
// carries each tool's authPolicyRef. Only a state without that list, as answered by a platform older than
// AIM 5.1, falls back to studio.tools, whose entries have no reference. An empty list is an answer: it does
// not fall back.
func ToolsStatus(state sdk.AimMcpProxyState) []mcpproxy.StudioToolStatus {
	if state.Tools != nil {
		if len(state.Tools) == 0 {
			return nil
		}
		tools := make([]mcpproxy.StudioToolStatus, 0, len(state.Tools))
		for _, tool := range state.Tools {
			tools = append(tools, mcpproxy.StudioToolStatus{
				Server:        valueOf(tool.Server),
				Tool:          valueOf(tool.Tool),
				Alias:         tool.Alias,
				EntityID:      valueOf(tool.EntityId),
				AuthPolicyRef: valueOf(tool.AuthPolicyRef),
			})
		}
		return tools
	}
	if state.Studio == nil || len(state.Studio.Tools) == 0 {
		return nil
	}
	tools := make([]mcpproxy.StudioToolStatus, 0, len(state.Studio.Tools))
	for _, tool := range state.Studio.Tools {
		tools = append(tools, mcpproxy.StudioToolStatus{
			Server:   tool.Server,
			Tool:     tool.Tool,
			Alias:    tool.Alias,
			EntityID: valueOf(tool.EntityId),
		})
	}
	return tools
}

func toStatusErrors(errs *sdk.AimErrors) status.Errors {
	if errs == nil {
		return status.Errors{}
	}
	return status.Errors{Severe: errs.Severe, Warning: errs.Warning}
}

func valueOf(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
