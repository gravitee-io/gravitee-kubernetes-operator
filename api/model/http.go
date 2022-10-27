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

const idleTimeout = 60000
const readTimeout = 10000
const connectTimeout = 5000
const maxConcurrentConnections = 100

func NewHttpClientOptions() *HttpClientOptions {
	return &HttpClientOptions{
		IdleTimeout:              idleTimeout,
		ConnectTimeout:           connectTimeout,
		KeepAlive:                true,
		ReadTimeout:              readTimeout,
		Pipelining:               false,
		MaxConcurrentConnections: maxConcurrentConnections,
		UseCompression:           true,
		FollowRedirects:          false,
		ClearTextUpgrade:         true,
		Version:                  Http1,
	}
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

type Request struct {
	Path     string       `json:"path,omitempty"`
	Method   HttpMethod   `json:"method,omitempty"`
	Headers  []HttpHeader `json:"headers,omitempty"`
	Body     string       `json:"body,omitempty"`
	FromRoot bool         `json:"fromRoot,omitempty"`
}

type Response struct {
	Assertions []string `json:"assertions,omitempty"`
}

func NewResponse() *Response {
	return &Response{
		Assertions: []string{"#response.status == 200"},
	}
}

type Step struct {
	Name     string   `json:"name,omitempty"`
	Request  Request  `json:"request,omitempty"`
	Response Response `json:"response,omitempty"`
}

func NewStep() *Step {
	return &Step{
		Name: "default-step",
	}
}

type HealthCheckService struct {
	Enabled  bool    `json:"enabled,omitempty"`
	Steps    []*Step `json:"steps,omitempty"`
	Schedule string  `json:"schedule,omitempty"`
}

func NewHealthCheckService() *HealthCheckService {
	return &HealthCheckService{Schedule: "health-check"}
}

type LoggingMode struct {
	Client bool `json:"client,omitempty"`
	Proxy  bool `json:"proxy,omitempty"`
}

var (
	NoLoggingMode   = LoggingMode{false, false}
	ClientMode      = LoggingMode{true, false}
	ProxyMode       = LoggingMode{false, true}
	ClientProxyMode = LoggingMode{true, true}
)

type LoggingScope struct {
	Request  bool `json:"request,omitempty"`
	Response bool `json:"response,omitempty"`
}

var (
	NoLoggingScope              = LoggingScope{false, false}
	RequestLoggingScope         = LoggingScope{true, false}
	ResponseLoggingScope        = LoggingScope{false, true}
	RequestResponseLoggingScope = LoggingScope{true, true}
)

type LoggingContent struct {
	Headers  bool `json:"headers,omitempty"`
	Payloads bool `json:"payloads,omitempty"`
}

var (
	NoLoggingContent              = LoggingContent{false, false}
	HeadersLoggingContent         = LoggingContent{true, false}
	PayloadsLoggingContent        = LoggingContent{false, true}
	HeadersPayloadsLoggingContent = LoggingContent{true, true}
)

type Logging struct {
	Mode      LoggingMode    `json:"mode,omitempty"`
	Scope     LoggingScope   `json:"scope,omitempty"`
	Content   LoggingContent `json:"content,omitempty"`
	Condition string         `json:"condition,omitempty"`
}

func NewLogging() *Logging {
	return &Logging{
		Mode:    NoLoggingMode,
		Scope:   NoLoggingScope,
		Content: NoLoggingContent,
	}
}

type FailoverCase string

type Failover struct {
	MaxAttempts  int            `json:"maxAttempts,omitempty"`
	RetryTimeout int64          `json:"retryTimeout,omitempty"`
	Cases        []FailoverCase `json:"cases,omitempty"`
}

const maxAttempts = 1
const retryTimeout = 10000

func NewFailover() *Failover {
	return &Failover{
		MaxAttempts:  maxAttempts,
		RetryTimeout: retryTimeout,
		Cases:        []FailoverCase{},
	}
}

type VirtualHost struct {
	Host               string `json:"host,omitempty"`
	Path               string `json:"path,omitempty"`
	OverrideEntrypoint bool   `json:"override_entrypoint,omitempty"`
}

type Proxy struct {
	VirtualHosts     []*VirtualHost   `json:"virtual_hosts,omitempty"`
	Groups           []*EndpointGroup `json:"groups,omitempty"`
	Failover         *Failover        `json:"failover,omitempty"`
	Cors             *Cors            `json:"cors,omitempty"`
	Logging          *Logging         `json:"logging,omitempty"`
	StripContextPath bool             `json:"strip_context_path,omitempty"`
	PreserveHost     bool             `json:"preserve_host,omitempty"`
}

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

type Plugin struct {
	Policy        string            `json:"policy,omitempty"`
	Resource      string            `json:"resource,omitempty"`
	Configuration *GenericStringMap `json:"configuration,omitempty"` // todo: check with David
}
