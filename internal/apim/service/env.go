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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/group"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/status"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

type Env struct {
	*client.Client
}

func NewEnv(client *client.Client) *Env {
	return &Env{Client: client}
}

// CreateGroup For tests purposes only.
func (svc *Env) CreateGroup(group *model.Group) error {
	url := svc.EnvV1Target("configuration").WithPath("groups")
	return svc.HTTP.Post(url.String(), group, group)
}

// CreateCategory For tests purposes only.
func (svc *Env) CreateCategory(category *model.Category) error {
	url := svc.EnvV1Target("configuration").WithPath("categories")
	return svc.HTTP.Post(url.String(), category, category)
}

func (svc *Env) DryRunImportGroup(grp *v1alpha1.Group) (*group.Status, error) {
	return svc.importGroup(grp, true)
}

func (svc *Env) ImportGroup(grp *v1alpha1.Group) (*group.Status, error) {
	return svc.importGroup(grp, false)
}

func (svc *Env) importGroup(grp *v1alpha1.Group, dryRun bool) (*group.Status, error) {
	url := svc.AutomationTarget("groups").WithQueryParam("dryRun", strconv.FormatBool(dryRun))

	setHridWithUUID := grp.Spec.ID != "" && !k8s.IsAutomationAPIManaged(grp)
	if setHridWithUUID {
		grp.Spec.HRID = grp.Spec.ID
		url = url.WithQueryParam("hridContainsUUID", strconv.FormatBool(true))
	}

	importStatus := struct {
		ID          string        `json:"id"`
		MemberCount uint          `json:"memberCount"`
		Errors      status.Errors `json:"errors,omitempty"`
	}{}
	if err := svc.HTTP.Put(url.String(), grp.Spec.Type, &importStatus); err != nil {
		return nil, err
	}

	if !setHridWithUUID {
		k8s.AddAutomationAPIManagedCondition(grp)
	}

	return &group.Status{
		ID:      importStatus.ID,
		Members: importStatus.MemberCount,
		Errors:  importStatus.Errors,
	}, nil
}

func (svc *Env) FindGroup(name string) (*model.Group, error) {
	url := svc.EnvV1Target("configuration").
		WithPath("groups").
		WithPath("_paged").
		WithQueryParam("query", name)

	paginatedGroup := new(model.PaginatedGroups)
	if err := svc.HTTP.Get(url.String(), paginatedGroup); err != nil {
		return nil, err
	}

	if paginatedGroup.Page.TotalElements == 0 {
		return nil, nil //nolint:nilnil // Returning nil, nil is intentional: not found is not an error condition
	}

	return &paginatedGroup.Data[0], nil
}

func (svc *Env) DeleteGroup(grp *v1alpha1.Group) error {
	id, hridContainsUUID := getGroupID(grp)
	url := svc.AutomationTarget("groups").
		WithPath(id).
		WithQueryParam("hridContainsUUID", strconv.FormatBool(hridContainsUUID))
	return svc.HTTP.Delete(url.String(), nil)
}

func getGroupID(grp *v1alpha1.Group) (string, bool) {
	if k8s.IsAutomationAPIManaged(grp) {
		return refs.NewNamespacedNameFromObject(grp).HRID(), false
	}
	return grp.GetID(), true
}

func (svc *Env) Get() (*model.Env, error) {
	env := new(model.Env)
	if err := svc.HTTP.Get(svc.URLs.EnvV2.String(), env); err != nil {
		return nil, err
	}
	return env, nil
}
