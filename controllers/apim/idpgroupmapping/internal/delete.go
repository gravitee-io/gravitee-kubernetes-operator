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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
)

func Delete(ctx context.Context, idpgroupmapping *v1alpha1.IDPGroupMapping) error {
	if !util.ContainsFinalizer(idpgroupmapping, core.IDPGroupMappingFinalizer) {
		return nil
	}

	ns := idpgroupmapping.Namespace

	apim, err := apim.FromContextRef(ctx, idpgroupmapping.ContextRef(), ns)
	if err != nil {
		return err
	}

	// Step 1, fetch the IDP configuration
	idpConfig, err := apim.Configuration.GetIDPConfiguration(idpgroupmapping.Spec.IDPID)
	if err != nil {
		return gerrors.NewControlPlaneError(err)
	}

	// Step 2, Extract the group mapping configuration from the IDP configuration
	groupMappingConfig := idpConfig.GroupMappings

	// Step 3, Resolve the group names to group IDs
	resolvedGroups := []string{}
	for _, group := range idpgroupmapping.Spec.Groups {
		// If the group is not an ID, we try to resolve it as a name
		resolvedGroup, err := apim.Env.FindGroup(group)
		
		if err != nil {
			return gerrors.NewControlPlaneError(err)
		}
		if resolvedGroup == nil {
			log.Error(ctx, nil, "unable to resolve group, skipping", "group", group, "mapping", idpgroupmapping.Name)
			continue
		}
		resolvedGroups = append(resolvedGroups, resolvedGroup.ID)
	}
	idpgroupmapping.Spec.Groups = resolvedGroups

	// Step 4,Check if the group mapping exists, if yes, remove it from the configuration
	if groupMappingConfig == nil || !containsGroupMapping(groupMappingConfig, idpgroupmapping.Spec.Condition, idpgroupmapping.Spec.Groups) {
		return nil
	} 

	// Step 5, generate a new group mapping configuration without the group mapping to delete
	newGroupMappingConfig := make([]model.GroupMapping, 0)
	for _, gm := range groupMappingConfig {
		if gm.Condition == idpgroupmapping.Spec.Condition && equalStringSlices(gm.Groups, idpgroupmapping.Spec.Groups) {
			continue
		}
		newGroupMappingConfig = append(newGroupMappingConfig, gm)
	}
	idpConfig.GroupMappings = newGroupMappingConfig
	log.Info(ctx, "GroupMapping", idpConfig)

	// Step 6, Update the IDP configuration with the new group mapping configuration
	return apim.Configuration.UpdateIDPConfiguration(idpConfig)
}
