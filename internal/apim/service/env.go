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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
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

func (svc *Env) DryRunImportGroup(spec *group.Type) (*group.Status, error) {
	return svc.importGroup(spec, true)
}

func (svc *Env) ImportGroup(spec *group.Type) (*group.Status, error) {
	return svc.importGroup(spec, false)
}

func (svc *Env) importGroup(spec *group.Type, dryRun bool) (*group.Status, error) {
	url := svc.EnvV2Target("groups").
		WithPath("_import").WithPath("crd").
		WithQueryParam("dryRun", strconv.FormatBool(dryRun))
	status := new(group.Status)
	if err := svc.HTTP.Put(url.String(), spec, status); err != nil {
		return nil, err
	}
	return status, nil
}

func (svc *Env) DeleteGroup(id string) error {
	url := svc.EnvV1Target("configuration").WithPath("groups").WithPath(id)
	return svc.HTTP.Delete(url.String(), nil)
}

func (svc *Env) Get() (*model.Env, error) {
	env := new(model.Env)
	if err := svc.HTTP.Get(svc.URLs.EnvV2.String(), env); err != nil {
		return nil, err
	}
	return env, nil
}
