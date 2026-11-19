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

package apiv4

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/apidefinition"
)

// ResolveRefs resolves the API's references for the reconcile: ApiResources, Shared Policy Groups, groups and the
// console notification are inlined, and template references are recorded. A missing group is reported on the
// ResolvedRefs condition, not as an error. While the API is being deleted it resolves nothing: Delete only needs
// the API's identity, and a reference deleted first must not leave the API in APIM.
func ResolveRefs(ctx context.Context, obj *v1alpha1.ApiV4Definition, _ string) error {
	if obj.IsBeingDeleted() {
		return nil
	}
	apidefinition.LogUnsyncedGroupRefs(ctx, obj, obj.GetGroupRefs())
	return apidefinition.PrepareV4SpecForAutomation(ctx, obj, true)
}

// ResolveRefsForAdmission is ResolveRefs for the admission webhook: it records no template references and logs
// nothing.
func ResolveRefsForAdmission(ctx context.Context, obj *v1alpha1.ApiV4Definition, _ string) error {
	return apidefinition.PrepareV4SpecForAutomation(ctx, obj, false)
}
