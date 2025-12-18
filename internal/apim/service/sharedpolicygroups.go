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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"strconv"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/sharedpolicygroups"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
)

const (
	spgsPath = "/shared-policy-groups"
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

	url := svc.EnvV2Target(spgsPath).WithPath(id)
	sg := new(model.SharedPolicyGroup)

	if err := svc.HTTP.Get(url.String(), sg); err != nil {
		return nil, err
	}

	return sg, nil
}

func (svc *SharedPolicyGroup) CreateOrUpdate(spg *v1alpha1.SharedPolicyGroup) (*sharedpolicygroups.Status, error) {
	return svc.createOrUpdate(spg, false)
}

func (svc *SharedPolicyGroup) DryRunCreateOrUpdate(spg *v1alpha1.SharedPolicyGroup) (*sharedpolicygroups.Status, error) {
	return svc.createOrUpdate(spg, true)
}

func (svc *SharedPolicyGroup) createOrUpdate(spg *v1alpha1.SharedPolicyGroup, dryRun bool) (*sharedpolicygroups.Status, error) {
	url := svc.AutomationTarget(spgsPath).
		WithQueryParam("dryRun", strconv.FormatBool(dryRun))

	usesHRID := spg.Spec.HRID != "" && spg.Status.CrossID == ""

	status := new(sharedpolicygroups.Status)

	// even if the resource was created with a previous version,
	// Automation API still need an HRID, so one is set.
	// API ignores it because a CrossID is set in that case see PopulateIDs().
	if !usesHRID {
		spg.Spec.HRID = spg.GetID()
		url = url.WithQueryParam("legacyID", strconv.FormatBool(true))
	}

	if err := svc.HTTP.Put(url.String(), spg.Spec, status); err != nil {
		return nil, err
	}

	status.UseHRID = usesHRID

	return status, nil
}

func (svc *SharedPolicyGroup) Delete(spg *v1alpha1.SharedPolicyGroup) error {
	id, legacy := getSPGID(spg)
	url := svc.AutomationTarget(spgsPath).WithPath(id).WithQueryParam("legacy", strconv.FormatBool(legacy))
	return svc.HTTP.Delete(url.String(), nil)
}

func getSPGID(spg *v1alpha1.SharedPolicyGroup) (string, bool) {
	var id string
	var legacy bool
	if spg.Status.UseHRID {
		id = refs.NewNamespacedNameFromObject(spg).HRID()
		legacy = false
	} else {
		id = spg.GetID()
		legacy = true
	}
	return id, legacy
}
