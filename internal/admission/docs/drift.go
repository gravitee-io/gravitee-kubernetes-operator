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
	apimdocs "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/docs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

// mergeDriftValidation compares through the SDK client of the page's parent context, not the APIM
// facade: the resolver builds that client, only once drift is known to be enabled, and hands no facade
// to the getter.
func mergeDriftValidation(
	ctx context.Context,
	oldDoc *v1alpha1.Documentation,
	newDoc *v1alpha1.Documentation,
	errs *errors.AdmissionErrors,
) {
	if !resolveTarget(ctx, newDoc, errs) {
		return
	}
	var client *apim.SDKClient
	resolveContext := func(ctx context.Context) (*apim.APIM, error) {
		var err error
		client, err = apimdocs.CreateAPIMClient(ctx, newDoc)
		return nil, err
	}
	getRemote := func(_ *apim.APIM, doc *v1alpha1.Documentation) (any, error) {
		dto, err := apimdocs.ToDTO(doc)
		if err != nil {
			return nil, err
		}
		return apimdocs.GetRemote(ctx, client, dto)
	}
	errs.MergeWith(drift.ValidateDriftWithContext(ctx, oldDoc, newDoc, resolveContext, resolveRefs, getRemote,
		drift.MapDTO(toDTO)))
}

func resolveRefs(context.Context, *v1alpha1.Documentation) error {
	return nil
}

// toDTO maps the CRD to the payload compared with the remote. getRemote maps the new object first and
// hands its error to the fetch failure policy: only the old object can reach here unmapped, and it is
// compared as an empty payload.
func toDTO(doc *v1alpha1.Documentation) apimdocs.DTO {
	dto, err := apimdocs.ToDTO(doc)
	if err != nil {
		return apimdocs.DTO{}
	}
	return dto
}
