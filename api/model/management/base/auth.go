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

package base

import "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"

type Auth struct {
	// The bearer token used to authenticate against the API Management instance
	// (must be generated from an admin account)
	BearerToken string `json:"bearerToken,omitempty"`
	// The Basic credentials used to authenticate against the API Management instance.
	Credentials *BasicAuth `json:"credentials,omitempty"`
	// A secret reference holding either a bearer token or the user name and password used for basic authentication
	SecretRef *refs.NamespacedName `json:"secretRef,omitempty"`
}

type BasicAuth struct {
	// +kubebuilder:validation:Required
	Username string `json:"username,omitempty"`
	// +kubebuilder:validation:Required
	Password string `json:"password,omitempty"`
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
