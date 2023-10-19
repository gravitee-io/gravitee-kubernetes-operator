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

package v1alpha1

import (
	"context"
	"fmt"

	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/convert/v2"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/convert/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1beta1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

func (api *ApiDefinition) ConvertTo(hub conversion.Hub) error {
	dst, err := castHupApiDefinition(hub)
	if err != nil {
		return err
	}
	dst.Spec.Api = *v4.FromV2(&api.Spec.Api)
	dst.Spec.Context = api.Spec.Context
	dst.Status = *tov1Beta1ApiDefinitionStatus(api)
	dst.ObjectMeta = api.ObjectMeta
	return nil
}

func (api *ApiDefinition) ConvertFrom(hub conversion.Hub) error {
	dst, err := castHupApiDefinition(hub)
	if err != nil {
		return err
	}
	api.Spec.Api = *v2.FromV4(&dst.Spec.Api)
	api.Spec.Context = dst.Spec.Context
	api.Status = *toV1Alpha1ApiDefinitionStatus(dst.Status)
	api.ObjectMeta = dst.ObjectMeta
	return nil
}

func castHupApiDefinition(hub conversion.Hub) (*v1beta1.ApiDefinition, error) {
	dst, ok := hub.(*v1beta1.ApiDefinition)
	if !ok {
		return nil, fmt.Errorf("unable to read hub version as v1beta1")
	}
	return dst, nil
}

func toV1Alpha1ApiDefinitionStatus(status v1beta1.ApiDefinitionStatus) *ApiDefinitionStatus {
	return &ApiDefinitionStatus{
		ID:                           status.ID,
		CrossID:                      status.CrossID,
		EnvID:                        status.EnvID,
		OrgID:                        status.OrgID,
		State:                        status.State,
		Status:                       ProcessingStatus(status.Status),
		DeprecatedStatus:             ProcessingStatus(status.Status),
		ObservedGeneration:           status.ObservedGeneration,
		DeprecatedObservedGeneration: status.ObservedGeneration,
	}
}

func tov1Beta1ApiDefinitionStatus(api *ApiDefinition) *v1beta1.ApiDefinitionStatus {
	return setStatusPlans(api, &v1beta1.ApiDefinitionStatus{
		ID:                 api.Status.ID,
		CrossID:            api.Status.CrossID,
		EnvID:              api.Status.EnvID,
		OrgID:              api.Status.OrgID,
		State:              api.Status.State,
		Status:             v1beta1.ProcessingStatus(api.Status.Status),
		ObservedGeneration: api.Status.ObservedGeneration,
	})
}

func setStatusPlans(
	api *ApiDefinition, status *v1beta1.ApiDefinitionStatus,
) *v1beta1.ApiDefinitionStatus {
	ref := api.Spec.Context
	if ref == nil {
		log.Global.Debug("Context not found, status plans will not be updated")
		return status
	}

	if api.Status.ID == "" {
		log.Global.Debug("API ID not found, status plans will not be updated")
		return status
	}

	k8s := k8s.GetClient()
	ctx := context.Background()
	mCtx := new(ManagementContext)

	if err := k8s.Get(ctx, ref.ToK8sType(), mCtx); err != nil {
		log.Global.Error(err, "Unable to get management context, status plans will not be updated")
		return status
	}

	apim, err := apim.FromContext(ctx, mCtx.Spec.Context)
	if err != nil {
		log.Global.Error(err, "Unable to get apim client, status plans will not be updated")
		return status
	}

	plans, pErr := apim.Plans.ListByAPI(api.Status.ID)
	if errors.IgnoreNotFound(pErr) != nil {
		log.Global.Error(pErr, "Unable to list plans, status plans will not be updated")
		return status
	}

	status.Plans = make(map[string]string)
	for _, plan := range plans.Data {
		status.Plans[plan.Name] = plan.ID
	}

	return status
}
