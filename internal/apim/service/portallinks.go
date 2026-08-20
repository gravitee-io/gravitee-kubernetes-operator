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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

type Links struct {
	*client.Client
}

func NewLinks(client *client.Client) *Links {
	return &Links{Client: client}
}

func (svc *Links) CreateOrUpdate(
	link *v1alpha1.PortalLink,
	prtl *v1alpha1.Portal,
) (portallink.Status, error) {
	return svc.createOrUpdate(link, prtl, false)
}

func (svc *Links) DryRunCreateOrUpdate(
	link *v1alpha1.PortalLink,
	prtl *v1alpha1.Portal,
) (portallink.Status, error) {
	return svc.createOrUpdate(link, prtl, true)
}

func (svc *Links) createOrUpdate(
	link *v1alpha1.PortalLink,
	prtl *v1alpha1.Portal,
	dryRun bool,
) (portallink.Status, error) {
	portalHrid := refs.NewNamespacedNameFromObject(prtl).HRID()
	url := svc.AutomationTarget("portals").
		WithPath(portalHrid).
		WithPath("links").
		WithQueryParam("dryRun", strconv.FormatBool(dryRun))

	dto := model.ToPortalLinkDTO(link.Spec.Type, refs.NewNamespacedNameFromObject(link).HRID())
	importStatus := &portallink.Status{}

	if err := svc.HTTP.Put(url.String(), dto, &importStatus); err != nil {
		return *importStatus, err
	}

	k8s.AddAutomationAPIManagedCondition(link)

	return *importStatus, nil
}

func (svc *Links) Delete(link *v1alpha1.PortalLink, prtl *v1alpha1.Portal) error {
	portalHrid := refs.NewNamespacedNameFromObject(prtl).HRID()
	linkHrid := refs.NewNamespacedNameFromObject(link).HRID()
	url := svc.AutomationTarget("portals").
		WithPath(portalHrid).
		WithPath("links").
		WithPath(linkHrid)
	return svc.HTTP.Delete(url.String(), nil)
}

func (svc *Links) GetByHRID(portalHrid, linkHrid string) (*model.PortalLinkState, error) {
	url := svc.AutomationTarget("portals").
		WithPath(portalHrid).
		WithPath("links").
		WithPath(linkHrid)
	link := new(model.PortalLinkState)
	if err := svc.HTTP.Get(url.String(), link); err != nil {
		return nil, err
	}
	return link, nil
}
