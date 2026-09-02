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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/portallink"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	gohttp "github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

// LinkParent identifies the owning resource of a PortalLink. A link is
// attached to exactly one of a portal or an API; the parent determines which
// automation endpoint the link is synced to. The HRID is computed from the
// resolved object here in the service layer.
type LinkParent struct {
	// Portal is the owning portal, or nil when attached to an API.
	Portal *v1alpha1.Portal
	// API is the owning API, or nil when attached to a portal.
	API core.ApiDefinitionObject
}

type Links struct {
	*client.Client
}

func NewLinks(client *client.Client) *Links {
	return &Links{Client: client}
}

func (svc *Links) CreateOrUpdate(
	link *v1alpha1.PortalLink,
	parent LinkParent,
) (portallink.Status, error) {
	return svc.createOrUpdate(link, parent, false)
}

func (svc *Links) DryRunCreateOrUpdate(
	link *v1alpha1.PortalLink,
	parent LinkParent,
) (portallink.Status, error) {
	return svc.createOrUpdate(link, parent, true)
}

func (svc *Links) createOrUpdate(
	link *v1alpha1.PortalLink,
	parent LinkParent,
	dryRun bool,
) (portallink.Status, error) {
	url := svc.linksTarget(parent).
		WithQueryParam("dryRun", strconv.FormatBool(dryRun))

	dto := model.ToPortalLinkDTO(link.Spec.Type, refs.NewNamespacedNameFromObject(link).HRID())
	importStatus := &portallink.Status{}

	if err := svc.HTTP.Put(url.String(), dto, &importStatus); err != nil {
		return *importStatus, err
	}

	k8s.AddAutomationAPIManagedCondition(link)

	return *importStatus, nil
}

func (svc *Links) Delete(link *v1alpha1.PortalLink, parent LinkParent) error {
	linkHrid := refs.NewNamespacedNameFromObject(link).HRID()
	url := svc.linksTarget(parent).WithPath(linkHrid)
	return svc.HTTP.Delete(url.String(), nil)
}

func (svc *Links) GetByHRID(parent LinkParent, linkHrid string) (*model.PortalLinkState, error) {
	url := svc.linksTarget(parent).WithPath(linkHrid)
	link := new(model.PortalLinkState)
	if err := svc.HTTP.Get(url.String(), link); err != nil {
		return nil, err
	}
	return link, nil
}

// linksTarget builds the links collection URL nested under the owning portal
// or API, computing the parent HRID from the resolved object.
func (svc *Links) linksTarget(parent LinkParent) *gohttp.URL {
	if parent.API != nil {
		apiRef := refs.NewNamespacedName(parent.API.GetNamespace(), parent.API.GetName())
		return svc.AutomationTarget("apis").
			WithPath(apiRef.HRID()).
			WithPath("links")
	}
	portalHrid := refs.NewNamespacedNameFromObject(parent.Portal).HRID()
	return svc.AutomationTarget("portals").
		WithPath(portalHrid).
		WithPath("links")
}
