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
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implie
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"context"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/apidefinition"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/apiv4"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
)

func CreateOrUpdate(ctx context.Context, apiDefinition client.Object) error {
	switch t := apiDefinition.(type) {
	case *v1alpha1.ApiDefinition:
		return createOrUpdateV2(ctx, t)
	case *v1alpha1.ApiV4Definition:
		return createOrUpdateV4(ctx, t)
	default:
		return fmt.Errorf("unknown type %T", t)
	}
}

func createOrUpdateV2(ctx context.Context, apiDefinition *v1alpha1.ApiDefinition) error {
	spec := &apiDefinition.Spec
	spec.EnsureDefinitionContext()

	if err := apidefinition.ResolveResources(ctx, spec.Resources, true); err != nil {
		return err
	}

	groups, err := apidefinition.ResolveGroupRefs(ctx, apiDefinition, apiDefinition.GetGroupRefs())
	if err != nil {
		return err
	}
	spec.Groups = append(spec.Groups, groups...)

	apiDefinition.PopulateIDs(nil, false)

	err = apidefinition.ResolveConsoleNotificationRefs(ctx, apiDefinition)
	if err != nil {
		return err
	}

	if !apiDefinition.HasContext() {
		if !spec.IsLocal {
			return gerrors.NewIllegalStateError(
				gerrors.NewUnrecoverableError("a context is required when setting local to false"))
		}
		if err := apidefinition.UpdateConfigMap(ctx, apiDefinition); err != nil {
			return err
		}
		apiDefinition.Status.State = spec.State
		apiDefinition.Status.ID = spec.ID
		return nil
	}

	log.Debug(ctx, "Syncing API definition with control plane", log.KeyValues(apiDefinition)...)

	apimClient, apimErr := apim.FromContextRef(ctx, spec.Context, apiDefinition.GetNamespace())
	if apimErr != nil {
		return apimErr
	}

	status, mgmtErr := apimClient.APIs.ImportV2(&spec.Api)
	if mgmtErr != nil {
		return gerrors.NewControlPlaneError(mgmtErr)
	}

	status.DeepCopyTo(&apiDefinition.Status.Status)

	if spec.IsLocal {
		retrieveMgmtPlanIds(spec, status)
		return apidefinition.UpdateConfigMap(ctx, apiDefinition)
	}

	log.Debug(ctx, "API successfully synced with control plane", log.KeyValues(apiDefinition)...)

	return nil
}

func createOrUpdateV4(ctx context.Context, apiDefinition *v1alpha1.ApiV4Definition) error {
	apidefinition.LogUnsyncedGroupRefs(ctx, apiDefinition, apiDefinition.GetGroupRefs())

	if err := apidefinition.PrepareV4SpecForAutomation(ctx, apiDefinition, true); err != nil {
		log.Error(ctx, err, "Unable to prepare API spec for automation", log.KeyValues(apiDefinition)...)
		return err
	}

	client, err := apiv4.CreateAPIMClient(ctx, apiDefinition)
	if err != nil {
		return err
	}
	dto, err := apiv4.ToDTO(apiDefinition)
	if err != nil {
		return err
	}

	if client != nil {
		log.Debug(ctx, "Syncing API definition with control plane", log.KeyValues(apiDefinition)...)
	}
	resp, err := apiv4.Upsert(ctx, client, dto)
	if accepted, ok := apiv4.AcceptedResponse(err); ok {
		// APIM accepted the API: keep its ID and plans in the status, then report the ConfigMap error.
		if stErr := apiv4.UpdateStatus(ctx, apiDefinition, accepted); stErr != nil {
			return stErr
		}
		return err
	}
	if err != nil {
		return gerrors.NewControlPlaneError(err)
	}
	if client != nil {
		log.Debug(ctx, "API successfully synced with control plane", log.KeyValues(apiDefinition)...)
	}
	return apiv4.UpdateStatus(ctx, apiDefinition, resp)
}

func retrieveMgmtPlanIds(spec *v1alpha1.ApiDefinitionV2Spec, status *base.Status) {
	plans := spec.Plans
	planIds := status.Plans

	for _, plan := range plans {
		plan.ID = planIds[plan.CrossID]
		plan.Api = &status.ID
	}
}
