// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package authzpdp

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/template"
)

func validateUpdate(
	ctx context.Context,
	oldObj *v1alpha1.AuthzPdp,
	newObj *v1alpha1.AuthzPdp,
) *errors.AdmissionErrors {
	// A metadata-only update, the controller's own finalizer and annotation writes included,
	// does not need the platform: only a spec change is dry-run. Hashed before validateContext,
	// which compiles templates into newObj while oldObj keeps them: hashed after, a templated spec
	// would always look changed.
	specChanged := oldObj.Spec.Hash() != newObj.Spec.Hash()

	errs := validateContext(ctx, newObj)
	if errs.IsSevere() {
		return errs
	}

	if specChanged {
		errs.MergeWith(validateDryRun(ctx, newObj))
		if errs.IsSevere() {
			return errs
		}
	}

	// A PDP the platform never received has nothing to drift from: comparing would apply the
	// remote-missing policy (deny by default) and refuse the very edit that fixes it.
	if oldObj.Status.ID == "" {
		return errs
	}

	// A strict change on a sibling of the same targetPdpId shows up as drift here, which is what
	// happened: the platform changed outside this resource.
	errs.MergeWith(drift.ValidateDrift(ctx, oldObj, newObj, resolveRefs, getRemoteAuthzPdp,
		drift.MapDTO(model.ToAuthzPdpDTO)))
	return errs
}

// resolveRefs compiles templates on the copy drift compares, so a templated name or path is
// compared resolved on both the old and the new side.
func resolveRefs(ctx context.Context, pdp *v1alpha1.AuthzPdp) error {
	return template.Compile(ctx, pdp, false)
}

// getRemoteAuthzPdp reads the platform's copy and returns the payload half only: the routing scope
// and discovery path are derived and never compared.
func getRemoteAuthzPdp(apimClient *apim.APIM, pdp *v1alpha1.AuthzPdp) (any, error) {
	hrid := refs.NewNamespacedNameFromObject(pdp).HRID()
	remote, err := apimClient.AuthzPdps.GetByHRID(hrid)
	if err != nil {
		return nil, err
	}
	return remote.AuthzPdpDTO, nil
}
