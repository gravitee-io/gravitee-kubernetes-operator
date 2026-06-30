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
	"k8s.io/apimachinery/pkg/runtime"
)

func mergeDriftValidation(
	ctx context.Context,
	oldObj runtime.Object,
	newObj runtime.Object,
	errs *errors.AdmissionErrors,
) {
	oldDoc, _ := oldObj.(*v1alpha1.Documentation)
	newDoc, _ := newObj.(*v1alpha1.Documentation)

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
			return *service.ToDocumentationDTO(doc)
		}),
	))
}

func documentationContextResolver(target *dryRunTarget) drift.ContextResolver {
	return func(ctx context.Context) (*apim.APIM, error) {
		return apim.FromContextRef(ctx, target.contextRef, target.contextNs)
	}
}

func resolveRefs(context.Context, runtime.Object) error {
	return nil
}

func getRemoteDocumentation(parent service.DocumentationParent) drift.RemoteObjectGetter {
	return func(apimClient *apim.APIM, o runtime.Object, admissionErrors *errors.AdmissionErrors) any {
		doc, _ := o.(*v1alpha1.Documentation)
		dto := service.ToDocumentationDTO(doc)
		remote, err := apimClient.Documentations.GetByHRID(parent, dto.HRID)
		if err != nil {
			admissionErrors.AddSeveref(
				"cannot fetch Documentation during drift detection from HRID %s: %s",
				dto.HRID, err.Error(),
			)
			return nil
		}
		return model.DocumentationDTO{
			HRID:     remote.HRID,
			Name:     remote.Name,
			PageType: remote.PageType,
			Content:  remote.Content,
			Location: remote.Location,
			Order:    remote.Order,
		}
	}
}
