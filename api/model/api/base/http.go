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

type HttpProxyType string
type SOCKSType string

const (
	Http   HttpProxyType = "HTTP"
	Socks4 SOCKSType     = "SOCKS4"
	Socks5 SOCKSType     = "SOCKS5"
)

type ProtocolVersion string

const (
	Http1 ProtocolVersion = "HTTP_1_1"
	Http2 ProtocolVersion = "HTTP_2"
)

// +kubebuilder:validation:Enum=GET;POST;PUT;PATCH;DELETE;OPTIONS;HEAD;CONNECT;TRACE;OTHER
type HttpMethod string

type Cors struct {
	Enabled                       bool     `json:"enabled"`
	AccessControlAllowOrigin      []string `json:"allowOrigin,omitempty"`
	AccessControlExposeHeaders    []string `json:"exposeHeaders,omitempty"`
	AccessControlMaxAge           int      `json:"maxAge"`
	AccessControlAllowCredentials bool     `json:"allowCredentials"`
	AccessControlAllowMethods     []string `json:"allowMethods,omitempty"`
	AccessControlAllowHeaders     []string `json:"allowHeaders,omitempty"`
	// +kubebuilder:default:=false
	RunPolicies bool `json:"runPolicies,omitempty"`
}

type HttpClientOptions struct {
	IdleTimeout              uint64          `json:"idleTimeout,omitempty"`
	ConnectTimeout           uint64          `json:"connectTimeout,omitempty"`
	KeepAlive                bool            `json:"keepAlive,omitempty"`
	ReadTimeout              uint64          `json:"readTimeout,omitempty"`
	Pipelining               bool            `json:"pipelining,omitempty"`
	MaxConcurrentConnections int             `json:"maxConcurrentConnections,omitempty"`
	UseCompression           bool            `json:"useCompression,omitempty"`
	FollowRedirects          bool            `json:"followRedirects,omitempty"`
	ClearTextUpgrade         bool            `json:"clearTextUpgrade,omitempty"`
	Version                  ProtocolVersion `json:"version,omitempty"`
}

type HttpClientSslOptions struct {
	TrustAll         bool        `json:"trustAll,omitempty"`
	HostnameVerifier bool        `json:"hostnameVerifier,omitempty"`
	TrustStore       *TrustStore `json:"trustStore,omitempty"`
	KeyStore         *KeyStore   `json:"keyStore,omitempty"`
}

type HttpProxy struct {
	Enabled        bool          `json:"enabled,omitempty"`
	UseSystemProxy bool          `json:"useSystemProxy,omitempty"`
	Host           string        `json:"host,omitempty"`
	Port           int           `json:"port,omitempty"`
	Username       string        `json:"username,omitempty"`
	Password       string        `json:"password,omitempty"`
	HttpProxyType  HttpProxyType `json:"type,omitempty"`
}

type HttpHeader struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
