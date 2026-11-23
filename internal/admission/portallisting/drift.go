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

package portallisting

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimportallisting "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portallisting"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

// mergeDriftValidation compares through the SDK client of the portal's context, not the APIM facade: the
// resolver builds that client, only once drift is known to be enabled, and hands no facade to the getter.
func mergeDriftValidation(
	ctx context.Context,
	oldListing *v1alpha1.PortalListing,
	newListing *v1alpha1.PortalListing,
	errs *errors.AdmissionErrors,
) {
	if validatePortal(ctx, newListing, errs) == nil {
		return
	}
	var client *apim.SDKClient
	resolveContext := func(ctx context.Context) (*apim.APIM, error) {
		var err error
		client, err = apimportallisting.CreateAPIMClient(ctx, newListing)
		return nil, err
	}
	getRemote := func(_ *apim.APIM, listing *v1alpha1.PortalListing) (any, error) {
		dto, err := apimportallisting.ToDTO(listing)
		if err != nil {
			return nil, err
		}
		return apimportallisting.GetRemote(ctx, client, dto)
	}
	errs.MergeWith(drift.ValidateDriftWithContext(ctx, oldListing, newListing, resolveContext, resolveRefs,
		getRemote, drift.MapDTO(toDTO)))
}

func resolveRefs(context.Context, *v1alpha1.PortalListing) error {
	return nil
}

// toDTO maps the CRD to the payload compared with the remote. getRemote maps the new object first and
// hands its error to the fetch failure policy: only the old object can reach here unmapped, and it is
// compared as an empty payload.
func toDTO(listing *v1alpha1.PortalListing) apimportallisting.DTO {
	dto, err := apimportallisting.ToDTO(listing)
	if err != nil {
		return apimportallisting.DTO{}
	}
	return dto
}
