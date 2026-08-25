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

package am

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
)

var _ core.ContextModel = &Context{}
var _ core.Auth = &Auth{}

// Context is a cloud-free, bearer-only connection to an Access Management instance.
type Context struct {
	// The URL of an Access Management API instance.
	// +kubebuilder:validation:Optional
	BaseUrl string `json:"baseUrl,omitempty"`
	// Allows to override the context path that will be appended to the baseURL.
	// This can be used when reverse proxying AM with URL rewrite.
	// +kubebuilder:validation:Optional
	Path *string `json:"path,omitempty"`
	// An existing organization id targeted by the context on the AM instance.
	// +kubebuilder:validation:Optional
	OrgID string `json:"organizationId,omitempty"`
	// An existing environment id targeted by the context within the organization.
	// +kubebuilder:validation:Optional
	EnvID string `json:"environmentId,omitempty"`
	// Auth is bearer-token only: an inline bearerToken or a secretRef holding key bearerToken.
	Auth *Auth `json:"auth,omitempty"`
}

// Auth authenticates against AM's Automation API. There is no credentials field:
// AM's Automation API is bearer-only.
type Auth struct {
	// The bearer token used to authenticate against the AM Automation API.
	// +kubebuilder:validation:Optional
	BearerToken string `json:"bearerToken,omitempty"`
	// A secret reference holding a "bearerToken" key.
	SecretRef *refs.NamespacedName `json:"secretRef,omitempty"`
}

func (c *Context) GetAuth() core.Auth {
	return c.Auth
}

func (c *Context) GetEnvID() string {
	return c.EnvID
}

func (c *Context) GetOrgID() string {
	return c.OrgID
}

func (c *Context) GetSecretRef() core.ObjectRef {
	return c.Auth.SecretRef
}

func (c *Context) GetURL() string {
	return c.BaseUrl
}

func (c *Context) GetPath() *string {
	return c.Path
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

func (in *Auth) GetBearerToken() string {
	return in.BearerToken
}

func (in *Auth) HasCredentials() bool {
	return false
}

func (in *Auth) GetCredentials() core.BasicAuth {
	return nil
}

func (in *Auth) GetSecretRef() core.ObjectRef {
	return in.SecretRef
}

func (in *Auth) SetSecretRef(ref core.ObjectRef) {
	nsm := refs.NewNamespacedName(ref.GetNamespace(), ref.GetName())
	in.SecretRef = &nsm
}

func (in *Auth) SetCredentials(_, _ string) {
	// AM auth is bearer-only.
}

func (in *Auth) SetToken(token string) {
	in.BearerToken = token
}
