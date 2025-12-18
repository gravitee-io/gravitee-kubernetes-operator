// Copyright (C) 2015 The Gravitee team (HTTP://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         HTTP://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"strconv"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/sharedpolicygroups"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
)

const (
	sharedPolicyGroupsPath = "/shared-policy-groups"
)

// SharedPolicyGroup brings support for managing gravitee.io APIM SharedPolicyGroups.
type SharedPolicyGroup struct {
	*client.Client
}

func NewSharedPolicyGroup(client *client.Client) *SharedPolicyGroup {
	return &SharedPolicyGroup{Client: client}
}

// GetByID For tests purposes only.
func (svc *SharedPolicyGroup) GetByID(id string) (*model.SharedPolicyGroup, error) {
	if id == "" {
		return nil, errors.NewNotFoundError()
	}

	url := svc.EnvV2Target(sharedPolicyGroupsPath).WithPath(id)
	sg := new(model.SharedPolicyGroup)

	if err := svc.HTTP.Get(url.String(), sg); err != nil {
		return nil, err
	}

	return sg, nil
}

// GetByHRID For test purposes only.
func (svc *SharedPolicyGroup) GetByHRID(hrid string) (*model.SharedPolicyGroup, error) {
	url := svc.AutomationTarget(sharedPolicyGroupsPath).WithPath(hrid)
	sg := new(model.SharedPolicyGroup)

	if err := svc.HTTP.Get(url.String(), sg); err != nil {
		return nil, err
	}

	return sg, nil
}

func (svc *SharedPolicyGroup) CreateOrUpdate(spg *v1alpha1.SharedPolicyGroup) (*sharedpolicygroups.Status, error) {
	return svc.createOrUpdate(spg, false)
}

func (svc *SharedPolicyGroup) DryRunCreateOrUpdate(
	spg *v1alpha1.SharedPolicyGroup,
) (*sharedpolicygroups.Status, error) {
	return svc.createOrUpdate(spg, true)
}

func (svc *SharedPolicyGroup) createOrUpdate(
	spg *v1alpha1.SharedPolicyGroup,
	dryRun bool,
) (*sharedpolicygroups.Status, error) {
	url := svc.AutomationTarget(sharedPolicyGroupsPath).
		WithQueryParam("dryRun", strconv.FormatBool(dryRun))

	setHridWithUUID := spg.GetID() != "" && !k8s.IsAutomationAPIManaged(spg)

	// If we are updating and CRD was created before upgrade to AutomationAPI
	// Then, set HRID to the ID and pass hridContainsUUID query param
	if setHridWithUUID {
		spg.Spec.HRID = spg.GetID()
		url = url.WithQueryParam("hridContainsUUID", strconv.FormatBool(true))
	}

	status := new(sharedpolicygroups.Status)

	if err := svc.HTTP.Put(url.String(), spg.Spec, status); err != nil {
		return nil, err
	}

	if !setHridWithUUID {
		k8s.AddAutomationAPIManagedCondition(spg)
	}

	return status, nil
}

func (svc *SharedPolicyGroup) Delete(spg *v1alpha1.SharedPolicyGroup) error {
	id, hridContainsUUID := getSPGID(spg)
	url := svc.AutomationTarget(sharedPolicyGroupsPath).
		WithPath(id).
		WithQueryParam("hridContainsUUID", strconv.FormatBool(hridContainsUUID))
	return svc.HTTP.Delete(url.String(), nil)
}

func getSPGID(spg *v1alpha1.SharedPolicyGroup) (string, bool) {
	if k8s.IsAutomationAPIManaged(spg) {
		return refs.NewNamespacedNameFromObject(spg).HRID(), false
	}
	return spg.GetID(), true
}
