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

package policygroups

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/status"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
)

type Status struct {
	// The organization ID, if a management context has been defined to sync with an APIM instance
	OrgID string `json:"organizationId,omitempty"`
	// The environment ID, if a management context has been defined to sync with an APIM instance
	EnvID string `json:"environmentId,omitempty"`
	// The Cross ID is used to identify an SharedPolicyGroup that has been promoted from one environment to another.
	CrossID string `json:"crossId,omitempty"`
	// The ID is used to identify an SharedPolicyGroup which is unique in any environment.
	ID string `json:"id,omitempty"`
	// The processing status of the SharedPolicyGroup.
	// The value is `Completed` if the sync with APIM succeeded, Failed otherwise.
	ProcessingStatus core.ProcessingStatus `json:"processingStatus,omitempty"`
	// When SharedPolicyGroup has been created regardless of errors, this field is
	// used to persist the error message encountered during admission
	Errors status.Errors `json:"errors,omitempty"`
}
