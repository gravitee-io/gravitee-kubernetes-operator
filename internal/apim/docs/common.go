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

// Package docs syncs Documentation resources with APIM through the APIM SDK.
package docs

import (
	"context"
	"fmt"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/navigation"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/mapper"
)

// DTO is the documentation page payload with the portal or API it belongs to, which goes in the request
// path. The parent is not sent and drift detection ignores it.
type DTO struct {
	sdk.Documentation
	Parent navigation.Parent `json:"-" drift:"ignore"`
}

// ToDTO maps a Documentation CRD to the Automation API payload. SDK defaults are not applied. A page that
// references both a portal and an API, or neither, is an error: no partial payload is built.
func ToDTO(obj *v1alpha1.Documentation) (DTO, error) {
	parent, err := navigation.ParentFromRefs(obj.Spec.Portal, obj.Spec.API, obj.GetNamespace())
	if err != nil {
		return DTO{}, fmt.Errorf("documentation [%s/%s]: %w", obj.GetNamespace(), obj.GetName(), err)
	}
	payload, err := mapper.MapViaJSON[sdk.Documentation](obj.Spec.Type)
	if err != nil {
		return DTO{}, err
	}
	payload.Hrid = refs.NewNamespacedNameFromObject(obj).HRID()
	return DTO{Documentation: payload, Parent: parent}, nil
}

// CreateAPIMClient builds the SDK client of the management context of the portal or API the page belongs
// to. Without a context, or while its portal is not synced, nothing was synced: the
// client is nil, Delete and DryRun then do nothing, and ResolveRefs refuses an upsert first.
func CreateAPIMClient(ctx context.Context, obj *v1alpha1.Documentation) (*apim.SDKClient, error) {
	parent, err := navigation.ResolveParent(ctx, obj.Spec.Portal, obj.Spec.API, obj.GetNamespace())
	if err != nil {
		return nil, err
	}
	if !parent.HasContext() || navigation.CheckPortalSynced(parent) != nil {
		return nil, nil
	}
	return apim.SDKFromContextRef(ctx, parent.ContextRef(), parent.GetNamespace())
}

// Response is what Upsert returns: the page state APIM answers, and the org and env it reports.
type Response struct {
	State sdk.DocumentationState
	OrgID string
	EnvID string
}

func newResponse(state sdk.DocumentationState) Response {
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
