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

package apidefinition

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

// ResolveV4Refs inlines the ApiResource and SharedPolicyGroup references of a v4 API, so that
// the admission dry-run and the reconciler send the same flows and resources to APIM.
// Admission passes updateMetadata=false to leave referenced Secrets and ConfigMaps untouched.
func ResolveV4Refs(ctx context.Context, api *v1alpha1.ApiV4Definition, updateMetadata bool) error {
	nsCtx := WithAPINamespace(ctx, api.Namespace)

	if err := ResolveResources(nsCtx, api.Spec.Resources, updateMetadata); err != nil {
		return err
	}

	return ResolveSharedPolicyGroups(nsCtx, &api.Spec, updateMetadata)
}
