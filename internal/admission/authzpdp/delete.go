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
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
)

// validateDelete names, per kind, the resources whose pdpRef still names the PDP through the PDP's
// ManagementContext. It does not call
// the platform's delete dry run: APIM 4.13 ignores dryRun on a DELETE and deletes the PDP.
func validateDelete(ctx context.Context, pdp *v1alpha1.AuthzPdp) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	ref := refs.NewNamespacedName(pdp.GetNamespace(), pdp.GetName())
	dependents, err := search.FindAuthzPdpDependents(ctx, pdp, search.AuthzPdpDependents())
	if err != nil {
		errs.AddWarningf("unable to check which resources target authz pdp [%s]: %s", ref.String(), err)
		return errs
	}

	scope := pdp.Status.RoutingScope
	if scope == "" {
		scope = "its routing scope"
	}
	for _, d := range dependents {
		errs.AddWarningf(
			"authz pdp [%s] is still targeted by %s [%s]: deleting it removes %s from their targets, "+
				"and each waits on ResolvedRefs=False until its pdpRef names another AuthzPdp",
			ref.String(), d.Kind, strings.Join(d.Names, ", "), scope,
		)
	}

	return errs
}
