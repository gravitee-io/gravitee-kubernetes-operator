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

// +kubebuilder:object:generate=true
package model

import (
	"net/http"
	"strings"
)

type ContextRef struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

type Context struct {
	// +kubebuilder:validation:Pattern=`^http(s?):\/\/.+$`
	BaseUrl string `json:"baseUrl"`
	// +kubebuilder:validation:Required
	OrgId string `json:"organizationId"`
	// +kubebuilder:validation:Required
	EnvId string `json:"environmentId"`
	// +kubebuilder:validation:Required
	Auth *Auth `json:"auth"`
}

type Auth struct {
	BearerToken string     `json:"bearerToken,omitempty"`
	Credentials *BasicAuth `json:"credentials,omitempty"`
}

type BasicAuth struct {
	// +kubebuilder:validation:Required
	Username string `json:"username,omitempty"`
	// +kubebuilder:validation:Required
	Password string `json:"password,omitempty"`
}

func (ctx Context) BuildUrl(path string) string {
	orgId, envId := ctx.OrgId, ctx.EnvId
	baseUrl := strings.TrimSuffix(ctx.BaseUrl, "/")
	url := baseUrl + "/management/organizations/" + orgId
	if envId != "" {
		url = url + "/environments/" + envId
	}
	return url + path
}

func (ctx Context) Authenticate(req *http.Request) {
	if ctx.Auth == nil {
		return
	}

	bearerToken := ctx.Auth.BearerToken
	if bearerToken != "" {
		req.Header.Add("Authorization", "Bearer "+bearerToken)
	} else if ctx.Auth.Credentials != nil {
		username := ctx.Auth.Credentials.Username
		password := ctx.Auth.Credentials.Password
		setBasicAuth(req, username, password)
	}
}

func setBasicAuth(request *http.Request, username, password string) {
	if username != "" {
		request.SetBasicAuth(username, password)
	}
}
