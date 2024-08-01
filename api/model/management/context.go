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
package management

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
)

var _ core.Auth = &Auth{}
var _ core.BasicAuth = &BasicAuth{}
var _ core.ContextModel = &Context{}

type Context struct {
	// The URL of a management API instance
	// +kubebuilder:validation:Pattern=`^http(s?):\/\/.+$`
	BaseUrl string `json:"baseUrl"`
	// An existing organization id targeted by the context on the management API instance.
	// +kubebuilder:validation:Required
	OrgID string `json:"organizationId"`
	// An existing environment id targeted by the context within the organization.
	// +kubebuilder:validation:Required
	EnvID string `json:"environmentId"`
	// Auth defines the authentication method used to connect to the API Management.
	// Can be either basic authentication credentials, a bearer token
	// or a reference to a kubernetes secret holding one of these two configurations.
	// +kubebuilder:validation:Required
	Auth *Auth `json:"auth"`
}

// GetAuth implements custom.Context.
func (c *Context) GetAuth() core.Auth {
	return c.Auth
}

// GetEnvID implements custom.Context.
func (c *Context) GetEnvID() string {
	return c.EnvID
}

// GetOrgID implements custom.Context.
func (c *Context) GetOrgID() string {
	return c.OrgID
}

// GetSecretRef implements custom.Context.
func (c *Context) GetSecretRef() core.ObjectRef {
	return c.Auth.SecretRef
}

// GetURL implements custom.Context.
func (c *Context) GetURL() string {
	return c.BaseUrl
}

type Auth struct {
	// The bearer token used to authenticate against the API Management instance
	// (must be generated from an admin account)
	BearerToken string `json:"bearerToken,omitempty"`
	// The Basic credentials used to authenticate against the API Management instance.
	Credentials *BasicAuth `json:"credentials,omitempty"`
	// A secret reference holding either a bearer token or the user name and password used for basic authentication
	SecretRef *refs.NamespacedName `json:"secretRef,omitempty"`
}

// GetBearerToken implements custom.Auth.
func (in *Auth) GetBearerToken() string {
	return in.BearerToken
}

// HasCredentials implements custom.Auth.
func (in *Auth) HasCredentials() bool {
	return in.Credentials != nil
}

// GetCredentials implements custom.Auth.
func (in *Auth) GetCredentials() core.BasicAuth {
	return in.Credentials
}

// GetSecretRef implements custom.Auth.
func (in *Auth) GetSecretRef() core.ObjectRef {
	return in.SecretRef
}

// SetCredentials implements custom.Auth.
func (in *Auth) SetCredentials(username string, password string) {
	in.Credentials = &BasicAuth{
		Username: username,
		Password: password,
	}
}

// SetToken implements custom.Auth.
func (in *Auth) SetToken(token string) {
	in.BearerToken = token
}

type BasicAuth struct {
	// +kubebuilder:validation:Required
	Username string `json:"username,omitempty"`
	// +kubebuilder:validation:Required
	Password string `json:"password,omitempty"`
}

// GetPassword implements custom.BasicAuth.
func (in *BasicAuth) GetPassword() string {
	return in.Password
}

// GetUsername implements custom.BasicAuth.
func (in *BasicAuth) GetUsername() string {
	return in.Username
}

func (c *Context) HasAuthentication() bool {
	return c.Auth != nil
}

func (c *Context) HasSecretRef() bool {
	if !c.HasAuthentication() {
		return false
	}

	return c.Auth.SecretRef != nil
}

func (c *Context) SecretRef() *refs.NamespacedName {
	if !c.HasSecretRef() {
		return nil
	}

	return c.Auth.SecretRef
}

func (c *Context) SetToken(token string) {
	if !c.HasAuthentication() {
		return
	}

	c.Auth.BearerToken = token
}

func (c *Context) SetCredentials(username, password string) {
	if !c.HasAuthentication() {
		return
	}

	c.Auth.Credentials = &BasicAuth{
		Username: username,
		Password: password,
	}
}
