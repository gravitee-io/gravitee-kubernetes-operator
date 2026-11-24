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

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimmcpproxy "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/mcpproxy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
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

	errs.MergeWith(validateDrift(ctx, oldObj, newObj))
	return errs
}

// validateDrift compares through the SDK client, not the APIM facade: the resolver builds that client,
// only once drift is known to be enabled, and hands no facade to the getter.
func validateDrift(ctx context.Context, oldObj, newObj *v1alpha1.McpProxy) *errors.AdmissionErrors {
	var client *apim.SDKClient
	resolveContext := func(ctx context.Context) (*apim.APIM, error) {
		var err error
		client, err = apimmcpproxy.CreateAPIMClient(ctx, newObj)
		return nil, err
	}
	getRemote := func(_ *apim.APIM, proxy *v1alpha1.McpProxy) (any, error) {
		dto, err := apimmcpproxy.ToDTO(proxy)
		if err != nil {
			return nil, err
		}
		return apimmcpproxy.GetRemote(ctx, client, dto)
	}
	return drift.ValidateDriftWithContext(ctx, oldObj, newObj, resolveContext, resolveRefs, getRemote,
		drift.MapDTO(toDTO))
}

// resolveRefs has nothing to resolve: the shared drift path compiles templates on both sides.
// Credential values are ignored by the payload's drift tags, whatever form the spec supplied them in.
func resolveRefs(context.Context, *v1alpha1.McpProxy) error {
	return nil
}

// toDTO maps the CRD to the payload compared with the remote. CEL pairs every union type with its block,
// so a mapping error is not expected. getRemote maps the new object first and hands its error to the fetch
// failure policy: only the old object can reach here unmapped, and it is compared as an empty payload.
func toDTO(proxy *v1alpha1.McpProxy) sdk.AimMcpProxy {
	dto, err := apimmcpproxy.ToDTO(proxy)
	if err != nil {
		return sdk.AimMcpProxy{}
	}
	return dto
}
