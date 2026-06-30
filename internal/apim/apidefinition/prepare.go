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

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

// PrepareV4SpecForAutomation resolves references and normalizes the spec the same
// way as before syncing to the APIM Automation API.
func PrepareV4SpecForAutomation(ctx context.Context, api *v1alpha1.ApiV4Definition) error {
	nsCtx := WithAPINamespace(ctx, api.Namespace)
	spec := &api.Spec

	if err := ResolveResources(nsCtx, spec.Resources); err != nil {
		return err
	}

	if err := ResolveSharedPolicyGroups(nsCtx, spec); err != nil {
		return err
	}

	if groups, err := ResolveGroupRefs(ctx, api, spec.GetGroupRefs()); err != nil {
		return err
	} else {
		spec.Groups = append(spec.Groups, groups...)
	}

	spec.DefinitionContext = v4.NewDefaultKubernetesContext().MergeWith(spec.DefinitionContext)

	if spec.Context != nil {
		return ResolveConsoleNotificationRefs(ctx, api)
	}

	return nil
}
