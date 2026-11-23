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

package application

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimapp "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

func mergeDriftValidation(
	ctx context.Context,
	oldApp *v1alpha1.Application,
	newApp *v1alpha1.Application,
	errs *errors.AdmissionErrors,
) {
	errs.MergeWith(drift.ValidateDrift(ctx, oldApp, newApp, resolveAppRefs, remoteAppGetter(ctx),
		drift.MapDTO(func(app *v1alpha1.Application) apimapp.DTO {
			// ToDTO only fails on a spec the dry run refuses as severe, and the dry run runs first.
			dto, _ := apimapp.ToDTO(app)
			return dto
		}),
	))
}

func resolveAppRefs(ctx context.Context, app *v1alpha1.Application) error {
	return apimapp.ResolveClientCertificates(ctx, app.Spec.Settings, app.GetNamespace(), app.GetName())
}

// remoteAppGetter fetches the application through the SDK client of the context the application references;
// the HTTP client drift hands over is not used.
func remoteAppGetter(ctx context.Context) drift.RemoteObjectGetter[*v1alpha1.Application] {
	return func(_ *apim.APIM, app *v1alpha1.Application) (any, error) {
		dto, err := apimapp.ToDTO(app)
		if err != nil {
			return nil, err
		}
		client, err := apimapp.CreateAPIMClient(ctx, app)
		if err != nil {
			return nil, err
		}
		return apimapp.GetRemote(ctx, client, dto)
	}
}
