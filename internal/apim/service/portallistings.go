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

package service

import (
	"strconv"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/portallisting"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

type Listings struct {
	*client.Client
}

func NewListings(client *client.Client) *Listings {
	return &Listings{Client: client}
}

func (svc *Listings) CreateOrUpdate(
	listing *v1alpha1.PortalListing,
	portalHrid string,
) (portallisting.Status, error) {
	return svc.createOrUpdate(listing, portalHrid, false)
}

func (svc *Listings) DryRunCreateOrUpdate(
	listing *v1alpha1.PortalListing,
	portalHrid string,
) (portallisting.Status, error) {
	return svc.createOrUpdate(listing, portalHrid, true)
}

func (svc *Listings) createOrUpdate(
	listing *v1alpha1.PortalListing,
	portalHrid string,
	dryRun bool,
) (portallisting.Status, error) {
	url := svc.AutomationTarget("portals").
		WithPath(portalHrid).
		WithPath("listings").
		WithQueryParam("dryRun", strconv.FormatBool(dryRun))

	dto := toPortalListingDTO(listing)
	importStatus := &portallisting.Status{}

	if err := svc.HTTP.Put(url.String(), dto, &importStatus); err != nil {
		return *importStatus, err
	}

	k8s.AddAutomationAPIManagedCondition(listing)

	return *importStatus, nil
}

func (svc *Listings) Delete(portalHrid, listingHrid string) error {
	url := svc.AutomationTarget("portals").
		WithPath(portalHrid).
		WithPath("listings").
		WithPath(listingHrid)
	return svc.HTTP.Delete(url.String(), nil)
}

// GetByHRID For test purposes only.
func (svc *Listings) GetByHRID(portalHrid, listingHrid string) (*model.PortalListingState, error) {
	url := svc.AutomationTarget("portals").
		WithPath(portalHrid).
		WithPath("listings").
		WithPath(listingHrid)
	listing := new(model.PortalListingState)
	if err := svc.HTTP.Get(url.String(), listing); err != nil {
		return nil, err
	}
	return listing, nil
}

func toPortalListingDTO(listing *v1alpha1.PortalListing) *model.PortalListingDTO {
	dto := &model.PortalListingDTO{
		HRID: refs.NewNamespacedNameFromObject(listing).HRID(),
		APIs: make([]model.PortalListingApiEntryDTO, 0, len(listing.Spec.APIs)),
	}

	for i := range listing.Spec.APIs {
		entry := &listing.Spec.APIs[i]
		ns := entry.Ref.Namespace
		if ns == "" {
			ns = listing.Namespace
		}
		apiRef := refs.NewNamespacedName(ns, entry.Ref.Name)
		dto.APIs = append(dto.APIs, model.PortalListingApiEntryDTO{
			ApiHrid:  apiRef.HRID(),
			Location: entry.Location,
			Order:    entry.Order,
		})
	}

	return dto
}
