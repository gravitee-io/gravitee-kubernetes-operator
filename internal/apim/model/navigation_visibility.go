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
	"strings"

	nav "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/navigation"
)

// WithResolvedVisibility returns a copy of the DTO in which every navigation entry
// carries the visibility APIM will store for it, so that drift detection can compare
// the field exactly instead of ignoring it.
//
// It is meant for the drift path only. The payload sent to APIM keeps an undeclared
// visibility omitted: APIM resolves it against the persisted parent, and sending an
// explicit PUBLIC under a PRIVATE parent would be rejected by its
// parent-must-be-public-if-child-is-public rule.
func (dto PortalDTO) WithResolvedVisibility() PortalDTO {
	if dto.Structure == nil || len(dto.Structure.TopNavbar) == 0 {
		return dto
	}
	declared := make(map[string]nav.Visibility, len(dto.Structure.TopNavbar))
	for _, entry := range dto.Structure.TopNavbar {
		if entry != nil {
			declared[normalizeNavigationPath(entry.Path)] = entry.Visibility
		}
	}
	resolved := make([]*NavigationEntryDTO, len(dto.Structure.TopNavbar))
	for i, entry := range dto.Structure.TopNavbar {
		if entry == nil {
			continue
		}
		copied := *entry
		copied.Visibility = resolveNavigationVisibility(copied.Path, declared)
		resolved[i] = &copied
	}
	dto.Structure = &NavigationStructureDTO{TopNavbar: resolved}
	return dto
}

// WithResolvedVisibility returns a copy of the DTO in which every portal navigation
// entry carries the visibility APIM will store for it. See
// [PortalDTO.WithResolvedVisibility] for why this is applied on the drift path only.
func (dto APIV4DTO) WithResolvedVisibility() APIV4DTO {
	if len(dto.PortalNavigation) == 0 {
		return dto
	}
	declared := make(map[string]nav.Visibility, len(dto.PortalNavigation))
	for _, entry := range dto.PortalNavigation {
		if entry != nil {
			declared[normalizeNavigationPath(entry.Path)] = entry.Visibility
		}
	}
	resolved := make([]*APIV4NavigationPathDTO, len(dto.PortalNavigation))
	for i, entry := range dto.PortalNavigation {
		if entry == nil {
			continue
		}
		copied := *entry
		copied.Visibility = resolveNavigationVisibility(copied.Path, declared)
		resolved[i] = &copied
	}
	dto.PortalNavigation = resolved
	return dto
}

// resolveNavigationVisibility mirrors PortalVisibility.resolve on the APIM side:
// the entry's own declaration wins, then the closest declared ancestor, then PUBLIC.
// Ancestors the spec does not list are folders APIM creates implicitly; they declare
// nothing, so they inherit too and the walk simply continues upwards.
func resolveNavigationVisibility(path string, declared map[string]nav.Visibility) nav.Visibility {
	segments := strings.Split(strings.Trim(normalizeNavigationPath(path), "/"), "/")
	for i := len(segments); i > 0; i-- {
		ancestor := "/" + strings.Join(segments[:i], "/")
		if visibility, ok := declared[ancestor]; ok && visibility != "" {
			return visibility
		}
	}
	return nav.Public
}

// normalizeNavigationPath drops a trailing slash so that "/a" and "/a/" describe the
// same node. The CRD only enforces the leading slash.
func normalizeNavigationPath(path string) string {
	trimmed := strings.TrimRight(path, "/")
	if trimmed == "" {
		return "/"
	}
	return trimmed
}
