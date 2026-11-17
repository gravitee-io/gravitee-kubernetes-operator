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

package internal

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	apimgroup "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/group"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

func CreateOrUpdate(ctx context.Context, group *v1alpha1.Group) error {
	client, err := apimgroup.CreateAPIMClient(ctx, group)
	if err != nil {
		return err
	}

	dto, err := apimgroup.ToDTO(group)
	if err != nil {
		return err
	}

	resp, err := apimgroup.Upsert(ctx, client, dto)
	if err != nil {
		return gerrors.NewControlPlaneError(err)
	}

	return apimgroup.UpdateStatus(ctx, group, resp)
}
