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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Status struct {
	// The Key of the security domain in the AM instance.
	// +kubebuilder:validation:Optional
	Key string `json:"key,omitempty"`
	// The organization Key defined in the AM context.
	// +kubebuilder:validation:Optional
	OrgID string `json:"organizationId,omitempty"`
	// The environment Key defined in the AM context.
	// +kubebuilder:validation:Optional
	EnvID string `json:"environmentId,omitempty"`
	// Conditions describe the current conditions of the security domain.
	//
	// Known condition types are:
	// * "Accepted"
	// * "ResolvedRefs"
	//
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=8
	// +kubebuilder:default={}
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

type DomainSubResourceStatus struct {
	Status `json:",inline"`
	// +kubebuilder:validation:Optional
	DomainKey string `json:"domainKey"`
}

type IdentityProviderStatus struct {
	DomainSubResourceStatus `json:",inline"`
	// The name of the identity provider in AM. For the system identity provider, the one AM builds.
	// +kubebuilder:validation:Optional
	Name string `json:"name,omitempty"`
	// The plugin type of the identity provider in AM. For the system identity provider, the one AM builds.
	// +kubebuilder:validation:Optional
	Type string `json:"type,omitempty"`
}

type CertificateStatus struct {
	DomainSubResourceStatus `json:",inline"`
	// The name of the certificate in AM. For the system certificate, the one AM builds.
	// +kubebuilder:validation:Optional
	Name string `json:"name,omitempty"`
	// The plugin type of the certificate in AM. For the system certificate, the one AM builds.
	// +kubebuilder:validation:Optional
	Type string `json:"type,omitempty"`
	// When the certificate expires, as reported by AM.
	// +kubebuilder:validation:Optional
	ExpiresAt *metav1.Time `json:"expiresAt,omitempty"`
}
