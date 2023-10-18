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
	"fmt"

	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/convert/v2"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/convert/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("apidefinition/conversion")

func (api *ApiDefinition) ConvertTo(hub conversion.Hub) error {
	log.Info("converting v1alpha1 to v1beta1")
	dst, err := toV1Beta1(hub)
	if err != nil {
		return err
	}
	dst.Spec.Api = *v4.FromV2(&api.Spec.Api)
	dst.Status = tov1Beta1Status(api.Status)
	dst.ObjectMeta = api.ObjectMeta
	return nil
}

func (api *ApiDefinition) ConvertFrom(hub conversion.Hub) error {
	log.Info("converting v1beta1 to v1alpha1")
	dst, err := toV1Beta1(hub)
	if err != nil {
		return err
	}
	api.Spec.Api = *v2.FromV4(&dst.Spec.Api)
	api.Status = toV1Alpha1Status(dst.Status)
	api.ObjectMeta = dst.ObjectMeta
	return nil
}

func toV1Beta1(hub conversion.Hub) (*v1beta1.ApiDefinition, error) {
	dst, ok := hub.(*v1beta1.ApiDefinition)
	if !ok {
		return nil, fmt.Errorf("unable to read hub version as v1beta1")
	}
	return dst, nil
}

func toV1Alpha1Status(status v1beta1.ApiDefinitionStatus) ApiDefinitionStatus {
	return ApiDefinitionStatus{
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

func tov1Beta1Status(status ApiDefinitionStatus) v1beta1.ApiDefinitionStatus {
	return v1beta1.ApiDefinitionStatus{
		ID:                 status.ID,
		CrossID:            status.CrossID,
		EnvID:              status.EnvID,
		OrgID:              status.OrgID,
		State:              status.State,
		Status:             v1beta1.ProcessingStatus(status.Status),
		ObservedGeneration: status.ObservedGeneration,
	}
}
