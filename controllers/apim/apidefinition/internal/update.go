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
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
)

type apiNamespaceKey string

const nsKey = apiNamespaceKey("api-namespace")

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
	cp := apiDefinition.DeepCopy()
	spec := &cp.Spec
	spec.EnsureDefinitionContext()

	if err := resolveResources(ctx, spec.Resources); err != nil {
		return err
	}

	if err := ResolveGroupRefs(ctx, cp); err != nil {
		return err
	}

	cp.PopulateIDs(nil)

	err := ResolveConsoleNotificationRefs(ctx, cp)
	if err != nil {
		return err
	}

	if !apiDefinition.HasContext() {
		if !spec.IsLocal {
			return gerrors.NewIllegalStateError(
				gerrors.NewUnrecoverableError("a context is required when setting local to false"))
		}
		if err := updateConfigMap(ctx, cp); err != nil {
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

	apiDefinition.Status = v1alpha1.ApiDefinitionStatus{
		Status: *status,
	}

	if spec.IsLocal {
		retrieveMgmtPlanIds(spec, status)
		return updateConfigMap(ctx, cp)
	}

	log.Debug(ctx, "API successfully synced with control plane", log.KeyValues(apiDefinition)...)

	return nil
}

func createOrUpdateV4(ctx context.Context, apiDefinition *v1alpha1.ApiV4Definition) error {
	cp := apiDefinition.DeepCopy()
	nsCtx := context.WithValue(ctx, nsKey, apiDefinition.Namespace)

	spec := &cp.Spec

	if err := resolveResources(ctx, spec.Resources); err != nil {
		log.Error(ctx, err, "Unable to resolve API resources from references", log.KeyValues(apiDefinition)...)
		return err
	}

	if err := resolveSharedPolicyGroups(nsCtx, spec); err != nil {
		log.Error(ctx, err, "Unable to resolve API resources from references", log.KeyValues(apiDefinition)...)
		return err
	}

	if err := ResolveGroupRefs(ctx, cp); err != nil {
		return err
	}

	spec.DefinitionContext = v4.NewDefaultKubernetesContext().MergeWith(spec.DefinitionContext)

	if spec.Context != nil {
		log.Debug(ctx, "Syncing API definition with control plane", log.KeyValues(apiDefinition)...)
		apimClient, err := apim.FromContextRef(ctx, spec.Context, apiDefinition.GetNamespace())
		if err != nil {
			return err
		}
		cp.PopulateIDs(apimClient.Context)

		err = ResolveConsoleNotificationRefs(ctx, cp)
		if err != nil {
			return err
		}

		status, err := apimClient.APIs.ImportV4(&spec.Api)

		if err != nil {
			return gerrors.NewControlPlaneError(err)
		}
		apiDefinition.Status.Status = *status
		log.Debug(ctx, "API successfully synced with control plane", log.KeyValues(apiDefinition)...)
	} else {
		cp.PopulateIDs(nil)
	}

	if spec.DefinitionContext.SyncFrom == v4.OriginKubernetes {
		return updateConfigMap(ctx, cp)
	}

	return nil
}

// Retrieve the plan ids from the API CRD status.
func retrieveMgmtPlanIds(spec *v1alpha1.ApiDefinitionV2Spec, status *base.Status) {
	plans := spec.Plans
	planIds := status.Plans

	for _, plan := range plans {
		plan.ID = planIds[plan.CrossID]
		plan.Api = &status.ID
	}
}
