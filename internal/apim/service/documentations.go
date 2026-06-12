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

package service

import (
	"strconv"

	documentation "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/docs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	gohttp "github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

// DocumentationParent identifies the owning resource of a documentation page.
// A page is attached to exactly one of a portal or an API; the parent
// determines which automation endpoint the page is synced to.
type DocumentationParent struct {
	// Portal is the HRID of the owning portal, or "" when attached to an API.
	Portal string
	// API is the HRID of the owning API, or "" when attached to a portal.
	API string
}

type Documentations struct {
	*client.Client
}

func NewDocumentations(client *client.Client) *Documentations {
	return &Documentations{Client: client}
}

func (svc *Documentations) CreateOrUpdate(
	doc *v1alpha1.Documentation,
	parent DocumentationParent,
) (documentation.Status, error) {
	return svc.createOrUpdate(doc, parent, false)
}

func (svc *Documentations) DryRunCreateOrUpdate(
	doc *v1alpha1.Documentation,
	parent DocumentationParent,
) (documentation.Status, error) {
	return svc.createOrUpdate(doc, parent, true)
}

func (svc *Documentations) createOrUpdate(
	doc *v1alpha1.Documentation,
	parent DocumentationParent,
	dryRun bool,
) (documentation.Status, error) {
	url := svc.documentationsTarget(parent).
		WithQueryParam("dryRun", strconv.FormatBool(dryRun))

	dto := toDocumentationDTO(doc)
	importStatus := &documentation.Status{}

	if err := svc.HTTP.Put(url.String(), dto, &importStatus); err != nil {
		return *importStatus, err
	}

	k8s.AddAutomationAPIManagedCondition(doc)

	return *importStatus, nil
}

func (svc *Documentations) Delete(parent DocumentationParent, docHrid string) error {
	url := svc.documentationsTarget(parent).WithPath(docHrid)
	return svc.HTTP.Delete(url.String(), nil)
}

// GetByHRID For test purposes only.
func (svc *Documentations) GetByHRID(parent DocumentationParent, docHrid string) (*model.DocumentationState, error) {
	url := svc.documentationsTarget(parent).WithPath(docHrid)
	doc := new(model.DocumentationState)
	if err := svc.HTTP.Get(url.String(), doc); err != nil {
		return nil, err
	}
	return doc, nil
}

// documentationsTarget builds the documentations collection URL nested under
// the owning portal or API.
func (svc *Documentations) documentationsTarget(parent DocumentationParent) *gohttp.URL {
	if parent.API != "" {
		return svc.AutomationTarget("apis").
			WithPath(parent.API).
			WithPath("documentations")
	}
	return svc.AutomationTarget("portals").
		WithPath(parent.Portal).
		WithPath("documentations")
}

func toDocumentationDTO(doc *v1alpha1.Documentation) *model.DocumentationDTO {
	return &model.DocumentationDTO{
		HRID:     refs.NewNamespacedNameFromObject(doc).HRID(),
		Name:     doc.Spec.Name,
		PageType: doc.Spec.PageType,
		Content:  doc.Spec.Content,
		Location: doc.Spec.Location,
		Order:    doc.Spec.Order,
	}
}
