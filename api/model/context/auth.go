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

package context

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
)

var _ core.Auth = &Auth{}
var _ core.BasicAuth = &BasicAuth{}

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
	in.SecretRef = refs.NewNamespacedNameFromObjectRef(ref)
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
