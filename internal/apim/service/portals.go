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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/portal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

type Portals struct {
	*client.Client
}

func NewPortals(client *client.Client) *Portals {
	return &Portals{Client: client}
}

func (svc *Portals) CreateOrUpdate(prtl *v1alpha1.Portal) (portal.Status, error) {
	return svc.createOrUpdate(prtl, false)
}

func (svc *Portals) DryRunCreateOrUpdate(prtl *v1alpha1.Portal) (portal.Status, error) {
	return svc.createOrUpdate(prtl, true)
}

func (svc *Portals) createOrUpdate(
	prtl *v1alpha1.Portal,
	dryRun bool,
) (portal.Status, error) {
	url := svc.AutomationTarget("portals").
		WithQueryParam("dryRun", strconv.FormatBool(dryRun))

	dto := &model.PortalDTO{
		HRID: refs.NewNamespacedNameFromObject(prtl).HRID(),
		Type: prtl.Spec.Type,
	}
	importStatus := &portal.Status{}

	if err := svc.HTTP.Put(url.String(), dto, &importStatus); err != nil {
		return *importStatus, err
	}

	k8s.AddAutomationAPIManagedCondition(prtl)

	return *importStatus, nil
}

func (svc *Portals) Delete(prtl *v1alpha1.Portal) error {
	hrid := refs.NewNamespacedNameFromObject(prtl).HRID()
	url := svc.AutomationTarget("portals").WithPath(hrid)
	return svc.HTTP.Delete(url.String(), nil)
}

// GetByHRID For test purposes only.
func (svc *Portals) GetByHRID(hrid string) (*model.PortalState, error) {
	url := svc.AutomationTarget("portals").WithPath(hrid)
	prtl := new(model.PortalState)
	if err := svc.HTTP.Get(url.String(), prtl); err != nil {
		return nil, err
	}
	return prtl, nil
}
