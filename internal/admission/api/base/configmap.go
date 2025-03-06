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

package base

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	"sigs.k8s.io/controller-runtime/pkg/client"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
Because this is the place where we have access to both
the old and the new object on updates, we use this hacky
side effect to delete the API definition config map for cases
when a switch is operated from syncing from the cluster to syncing
from the data store. The error is ignored and we do best effort here,
because some users might create RBAC by themselves and deny access to config map.
*/
func DeleteDefinitionConfigMapIfNeeded(
	ctx context.Context,
	oldAPI core.ApiDefinitionObject,
	newAPi core.ApiDefinitionObject,
) {
	if !oldAPI.IsSyncFromManagement() && newAPi.IsSyncFromManagement() {
		log.Debug(ctx, "deleting configmap following switch in sync mode")
		configMap := &coreV1.ConfigMap{
			ObjectMeta: metaV1.ObjectMeta{
				Name:      oldAPI.GetName(),
				Namespace: oldAPI.GetNamespace(),
			},
		}
		err := client.IgnoreNotFound(k8s.GetClient().Delete(ctx, configMap))
		if err != nil {
			log.Debug(ctx, err.Error())
		}
	}
}
