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
	Condition        string          `json:"condition,omitempty"`
	MessageCondition string          `json:"messageCondition,omitempty"`
	Content          *LoggingContent `json:"content,omitempty"`
	Mode             *LoggingMode    `json:"mode,omitempty"`
	Phase            *LoggingPhase   `json:"phase,omitempty"`
}

type LoggingPhase struct {
	Request  bool `json:"request"`
	Response bool `json:"response"`
}

func NewLoggingPhase(request, response bool) *LoggingPhase {
	return &LoggingPhase{
		Request:  request,
		Response: response,
	}
}

type LoggingMode struct {
	Entrypoint bool `json:"entrypoint"`
	Endpoint   bool `json:"endpoint"`
}

func NewLoggingMode(entrypoint, endpoint bool) *LoggingMode {
	return &LoggingMode{
		Entrypoint: entrypoint,
		Endpoint:   endpoint,
	}
}

type LoggingContent struct {
	Headers         bool `json:"headers"`
	MessageHeaders  bool `json:"messageHeaders"`
	Payload         bool `json:"payload"`
	MessagePayload  bool `json:"messagePayload"`
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
	Enabled  bool      `json:"enabled"`
	Sampling *Sampling `json:"sampling,omitempty"`
	Logging  *Logging  `json:"logging,omitempty"`
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
	Type  SamplingType `json:"type"`
	Value string       `json:"value"`
}
