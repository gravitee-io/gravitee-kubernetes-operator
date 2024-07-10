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

package model

import (
	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
)

func ToApiImport(spec *v2.Api) *ApiImport {
	apiImport := &ApiImport{Api: spec}
	pages := make([]*PageImport, 0)
	for _, p := range spec.Pages {
		pages = append(pages, &PageImport{
			Page: p,
			DefinitionContext: &v2.DefinitionContext{
				Origin: v2.OriginKubernetes,
			},
		})
	}
	spec.Pages = nil
	apiImport.Pages = pages
	apiImport.DisableMembershipNotifications = !spec.NotifyMembers
	return apiImport
}
