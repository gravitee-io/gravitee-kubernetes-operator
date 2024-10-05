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
	// The URL of a management API instance.
	// This is optional when this context targets Gravitee Cloud otherwise it is required.
	// +kubebuilder:validation:Optional
	BaseUrl string `json:"baseUrl,omitempty"`
	// An existing organization id targeted by the context on the management API instance.
	// This is optional when this context targets Gravitee Cloud otherwise it is required.
	// +kubebuilder:validation:Optional
	OrgID string `json:"organizationId,omitempty"`
	// An existing environment id targeted by the context within the organization.
	// This is optional when this context targets Gravitee Cloud
	// and your cloud token contains only one environment ID, otherwise it is required.
	// +kubebuilder:validation:Optional
	EnvID string `json:"environmentId,omitempty"`
	// Auth defines the authentication method used to connect to the API Management.
	// Can be either basic authentication credentials, a bearer token
	// or a reference to a kubernetes secret holding one of these two configurations.
	// This is optional when this context targets Gravitee Cloud.
	Auth *Auth `json:"auth,omitempty"`
	// Cloud when set (token or secretRef) this context will target Gravitee Cloud.
	// BaseUrl will be defaulted from token data if not set,
	// Auth is defaulted to use the token (bearerToken),
	// OrgID is extracted from the token,
	// EnvID is defaulted when the token contains exactly one environment.
	Cloud *Cloud `json:"cloud,omitempty"`
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

// HasCloud implements custom.Context.
func (c *Context) HasCloud() bool {
	return c.Cloud != nil
}

// GetCloud implements custom.Context.
func (c *Context) GetCloud() core.Cloud {
	return c.Cloud
}

func (c *Context) ConfigureCloud(url string, orgID string, envID string) {
	c.BaseUrl = url
	c.OrgID = orgID
	c.EnvID = envID

	if !c.HasAuthentication() {
		c.Auth = &Auth{}
	}

	// override Auth to be bearer token
	if c.GetCloud().HasSecretRef() {
		c.GetAuth().SetSecretRef(c.GetCloud().GetSecretRef())
		c.Auth.BearerToken = ""
	} else {
		c.GetAuth().SetToken(c.GetCloud().GetToken())
		c.Auth.SecretRef = nil
	}
	c.Auth.Credentials = nil
}

type Auth struct {
	// The bearer token used to authenticate against the API Management instance
	// (must be generated from an admin account)
	// +kubebuilder:validation:Optional
	BearerToken string `json:"bearerToken,omitempty"`
	// The Basic credentials used to authenticate against the API Management instance.
	Credentials *BasicAuth `json:"credentials,omitempty"`
	// A secret reference holding either a "bearerToken" key for bearer token authentication
	// or "username" and "password" keys for basic authentication
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

// SetSecretRef implements custom.Auth.
func (in *Auth) SetSecretRef(ref core.ObjectRef) {
	nsm := refs.NewNamespacedName(ref.GetNamespace(), ref.GetName())
	in.SecretRef = &nsm
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

type Cloud struct {
	// Token plain text Gravitee cloud token (JWT)
	// +kubebuilder:validation:Optional
	Token string `json:"token,omitempty"`
	// SecretRef secret reference holding the Gravitee cloud token in the "cloudToken" key
	// +kubebuilder:validation:Optional
	SecretRef *refs.NamespacedName `json:"secretRef,omitempty"`
}

func (c *Cloud) HasSecretRef() bool {
	return c.SecretRef != nil && c.SecretRef.GetName() != ""
}

func (c *Cloud) GetSecretRef() core.ObjectRef {
	return c.SecretRef
}

func (c *Cloud) GetToken() string {
	return c.Token
}

func (c *Cloud) IsEnabled() bool {
	return c.Token != "" || c.HasSecretRef()
}
