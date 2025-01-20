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

package base

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

func validatePages(api core.ApiDefinitionObject) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	pages := make(map[string]*base.Page)

	switch t := api.(type) {
	case *v1alpha1.ApiDefinition:
		if t.Spec.Pages != nil {
			for k, v := range *t.Spec.Pages {
				pages[k] = v.Page
			}
		}
	case *v1alpha1.ApiV4Definition:
		if t.Spec.Pages != nil {
			for k, v := range *t.Spec.Pages {
				pages[k] = v.Page
			}
		}
	}

	for name, page := range pages {
		if page.Parent != nil && !parentFound(pages, page.Parent) {
			errs.AddSeveref("can not apply API [%s]. Parent page [%s] can not be found for page [%s]",
				api.GetName(),
				*page.Parent,
				name,
			)
		}
	}

	return errs
}

func parentFound(pages map[string]*base.Page, parent *string) bool {
	if _, ok := pages[*parent]; !ok {
		return false
	}

	return true
}
