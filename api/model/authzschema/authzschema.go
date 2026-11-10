// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package authzschema

import "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"

// Type defines the specification of an AuthzSchema resource: a GAPL schema document of
// Authorization Management and the PDP gateway that enforces it.
type Type struct {
	// Display name of the schema document.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=255
	Name string `json:"name"`
	// The PDP gateway that enforces the document. One document per gateway: a target another
	// schema holds is refused. Omitted, the document is stored and published nowhere. Deleting
	// the AuthzPdp a pdpRef names removes the target on the platform.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule="!has(self.scope) || self.scope != '*'",message="scope cannot be '*'"
	Target *authzpdp.Target `json:"target,omitempty"`
	// The GAPL schema text, replaced whole on every apply.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Schema string `json:"schema"`
}
