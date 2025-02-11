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

package v4

type Logging struct {
	// The logging condition. This field is evaluated for HTTP requests and supports EL expressions.
	Condition string `json:"condition,omitempty"`

	// The logging message condition. This field is evaluated for messages and supports EL expressions.
	MessageCondition string `json:"messageCondition,omitempty"`

	// Defines which component of the request should be included in the log payload.
	Content *LoggingContent `json:"content,omitempty"`

	// The logging mode defines which "hop" of the request roundtrip
	// should be included in the log payload.
	// This can be either the inbound request to the gateway,
	// the request issued by the gateway to the upstream service, or both.
	Mode *LoggingMode `json:"mode,omitempty"`

	// Defines which phase of the request roundtrip
	// should be included in the log payload.
	// This can be either the request phase, the response phase, or both.
	Phase *LoggingPhase `json:"phase,omitempty"`
}

type LoggingPhase struct {
	// Should the request phase of the request roundtrip be included in the log payload or not ?
	Request bool `json:"request"`

	// Should the response phase of the request roundtrip be included in the log payload or not ?
	Response bool `json:"response"`
}

func NewLoggingPhase(request, response bool) *LoggingPhase {
	return &LoggingPhase{
		Request:  request,
		Response: response,
	}
}

type LoggingMode struct {
	// If true, the inbound request to the gateway will be included in the log payload
	Entrypoint bool `json:"entrypoint"`

	// If true, the request to the upstream service will be included in the log payload
	Endpoint bool `json:"endpoint"`
}

func NewLoggingMode(entrypoint, endpoint bool) *LoggingMode {
	return &LoggingMode{
		Entrypoint: entrypoint,
		Endpoint:   endpoint,
	}
}

type LoggingContent struct {
	// Should HTTP headers be logged or not ?
	Headers bool `json:"headers"`

	// Should message headers be logged or not ?
	MessageHeaders bool `json:"messageHeaders"`

	// Should HTTP payloads be logged or not ?
	Payload bool `json:"payload"`

	// Should message payloads be logged or not ?
	MessagePayload bool `json:"messagePayload"`

	// Should message metadata be logged or not ?
	MessageMetadata bool `json:"messageMetadata"`
}

func NewLoggingContent(
	headers, messageHeaders, payload, messagePayload, messageMetadata bool,
) *LoggingContent {
	return &LoggingContent{
		Headers:         headers,
		MessageHeaders:  messageHeaders,
		Payload:         payload,
		MessagePayload:  messagePayload,
		MessageMetadata: messageMetadata,
	}
}

type Analytics struct {
	// +kubebuilder:default:=true
	// Analytics Enabled or not?
	Enabled bool `json:"enabled"`

	// Analytics Sampling
	Sampling *Sampling `json:"sampling,omitempty"`

	// Analytics Logging
	Logging *Logging `json:"logging,omitempty"`

	// Analytics Tracing
	Tracing *Tracing `json:"tracing,omitempty"`
}

func NewAnalytics() *Analytics {
	return &Analytics{
		Enabled: true,
	}
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

type Tracing struct {
	// Specify if Tracing is Enabled or not
	Enabled *bool `json:"enabled,omitempty"`

	// Specify if Tracing is Verbose or not
	Verbose *bool `json:"verbose,omitempty"`
}
