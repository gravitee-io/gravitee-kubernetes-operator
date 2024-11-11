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

// +kubebuilder:validation:Enum=PEM;PKCS12;JKS;
type KeyStoreType string

type TrustStore struct {
	// The TrustStore type to use (possible values are PEM, PKCS12, JKS)
	TrustStoreType KeyStoreType `json:"type,omitempty"`
	// +kubebuilder:validation:Optional
	Path *string `json:"path,omitempty"`
	// The base64 encoded trustStore content, if not relying on a path to a file
	// +kubebuilder:validation:Optional
	Content *string `json:"content,omitempty"`
	// TrustStore password (Not applicable for PEM TrustStore)
	// +kubebuilder:validation:Optional
	Password *string `json:"password,omitempty"`
}

type KeyStore struct {
	// The KeyStore type to use (possible values are PEM, PKCS12, JKS)
	KeyStoreType KeyStoreType `json:"type,omitempty"`
	// KeyStore path
	// +kubebuilder:validation:Optional
	Path *string `json:"path,omitempty"`
	// The base64 encoded trustStore content, if not relying on a path to a file
	// +kubebuilder:validation:Optional
	Content *string `json:"content,omitempty"`
	// +kubebuilder:validation:Optional
	Password *string `json:"password,omitempty"`

	// KeyStore key path (Only applicable for PEM KeyStore)
	// +kubebuilder:validation:Optional
	KeyPath *string `json:"keyPath,omitempty"`
	// The base64 encoded trustStore content, if not relying on a path to a file
	// (Only applicable for PEM KeyStore)
	// +kubebuilder:validation:Optional
	KeyContent *string `json:"keyContent,omitempty"`
	// KeyStore cert path (Only applicable for PEM KeyStore)
	// +kubebuilder:validation:Optional
	CertPath *string `json:"certPath,omitempty"`
	// KeyStore cert content (Only applicable for PEM KeyStore)
	// +kubebuilder:validation:Optional
	CertContent *string `json:"certContent,omitempty"`
}
