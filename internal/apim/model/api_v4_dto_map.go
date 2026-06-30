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

func mapAPIV4Plans(plans *map[string]*v4.Plan) []*APIV4Plan {
	if plans == nil {
		return nil
	}

	result := make([]*APIV4Plan, 0, len(*plans))
	for hrid, plan := range *plans {
		if plan == nil {
			continue
		}
		mapped := mapViaJSON[*APIV4Plan](plan)
		if mapped == nil {
			continue
		}
		mapped.HRID = hrid
		result = append(result, mapped)
	}

	return result
}

func mapAPIV4Pages(pages *map[string]*v4.Page) []*APIV4Page {
	if pages == nil {
		return nil
	}

	result := make([]*APIV4Page, 0, len(*pages))
	for hrid, page := range *pages {
		if page == nil {
			continue
		}
		mapped := mapViaJSON[*APIV4Page](page)
		if mapped == nil {
			continue
		}
		mapped.HRID = hrid
		result = append(result, mapped)
	}

	return result
}
