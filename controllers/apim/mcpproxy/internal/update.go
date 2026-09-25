// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"context"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/mcpproxy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

func CreateOrUpdate(ctx context.Context, proxy *v1alpha1.McpProxy) error {
	if !proxy.HasContext() {
		return gerrors.NewIllegalStateError(
			fmt.Errorf("mcp proxy [%s] has no management context", proxy.GetName()),
		)
	}

	apimClient, err := apim.FromContextRef(ctx, proxy.ContextRef(), proxy.GetNamespace())
	if err != nil {
		return err
	}

	state, err := apimClient.McpProxies.CreateOrUpdate(proxy)
	if err != nil {
		// A 400 is how the platform refuses an apply. Admission ran the same checks as a dry run,
		// so a refusal here comes from a state that changed since (a tool gone upstream, a context
		// path taken): keep it recoverable (backoff) rather than terminal until the next spec change.
		if gerrors.IsBadRequest(err) {
			if findings, refused := model.AutomationRefusal(err); refused {
				proxy.Status.Errors = findings
			}
			return fmt.Errorf("mcp proxy [%s] was refused by the automation api: %w", proxy.GetName(), err)
		}
		return gerrors.NewControlPlaneError(err)
	}

	// Setting fields by fields to keep the rest (conditions) intact
	proxy.Status.ID = state.ID
	proxy.Status.OrgID = state.OrgID
	proxy.Status.EnvID = state.EnvID
	proxy.Status.HRID = refs.NewNamespacedNameFromObject(proxy).HRID()
	proxy.Status.Errors = state.Errors
	proxy.Status.State = state.State
	proxy.Status.Tools = toolsStatus(state.Studio)

	return nil
}

func toolsStatus(studio *model.McpProxyStudioDTO) []mcpproxy.StudioToolStatus {
	if studio == nil || len(studio.Tools) == 0 {
		return nil
	}
	tools := make([]mcpproxy.StudioToolStatus, 0, len(studio.Tools))
	for _, tool := range studio.Tools {
		tools = append(tools, mcpproxy.StudioToolStatus{
			Server:   tool.Server,
			Tool:     tool.Tool,
			Alias:    tool.Alias,
			EntityID: tool.EntityID,
		})
	}
	return tools
}
