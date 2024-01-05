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
	// Indicate if the cors enabled or not
	Enabled bool `json:"enabled"`

	// Access Control -  List of Allowed origins
	AccessControlAllowOrigin []string `json:"allowOrigin,omitempty"`

	// Access Control - List of Exposed Headers
	AccessControlExposeHeaders []string `json:"exposeHeaders,omitempty"`

	// Access Control -  Max age
	AccessControlMaxAge int `json:"maxAge"`

	// Access Control - Allow credentials or not
	AccessControlAllowCredentials bool `json:"allowCredentials"`

	// Access Control - List of allowed methods
	AccessControlAllowMethods []string `json:"allowMethods,omitempty"`

	// Access Control - List of allowed headers
	AccessControlAllowHeaders []string `json:"allowHeaders,omitempty"`

	// +kubebuilder:default:=false
	// Run policies or not
	RunPolicies bool `json:"runPolicies,omitempty"`
}

type HttpClientOptions struct {
	//  Idle Timeout for the http connection
	IdleTimeout uint64 `json:"idleTimeout,omitempty"`

	// Connection timeout of the http connection
	ConnectTimeout uint64 `json:"connectTimeout,omitempty"`

	// Should keep alive be used for the HTTP connection ?
	KeepAlive bool `json:"keepAlive,omitempty"`

	// Read timeout
	ReadTimeout uint64 `json:"readTimeout,omitempty"`

	// Should HTTP/1.1 pipelining be used for the connection or not ?
	Pipelining bool `json:"pipelining,omitempty"`

	// HTTP max concurrent connections
	MaxConcurrentConnections int `json:"maxConcurrentConnections,omitempty"`

	// Should compression be used or not ?
	UseCompression bool `json:"useCompression,omitempty"`

	// Should HTTP redirects be followed or not ?
	FollowRedirects bool `json:"followRedirects,omitempty"`

	// Should HTTP/2 clear text upgrade be used or not ?
	ClearTextUpgrade bool `json:"clearTextUpgrade,omitempty"`

	// HTTP Protocol Version (Possible values Http1 or Http2)
	Version ProtocolVersion `json:"version,omitempty"`
}

type HttpClientSslOptions struct {
	// Whether to trust all issuers or not
	TrustAll bool `json:"trustAll,omitempty"`

	// Verify Hostname when establishing connection
	HostnameVerifier bool `json:"hostnameVerifier,omitempty"`

	// TrustStore type (possible values PEM, PKCS12, JKS)
	TrustStore *TrustStore `json:"trustStore,omitempty"`

	// KeyStore type (possible values PEM, PKCS12, JKS)
	KeyStore *KeyStore `json:"keyStore,omitempty"`
}

type HttpProxy struct {
	// Specifies that the HTTP connection will be established through a proxy
	Enabled bool `json:"enabled,omitempty"`

	// If true, the proxy defined at the system level will be used
	UseSystemProxy bool `json:"useSystemProxy,omitempty"`

	// Proxy host name
	Host string `json:"host,omitempty"`

	// The HTTP proxy port
	Port int `json:"port,omitempty"`

	// The HTTP proxy username (if the proxy requires authentication)
	Username string `json:"username,omitempty"`

	// The HTTP proxy password (if the proxy requires authentication)
	Password string `json:"password,omitempty"`

	// The HTTP proxy type (possible values Http, Socks4, Socks5)
	HttpProxyType HttpProxyType `json:"type,omitempty"`
}

type HttpHeader struct {
	// The HTTP header name
	Name string `json:"name,omitempty"`

	// The HTTP header value
	Value string `json:"value,omitempty"`
}