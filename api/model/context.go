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

type Management struct {
	// The URL of a management API instance
	// +kubebuilder:validation:Pattern=`^http(s?):\/\/.+$`
	BaseUrl string `json:"baseUrl"`
	// An existing organization id targeted by the context on the management API instance.
	// +kubebuilder:validation:Required
	OrgId string `json:"organizationId"`
	// An existing environment id targeted by the context within the organization.
	// +kubebuilder:validation:Required
	EnvId string `json:"environmentId"`
	// Auth defines the authentication method used to connect to the API Management.
	// Can be either basic authentication credentials, a bearer token
	// or a reference to a kubernetes secret holding one of these two configurations.
	// +kubebuilder:validation:Required
	Auth *Auth `json:"auth"`
}

type Auth struct {
	// The bearer token used to authenticate against the API Management instance
	// (must be generated from an admin account)
	BearerToken string `json:"bearerToken,omitempty"`
	// The Basic credentials used to authenticate against the API Management instance.
	Credentials *BasicAuth `json:"credentials,omitempty"`
	// A secret reference holding either a bearer token or the user name and password used for basic authentication
	SecretRef *NamespacedName `json:"secretRef,omitempty"`
}

type BasicAuth struct {
	// +kubebuilder:validation:Required
	Username string `json:"username,omitempty"`
	// +kubebuilder:validation:Required
	Password string `json:"password,omitempty"`
}

func (m *Management) HasAuthentication() bool {
	return m.Auth != nil
}

func (m *Management) HasSecretRef() bool {
	if !m.HasAuthentication() {
		return false
	}

	return m.Auth.SecretRef != nil
}

func (m *Management) SecretRef() *NamespacedName {
	if !m.HasSecretRef() {
		return nil
	}

	return m.Auth.SecretRef
}

func (m *Management) SetToken(token string) {
	if !m.HasAuthentication() {
		return
	}

	m.Auth.BearerToken = token
}

func (m *Management) SetCredentials(username, password string) {
	if !m.HasAuthentication() {
		return
	}

	m.Auth.Credentials = &BasicAuth{
		Username: username,
		Password: password,
	}
}
