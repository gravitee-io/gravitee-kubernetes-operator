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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
)

func CreateOrUpdate(ctx context.Context, idpgroupmapping *v1alpha1.IDPGroupMapping) error {
	ns := idpgroupmapping.Namespace
	spec := idpgroupmapping.Spec

	apim, err := apim.FromContextRef(ctx, idpgroupmapping.ContextRef(), ns)
	if err != nil {
		return err
	}

	idpgroupmapping.PopulateIDs(apim.Context)

	// Step 1, fetch the IDP configuration
	idpConfig, err := apim.Configuration.GetIDPConfiguration(idpgroupmapping.Spec.IDPID)
	if err != nil {
		return gerrors.NewControlPlaneError(err)
	}

	// Step 2, Extract the group mapping configuration from the IDP configuration
	groupMappingConfig := idpConfig.GroupMappings
	updated := false

	// Step 3, Resolve the group names to group IDs
	resolvedGroups := []string{}
	for _, group := range spec.Groups {
		// If the group is not an ID, we try to resolve it as a name
		resolvedGroup, err := apim.Env.FindGroup(group)
		
		if err != nil {
			return gerrors.NewControlPlaneError(err)
		}
		if resolvedGroup == nil || resolvedGroup.ID == "" {
			log.Error(ctx, nil, "unable to resolve group, skipping", "group", group, "mapping", idpgroupmapping.Name)
			continue
		}
		resolvedGroups = append(resolvedGroups, resolvedGroup.ID)
	}
	spec.Groups = resolvedGroups

	// Step 4, Check if the group mapping configuration is nil, if yes, initialize it
	if groupMappingConfig == nil {
		groupMappingConfig = []model.GroupMapping{
			{
				Condition: spec.Condition,
				Groups:    spec.Groups,
			},
		}
		updated = true
	} else if !containsGroupMapping(groupMappingConfig, spec.Condition, spec.Groups) {
		groupMappingConfig = append(groupMappingConfig, model.GroupMapping{
			Condition: spec.Condition,
			Groups:    spec.Groups,
		})
		updated = true
	}

	// Step 5, Update the IDP configuration with the new group mapping configuration if needed
	if updated {
		idpConfig.GroupMappings = groupMappingConfig
		if err := apim.Configuration.UpdateIDPConfiguration(idpConfig); err != nil {
			return gerrors.NewControlPlaneError(err)
		}
	}

	idpgroupmapping.Status.ID = idpConfig.ID
	idpgroupmapping.Status.OrgID = apim.Context.GetOrgID()
	idpgroupmapping.Status.EnvID = apim.Context.GetEnvID()

	return nil
}