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

package portal

import (
	"context"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimportal "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

func validateUpdate(ctx context.Context, oldObj *v1alpha1.Portal, newObj *v1alpha1.Portal) *errors.AdmissionErrors {
	errs := validateCreate(ctx, newObj)
	if errs.IsSevere() {
		return errs
	}
	errs.MergeWith(validateDrift(ctx, oldObj, newObj))
	return errs
}

// validateDrift compares through the SDK client, not the APIM facade: the resolver builds that client,
// only once drift is known to be enabled, and hands no facade to the getter.
func validateDrift(ctx context.Context, oldObj, newObj *v1alpha1.Portal) *errors.AdmissionErrors {
	var client *apim.SDKClient
	resolveContext := func(ctx context.Context) (*apim.APIM, error) {
		var err error
		client, err = apimportal.CreateAPIMClient(ctx, newObj)
		return nil, err
	}
	getRemote := func(_ *apim.APIM, prtl *v1alpha1.Portal) (any, error) {
		dto, err := apimportal.ToDTO(prtl)
		if err != nil {
			return nil, err
		}
		remote, err := apimportal.GetRemote(ctx, client, dto)
		if err != nil {
			return nil, err
		}
		return apimportal.DriftDTO(remote), nil
	}
	return drift.ValidateDriftWithContext(ctx, oldObj, newObj, resolveContext, resolveRefs, getRemote,
		drift.MapDTO(toDTO))
}

func resolveRefs(context.Context, *v1alpha1.Portal) error {
	return nil
}

// toDTO maps the CRD to the payload compared with the remote, navigation visibility resolved as on the
// remote side. getRemote maps the new object first and hands its error to the fetch failure policy: only
// the old object can reach here unmapped, and it is compared as an empty payload.
func toDTO(prtl *v1alpha1.Portal) sdk.Portal {
	dto, err := apimportal.ToDTO(prtl)
	if err != nil {
		return sdk.Portal{}
	}
	return apimportal.DriftDTO(dto)
}
