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
	appResolve "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

func mergeDriftValidation(
	ctx context.Context,
	oldApp *v1alpha1.Application,
	newApp *v1alpha1.Application,
	errs *errors.AdmissionErrors,
) {
	errs.MergeWith(drift.ValidateDrift(ctx, oldApp, newApp, resolveAppRefs, getRemoteApp,
		drift.MapDTO(func(app *v1alpha1.Application) model.ApplicationDTO {
			return model.ToApplicationDTO(app.Spec)
		}),
	))
}

func resolveAppRefs(ctx context.Context, app *v1alpha1.Application) error {
	return appResolve.ResolveClientCertificates(ctx, app.Spec.Settings, app.GetNamespace(), app.GetName())
}

func getRemoteApp(apimClient *apim.APIM, app *v1alpha1.Application) (any, error) {
	automation := k8s.IsAutomationAPIManaged(app)
	app.PopulateIDs(apimClient.Context, automation)
	if automation {
		remoteApp, err := apimClient.Applications.GetByHRID(app.Spec.HRID)
		if err != nil {
			return nil, err
		}
		return *remoteApp, nil
	}
	remoteApp, err := apimClient.Applications.GetByID(app.Spec.ID)
	if err != nil {
		return nil, err
	}
	return *remoteApp, nil
}
