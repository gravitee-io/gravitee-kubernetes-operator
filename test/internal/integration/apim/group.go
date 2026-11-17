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
	apimgroup "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/group"
)

// GetGroup gets a group state by HRID, server-managed fields (id, member count) included.
func (a *APIM) GetGroup(ctx context.Context, hrid string) (sdk.GroupState, error) {
	return apimgroup.GetState(ctx, a.SDK, apimgroup.DTO{Group: sdk.Group{Hrid: hrid}})
}

// ImportGroup creates or updates a group in APIM outside the operator.
func (a *APIM) ImportGroup(ctx context.Context, grp *v1alpha1.Group) error {
	dto, err := apimgroup.ToDTO(grp)
	if err != nil {
		return err
	}
	_, err = apimgroup.Upsert(ctx, a.SDK, dto)
	return err
}
