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

package dictionary

// DictionaryType is the type of dictionary.
// MANUAL is to be updated manually, DYNAMIC is updated and deployed automatically.
// +kubebuilder:validation:Enum=MANUAL;DYNAMIC;
type DictionaryType string

const (
	ManualType  DictionaryType = "MANUAL"
	DynamicType DictionaryType = "DYNAMIC"
)

// Type defines the specification of a dictionary resource.
// Dictionaries can be used in Gravitee EL expressions:
// `{#dictionaries['hrid']['property key']}`.
type Type struct {
	// Display name of the dictionary.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// Detailed description of the dictionary.
	// +kubebuilder:validation:Optional
	Description string `json:"description,omitempty"`
	// If true, a MANUAL dictionary is deployed in the gateway and a DYNAMIC dictionary is started.
	// Setting this back to false will stop or undeploy the dictionary.
	// +kubebuilder:validation:Required
	Deployed bool `json:"deployed"`
	// +kubebuilder:validation:Required
	DictionaryType DictionaryType `json:"type"`
	// Manual dictionary specification. Required when type is MANUAL, forbidden when type is DYNAMIC.
	// +kubebuilder:validation:Optional
	Manual *ManualSpec `json:"manual,omitempty"`
	// Dynamic dictionary specification. Required when type is DYNAMIC, forbidden when type is MANUAL.
	// +kubebuilder:validation:Optional
	Dynamic *DynamicSpec `json:"dynamic,omitempty"`
}

// ManualSpec defines a manual dictionary with static key/value properties.
type ManualSpec struct {
	// Key/value pairs that constitute the dictionary data.
	// +kubebuilder:validation:Required
	Properties map[string]string `json:"properties"`
}

// DynamicSpec defines a dynamic dictionary populated from an external provider on a schedule.
type DynamicSpec struct {
	// HTTP provider configuration for fetching dictionary data.
	// +kubebuilder:validation:Required
	Provider *Provider `json:"provider"`
	// Renewal schedule controlling how often the provider is polled.
	// +kubebuilder:validation:Required
	Trigger *Trigger `json:"trigger"`
}

// Provider defines the HTTP provider configuration for a DYNAMIC dictionary.
type Provider struct {
	// Type of dictionary provider.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=HTTP;
	ProviderType string `json:"type"`
	// URL of the provider to fetch data from.
	// +kubebuilder:validation:Required
	URL string `json:"url"`
	// HTTP method used to call the provider.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=GET;POST;PUT;PATCH;DELETE;HEAD;OPTIONS;TRACE;CONNECT;
	Method string `json:"method"`
	// JOLT specification to transform the returned payload into dictionary entries.
	// +kubebuilder:validation:Required
	Specification string `json:"specification"`
	// Optional request payload sent to the provider.
	// +kubebuilder:validation:Optional
	Body string `json:"body,omitempty"`
	// Use the system proxy for outbound requests to the provider.
	// +kubebuilder:validation:Optional
	UseSystemProxy bool `json:"useSystemProxy,omitempty"`
	// HTTP headers sent with the provider request.
	// +kubebuilder:validation:Optional
	Headers []ProviderHeader `json:"headers,omitempty"`
}

// ProviderHeader is an HTTP header sent with a provider request.
type ProviderHeader struct {
	// Header name.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// Header value.
	// +kubebuilder:validation:Required
	Value string `json:"value"`
}

// TriggerUnit is the time unit for a dictionary trigger schedule.
// +kubebuilder:validation:Enum=MICROSECONDS;MILLISECONDS;SECONDS;MINUTES;HOURS;DAYS;
type TriggerUnit string

const (
	MicrosecondsUnit TriggerUnit = "MICROSECONDS"
	MillisecondsUnit TriggerUnit = "MILLISECONDS"
	SecondsUnit      TriggerUnit = "SECONDS"
	MinutesUnit      TriggerUnit = "MINUTES"
	HoursUnit        TriggerUnit = "HOURS"
	DaysUnit         TriggerUnit = "DAYS"
)

// Trigger defines the renewal configuration for a DYNAMIC dictionary.
type Trigger struct {
	// Polling interval value (used with Unit to define the schedule).
	// +kubebuilder:validation:Required
	Rate int64 `json:"rate"`
	// Time unit for the polling interval.
	// +kubebuilder:validation:Required
	Unit TriggerUnit `json:"unit"`
}
