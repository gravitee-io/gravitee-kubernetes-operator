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

// Package portaltheme syncs PortalTheme resources with APIM through the APIM SDK.
package portaltheme

import (
	"context"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/mapper"
)

// ToDTO maps a PortalTheme CRD to the Automation API payload. SDK defaults are not applied.
func ToDTO(obj *v1alpha1.PortalTheme) (sdk.PortalTheme, error) {
	dto, err := mapper.MapViaJSON[sdk.PortalTheme](obj.Spec.Type)
	if err != nil {
		return sdk.PortalTheme{}, err
	}
	dto.Hrid = refs.NewNamespacedNameFromObject(obj).HRID()
	return dto, nil
}

// CreateAPIMClient builds the SDK client of the management context the theme references.
func CreateAPIMClient(ctx context.Context, obj *v1alpha1.PortalTheme) (*apim.SDKClient, error) {
	return apim.SDKFromContextRef(ctx, obj.ContextRef(), obj.GetNamespace())
}

// Response is what Upsert returns: the theme state APIM answers, and the org and env it reports.
type Response struct {
	State sdk.PortalThemeState
	OrgID string
	EnvID string
}

func newResponse(state sdk.PortalThemeState) Response {
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
