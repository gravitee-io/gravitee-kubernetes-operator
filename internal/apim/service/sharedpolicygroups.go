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

	spg "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/policygroups"
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

func (svc *SharedPolicyGroup) CreateOrUpdate(spec *spg.SharedPolicyGroup) (*spg.Status, error) {
	return svc.createOrUpdate(spec, false)
}

func (svc *SharedPolicyGroup) DryRunCreateOrUpdate(spec *spg.SharedPolicyGroup) (*spg.Status, error) {
	return svc.createOrUpdate(spec, true)
}

func (svc *SharedPolicyGroup) createOrUpdate(spec *spg.SharedPolicyGroup, dryRun bool) (*spg.Status, error) {
	url := svc.EnvV2Target(sharedPolicyGroupsPath).
		WithPath("_import/crd").
		WithQueryParam("dryRun", strconv.FormatBool(dryRun))

	status := new(spg.Status)

	if err := svc.HTTP.Put(url.String(), spec, status); err != nil {
		return nil, err
	}

	return status, nil
}

func (svc *SharedPolicyGroup) Delete(id string) error {
	url := svc.EnvV2Target(sharedPolicyGroupsPath).WithPath(id)
	return svc.HTTP.Delete(url.String(), nil)
}
