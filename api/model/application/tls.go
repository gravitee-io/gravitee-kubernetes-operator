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

package application

// CertificateRef is a reference to a Secret or ConfigMap containing a client certificate.
type CertificateRef struct {
	// Kind of the referenced resource. Defaults to "secrets".
	// +kubebuilder:validation:Enum=secrets;configmaps
	// +kubebuilder:default="secrets"
	Kind string `json:"kind,omitempty"`
	// Name of the referenced Secret or ConfigMap.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// Key in the referenced Secret or ConfigMap data. Defaults to "tls.crt".
	// +kubebuilder:default="tls.crt"
	Key string `json:"key,omitempty"`
	// Namespace of the referenced resource. Defaults to the Application namespace.
	// +kubebuilder:validation:Optional
	Namespace string `json:"namespace,omitempty"`
}

// ClientCertificate represents a client certificate for mTLS plans.
// Either content or ref must be set, but not both.
type ClientCertificate struct {
	// Name is an optional label for this certificate.
	// Defaults to the application name suffixed with the certificate index.
	// +kubebuilder:validation:Optional
	Name string `json:"name,omitempty"`
	// Content is the certificate inlined (PEM or Base64) or a template [[ ]] notation.
	// +kubebuilder:validation:Optional
	Content string `json:"content,omitempty"`
	// Ref is a reference to a Secret or ConfigMap containing the certificate.
	// +kubebuilder:validation:Optional
	Ref *CertificateRef `json:"ref,omitempty"`
	// StartsAt is the optional start date of the certificate validity (RFC3339).
	// +kubebuilder:validation:Optional
	StartsAt string `json:"startsAt,omitempty"`
	// EndsAt is the optional end date of the certificate validity (RFC3339).
	// +kubebuilder:validation:Optional
	EndsAt string `json:"endsAt,omitempty"`
	// Encoded indicates whether the content is base64 encoded.
	// If true, the content will be decoded before being sent to APIM.
	// +kubebuilder:validation:Optional
	Encoded bool `json:"encoded,omitempty"`
}
