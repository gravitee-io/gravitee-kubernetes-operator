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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hrid"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

// preCheck validates the derived HRID, then what the manifest alone can tell about a studio.
func preCheck(_ context.Context, proxy *v1alpha1.McpProxy) *errors.AdmissionErrors {
	if err := hrid.Validate(refs.NewNamespacedNameFromObject(proxy).HRID()); err != nil {
		errs := errors.NewAdmissionErrors()
		errs.AddSevere(err.Error())
		return errs
	}
	return validateStudio(proxy)
}

// deferDryRun skips the dry run while a CatalogMcpServer the studio selects is missing or not synced: the
// platform would refuse the proxy, a normal state when a directory is applied at once. The controller waits
// for the server before applying.
func deferDryRun(ctx context.Context, proxy *v1alpha1.McpProxy) string {
	if err := dynamic.AssertCatalogMcpServersSynced(ctx, proxy); err != nil {
		return err.Error() + ": the proxy is checked against the platform once it is"
	}
	return ""
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
