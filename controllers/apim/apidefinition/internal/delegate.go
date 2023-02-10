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

	"github.com/go-logr/logr"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	coreV1 "k8s.io/api/core/v1"
	k8s "sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	separator = "/"

	bearerTokenSecretKey = "bearerToken"
	usernameSecretKey    = "username"
	passwordSecretKey    = "password"
)

type Delegate struct {
	ctx  context.Context
	k8s  k8s.Client
	log  logr.Logger
	apim *apim.APIM
}

func NewDelegate(ctx context.Context, k8s k8s.Client, log logr.Logger) *Delegate {
	return &Delegate{
		ctx, k8s, log, nil,
	}
}

func (d *Delegate) ResolveContext(api *gio.ApiDefinition) error {
	managementContext := new(gio.ManagementContext)

	ref := api.Spec.Context
	ns := ref.ToK8sType()

	d.log.Info("Resolving API context", "namespace", ref.Namespace, "name", ref.Name)

	if err := d.k8s.Get(d.ctx, ns, managementContext); err != nil {
		return err
	}

	if err := d.resolveContextSecrets(managementContext); err != nil {
		return err
	}

	apim, err := apim.FromContext(d.ctx, managementContext.Spec.Context)
	if err != nil {
		return err
	}

	d.apim = apim
	return nil
}

func (d *Delegate) HasContext() bool {
	return d.apim != nil
}

func (d *Delegate) AddDeletionFinalizer(api *gio.ApiDefinition) {
	if api.IsMissingDeletionFinalizer() {
		util.AddFinalizer(api, keys.ApiDefinitionDeletionFinalizer)
		if err := d.k8s.Update(d.ctx, api); err != nil {
			d.log.Error(err, "Unable to add deletion finalizer to API definition")
		}
	}
}

func (d *Delegate) resolveContextSecrets(context *gio.ManagementContext) error {
	management := context.Spec

	if management.HasSecretRef() {
		secret := new(coreV1.Secret)

		secretKey := management.SecretRef().ToK8sType()
		secretKey.Namespace = getSecretNamespace(context)

		if err := d.k8s.Get(d.ctx, secretKey, secret); err != nil {
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

func getSecretNamespace(context *gio.ManagementContext) string {
	secretRef := context.Spec.SecretRef()
	if secretRef.Namespace != "" {
		return secretRef.Namespace
	}
	return context.Namespace
}
