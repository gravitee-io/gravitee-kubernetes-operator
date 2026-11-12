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

// ToAPIV4DTO maps the spec to the payload sent to APIM. APIM numbers the plans from their position
// in the list, so they are sent in the order they should be stored.
func ToAPIV4DTO(api *v4.Api) APIV4DTO {
	if api == nil || api.V4BaseApi == nil {
		return APIV4DTO{}
	}

	dto := mapViaJSON[APIV4DTO](api.V4BaseApi)
	dto.Pages = mapAPIV4Pages(api.Pages)
	dto = dto.Canonicalize()
	dto.Plans = inSendOrder(mapAPIV4Plans(api.Plans))

	return dto
}

// ToAPIV4DriftDTO maps the spec the way drift compares it with the remote API, plans aligned by HRID.
// When every plan sets an order, each plan's order is its rank in the send order, which is the order
// APIM stores. Otherwise the orders are cleared, which leaves the plan order out of the comparison.
func ToAPIV4DriftDTO(api *v4.Api) APIV4DTO {
	dto := ToAPIV4DTO(api)
	if slices.ContainsFunc(dto.Plans, func(plan *APIV4PlanDTO) bool { return plan.Order == nil }) {
		for _, plan := range dto.Plans {
			plan.Order = nil
		}
	} else {
		setOrderFromPosition(dto.Plans)
	}
	return dto.Canonicalize()
}

// inSendOrder sorts plans by order, unset orders last, then by HRID.
func inSendOrder(plans []*APIV4PlanDTO) []*APIV4PlanDTO {
	return slices.SortedStableFunc(slices.Values(plans), func(a, b *APIV4PlanDTO) int {
		return cmp.Or(compareOrder(a.Order, b.Order), cmp.Compare(a.HRID, b.HRID))
	})
}

func compareOrder(a, b *int) int {
	switch {
	case a == nil && b == nil:
		return 0
	case a == nil:
		return 1
	case b == nil:
		return -1
	default:
		return cmp.Compare(*a, *b)
	}
}

// setOrderFromPosition numbers plans 1…n in list order, as APIM does with the list it receives.
func setOrderFromPosition(plans []*APIV4PlanDTO) {
	for i, plan := range plans {
		order := i + 1
		plan.Order = &order
	}
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

	return result
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
	return result
}

// Canonicalize returns a copy with plans and pages sorted by HRID, so that a mapped CRD and a
// remote state compare item by item. The CRD declares them as maps; the plan order is compared
// through each plan's order field, not through its position.
func (dto APIV4DTO) Canonicalize() APIV4DTO {
	dto.Plans = slices.SortedStableFunc(slices.Values(dto.Plans), func(a, b *APIV4PlanDTO) int {
		return cmp.Compare(a.HRID, b.HRID)
	})
	dto.Pages = slices.SortedStableFunc(slices.Values(dto.Pages), func(a, b *APIV4PageDTO) int {
		return cmp.Compare(a.HRID, b.HRID)
	})
	return dto
}
