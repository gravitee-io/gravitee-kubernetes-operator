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

package apidefinition

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
)

// AdoptLegacyID re-attaches an API that GKO 4.11 created to its CR when the CR comes back without status,
// for example after a delete and re-apply during an upgrade. GKO 4.11 derived the API id from the cross id,
// organization and environment; APIM only keeps that id when the API is written with hridContainsUUID.
// When APIM holds an API with that id and this CR's cross id, the status gets the id, so that PopulateIDs
// takes the legacy branch. Otherwise the API is new and is synced by HRID, as before.
func AdoptLegacyID(ctx context.Context, apimClient *apim.APIM, api *v1alpha1.ApiV4Definition) error {
	if api.Status.ID != "" || api.Spec.ID != "" || k8s.IsAutomationAPIManaged(api) {
		return nil
	}

	id := api.LegacyID(apimClient.Context)
	remote, err := apimClient.APIs.GetV4WithUUID(id)
	if err != nil {
		var serverErr gerrors.ServerError
		if errors.As(err, &serverErr) && serverErr.StatusCode < http.StatusInternalServerError {
			return nil
		}
		return fmt.Errorf("looking up the API GKO 4.11 created with id %s: %w", id, err)
	}

	if remote.V4BaseApi == nil || remote.ApiBase == nil || remote.CrossID != api.LegacyCrossID() {
		return nil
	}

	api.Status.ID = id
	api.Status.CrossID = remote.CrossID
	log.Info(ctx, "Re-attached API created before GKO 4.12", log.KeyValues(api, "id", id)...)
	return nil
}
