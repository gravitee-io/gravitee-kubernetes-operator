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

package group

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimgroup "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/group"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

func validateUpdate(ctx context.Context, oldObj *v1alpha1.Group, newObj *v1alpha1.Group) *errors.AdmissionErrors {
	errs := validateCreate(ctx, newObj)
	if errs.IsSevere() {
		return errs
	}

	// Group talks to APIM through its SDK client, not the APIM facade: the resolver builds that client,
	// only once drift is known to be enabled, and hands no facade to the getter.
	var client *apim.SDKClient
	resolveContext := func(ctx context.Context) (*apim.APIM, error) {
		var err error
		client, err = apimgroup.CreateAPIMClient(ctx, newObj)
		return nil, err
	}
	getRemote := func(_ *apim.APIM, grp *v1alpha1.Group) (any, error) {
		dto, err := apimgroup.ToDTO(grp)
		if err != nil {
			return nil, err
		}
		return apimgroup.GetRemote(ctx, client, dto)
	}
	errs.MergeWith(drift.ValidateDriftWithContext(ctx, oldObj, newObj, resolveContext, resolveRefs, getRemote,
		drift.MapDTO(toDTO)))
	return errs
}

func resolveRefs(context.Context, *v1alpha1.Group) error {
	return nil
}

// toDTO maps the CRD to the payload compared with the remote. The same spec already went through the
// dry run in validateCreate, so a mapping error is not expected here; should one occur, the empty
// payload is reported as drift rather than silently admitted.
func toDTO(grp *v1alpha1.Group) apimgroup.DTO {
	dto, err := apimgroup.ToDTO(grp)
	if err != nil {
		return apimgroup.DTO{}
	}
	return dto.WithoutIgnoredMembers()
}
