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

package apim

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
)

// APIs brings support for managing gravitee.io APIM APIs.
type Export struct {
	*client.Client
}

func NewExport(client *client.Client) *Export {
	return &Export{Client: client}
}

func (svc *Export) V2Api(id string) (*v1alpha1.ApiDefinition, error) {
	url := svc.EnvV1Target("apis").WithPath(id).WithPath("/crd")
	exported := new(v1alpha1.ApiDefinition)
	if err := svc.HTTP.GetYAML(url.String(), &exported); err != nil {
		return nil, err
	}
	return exported, nil
}

func (svc *Export) V4Api(id string) (*v1alpha1.ApiV4Definition, error) {
	url := svc.EnvV2Target("apis").WithPath(id).WithPath("/_export/crd")
	exported := new(v1alpha1.ApiV4Definition)
	if err := svc.HTTP.GetYAML(url.String(), &exported); err != nil {
		return nil, err
	}
	return exported, nil
}
