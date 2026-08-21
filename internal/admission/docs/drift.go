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

package docs

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/service"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

func mergeDriftValidation(
	ctx context.Context,
	oldDoc *v1alpha1.Documentation,
	newDoc *v1alpha1.Documentation,
	errs *errors.AdmissionErrors,
) {
	var target *dryRunTarget
	if newDoc.IsPortalDoc() {
		target = resolvePortalTarget(ctx, newDoc, errs)
	} else {
		target = resolveApiTarget(ctx, newDoc, errs)
	}
	if errs.IsSevere() || target == nil {
		return
	}

	errs.MergeWith(drift.ValidateDriftWithContext(ctx, oldDoc, newDoc,
		documentationContextResolver(target),
		resolveRefs,
		getRemoteDocumentation(target.parent),
		drift.MapDTO(func(doc *v1alpha1.Documentation) model.DocumentationDTO {
			return model.ToDocumentationDTO(doc)
		}),
	))
}

func documentationContextResolver(target *dryRunTarget) drift.ContextResolver {
	return func(ctx context.Context) (*apim.APIM, error) {
		return apim.FromContextRef(ctx, target.contextRef, target.contextNs)
	}
}

func resolveRefs(context.Context, *v1alpha1.Documentation) error {
	return nil
}

func getRemoteDocumentation(parent service.DocumentationParent) drift.RemoteObjectGetter[*v1alpha1.Documentation] {
	return func(apimClient *apim.APIM, doc *v1alpha1.Documentation) (any, error) {
		dto := model.ToDocumentationDTO(doc)
		remote, err := apimClient.Documentations.GetByHRID(parent, dto.HRID)
		if err != nil {
			return nil, err
		}
		return model.DocumentationDTO{
			HRID:     remote.HRID,
			Name:     remote.Name,
			PageType: remote.PageType,
			Content:  remote.Content,
			Location: remote.Location,
			Order:    remote.Order,
			Area:     remote.Area,
		}, nil
	}
}
