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

package apim

import (
	"context"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	spg "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/sharedpolicygroup"
)

// SharedPolicyGroups reads and writes shared policy groups in APIM outside the operator.
type SharedPolicyGroups struct {
	sdk *apim.SDKClient
}

// SharedPolicyGroupState is what the specs check of the state APIM answers, its cross id as the CR status holds it.
type SharedPolicyGroupState struct {
	Name    string
	CrossID string
}

// GetByHRID gets a shared policy group state by HRID.
func (s *SharedPolicyGroups) GetByHRID(hrid string) (*SharedPolicyGroupState, error) {
	state, err := spg.GetState(context.Background(), s.sdk, spg.DTO{SharedPolicyGroup: sdk.SharedPolicyGroup{Hrid: hrid}})
	if err != nil {
		return nil, err
	}
	crossID := ""
	if state.CrossId != nil {
		crossID = *state.CrossId
	}
	return &SharedPolicyGroupState{Name: state.Name, CrossID: crossID}, nil
}

// CreateOrUpdate creates or updates a shared policy group in APIM as the operator would send it.
func (s *SharedPolicyGroups) CreateOrUpdate(obj *v1alpha1.SharedPolicyGroup) (spg.Response, error) {
	dto, err := spg.ToDTO(obj)
	if err != nil {
		return spg.Response{}, err
	}
	return spg.Upsert(context.Background(), s.sdk, dto)
}
