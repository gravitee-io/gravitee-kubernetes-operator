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

package authzpdp

// Endpoint is the AuthZEN endpoint of the PDP: an API listening on path and on the matching
// AuthZEN discovery path.
type Endpoint struct {
	// Context path of the endpoint: unique among the APIs of the environment, outside
	// /.well-known/authzen-configuration, and written the way APIM stores it, without a trailing
	// slash. A changed path replaces the endpoint.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Path string `json:"path"`
}

// Type defines the specification of an AuthzPdp resource: a PDP gateway of Authorization
// Management, the engine it runs, the gateways it runs on and its AuthZEN endpoint.
type Type struct {
	// Display name of the PDP gateway.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=255
	Name string `json:"name"`
	// The PDP engine the gateway runs. Every AuthzPdp of one targetPdpId shares the engine and its
	// strict value. Immutable.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^[A-Za-z0-9_.-]+$`
	// +kubebuilder:validation:MaxLength=255
	TargetPdpID string `json:"targetPdpId"`
	// Sharding tag of the gateways that run the PDP; omit it for an untagged gateway. Immutable.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`^-?[a-z0-9]+(-[a-z0-9]+)*$`
	// +kubebuilder:validation:MaxLength=64
	Tag string `json:"tag,omitempty"`
	// Whether the engine checks policies and entities against its schema. Defaults to false.
	// Shared by every PDP of the targetPdpId: on create it must agree with the registered ones, and
	// a change aligns them all with a warning naming each. Two AuthzPdp resources of one
	// targetPdpId that disagree overwrite each other on every reconcile.
	// +kubebuilder:validation:Optional
	Strict bool `json:"strict,omitempty"`
	// The AuthZEN endpoint of the PDP. Omitted, the PDP has none, and an existing one is deleted.
	// +kubebuilder:validation:Optional
	Endpoint *Endpoint `json:"endpoint,omitempty"`
}
