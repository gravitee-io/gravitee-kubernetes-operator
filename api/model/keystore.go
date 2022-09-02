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

package model

type KeyStoreType int

const (
	PEM KeyStoreType = iota
	PKCS12
	JKS
)

type TrustStore struct {
	TrustStoreType KeyStoreType `json:"type,omitempty"`
}

type KeyStore struct {
	KeyStoreType KeyStoreType `json:"type,omitempty"`
}
type PEMTrustStore struct {
	Type    KeyStoreType `json:"type,omitempty"`
	Path    string       `json:"path,omitempty"`
	Content string       `json:"content,omitempty"`
}

type PKCS12TrustStore struct {
	Type     KeyStoreType `json:"type,omitempty"`
	Path     string       `json:"path,omitempty"`
	Content  string       `json:"content,omitempty"`
	Password string       `json:"password,omitempty"`
}
type JKSTrustStore struct {
	Type     KeyStoreType `json:"type,omitempty"`
	Path     string       `json:"path,omitempty"`
	Content  string       `json:"content,omitempty"`
	Password string       `json:"password,omitempty"`
}
type PEMKeyStore struct {
	Type        KeyStoreType `json:"type,omitempty"`
	KeyPath     string       `json:"keyPath,omitempty"`
	KeyContent  string       `json:"keyContent,omitempty"`
	CertPath    string       `json:"certPath,omitempty"`
	CertContent string       `json:"certContent,omitempty"`
}

type PKCS12KeyStore struct {
	Type     KeyStoreType `json:"type,omitempty"`
	Path     string       `json:"path,omitempty"`
	Content  string       `json:"content,omitempty"`
	Password string       `json:"password,omitempty"`
}

type JKSKeyStore struct {
	Type     KeyStoreType `json:"type,omitempty"`
	Path     string       `json:"path,omitempty"`
	Content  string       `json:"content,omitempty"`
	Password string       `json:"password,omitempty"`
}
