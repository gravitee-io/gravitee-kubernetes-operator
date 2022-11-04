/*
 * Copyright (C) 2015 The Gravitee team (http://gravitee.io)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha1

import (
	"net/http"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ManagementContext represents the configuration for a Management API.
// +kubebuilder:object:generate=true
type ManagementContextSpec struct {
	model.Context `json:",inline"`
}

// ManagementContextStatus defines the observed state of ManagementContext.
type ManagementContextStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.spec.name`
// +kubebuilder:printcolumn:name="BaseUrl",type=string,JSONPath=`.spec.baseUrl`
// +kubebuilder:resource:shortName=graviteecontexts
// ManagementContext is the Schema for the Management API.
type ManagementContext struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ManagementContextSpec   `json:"spec,omitempty"`
	Status ManagementContextStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// ManagementContextList contains a list of ManagementContext.
type ManagementContextList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ManagementContext `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ManagementContext{}, &ManagementContextList{})
}

func (r *ManagementContext) HasAuthentication() bool {
	return r.Spec.Auth != nil
}

func (r *ManagementContext) HasSecretRef() bool {
	if !r.HasAuthentication() {
		return false
	}

	return r.Spec.Auth.SecretRef != nil
}

func (r *ManagementContext) Authenticate(req *http.Request) {
	if !r.HasAuthentication() {
		return
	}

	bearerToken := r.Spec.Auth.BearerToken
	basicAuth := r.Spec.Auth.Credentials

	if bearerToken != "" {
		setBearerToken(req, bearerToken)
	} else if basicAuth != nil {
		setBasicAuth(req, basicAuth)
	}
}

func setBearerToken(request *http.Request, token string) {
	if token != "" {
		request.Header.Add("Authorization", "Bearer "+token)
	}
}

func setBasicAuth(request *http.Request, auth *model.BasicAuth) {
	if auth != nil && auth.Username != "" {
		request.SetBasicAuth(auth.Username, auth.Password)
	}
}
