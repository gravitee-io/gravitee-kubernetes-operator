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

// Org brings support for managing gravitee.io APIM support for organization level operations.
// This service is used for testing purposes only and not initialized by the operator manager.
type Org struct {
	*client.Client
}

func NewOrg(client *client.Client) *Org {
	return &Org{Client: client}
}

func (svc *Org) CreateUser(user *model.User) error {
	url := svc.OrgTarget("users")
	return svc.HTTP.Post(url.String(), user, user)
}
