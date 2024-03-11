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

type Member struct {
	// +kubebuilder:validation:Required
	// Member ID
	Id string `json:"id,omitempty"`
	// +kubebuilder:validation:Optional
	// Member email
	Email string `json:"email,omitempty"`
	// +kubebuilder:validation:Optional
	// Member Display Name
	DisplayName string `json:"displayName,omitempty"`
	// +kubebuilder:validation:Required
	// Member type
	Type string `json:"type,omitempty"`
}
