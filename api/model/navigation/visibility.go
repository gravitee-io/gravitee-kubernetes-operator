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

package navigation

// Visibility tells whether a navigation entry of the next-gen portal is shown
// to anonymous visitors. Left unset, APIM resolves it: the entry inherits the
// visibility of its parent folder, and a root entry defaults to PUBLIC.
// +kubebuilder:validation:Enum=PUBLIC;PRIVATE
type Visibility string

const (
	Public  Visibility = "PUBLIC"
	Private Visibility = "PRIVATE"
)
