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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	apim "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	coreV1 "k8s.io/api/core/v1"
	k8s "sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	separator           = "/"
	defaultPlanSecurity = "KEY_LESS"
	defaultPlanStatus   = "PUBLISHED"
	defaultPlanName     = "G.K.O. Default"
	origin              = "kubernetes"
	mode                = "fully_managed"
)

type Delegate struct {
	ctx      context.Context
	k8s      k8s.Client
	log      logr.Logger
	contexts []DelegateContext
}

func NewDelegate(ctx context.Context, k8s k8s.Client, log logr.Logger) *Delegate {
	return &Delegate{
		ctx, k8s, log, make([]DelegateContext, 0),
	}
}

func (d *Delegate) ResolveContexts(api *gio.ApiDefinition) {
	contexts := api.Spec.Contexts
	for _, ref := range contexts {
		context, err := d.resolveContext(ref)

		if err != nil {
			d.log.Error(err, "Unable to resolve context ", "namespace", ref.Namespace, "name", ref.Name)
			continue
		}

		d.addContext(context)

		if api.IsMissingDeletionFinalizer() {
			util.AddFinalizer(api, keys.ApiDefinitionDeletionFinalizer)
			if err = d.k8s.Update(d.ctx, api); err != nil {
				d.log.Error(err, "Unable to add deletion finalizer to API definition", "namespace", api.Namespace, "name", api.Name)
			}
		}
	}
}

func (d *Delegate) HasContext() bool {
	return len(d.contexts) > 0
}

func (d *Delegate) resolveContext(
	ref model.NamespacedName,
) (*gio.ApiContext, error) {
	apiContext := new(gio.ApiContext)
	ns := ref.ToK8sType()

	d.log.Info("Resolving API context", "namespace", ref.Namespace, "name", ref.Name)

	if err := d.k8s.Get(d.ctx, ns, apiContext); err != nil {
		return nil, err
	}

	if err := d.resolveContextSecrets(apiContext); err != nil {
		return nil, err
	}

	return apiContext, nil
}

func (d *Delegate) resolveContextSecrets(context *gio.ApiContext) error {
	management := context.Spec.Management

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

func (d *Delegate) addContext(apiContext *gio.ApiContext) {
	spec := apiContext.Spec

	context := DelegateContext{
		Values:   spec.Values,
		Location: apiContext.Namespace + separator + apiContext.Name,
	}

	if spec.Management == nil {
		d.contexts = append(d.contexts, context)
		return
	}

	client, err := apim.NewClient(d.ctx, spec.Management)

	if err != nil {
		d.log.Error(err, "Unable to create management API client")
	} else {
		context.Client = client
	}

	d.contexts = append(d.contexts, context)
}

func getSecretNamespace(context *gio.ApiContext) string {
	secretRef := context.Spec.Management.SecretRef()
	if secretRef.Namespace != "" {
		return secretRef.Namespace
	}
	return context.Namespace
}
