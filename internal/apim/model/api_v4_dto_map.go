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
	"cmp"
	"slices"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
)

func ToAPIV4DTO(api *v4.Api) APIV4DTO {
	if api == nil || api.V4BaseApi == nil {
		return APIV4DTO{}
	}

	dto := mapViaJSON[APIV4DTO](api.V4BaseApi)
	dto.Plans = mapAPIV4Plans(api.Plans)
	dto.Pages = mapAPIV4Pages(api.Pages)

	return dto
}

func mapAPIV4Plans(plans *map[string]*v4.Plan) []*APIV4PlanDTO {
	if plans == nil {
		return nil
	}

	result := make([]*APIV4PlanDTO, 0, len(*plans))
	for _, plan := range *plans {
		if plan == nil {
			continue
		}
		mapped := mapViaJSON[*APIV4PlanDTO](plan)
		if mapped == nil {
			continue
		}
		result = append(result, mapped)
	}

	return slices.SortedFunc(slices.Values(result), func(a, b *APIV4PlanDTO) int {
		return cmp.Compare(a.HRID, b.HRID)
	})
}

func mapAPIV4Pages(pages *map[string]*v4.Page) []*APIV4PageDTO {
	if pages == nil {
		return nil
	}

	result := make([]*APIV4PageDTO, 0, len(*pages))
	for _, page := range *pages {
		if page == nil {
			continue
		}
		mapped := mapViaJSON[*APIV4PageDTO](page)
		if mapped == nil {
			continue
		}
		result = append(result, mapped)
	}
	return slices.SortedFunc(slices.Values(result), func(a, b *APIV4PageDTO) int {
		return cmp.Compare(a.HRID, b.HRID)
	})
}
