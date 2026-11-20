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
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/apidefinition"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
)

func updateConfigMap(ctx context.Context, api *v1alpha1.ApiDefinition) error {
	if api.IsStopped() {
		if err := apidefinition.DeleteConfigMap(ctx, api.GetName(), api.GetNamespace()); err != nil {
			log.Error(ctx, err, "Unable to delete config map for API definition", log.KeyValues(api)...)
			return err
		}
		return nil
	}
	if err := saveConfigMap(ctx, api); err != nil {
		log.Error(ctx, err, "Unable to create or update config map for API definition", log.KeyValues(api)...)
		return err
	}
	return nil
}

func saveConfigMap(ctx context.Context, api *v1alpha1.ApiDefinition) error {
	definition, err := json.Marshal(&api.Spec)
	if err != nil {
		return err
	}
	gvk := api.GetObjectKind().GroupVersionKind()
	return apidefinition.SaveConfigMap(ctx, apidefinition.GatewayConfigMap{
		Name:      api.GetName(),
		Namespace: api.GetNamespace(),
		// Owned by the API, so that it is deleted with it.
		Owner: metav1.OwnerReference{
			Kind:       gvk.Kind,
			Name:       api.GetName(),
			APIVersion: gvk.GroupVersion().String(),
			UID:        api.GetUID(),
		},
		OrgID:      api.GetOrgID(),
		EnvID:      api.GetEnvID(),
		Definition: definition,
	})
}
