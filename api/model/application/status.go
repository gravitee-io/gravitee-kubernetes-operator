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

package application

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/status"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Status struct {
	// The organization ID, if a management context has been defined to sync with an APIM instance
	OrgID string `json:"organizationId,omitempty"`
	// The environment ID, if a management context has been defined to sync with an APIM instance
	EnvID string `json:"environmentId,omitempty"`
	// The ID of the Application, if a management context has been defined to sync with an APIM instance
	ID string `json:"id,omitempty"`
	// Conditions describe the current conditions of the Application.
	//
	// Known condition types are:
	// * "Accepted"
	// * "ResolvedRefs"
	//
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=8
	// +kubebuilder:default={}
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// The processing status of the Application. *** DEPRECATED ***
	// The value is `Completed` if the sync with APIM succeeded, Failed otherwise.
	ProcessingStatus core.ProcessingStatus `json:"processingStatus,omitempty"`
	// The number of subscriptions that reference the application
	SubscriptionCount uint `json:"subscriptions,omitempty"`
	// When application has been created regardless of errors, this field is
	// used to persist the error message encountered during admission
	Errors status.Errors `json:"errors,omitempty"`
}
