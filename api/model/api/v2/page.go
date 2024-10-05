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

import "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"

type Page struct {
	*base.Page `json:",inline"`
	// If the page is private, defines a set of user groups with access
	// +kubebuilder:validation:Optional
	AccessControls []base.AccessControl `json:"accessControls"`
	// if true, the references defined in the accessControls list will be
	// denied access instead of being granted
	// +kubebuilder:validation:Optional
	ExcludedAccessControl bool `json:"excludedAccessControls"`
}
