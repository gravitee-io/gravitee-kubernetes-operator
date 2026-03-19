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

// Configuration brings support for managing gravitee.io APIM support for organization level operations.
// This service is used for testing purposes only and not initialized by the operator manager.
type Configuration struct {
	*client.Client
}

func NewConfiguration(client *client.Client) *Configuration {
	return &Configuration{Client: client}
}

func (svc *Configuration) GetIDPConfiguration(idpId string) (*model.IDPConfiguration, error) {
	url := svc.ConfigurationTarget("identities").WithPath(idpId)
	idpConfig := new(model.IDPConfiguration)
	if err := svc.HTTP.Get(url.String(), idpConfig); err != nil {
		return nil, err
	}
	return idpConfig, nil
}

func (svc *Configuration) UpdateIDPConfiguration(idpConfiguration *model.IDPConfiguration) error {
	url := svc.ConfigurationTarget("identities").WithPath(idpConfiguration.ID)
	return svc.HTTP.Put(url.String(), idpConfiguration, idpConfiguration)
}