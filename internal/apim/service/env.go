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

func (svc *Env) Get() (*model.Env, error) {
	env := new(model.Env)
	if err := svc.HTTP.Get(svc.URLs.EnvV2.String(), env); err != nil {
		return nil, err
	}
	return env, nil
}
