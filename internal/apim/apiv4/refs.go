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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/apidefinition"
)

// ResolveRefs resolves the API's references for the reconcile: ApiResources, Shared Policy Groups, groups and the
// console notification are inlined, and template references are recorded. A missing group is reported on the
// ResolvedRefs condition, not as an error. While the API is being deleted it resolves nothing: Delete only needs
// the API's identity, and a reference deleted first must not leave the API in APIM. A CR re-created without status
// gets back the id of the API GKO 4.11 created for it.
func ResolveRefs(ctx context.Context, obj *v1alpha1.ApiV4Definition, _ string) error {
	if obj.IsBeingDeleted() {
		return nil
	}
	apidefinition.LogUnsyncedGroupRefs(ctx, obj, obj.GetGroupRefs())
	if err := apidefinition.PrepareV4SpecForAutomation(ctx, obj, true); err != nil {
		return err
	}
	return adoptLegacyID(ctx, obj)
}

// ResolveRefsForAdmission is ResolveRefs for the admission webhook: it records no template references and logs
// nothing.
func ResolveRefsForAdmission(ctx context.Context, obj *v1alpha1.ApiV4Definition, _ string) error {
	if err := apidefinition.PrepareV4SpecForAutomation(ctx, obj, false); err != nil {
		return err
	}
	return adoptLegacyID(ctx, obj)
}

// adoptLegacyID gives a CR re-created without status the id of the API GKO 4.11 created for it, so that ToDTO
// sends it by UUID (GKO-3474). The id lands on the spec of the copy being synced, which the lifecycle never saves,
// and the lookup skips a CR being deleted: neither a failed PUT nor a delete before the first sync leads a delete
// to the API the CR would have re-attached to.
func adoptLegacyID(ctx context.Context, obj *v1alpha1.ApiV4Definition) error {
	if !obj.HasContext() || !apidefinition.MayHaveLegacyID(obj) {
		return nil
	}
	client, err := apim.FromContextRef(ctx, obj.ContextRef(), obj.GetNamespace())
	if err != nil {
		return err
	}
	return apidefinition.AdoptLegacyID(ctx, client, obj)
}
