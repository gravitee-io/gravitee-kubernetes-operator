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

package model

import (
	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/group"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/status"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/mapper"
)

// ToGroup maps a Group CRD to the Automation API group payload. SDK defaults are not applied:
// they differ from the CRD's and would defeat the drift equivalences. The CRD id is not part of
// the payload; a legacy group sends it as its hrid (see the Env service).
func ToGroup(grp *v1alpha1.Group) (sdk.Group, error) {
	return mapper.MapViaJSON[sdk.Group](grp.Spec.Type)
}

// GroupFromState decodes a group state, as returned by a GET or a PUT, into the payload type,
// so that drift compares two values of the same type. State-only fields are dropped.
func GroupFromState(state sdk.GroupState) (sdk.Group, error) {
	return mapper.MapViaJSON[sdk.Group](state)
}

// ToGroupStatus maps the group state APIM answers to the CR status fields it sets.
func ToGroupStatus(state sdk.GroupState) *group.Status {
	result := &group.Status{}
	if state.Id != nil {
		result.ID = *state.Id
	}
	if state.MemberCount != nil && *state.MemberCount > 0 {
		result.Members = uint(*state.MemberCount)
	}
	if state.Errors != nil {
		result.Errors = status.Errors{
			Severe:  state.Errors.Severe,
			Warning: state.Errors.Warning,
		}
	}
	return result
}
