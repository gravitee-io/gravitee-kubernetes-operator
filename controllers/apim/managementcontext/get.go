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

package managementcontext

import (
	"context"
	"fmt"

	log "github.com/go-logr/logr"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Get(
	ctx context.Context,
	k8sClient client.Client,
	log log.Logger,
	contextRef *model.ContextRef,
) (*gio.ManagementContext, error) {
	apimContext := new(gio.ManagementContext)
	ns := types.NamespacedName{Name: contextRef.Name, Namespace: contextRef.Namespace}

	log.Info("Looking for context from", "namespace", contextRef.Namespace, "name", contextRef.Name)

	if err := k8sClient.Get(ctx, ns, apimContext); err != nil {
		return nil, err
	}

	if apimContext.HasSecretRef() {
		secret := new(coreV1.Secret)
		secretName := apimContext.Spec.Auth.SecretRef.Name
		secretNameSpace := apimContext.Spec.Auth.SecretRef.Namespace
		secretKey := types.NamespacedName{Name: secretName, Namespace: secretNameSpace}

		if err := k8sClient.Get(ctx, secretKey, secret); err != nil {
			return nil, err
		}

		bearerToken, ok := secret.StringData["token"]

		if !ok {
			return nil, fmt.Errorf("token not found in secret %s/%s", secretNameSpace, secretName)
		}

		apimContext.Spec.Auth.BearerToken = bearerToken
	}

	return apimContext, nil
}
