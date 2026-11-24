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

package apim

import (
	"context"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	apimdocs "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/docs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/navigation"
	apimportal "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portal"
	apimportallink "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portallink"
	apimportallisting "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portallisting"
	apimportaltheme "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portaltheme"
)

// GetPortalTheme gets a portal theme state by HRID, server-managed fields included.
func (a *APIM) GetPortalTheme(ctx context.Context, hrid string) (sdk.PortalThemeState, error) {
	return apimportaltheme.GetState(ctx, a.SDK, sdk.PortalTheme{Hrid: hrid})
}

// ImportPortalTheme creates or updates a portal theme in APIM outside the operator.
func (a *APIM) ImportPortalTheme(ctx context.Context, thm *v1alpha1.PortalTheme) error {
	dto, err := apimportaltheme.ToDTO(thm)
	if err != nil {
		return err
	}
	_, err = apimportaltheme.Upsert(ctx, a.SDK, dto)
	return err
}

// GetPortal gets a portal state by HRID, server-managed fields included.
func (a *APIM) GetPortal(ctx context.Context, hrid string) (sdk.PortalState, error) {
	return apimportal.GetState(ctx, a.SDK, sdk.Portal{Hrid: hrid})
}

// ImportPortal creates or updates a portal in APIM outside the operator.
func (a *APIM) ImportPortal(ctx context.Context, prtl *v1alpha1.Portal) error {
	dto, err := apimportal.ToDTO(prtl)
	if err != nil {
		return err
	}
	_, err = apimportal.Upsert(ctx, a.SDK, dto)
	return err
}

// GetPortalLink gets a link state by its parent and HRID, server-managed fields included.
func (a *APIM) GetPortalLink(ctx context.Context, parent navigation.Parent, hrid string) (sdk.PortalLinkState, error) {
	return apimportallink.GetState(ctx, a.SDK, apimportallink.DTO{PortalLink: sdk.PortalLink{Hrid: hrid}, Parent: parent})
}

// ImportPortalLink creates or updates a link in APIM outside the operator, under the parent it references.
func (a *APIM) ImportPortalLink(ctx context.Context, link *v1alpha1.PortalLink) error {
	dto, err := apimportallink.ToDTO(link)
	if err != nil {
		return err
	}
	_, err = apimportallink.Upsert(ctx, a.SDK, dto)
	return err
}

// GetPortalListing gets a listing state by its portal HRID and its HRID, server-managed fields included.
func (a *APIM) GetPortalListing(ctx context.Context, portalHrid, hrid string) (sdk.PortalListingState, error) {
	return apimportallisting.GetState(ctx, a.SDK,
		apimportallisting.DTO{PortalListing: sdk.PortalListing{Hrid: hrid}, PortalHrid: portalHrid})
}

// ImportPortalListing creates or updates a listing in APIM outside the operator, under its portal.
func (a *APIM) ImportPortalListing(ctx context.Context, listing *v1alpha1.PortalListing) error {
	dto, err := apimportallisting.ToDTO(listing)
	if err != nil {
		return err
	}
	_, err = apimportallisting.Upsert(ctx, a.SDK, dto)
	return err
}

// GetDocumentation gets a documentation page state by its parent and HRID, server-managed fields included.
func (a *APIM) GetDocumentation(
	ctx context.Context, parent navigation.Parent, hrid string,
) (sdk.DocumentationState, error) {
	return apimdocs.GetState(ctx, a.SDK,
		apimdocs.DTO{Documentation: sdk.Documentation{Hrid: hrid}, Parent: parent})
}

// ImportDocumentation creates or updates a documentation page in APIM outside the operator, under the parent
// it references.
func (a *APIM) ImportDocumentation(ctx context.Context, doc *v1alpha1.Documentation) error {
	dto, err := apimdocs.ToDTO(doc)
	if err != nil {
		return err
	}
	_, err = apimdocs.Upsert(ctx, a.SDK, dto)
	return err
}

// DeletePortal deletes a portal in APIM outside the operator. A portal already gone counts as deleted.
func (a *APIM) DeletePortal(ctx context.Context, hrid string) error {
	return apimportal.Delete(ctx, a.SDK, sdk.Portal{Hrid: hrid})
}

// DeletePortalLink deletes a link in APIM outside the operator. A link already gone counts as deleted.
func (a *APIM) DeletePortalLink(ctx context.Context, parent navigation.Parent, hrid string) error {
	return apimportallink.Delete(ctx, a.SDK, apimportallink.DTO{PortalLink: sdk.PortalLink{Hrid: hrid}, Parent: parent})
}
