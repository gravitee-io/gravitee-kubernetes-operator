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

package catalogmcpserver

import (
	"context"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
)

func validateDelete(ctx context.Context, srv *v1alpha1.CatalogMcpServer) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	ref := refs.NewNamespacedName(srv.GetNamespace(), srv.GetName())
	proxies := &v1alpha1.McpProxyList{}
	if err := search.FindByFieldReferencing(ctx, search.McpProxyCatalogMcpServerField, ref, proxies); err != nil {
		errs.AddWarningf("unable to check which mcp proxies select catalog mcp server [%s]: %s", ref.String(), err)
		return errs
	}

	if len(proxies.Items) == 0 {
		return errs
	}

	names := make([]string, 0, len(proxies.Items))
	for i := range proxies.Items {
		names = append(names, refs.NewNamespacedNameFromObject(&proxies.Items[i]).String())
	}
	errs.AddWarningf(
		"catalog mcp server [%s] is still selected by McpProxy [%s]",
		ref.String(), strings.Join(names, ", "),
	)

	return errs
}
