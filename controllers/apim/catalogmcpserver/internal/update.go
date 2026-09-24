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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

func CreateOrUpdate(ctx context.Context, srv *v1alpha1.CatalogMcpServer) error {
	if !srv.HasContext() {
		return gerrors.NewIllegalStateError(
			fmt.Errorf("catalog mcp server [%s] has no management context", srv.GetName()),
		)
	}

	apimClient, err := apim.FromContextRef(ctx, srv.ContextRef(), srv.GetNamespace())
	if err != nil {
		return err
	}

	state, err := apimClient.CatalogMcpServers.CreateOrUpdate(srv)
	if err != nil {
		// A 400 is how the platform refuses an apply, a failed discovery of the upstream included,
		// which is often transient: keep it recoverable (backoff) rather than terminal until the
		// next spec change.
		if gerrors.IsBadRequest(err) {
			if findings, refused := model.CatalogMcpServerRefusal(err); refused {
				srv.Status.Errors = findings
			}
			return fmt.Errorf("catalog mcp server [%s] was refused by the automation api: %w", srv.GetName(), err)
		}
		return gerrors.NewControlPlaneError(err)
	}

	// Setting fields by fields to keep the rest (conditions) intact
	srv.Status.ID = state.ID
	srv.Status.OrgID = state.OrgID
	srv.Status.EnvID = state.EnvID
	srv.Status.HRID = refs.NewNamespacedNameFromObject(srv).HRID()
	srv.Status.Errors = state.Errors
	srv.Status.Discovered = state.Discovered

	return nil
}
