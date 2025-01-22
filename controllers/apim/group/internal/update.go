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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
)

func CreateOrUpdate(ctx context.Context, group *v1alpha1.Group) error {
	ns := group.Namespace
	spec := group.Spec

	apim, err := apim.FromContextRef(ctx, group.ContextRef(), ns)
	if err != nil {
		return err
	}

	group.PopulateIDs(apim.Context)

	status, err := apim.Env.ImportGroup(spec.Type)
	if err != nil {
		return err
	}

	group.Status.ID = status.ID
	group.Status.OrgID = apim.Context.GetOrgID()
	group.Status.EnvID = apim.Context.GetEnvID()
	group.Status.Members = status.Members

	return nil
}
