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

package model

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/application"
)

type Application struct {
	Id          string               `json:"id,omitempty"`
	Name        string               `json:"name,omitempty"`
	Status      string               `json:"status,omitempty"`
	Description string               `json:"description,omitempty"`
	Settings    *application.Setting `json:"settings,omitempty"`
	Background  string               `json:"background,omitempty"`
	Domain      string               `json:"domain,omitempty"`
	Groups      []string             `json:"groups,omitempty"`
	Picture     string               `json:"picture,omitempty"`
	AppKeyMode  *application.KeyMode `json:"app_key_mode,omitempty"`
	AppType     string               `json:"type,omitempty"`
}

type ApplicationMetaData struct {
	Name          string                      `json:"name"`
	ApplicationId string                      `json:"applicationId,omitempty"`
	Value         string                      `json:"value,omitempty"`
	DefaultValue  string                      `json:"defaultValue,omitempty"`
	Format        *application.MetaDataFormat `json:"format,omitempty"`
	Hidden        bool                        `json:"hidden,omitempty"`
	Key           string                      `json:"key,omitempty"`
}
