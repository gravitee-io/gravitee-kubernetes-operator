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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	apimmcpproxy "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/mcpproxy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

func CreateOrUpdate(ctx context.Context, proxy *v1alpha1.McpProxy) error {
	if !proxy.HasContext() {
		return gerrors.NewIllegalStateError(
			fmt.Errorf("mcp proxy [%s] has no management context", proxy.GetName()),
		)
	}

	client, err := apimmcpproxy.CreateAPIMClient(ctx, proxy)
	if err != nil {
		return err
	}

	dto, err := apimmcpproxy.ToDTO(proxy)
	if err != nil {
		return err
	}

	resp, err := apimmcpproxy.Upsert(ctx, client, dto)
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

	return apimmcpproxy.UpdateStatus(ctx, proxy, resp)
}
