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

// Package portal syncs Portal resources with APIM through the APIM SDK.
package portal

import (
	"context"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/navigation"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/mapper"
)

// ToDTO maps a Portal CRD to the Automation API payload. SDK defaults are not applied. The active theme is
// the HRID of the referenced theme, nil without one: APIM then deactivates the portal's theme. The theme is
// not read here; reconcile refuses a portal whose theme it cannot resolve.
func ToDTO(obj *v1alpha1.Portal) (sdk.Portal, error) {
	dto, err := mapper.MapViaJSON[sdk.Portal](obj.Spec.Type)
	if err != nil {
		return sdk.Portal{}, err
	}
	dto.Hrid = refs.NewNamespacedNameFromObject(obj).HRID()
	dto.ActiveThemeHrid = nil
	if hrid := obj.ActiveThemeHRID(); hrid != "" {
		dto.ActiveThemeHrid = &hrid
	}
	return dto, nil
}

// DriftDTO returns a copy of the payload in which every top navbar entry carries the visibility APIM
// stores for it, so drift compares the field exactly. It is applied to both sides of the comparison and
// never to the payload sent: see navigation.WithResolvedTopNavbar. The input is never written to.
func DriftDTO(dto sdk.Portal) sdk.Portal {
	if dto.Structure == nil || len(dto.Structure.TopNavbar) == 0 {
		return dto
	}
	dto.Structure = &sdk.PortalNavigationStructure{TopNavbar: navigation.WithResolvedTopNavbar(dto.Structure.TopNavbar)}
	return dto
}

// CreateAPIMClient builds the SDK client of the management context the portal references.
// Without a context nothing was synced: the client is nil, Delete and DryRun then do nothing, and
// ResolveRefs refuses an upsert first.
func CreateAPIMClient(ctx context.Context, obj *v1alpha1.Portal) (*apim.SDKClient, error) {
	if !obj.HasContext() {
		return nil, nil
	}
	return apim.SDKFromContextRef(ctx, obj.ContextRef(), obj.GetNamespace())
}

// Response is what Upsert returns: the portal state APIM answers, and the org and env it reports.
type Response struct {
	State sdk.PortalState
	OrgID string
	EnvID string
}

func newResponse(state sdk.PortalState) Response {
	return Response{
		State: state,
		OrgID: utils.SafeDereference(state.OrganizationId),
		EnvID: utils.SafeDereference(state.EnvironmentId),
	}
}

func (r Response) GetOrgID() string {
	return r.OrgID
}

func (r Response) GetEnvID() string {
	return r.EnvID
}
