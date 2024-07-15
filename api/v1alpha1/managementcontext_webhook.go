/*
 * Copyright (C) 2015 The Gravitee team (http://gravitee.io)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha1

import (
	"context"
	"net"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/uuid"
	errors "github.com/pkg/errors"
	runtime "k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var _ webhook.Defaulter = &ManagementContext{}
var _ webhook.Validator = &ManagementContext{}

func (context *ManagementContext) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(context).
		Complete()
}

func (context *ManagementContext) Default() {}

func (context *ManagementContext) ValidateCreate() (admission.Warnings, error) {
	if err := resolveManagementContext(context); err != nil {
		return admission.Warnings{}, err
	}

	return admission.Warnings{}, nil
}

func (context *ManagementContext) ValidateUpdate(_ runtime.Object) (admission.Warnings, error) {
	return admission.Warnings{}, nil
}

func (*ManagementContext) ValidateDelete() (admission.Warnings, error) {
	return admission.Warnings{}, nil
}

func resolveManagementContext(ctx *ManagementContext) error {
	urLs, _ := client.NewURLs(ctx.Spec.BaseUrl, ctx.Spec.OrgId, ctx.Spec.EnvId)

	httpClient := http.NewClient(context.Background(), nil)
	cli := client.Client{
		HTTP: httpClient,
		URLs: urLs,
	}

	api := make(map[string]interface{})
	err := httpClient.Get(cli.EnvV1Target("apis").WithPath(uuid.NewV4String()).String(), api)

	var opError *net.OpError
	if errors.As(err, &opError) {
		return err
	}

	return nil
}
