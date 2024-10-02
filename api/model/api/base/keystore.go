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

type KeyStoreType string

const (
	PEM KeyStoreType = "PEM"
	PKCS12
	JKS
)

type TrustStore struct {
	// The TrustStore type to use (possible values are PEM, PKCS12, JKS)
	// +kubebuilder:validation:Optional
	TrustStoreType KeyStoreType `json:"type"`
}

type KeyStore struct {
	// The KeyStore type to use (possible values are PEM, PKCS12, JKS)
	// +kubebuilder:validation:Optional
	KeyStoreType KeyStoreType `json:"type"`
}
type PEMTrustStore struct {
	// The TrustStore type (should be set to PEM in that case)
	// +kubebuilder:validation:Optional
	Type KeyStoreType `json:"type"`
	// The path to the TrustStore
	// +kubebuilder:validation:Optional
	Path string `json:"path,omitempty"`
	// The base64 encoded trustStore content, if not relying on a path to a file
	// +kubebuilder:validation:Optional
	Content string `json:"content,omitempty"`
}

type PKCS12TrustStore struct {
	// The trustStore type (should be set to PKCS12 in that case)
	Type KeyStoreType `json:"type,omitempty"`
	// The TrustStore path
	// +kubebuilder:validation:Optional
	Path string `json:"path,omitempty"`
	// The base64 encoded trustStore content, if not relying on a path to a file
	// +kubebuilder:validation:Optional
	Content string `json:"content,omitempty"`
	// TrustStore password
	// +kubebuilder:validation:Optional
	Password string `json:"password,omitempty"`
}
type JKSTrustStore struct {
	// The TrustStore type (should be JKS in that case)
	Type KeyStoreType `json:"type,omitempty"`
	// TrustStore path
	// +kubebuilder:validation:Optional
	Path string `json:"path,omitempty"`
	// The base64 encoded trustStore content, if not relying on a path to a file
	// +kubebuilder:validation:Optional
	Content string `json:"content,omitempty"`
	// TrustStore password
	// +kubebuilder:validation:Optional
	Password string `json:"password,omitempty"`
}
type PEMKeyStore struct {
	// KeyStore type (should be PEM in that case)
	Type KeyStoreType `json:"type,omitempty"`
	// KeyStore key path
	// +kubebuilder:validation:Optional
	KeyPath string `json:"keyPath,omitempty"`
	// The base64 encoded trustStore content, if not relying on a path to a file
	// +kubebuilder:validation:Optional
	KeyContent string `json:"keyContent,omitempty"`
	// KeyStore cert path
	// +kubebuilder:validation:Optional
	CertPath string `json:"certPath,omitempty"`
	// KeyStore cert content
	// +kubebuilder:validation:Optional
	CertContent string `json:"certContent,omitempty"`
}

type PKCS12KeyStore struct {
	// KeyStore type (should be PKCS12 in that case)
	Type KeyStoreType `json:"type,omitempty"`
	// KeyStore path
	// +kubebuilder:validation:Optional
	Path string `json:"path,omitempty"`
	// The base64 encoded trustStore content, if not relying on a path to a file
	// +kubebuilder:validation:Optional
	Content string `json:"content,omitempty"`
	// +kubebuilder:validation:Optional
	Password string `json:"password,omitempty"`
}

type JKSKeyStore struct {
	Type KeyStoreType `json:"type,omitempty"`
	// +kubebuilder:validation:Optional
	Path string `json:"path,omitempty"`
	// The base64 encoded trustStore content, if not relying on a path to a file
	// +kubebuilder:validation:Optional
	Content string `json:"content,omitempty"`
	// KeyStore password
	// +kubebuilder:validation:Optional
	Password string `json:"password,omitempty"`
}
