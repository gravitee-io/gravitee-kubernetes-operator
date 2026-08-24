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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/portaltheme"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

type PortalThemes struct {
	*client.Client
}

func NewPortalThemes(client *client.Client) *PortalThemes {
	return &PortalThemes{Client: client}
}

func (svc *PortalThemes) CreateOrUpdate(thm *v1alpha1.PortalTheme) (portaltheme.Status, error) {
	return svc.createOrUpdate(thm, false)
}

func (svc *PortalThemes) DryRunCreateOrUpdate(thm *v1alpha1.PortalTheme) (portaltheme.Status, error) {
	return svc.createOrUpdate(thm, true)
}

func (svc *PortalThemes) createOrUpdate(
	thm *v1alpha1.PortalTheme,
	dryRun bool,
) (portaltheme.Status, error) {
	url := svc.AutomationTarget("portal-themes").
		WithQueryParam("dryRun", strconv.FormatBool(dryRun))

	dto := model.ToPortalThemeDTO(thm)
	importStatus := &portaltheme.Status{}

	if err := svc.HTTP.Put(url.String(), dto, &importStatus); err != nil {
		return *importStatus, err
	}

	k8s.AddAutomationAPIManagedCondition(thm)

	return *importStatus, nil
}

func (svc *PortalThemes) Delete(thm *v1alpha1.PortalTheme) error {
	hrid := refs.NewNamespacedNameFromObject(thm).HRID()
	url := svc.AutomationTarget("portal-themes").WithPath(hrid)
	return svc.HTTP.Delete(url.String(), nil)
}

// GetByHRID For test purposes only.
func (svc *PortalThemes) GetByHRID(hrid string) (*model.PortalThemeState, error) {
	url := svc.AutomationTarget("portal-themes").WithPath(hrid)
	thm := new(model.PortalThemeState)
	if err := svc.HTTP.Get(url.String(), thm); err != nil {
		return nil, err
	}
	return thm, nil
}
