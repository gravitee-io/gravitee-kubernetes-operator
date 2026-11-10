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

// validateDelete names the RESOURCE entities of the entity's ManagementContext its delete deletes on
// the platform (search.FindAuthzEntityCascadeDeletes). The PRINCIPAL ones it detaches are applied
// again as declared: nothing the user sees changes, so they are not named. It does not call the
// platform's delete dry run, which would also name entities other sources created: APIM 4.13
// ignores dryRun on a DELETE and deletes.
func validateDelete(ctx context.Context, entity *v1alpha1.AuthzEntity) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	ref := refs.NewNamespacedName(entity.GetNamespace(), entity.GetName())
	deletes, err := search.FindAuthzEntityCascadeDeletes(ctx, entity)
	if err != nil {
		errs.AddWarningf("could not list the children of authz entity [%s]: %s", ref.String(), err)
		return errs
	}
	if len(deletes) == 0 {
		return errs
	}

	errs.AddWarningf(
		"deleting authz entity [%s] also deletes AuthzEntity [%s] on the platform, "+
			"which refuses them while their parent is missing",
		ref.String(), strings.Join(deletes, ", "),
	)
	return errs
}
