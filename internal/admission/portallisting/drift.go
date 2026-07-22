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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/service"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	"k8s.io/apimachinery/pkg/runtime"
)

func mergeDriftValidation(
	ctx context.Context,
	oldObj runtime.Object,
	newObj runtime.Object,
	errs *errors.AdmissionErrors,
) {
	oldListing, _ := oldObj.(*v1alpha1.PortalListing)
	newListing, _ := newObj.(*v1alpha1.PortalListing)

	prtl, err := dynamic.ResolvePortal(ctx, newListing.GetPortalRef(), newListing.GetNamespace())
	if err != nil {
		errs.AddSeveref(
			"portal listing [%s] references portal [%v] that can't be resolved",
			newListing.GetName(), newListing.GetPortalRef(),
		)
		return
	}
	if !prtl.HasContext() {
		errs.AddSeveref(
			"referenced portal [%v] has no management context (spec.contextRef)",
			newListing.GetPortalRef(),
		)
		return
	}

	errs.MergeWith(drift.ValidateDriftWithContext(ctx, oldListing, newListing,
		func(ctx context.Context) (*apim.APIM, error) {
			return apim.FromContextRef(ctx, prtl.ContextRef(), prtl.GetNamespace())
		},
		resolveRefs,
		getRemotePortalListing(prtl),
		drift.MapDTO(func(listing *v1alpha1.PortalListing) model.PortalListingDTO {
			return *service.ToPortalListingDTO(listing)
		}),
	))
}

func resolveRefs(context.Context, runtime.Object) error {
	return nil
}

func getRemotePortalListing(prtl *v1alpha1.Portal) drift.RemoteObjectGetter {
	return func(apimClient *apim.APIM, o runtime.Object, admissionErrors *errors.AdmissionErrors) any {
		listing, _ := o.(*v1alpha1.PortalListing)
		dto := service.ToPortalListingDTO(listing)
		portalHrid := refs.NewNamespacedNameFromObject(prtl).HRID()
		remote, err := apimClient.Listings.GetByHRID(portalHrid, dto.HRID)
		if err != nil {
			admissionErrors.AddSeveref(
				"cannot fetch PortalListing during drift detection from portal HRID %s and listing HRID %s: %s",
				portalHrid, dto.HRID, err.Error(),
			)
			return nil
		}
		return model.PortalListingDTO{
			HRID: remote.HRID,
			APIs: remote.APIs,
		}
	}
}
