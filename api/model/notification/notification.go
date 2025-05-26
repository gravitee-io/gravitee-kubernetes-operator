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

package notification

import "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"

const TargetConsole Target = "console"
const EventTypeAPI EventType = "api"

// ApiEvent defines the events that can be sent to the console.
// +kubebuilder:validation:Enum=APIKEY_EXPIRED;APIKEY_RENEWED;APIKEY_REVOKED;SUBSCRIPTION_NEW;SUBSCRIPTION_ACCEPTED;SUBSCRIPTION_CLOSED;SUBSCRIPTION_PAUSED;SUBSCRIPTION_RESUMED;SUBSCRIPTION_REJECTED;SUBSCRIPTION_TRANSFERRED;SUBSCRIPTION_FAILED;NEW_SUPPORT_TICKET;API_STARTED;API_STOPPED;API_UPDATED;API_DEPLOYED;NEW_RATING;NEW_RATING_ANSWER;MESSAGE;ASK_FOR_REVIEW;REVIEW_OK;REQUEST_FOR_CHANGES;API_DEPRECATED;NEW_SPEC_GENERATED
type ApiEvent string

// Target defines the target of the notification.
// +kubebuilder:validation:Enum=console;
type Target string

// EventType defines the subject of those events.
// +kubebuilder:validation:Enum=api;
type EventType string

type Type struct {
	// Target of the notification: "console" is for notifications in Gravitee console UI.
	// For each target there is an attribute of the same name to configure it.
	// +kubebuilder:validation:Required
	// +kubebuilder:default="console"
	Target Target `json:"target"`
	// EventType defines the subject of those events.
	// Notification can be used in API or Applications, each of those have different events.
	// An attribute starting with eventType value exists in the target configuration
	// to configure events: < eventType >Events (e.g apiEvents)
	// +kubebuilder:validation:Required
	// +kubebuilder:default="api"
	EventType EventType `json:"eventType"`
	// Console is used when the target value is "console" and is meant
	// to configure Gravitee console UI notifications.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	Console Console `json:"console"`
}

type Console struct {
	// List of group references associated with this console notification.
	// These groups are references to gravitee.io/Group custom resources created on the cluster.
	// All members of those groups will receive a notification for the defined events.
	// +kubebuilder:validation:Optional
	GroupRefs []refs.NamespacedName `json:"groupRefs"`
	// List events that will trigger a notification for an API. Recipients are the API primary owner
	// and all members of groups referenced in groupRefs
	// Notification spec attribute eventType must be set to "api".
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:items:UniqueItems=true
	APIEvents []ApiEvent `json:"apiEvents"`
	// List of groups associated with the API.
	// These groups are id to existing groups in APIM.
	// +kubebuilder:validation:Optional
	Groups []string `json:"groups"`
}

func (c *Console) APIEventsAsString() []string {
	result := make([]string, 0)
	for _, event := range c.APIEvents {
		result = append(result, string(event))
	}
	return result
}
