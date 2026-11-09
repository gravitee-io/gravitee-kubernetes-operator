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

// AuthzSchemaDTO is the automation API wire representation of an AuthzSchema (PUT /authz/schemas).
// The platform stores schema and target as sent, so both compare exactly.
type AuthzSchemaDTO struct {
	HRID   string `json:"hrid,omitempty" drift:"ignore"`
	Name   string `json:"name"`
	Target string `json:"target,omitempty"`
	Schema string `json:"schema"`
}

// AuthzSchemaState is what a PUT or a GET answers: the payload and the envelope.
type AuthzSchemaState struct {
	AuthzSchemaDTO `json:",inline"`
	ID             string        `json:"id,omitempty"`
	OrgID          string        `json:"organizationId,omitempty"`
	EnvID          string        `json:"environmentId,omitempty"`
	Errors         status.Errors `json:"errors,omitempty"`
}

// ToAuthzSchemaDTO maps a schema whose refs dynamic.ResolveAuthzSchemaRefs has resolved onto the
// wire payload: its target is a scope, or absent. It is the single mapping used by the sync path,
// the dry run and drift detection, and reads nothing from the cluster. An unresolved pdpRef would
// go out as no target: callers resolve first.
func ToAuthzSchemaDTO(crd *v1alpha1.AuthzSchema) AuthzSchemaDTO {
	target := ""
	if crd.Spec.Target != nil {
		target = crd.Spec.Target.Scope
	}
	return AuthzSchemaDTO{
		HRID:   refs.NewNamespacedNameFromObject(crd).HRID(),
		Name:   crd.Spec.Name,
		Target: target,
		Schema: crd.Spec.Schema,
	}
}
