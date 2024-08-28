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
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
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
	cp := apiDefinition.DeepCopy()
	spec := &cp.Spec
	spec.EnsureDefinitionContext()

	if err := resolveResources(ctx, spec.Resources); err != nil {
		return err
	}

	cp.PopulateIDs(nil)

	if !apiDefinition.HasContext() {
		if !spec.IsLocal {
			return errors.NewUnrecoverableError("a context is required when setting local to false")
		}
		if err := updateConfigMap(ctx, cp); err != nil {
			return err
		}
		apiDefinition.Status.State = spec.State
		apiDefinition.Status.ID = spec.ID
		return nil
	}

	log.FromContext(ctx).Info("Syncing API with APIM")

	apimClient, apimErr := apim.FromContextRef(ctx, spec.Context, apiDefinition.GetNamespace())
	if apimErr != nil {
		return apimErr
	}

	status, mgmtErr := apimClient.APIs.ImportV2(&spec.Api)
	if mgmtErr != nil {
		return errors.NewContextError(mgmtErr)
	}

	apiDefinition.Status = v1alpha1.ApiDefinitionStatus{
		Status: *status,
	}

	if spec.IsLocal {
		return updateConfigMap(ctx, cp)
	}

	if err := deleteConfigMap(ctx, apiDefinition); err != nil {
		return err
	}

	return nil
}

func createOrUpdateV4(ctx context.Context, apiDefinition *v1alpha1.ApiV4Definition) error {
	cp := apiDefinition.DeepCopy()

	spec := &cp.Spec

	if err := resolveResources(ctx, spec.Resources); err != nil {
		log.FromContext(ctx).Error(err, "Unable to resolve API resources from references")
		return err
	}

	spec.DefinitionContext = v4.NewDefaultKubernetesContext().MergeWith(spec.DefinitionContext)

	if spec.Context != nil {
		log.FromContext(ctx).Info("Syncing API with APIM")
		apimClient, err := apim.FromContextRef(ctx, spec.Context, apiDefinition.GetNamespace())
		if err != nil {
			return err
		}
		cp.PopulateIDs(apimClient.Context)
		status, err := apimClient.APIs.ImportV4(&spec.Api)
		if err != nil {
			return err
		}
		apiDefinition.Status.Status = *status
		log.FromContext(ctx).WithValues("id", spec.ID).Info("API successfully synced with APIM")
	} else {
		cp.PopulateIDs(nil)
	}

	if spec.DefinitionContext.SyncFrom == v4.OriginManagement || spec.State == base.StateStopped {
		log.FromContext(ctx).Info(
			"Deleting config map as API is not managed by operator or is stopped",
			"syncFrom", spec.DefinitionContext.SyncFrom,
			"state", spec.State,
		)
		if err := deleteConfigMap(ctx, cp); err != nil {
			return err
		}
	} else {
		log.FromContext(ctx).Info("Saving config map")
		if err := saveConfigMap(ctx, cp); err != nil {
			return err
		}
	}
	return nil
}
