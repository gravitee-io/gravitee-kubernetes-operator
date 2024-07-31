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

import "github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"

type Status struct {
	// The organization ID, if a management context has been defined to sync with an APIM instance
	OrgID string `json:"organizationId,omitempty"`
	// The environment ID, if a management context has been defined to sync with an APIM instance
	EnvID string `json:"environmentId,omitempty"`
	// The ID of the API definition in the Gravitee API Management instance (if an API context has been configured).
	ID string `json:"id,omitempty"`
	// The Cross ID is used to identify an API that has been promoted from one environment to another.
	CrossID string `json:"crossId,omitempty"`
	// The processing status of the API definition.
	ProcessingStatus core.ProcessingStatus `json:"processingStatus,omitempty"`
	// The state of the API. Can be either STARTED or STOPPED.
	State ApiState `json:"state,omitempty"`
	// This is the object generation observed during the latest reconcile.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}
