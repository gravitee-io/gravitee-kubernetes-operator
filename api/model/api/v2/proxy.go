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

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
)

type VirtualHost struct {
	// Host name
	// +kubebuilder:validation:Optional
	Host string `json:"host,omitempty"`
	// Path
	// +kubebuilder:validation:Optional
	Path string `json:"path,omitempty"`
	// Indicate if Entrypoint should be overridden or not
	// +kubebuilder:validation:Optional
	OverrideEntrypoint bool `json:"override_entrypoint,omitempty"`
}

func NewVirtualHost(host, path string) *VirtualHost {
	return &VirtualHost{
		Host: host,
		Path: path,
	}
}

type Proxy struct {
	// list of Virtual hosts fot the proxy
	// +kubebuilder:validation:Optional
	VirtualHosts []*VirtualHost `json:"virtual_hosts"`
	// List of endpoint groups of the proxy
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	Groups []*EndpointGroup `json:"groups"`
	// Proxy Failover
	Failover *Failover `json:"failover,omitempty"`
	// Proxy Cors
	Cors *base.Cors `json:"cors,omitempty"`
	// Logging
	Logging *Logging `json:"logging,omitempty"`
	// Strip Context Path
	// +kubebuilder:validation:Optional
	StripContextPath *bool `json:"strip_context_path,omitempty"`
	// Preserve Host
	// +kubebuilder:validation:Optional
	PreserveHost *bool `json:"preserve_host,omitempty"`
}
