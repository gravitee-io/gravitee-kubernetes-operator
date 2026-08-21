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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

func mergeDriftValidation(
	ctx context.Context,
	oldLink *v1alpha1.PortalLink,
	newLink *v1alpha1.PortalLink,
	errs *errors.AdmissionErrors,
) {
	prtl, err := dynamic.ResolvePortal(ctx, newLink.GetPortalRef(), newLink.GetNamespace())
	if err != nil {
		errs.AddSeveref(
			"portal link [%s] references portal [%v] that can't be resolved",
			newLink.GetName(), newLink.GetPortalRef(),
		)
		return
	}
	if !prtl.HasContext() {
		errs.AddSeveref(
			"referenced portal [%v] has no management context (spec.contextRef)",
			newLink.GetPortalRef(),
		)
		return
	}

	errs.MergeWith(drift.ValidateDriftWithContext(ctx, oldLink, newLink,
		func(ctx context.Context) (*apim.APIM, error) {
			return apim.FromContextRef(ctx, prtl.ContextRef(), prtl.GetNamespace())
		},
		resolveRefs,
		getRemotePortalLink(prtl),
		drift.MapDTO(func(link *v1alpha1.PortalLink) model.PortalLinkDTO {
			return model.ToPortalLinkDTO(link.Spec.Type, refs.NewNamespacedNameFromObject(link).HRID())
		}),
	))
}

func resolveRefs(context.Context, *v1alpha1.PortalLink) error {
	return nil
}

func getRemotePortalLink(prtl *v1alpha1.Portal) drift.RemoteObjectGetter[*v1alpha1.PortalLink] {
	return func(apimClient *apim.APIM, link *v1alpha1.PortalLink) (any, error) {
		portalHrid := refs.NewNamespacedNameFromObject(prtl).HRID()
		linkHrid := refs.NewNamespacedNameFromObject(link).HRID()
		remote, err := apimClient.Links.GetByHRID(portalHrid, linkHrid)
		if err != nil {
			return nil, err
		}
		return model.PortalLinkDTO{
			HRID:     remote.HRID,
			Name:     remote.Name,
			Href:     remote.Href,
			Location: remote.Location,
			Order:    remote.Order,
		}, nil
	}
}
