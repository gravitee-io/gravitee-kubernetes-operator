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

package portallink

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimportallink "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portallink"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

// mergeDriftValidation compares through the SDK client of the link's parent context, not the APIM
// facade: the resolver builds that client, only once drift is known to be enabled, and hands no facade
// to the getter.
func mergeDriftValidation(
	ctx context.Context,
	oldLink *v1alpha1.PortalLink,
	newLink *v1alpha1.PortalLink,
	errs *errors.AdmissionErrors,
) {
	if !resolveTarget(ctx, newLink, errs) {
		return
	}
	var client *apim.SDKClient
	resolveContext := func(ctx context.Context) (*apim.APIM, error) {
		var err error
		client, err = apimportallink.CreateAPIMClient(ctx, newLink)
		return nil, err
	}
	getRemote := func(_ *apim.APIM, link *v1alpha1.PortalLink) (any, error) {
		dto, err := apimportallink.ToDTO(link)
		if err != nil {
			return nil, err
		}
		return apimportallink.GetRemote(ctx, client, dto)
	}
	errs.MergeWith(drift.ValidateDriftWithContext(ctx, oldLink, newLink, resolveContext, resolveRefs, getRemote,
		drift.MapDTO(toDTO)))
}

func resolveRefs(context.Context, *v1alpha1.PortalLink) error {
	return nil
}

// toDTO maps the CRD to the payload compared with the remote. getRemote maps the new object first and
// hands its error to the fetch failure policy: only the old object can reach here unmapped, and it is
// compared as an empty payload.
func toDTO(link *v1alpha1.PortalLink) apimportallink.DTO {
	dto, err := apimportallink.ToDTO(link)
	if err != nil {
		return apimportallink.DTO{}
	}
	return dto
}
