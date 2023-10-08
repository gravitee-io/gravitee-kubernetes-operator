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
	TrustStoreType KeyStoreType `json:"type,omitempty"`
}

type KeyStore struct {
	// The KeyStore type to use (possible values are PEM, PKCS12, JKS)
	KeyStoreType KeyStoreType `json:"type,omitempty"`
}
type PEMTrustStore struct {
	// The TrustStore type (should be set to PEM in that case)
	Type KeyStoreType `json:"type,omitempty"`

	// The path to the TrustStore
	Path string `json:"path,omitempty"`

	// The base64 encoded trustStore content, if not relying on a path to a file
	Content string `json:"content,omitempty"`
}

type PKCS12TrustStore struct {
	// // The trustStore type (should be set to PKCS12 in that case)
	Type KeyStoreType `json:"type,omitempty"`

	// The TrustStore path
	Path string `json:"path,omitempty"`

	// The base64 encoded trustStore content, if not relying on a path to a file
	Content string `json:"content,omitempty"`

	// TrustStore password
	Password string `json:"password,omitempty"`
}
type JKSTrustStore struct {
	// The TrustStore type (should be JKS in that case)
	Type KeyStoreType `json:"type,omitempty"`

	// TrustStore path
	Path string `json:"path,omitempty"`

	// The base64 encoded trustStore content, if not relying on a path to a file
	Content string `json:"content,omitempty"`

	// TrustStore password
	Password string `json:"password,omitempty"`
}
type PEMKeyStore struct {
	// KeyStore type (should be PEM in that case)
	Type KeyStoreType `json:"type,omitempty"`

	// KeyStore key path
	KeyPath string `json:"keyPath,omitempty"`

	// The base64 encoded trustStore content, if not relying on a path to a file
	KeyContent string `json:"keyContent,omitempty"`

	// KeyStore cert path
	CertPath string `json:"certPath,omitempty"`

	// KeyStore cert content
	CertContent string `json:"certContent,omitempty"`
}

type PKCS12KeyStore struct {
	// KeyStore type (should be PKCS12 in that case)
	Type KeyStoreType `json:"type,omitempty"`

	// KeyStore path
	Path string `json:"path,omitempty"`

	// The base64 encoded trustStore content, if not relying on a path to a file
	Content string `json:"content,omitempty"`

	// KeyStore password
	Password string `json:"password,omitempty"`
}

type JKSKeyStore struct {
	// KeyStore type (should be JKS in that case)
	Type KeyStoreType `json:"type,omitempty"`

	// KeyStore path
	Path string `json:"path,omitempty"`

	// The base64 encoded trustStore content, if not relying on a path to a file
	Content string `json:"content,omitempty"`

	// KeyStore password
	Password string `json:"password,omitempty"`
}
