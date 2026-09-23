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
	oldObj *v1alpha1.CatalogMcpServer,
	newObj *v1alpha1.CatalogMcpServer,
) *errors.AdmissionErrors {
	errs := validateCreate(ctx, newObj)
	if errs.IsSevere() {
		return errs
	}
	errs.MergeWith(drift.ValidateDrift(ctx, oldObj, newObj, resolveRefs, getRemoteCatalogMcpServer,
		drift.MapDTO(model.ToCatalogMcpServerDTO)))
	return errs
}

// resolveRefs compiles templates on the copy drift compares, so a templated endpoint or token
// URL is compared resolved on both the old and the new side. Credential values are ignored by
// the DTO's drift tags, whatever form the spec supplied them in.
func resolveRefs(ctx context.Context, srv *v1alpha1.CatalogMcpServer) error {
	return template.Compile(ctx, srv, false)
}

// getRemoteCatalogMcpServer reads the platform's copy and returns the payload half only: what
// the platform discovered lives in the state and is never compared.
func getRemoteCatalogMcpServer(apimClient *apim.APIM, srv *v1alpha1.CatalogMcpServer) (any, error) {
	hrid := refs.NewNamespacedNameFromObject(srv).HRID()
	remote, err := apimClient.CatalogMcpServers.GetByHRID(hrid)
	if err != nil {
		return nil, err
	}
	return remote.CatalogMcpServerDTO, nil
}
