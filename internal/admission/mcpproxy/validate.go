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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/ctxref"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hrid"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

// validateCreate runs the checks in order: the context, the derived HRID, the studio references,
// then the platform dry run. Schema and CEL own the spec-only rules (entityId grammar, union
// pairing, mode and the immutable fields).
func validateCreate(ctx context.Context, proxy *v1alpha1.McpProxy) *errors.AdmissionErrors {
	errs := validateContext(ctx, proxy)
	if errs.IsSevere() {
		return errs
	}

	if err := hrid.Validate(refs.NewNamespacedNameFromObject(proxy).HRID()); err != nil {
		errs.AddSevere(err.Error())
		return errs
	}

	errs.MergeWith(validateStudio(proxy))
	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(validatePlatform(ctx, proxy))
	return errs
}

// validateContext requires a context that resolves. Resolving also compiles templates, so a dry
// run that follows carries resolved credentials.
func validateContext(ctx context.Context, proxy *v1alpha1.McpProxy) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	if !proxy.HasContext() {
		errs.AddSevere("a management context reference (spec.contextRef) is required")
		return errs
	}

	errs.Add(ctxref.Validate(ctx, proxy))
	return errs
}

// validateStudio checks what the manifest alone can tell about a studio: every reference names a
// CatalogMcpServer, and every server a selected tool comes from has a credential entry.
func validateStudio(proxy *v1alpha1.McpProxy) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	studio := proxy.Spec.Studio
	if studio == nil {
		return errs
	}

	for i := range studio.Tools {
		validateServerKind(errs, "studio.tools", studio.Tools[i].ServerRef)
	}
	for i := range studio.UpstreamAuth {
		validateServerKind(errs, "studio.upstreamAuth", studio.UpstreamAuth[i].ServerRef)
	}

	for _, server := range proxy.Spec.ServersWithoutUpstreamAuth(proxy.GetNamespace()) {
		errs.AddSeveref(
			"studio.upstreamAuth: server [%s] has no entry; declare the credential the gateway "+
				"presents to it, or type NONE when it needs none",
			server.String(),
		)
	}

	return errs
}

func validateServerKind(errs *errors.AdmissionErrors, path string, ref refs.NamespacedName) {
	if ref.Kind != "" && dynamic.ResourceFromKind(ref.Kind) != core.CRDCatalogMcpServerResource {
		errs.AddSeveref(
			"%s [%s]: serverRef kind must be CatalogMcpServer, got [%s]",
			path, ref.String(), ref.Kind,
		)
	}
}

// validatePlatform dry-runs the declaration. A studio selecting a catalog server that does not
// exist or has not synced yet cannot be checked by the platform, which would refuse it: that is
// a normal state when a directory is applied at once, so it is a warning and the dry run is
// skipped. The controller waits for the server before applying.
func validatePlatform(ctx context.Context, proxy *v1alpha1.McpProxy) *errors.AdmissionErrors {
	if err := dynamic.AssertCatalogMcpServersSynced(ctx, proxy); err != nil {
		errs := errors.NewAdmissionErrors()
		errs.AddWarningf("%s: the proxy is checked against the platform once it is", err)
		return errs
	}
	return validateDryRun(ctx, proxy)
}

func validateDryRun(ctx context.Context, proxy *v1alpha1.McpProxy) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	cp := proxy.DeepCopy()

	apimClient, err := apim.FromContextRef(ctx, cp.ContextRef(), cp.GetNamespace())
	if err != nil {
		errs.AddSevere(err.Error())
		return errs
	}

	// A dry run answers 200 with its findings; a refusal body is only expected from a real apply,
	// but is read the same way should the platform answer one.
	state, err := apimClient.McpProxies.DryRunCreateOrUpdate(cp)
	if err != nil {
		errs.MergeWith(refusalToAdmissionErrors(err))
		return errs
	}

	errs.MergeWith(errors.NewAdmissionErrorsFromStatus(state.Errors))

	return errs
}

// refusalToAdmissionErrors reports a platform refusal finding by finding; anything else is
// reported as the request error.
func refusalToAdmissionErrors(err error) *errors.AdmissionErrors {
	if findings, refused := model.AutomationRefusal(err); refused {
		return errors.NewAdmissionErrorsFromStatus(findings)
	}

	errs := errors.NewAdmissionErrors()
	errs.AddSevere(err.Error())
	return errs
}
