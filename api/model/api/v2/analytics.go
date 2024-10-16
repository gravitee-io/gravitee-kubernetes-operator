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

package v2

// +kubebuilder:validation:Enum=NONE;CLIENT;PROXY;CLIENT_PROXY
type LoggingMode string

const (
	NoLoggingMode   = "NONE"
	ClientMode      = "CLIENT"
	ProxyMode       = "PROXY"
	ClientProxyMode = "CLIENT_PROXY"
)

// +kubebuilder:validation:Enum=NONE;REQUEST;RESPONSE;REQUEST_RESPONSE
type LoggingScope string

const (
	NoLoggingScope              = "NONE"
	RequestLoggingScope         = "REQUEST"
	ResponseLoggingScope        = "RESPONSE"
	RequestResponseLoggingScope = "REQUEST_RESPONSE"
)

// +kubebuilder:validation:Enum=NONE;HEADERS;PAYLOADS;HEADERS_PAYLOADS
type LoggingContent string

const (
	NoLoggingContent              = "NONE"
	HeadersLoggingContent         = "HEADERS"
	PayloadsLoggingContent        = "PAYLOADS"
	HeadersPayloadsLoggingContent = "HEADERS_PAYLOADS"
)

type Logging struct {
	// The logging mode.
	// CLIENT identifies the inbound request issued to the gateway,
	// while PROXY identifies the request issued to the upstream service.
	Mode LoggingMode `json:"mode,omitempty"`
	// The logging scope (which phase of the request roundtrip should be included in each log entry.
	Scope LoggingScope `json:"scope,omitempty"`
	// Which part of the request/response should be logged ?
	Content LoggingContent `json:"content,omitempty"`
	// The logging condition (supports EL expressions)
	// +kubebuilder:validation:Optional
	Condition string `json:"condition,omitempty"`
}

type Analytics struct {
	// Analytics Enabled or not?
	Enabled bool `json:"enabled"`
	// Analytics Sampling
	Sampling *Sampling `json:"sampling,omitempty"`
	// Analytics Logging
	Logging *Logging `json:"logging,omitempty"`
}

type SamplingType string

const (
	ProbabilitySamplingType = SamplingType("PROBABILITY")
	TemporalSamplingType    = SamplingType("TEMPORAL")
	CountSamplingType       = SamplingType("COUNT")
	QuerySamplingType       = SamplingType("QUERY")
)

type Sampling struct {
	// The sampling type to use
	Type SamplingType `json:"type"`
	// Sampling Value
	Value string `json:"value"`
}
