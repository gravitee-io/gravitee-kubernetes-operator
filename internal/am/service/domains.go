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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/client"
)

const domainsPath = "domains"

type Domains struct {
	*client.Client
}

func NewDomains(c *client.Client) *Domains {
	return &Domains{Client: c}
}

// Probe checks that AM's Automation API is reachable at this org/env.
// A 200 is the whole contract: the body is discarded.
func (svc *Domains) Probe() error {
	url := svc.AutomationTarget(domainsPath).WithQueryParam("size", "1")
	return svc.HTTP.Get(url.String(), nil)
}
