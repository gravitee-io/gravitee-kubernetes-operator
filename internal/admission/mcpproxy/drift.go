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

package mcpproxy

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
	oldObj *v1alpha1.McpProxy,
	newObj *v1alpha1.McpProxy,
) *errors.AdmissionErrors {
	errs := validateContext(ctx, newObj)
	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(validateStudio(newObj))
	if errs.IsSevere() {
		return errs
	}

	// A metadata-only update, the controller's own finalizer and annotation writes included,
	// does not need the platform: only a spec change is dry-run.
	if oldObj.Spec.Hash() != newObj.Spec.Hash() {
		errs.MergeWith(validatePlatform(ctx, newObj))
		if errs.IsSevere() {
			return errs
		}
	}

	// A proxy the platform never received, a studio still waiting for its servers for instance, has
	// nothing to drift from: comparing would apply the remote-missing policy (deny by default) and
	// refuse the very edit that unblocks it, a serverRef typo fix included.
	if oldObj.Status.ID == "" {
		return errs
	}

	errs.MergeWith(drift.ValidateDrift(ctx, oldObj, newObj, resolveRefs, getRemoteMcpProxy,
		drift.MapDTO(model.ToMcpProxyDTO)))
	return errs
}

// resolveRefs compiles templates on the copy drift compares, so a templated URL or step
// configuration is compared resolved on both the old and the new side. Credential values are
// ignored by the DTO's drift tags, whatever form the spec supplied them in.
func resolveRefs(ctx context.Context, proxy *v1alpha1.McpProxy) error {
	return template.Compile(ctx, proxy, false)
}

// getRemoteMcpProxy reads the platform's copy and returns the payload half, in the order the
// mapper produces.
func getRemoteMcpProxy(apimClient *apim.APIM, proxy *v1alpha1.McpProxy) (any, error) {
	hrid := refs.NewNamespacedNameFromObject(proxy).HRID()
	remote, err := apimClient.McpProxies.GetByHRID(hrid)
	if err != nil {
		return nil, err
	}
	remote.Canonicalize()
	return remote.McpProxyDTO, nil
}
