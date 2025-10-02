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

// +kubebuilder:validation:Enum=HTTP_1_1;HTTP_2;
type ProtocolVersion string

// +kubebuilder:validation:Enum=GET;POST;PUT;PATCH;DELETE;OPTIONS;HEAD;CONNECT;TRACE;OTHER
type HttpMethod string

type Cors struct {
	// Indicate if the cors enabled or not
	Enabled bool `json:"enabled"`
	// Access Control -  List of Allowed origins
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	AccessControlAllowOrigin []string `json:"allowOrigin"`
	// Access Control - List of Exposed Headers
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	AccessControlExposeHeaders []string `json:"exposeHeaders"`
	// Access Control -  Max age
	AccessControlMaxAge int `json:"maxAge"`
	// Access Control - Allow credentials or not
	AccessControlAllowCredentials bool `json:"allowCredentials"`
	// Access Control - List of allowed methods
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	AccessControlAllowMethods []string `json:"allowMethods"`
	// Access Control - List of allowed headers
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	AccessControlAllowHeaders []string `json:"allowHeaders"`
	// +kubebuilder:default:=false
	// Run policies or not
	RunPolicies bool `json:"runPolicies,omitempty"`
}

type HttpClientOptions struct {
	//  Idle Timeout for the http connection
	// +kubebuilder:validation:Optional
	IdleTimeout *uint64 `json:"idleTimeout,omitempty"`
	// Connection timeout of the http connection
	// +kubebuilder:validation:Optional
	ConnectTimeout *uint64 `json:"connectTimeout,omitempty"`
	// +kubebuilder:validation:Optional (default: true)
	KeepAlive *bool `json:"keepAlive"`
	// +kubebuilder:validation:Optional
	// Should keep alive be used for the HTTP connection ? (default: 30000)
	KeepAliveTimeout *uint64 `json:"keepAliveTimeout"`
	// Read timeout
	// +kubebuilder:validation:Optional
	ReadTimeout *uint64 `json:"readTimeout,omitempty"`
	// Should HTTP/1.1 pipelining be used for the connection or not ? (default: false)
	Pipelining *bool `json:"pipelining"`
	// HTTP max concurrent connections
	// +kubebuilder:validation:Optional
	MaxConcurrentConnections *int `json:"maxConcurrentConnections,omitempty"`
	// Should compression be used or not ? (default: false)
	UseCompression *bool `json:"useCompression"`
	// +kubebuilder:validation:Optional
	// Propagate Client Accept-Encoding header (default: false)
	PropagateClientAcceptEncoding *bool `json:"propagateClientAcceptEncoding"`
	// Should HTTP redirects be followed or not ?
	FollowRedirects *bool `json:"followRedirects"`
	// +kubebuilder:validation:Optional
	// Should HTTP/2 clear text upgrade be used or not ? (default: true)
	ClearTextUpgrade *bool `json:"clearTextUpgrade"`
	// HTTP Protocol Version (Possible values Http1 or Http2) (default: HTTP_1_1)
	ProtocolVersion *ProtocolVersion `json:"version,omitempty"`
}

type HttpClientSslOptions struct {
	// Whether to trust all issuers or not (default: false)
	TrustAll *bool `json:"trustAll,omitempty"`
	// Verify Hostname when establishing connection
	HostnameVerifier *bool `json:"hostnameVerifier,omitempty"`
	// TrustStore type (possible values PEM, PKCS12, JKS)
	TrustStore *TrustStore `json:"trustStore,omitempty"`
	// KeyStore type (possible values PEM, PKCS12, JKS)
	KeyStore *KeyStore `json:"keyStore,omitempty"`
	// Http headers
	Headers []*HttpHeader `json:"headers,omitempty"`
}

type HttpProxy struct {
	// Specifies that the HTTP connection will be established through a proxy (default: false)
	Enabled *bool `json:"enabled,omitempty"`
	// If true, the proxy defined at the system level will be used (default: false)
	UseSystemProxy *bool `json:"useSystemProxy,omitempty"`
	// Proxy host name
	// +kubebuilder:validation:Optional
	Host *string `json:"host,omitempty"`
	// The HTTP proxy port
	// +kubebuilder:validation:Optional
	Port *int `json:"port,omitempty"`
	// The HTTP proxy username (if the proxy requires authentication)
	// +kubebuilder:validation:Optional
	Username *string `json:"username,omitempty"`
	// The HTTP proxy password (if the proxy requires authentication)
	// +kubebuilder:validation:Optional
	Password *string `json:"password,omitempty"`
	// The HTTP proxy type (possible values Http, Socks4, Socks5)
	HttpProxyType HttpProxyType `json:"type,omitempty"`
}

type HttpHeader struct {
	// The HTTP header name
	// +kubebuilder:validation:Optional
	Name *string `json:"name,omitempty"`
	// The HTTP header value
	// +kubebuilder:validation:Optional
	Value *string `json:"value,omitempty"`
}
