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

type Member struct {
	// Member user ID
	Id string `json:"id"`

	// Member source ID
	DisplayName string `json:"displayName,omitempty"`

	// The API role associated with this Member
	// +kubebuilder:default:={}
	Roles []Role `json:"roles,omitempty"`
}

type Role struct {
	// Name of the role (USER, REVIEWER ...)
	// +kubebuilder:default:=`USER`
	Name string `json:"name,omitempty"`

	// Role scope, by default it is API scope
	// +kubebuilder:default:=`API`
	Scope string `json:"scope,omitempty"`
}
