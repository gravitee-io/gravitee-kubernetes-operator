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

package model

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/status"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

// AuthzPolicyDTO is the automation API wire representation of an AuthzPolicy (PUT /authz/policies).
// The platform stores and answers policyText verbatim, so it compares exactly; it sorts targets
// and answers [] when there are none, so they do not compare by position.
type AuthzPolicyDTO struct {
	HRID       string   `json:"hrid,omitempty" drift:"ignore"`
	EntityID   string   `json:"entityId"`
	Name       string   `json:"name"`
	PolicyType string   `json:"policyType"`
	PolicyText string   `json:"policyText"`
	Targets    []string `json:"targets,omitempty" drift:"set"`
}

// AuthzPolicyState is what a PUT or a GET answers: the payload, the envelope and the deployment.
// PolicyStatus sits beside the payload, so drift detection does not compare it: every apply
// deploys, and no spec can declare a disabled policy.
type AuthzPolicyState struct {
	AuthzPolicyDTO `json:",inline"`
	ID             string        `json:"id,omitempty"`
	OrgID          string        `json:"organizationId,omitempty"`
	EnvID          string        `json:"environmentId,omitempty"`
	Errors         status.Errors `json:"errors,omitempty"`
	PolicyStatus   string        `json:"policyStatus,omitempty"`
}

// ToAuthzPolicyDTO maps a policy whose refs dynamic.ResolveAuthzPolicyRefs has resolved onto the
// wire payload: each target is a scope, sent in declaration order. It is the single mapping used
// by the sync path, the dry run and drift detection, and reads nothing from the cluster.
func ToAuthzPolicyDTO(crd *v1alpha1.AuthzPolicy) AuthzPolicyDTO {
	var targets []string
	for _, target := range crd.Spec.Targets {
		targets = append(targets, target.Scope)
	}
	return AuthzPolicyDTO{
		HRID:       refs.NewNamespacedNameFromObject(crd).HRID(),
		EntityID:   crd.Spec.EntityID,
		Name:       crd.Spec.Name,
		PolicyType: crd.Spec.PolicyType,
		PolicyText: crd.Spec.PolicyText,
		Targets:    targets,
	}
}
