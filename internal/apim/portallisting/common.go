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

// Package portallisting syncs PortalListing resources with APIM through the APIM SDK.
package portallisting

import (
	"context"
	"fmt"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

// DTO is the listing payload with the HRID of the portal it belongs to, which goes in the request path.
// The portal HRID is not sent and drift detection ignores it.
type DTO struct {
	sdk.PortalListing
	PortalHrid string `json:"-" drift:"ignore"`
}

// ToDTO maps a PortalListing CRD to the Automation API payload, entry by entry. An API ref without
// namespace points to the listing's namespace, as resolving it would. The APIs are not resolved here.
func ToDTO(obj *v1alpha1.PortalListing) (DTO, error) {
	apis := make([]sdk.PortalListingApiEntry, 0, len(obj.Spec.APIs))
	for i := range obj.Spec.APIs {
		entry := &obj.Spec.APIs[i]
		apis = append(apis, sdk.PortalListingApiEntry{
			ApiHrid:    hridOf(entry.Ref, obj.GetNamespace()),
			Location:   entry.Location,
			Order:      toIntPtr(entry.Order),
			Visibility: toVisibility(entry.Visibility),
		})
	}
	return DTO{
		PortalListing: sdk.PortalListing{
			Hrid: refs.NewNamespacedNameFromObject(obj).HRID(),
			Apis: apis,
		},
		PortalHrid: hridOf(obj.Spec.Portal, obj.GetNamespace()),
	}, nil
}

func hridOf(ref refs.NamespacedName, ns string) string {
	if ref.IsMissingNamespace() {
		ref.SetNamespace(ns)
	}
	return ref.HRID()
}

func toIntPtr(v *int32) *int {
	if v == nil {
		return nil
	}
	i := int(*v)
	return &i
}

func toVisibility[V ~string](v *V) *sdk.PortalVisibility {
	if v == nil || *v == "" {
		return nil
	}
	visibility := sdk.PortalVisibility(*v)
	return &visibility
}

// CreateAPIMClient builds the SDK client of the management context of the portal the listing belongs to.
// A portal without a context is refused.
func CreateAPIMClient(ctx context.Context, obj *v1alpha1.PortalListing) (*apim.SDKClient, error) {
	prtl, err := dynamic.ResolvePortal(ctx, obj.GetPortalRef(), obj.GetNamespace())
	if err != nil {
		return nil, err
	}
	if !prtl.HasContext() {
		return nil, fmt.Errorf("portal [%s] has no management context", prtl.GetName())
	}
	return apim.SDKFromContextRef(ctx, prtl.ContextRef(), prtl.GetNamespace())
}

// Response is what Upsert returns: the listing state APIM answers, and the org and env it reports.
type Response struct {
	State sdk.PortalListingState
	OrgID string
	EnvID string
}

func newResponse(state sdk.PortalListingState) Response {
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
