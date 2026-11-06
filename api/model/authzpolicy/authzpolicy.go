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

package authzpolicy

import "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"

// Type defines the specification of an AuthzPolicy resource: one GAPL policy of Authorization
// Management and the PDP gateways it is deployed to.
type Type struct {
	// Id of the policy in the engine and in the Console: policy. then one or more segments
	// separated by a dot. One policy of the environment holds it: an id another policy holds is
	// refused, naming that policy. Immutable.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^policy(\.[A-Za-z0-9_:-]+)+$`
	// +kubebuilder:validation:MaxLength=255
	EntityID string `json:"entityId"`
	// Display name of the policy in the Console.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=255
	Name string `json:"name"`
	// What the policy governs, as the Console groups policies. Immutable.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=MCP;AGENT;MODEL;API;CUSTOM
	PolicyType string `json:"policyType"`
	// The GAPL text, one or more permit and forbid statements, replaced whole on every apply. A
	// text that does not parse, or a template (?principal, ?resource), is refused.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	PolicyText string `json:"policyText"`
	// The PDP gateways the policy is deployed to. Empty, it is deployed to the default engine on
	// every gateway, the only engine an MCP proxy asks. Deleting the AuthzPdp a pdpRef names removes
	// its scope here; the platform refuses that delete when the policy would be left with no target.
	// +kubebuilder:validation:Optional
	Targets []authzpdp.Target `json:"targets,omitempty"`
}
