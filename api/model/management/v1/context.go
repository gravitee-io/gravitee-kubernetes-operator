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

// +kubebuilder:object:generate=true
package v1

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/management/base"
)

type Context struct {
	base.Context `json:",inline"`

	// The Gravitee APIM organization targeted by the management context.
	// +kubebuilder:validation:Required
	OrgID string `json:"organizationId"`
}
