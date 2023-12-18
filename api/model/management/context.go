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

type Context struct {
	// The baseURL of a management API instance
	// +kubebuilder:validation:Pattern=`^http(s?):\/\/.+$`
	BaseURL string `json:"baseUrl"`

	// An existing environment id targeted by the context within the organization.
	// +kubebuilder:validation:Required
	EnvID string `json:"environmentId"`

	// The Gravitee APIM organization targeted by the management context.
	// +kubebuilder:validation:Required
	OrgID string `json:"organizationId"`

	// Auth defines the authentication method used to connect to the API Management.
	// Can be either basic authentication credentials, a bearer token
	// or a reference to a kubernetes secret holding one of these two configurations.
	// +kubebuilder:validation:Required
	Auth *Auth `json:"auth"`
}
