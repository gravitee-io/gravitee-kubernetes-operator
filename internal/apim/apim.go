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

package apim

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/management"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/service"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/template"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/types/k8s/custom"

	coreV1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	bearerTokenSecretKey = "bearerToken"
	usernameSecretKey    = "username"
	passwordSecretKey    = "password"
)

// APIM wraps services needed to sync resources with a given environment on a Gravitee.io APIM instance.
type APIM struct {
	APIs         *service.APIs
	Applications *service.Applications

	Context *management.Context
}

// EnvID returns the environment ID of the current managed APIM instance.
func (apim *APIM) EnvID() string {
	return apim.Context.EnvId
}

// OrgID returns the organization ID of the current managed APIM instance.
func (apim *APIM) OrgID() string {
	return apim.Context.OrgId
}

// FromContext returns a new APIM instance from a given reconcile context and management context.
func FromContext(ctx context.Context, managementContext *management.Context) (*APIM, error) {
	orgID, envID := managementContext.OrgId, managementContext.EnvId
	urls, err := client.NewURLs(managementContext.BaseUrl, orgID, envID)
	if err != nil {
		return nil, err
	}

	client := &client.Client{
		HTTP: http.NewClient(ctx, toHttpAuth(managementContext)),
		URLs: urls,
	}

	return &APIM{
		APIs:         service.NewAPIs(client),
		Applications: service.NewApplications(client),
		Context:      managementContext,
	}, nil
}

func FromContextRef(ctx context.Context, ref custom.ResourceRef) (*APIM, error) {
	managementContext, err := resolveManagementContext(ctx, ref)
	if err != nil {
		return nil, err
	}
	return FromContext(ctx, managementContext)
}

func resolveManagementContext(
	ctx context.Context,
	ref custom.ResourceRef,
) (*management.Context, error) {
	managementContext := new(v1alpha1.ManagementContext)

	nsm := ref.NamespacedName()

	log.FromContext(ctx).Info("Resolving management context", "namespace", nsm.Namespace, "name", nsm.Name)

	if err := k8s.GetClient().Get(ctx, nsm, managementContext); err != nil {
		return nil, err
	}

	if err := template.Compile(ctx, managementContext); err != nil {
		return nil, err
	}

	if err := resolveContextSecrets(ctx, managementContext); err != nil {
		return nil, err
	}

	return managementContext.Spec.Context, nil
}

func resolveContextSecrets(ctx context.Context, context *v1alpha1.ManagementContext) error {
	management := context.Spec

	if management.HasSecretRef() {
		secret := new(coreV1.Secret)

		secretKey := management.SecretRef().NamespacedName()
		secretKey.Namespace = getSecretNamespace(context)

		if err := k8s.GetClient().Get(ctx, secretKey, secret); err != nil {
			return err
		}

		bearerToken := string(secret.Data[bearerTokenSecretKey])
		username := string(secret.Data[usernameSecretKey])
		password := string(secret.Data[passwordSecretKey])

		management.SetToken(bearerToken)
		management.SetCredentials(username, password)
	}

	return nil
}

func getSecretNamespace(context *v1alpha1.ManagementContext) string {
	secretRef := context.Spec.SecretRef()
	if secretRef.Namespace != "" {
		return secretRef.Namespace
	}
	return context.Namespace
}
