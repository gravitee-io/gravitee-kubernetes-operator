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

package authzentity

import (
	"context"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
)

// validateDelete names the AuthzEntity resources that declare the entity as a parent through its
// ManagementContext: the platform deletes the RESOURCE ones and detaches the PRINCIPAL ones, and
// no event tells their controller. It does not call the platform's delete dry run, which would
// also name entities other sources created: APIM 4.13 ignores dryRun on a DELETE and deletes.
func validateDelete(ctx context.Context, entity *v1alpha1.AuthzEntity) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	ref := refs.NewNamespacedName(entity.GetNamespace(), entity.GetName())
	children, err := search.FindAuthzEntityChildren(ctx, entity)
	if err != nil {
		errs.AddWarningf("unable to check which resources declare authz entity [%s] as a parent: %s", ref.String(), err)
		return errs
	}
	if len(children) == 0 {
		return errs
	}

	errs.AddWarningf(
		"authz entity [%s] is a parent of AuthzEntity [%s]: deleting it detaches each PRINCIPAL among "+
			"them and deletes each RESOURCE on the platform, and they differ from their manifests until "+
			"their spec next changes",
		ref.String(), strings.Join(children, ", "),
	)
	return errs
}
