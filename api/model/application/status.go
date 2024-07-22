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

import "github.com/gravitee-io/gravitee-kubernetes-operator/pkg/types/k8s/custom"

type Status struct {
	// The organization ID, if a management context has been defined to sync with an APIM instance
	OrgID string `json:"organizationId,omitempty"`
	// The environment ID, if a management context has been defined to sync with an APIM instance
	EnvID string `json:"environmentId,omitempty"`
	// The ID of the Application, if a management context has been defined to sync with an APIM instance
	ID string `json:"id,omitempty"`
	// The processing status of the Application.
	// The value is `Completed` if the sync with APIM succeeded, Failed otherwise.
	ProcessingStatus custom.ProcessingStatus `json:"processingStatus,omitempty"`
	// This is the object generation observed during the latest reconcile.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}
