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

package apiv4

import (
	"cmp"
	"slices"
	"strings"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
)

// DriftDTO applies the normalisations the drift comparison needs, to the CR and the remote side alike. Drift
// path only: the payload sent to APIM keeps the values as they are (AGENTS.md, Drift detection).
//
//   - plans and pages are compared in hrid order (the plan order is compared through PlanOrder);
//   - portal navigation visibility is resolved from the closest declared ancestor, as APIM does;
//   - listener paths lose overrideAccess, which APIM does not store as sent;
//   - KAFKA entrypoints lose their qos, which APIM does not store for a native entrypoint.
//
// It never mutates the slices or maps of dto (a remote DTO shares them with the fetched state).
func DriftDTO(dto DTO) DTO {
	dto.Plans = slices.SortedStableFunc(slices.Values(dto.Plans), func(a, b sdk.PlanV4) int {
		return cmp.Compare(a.Hrid, b.Hrid)
	})
	dto.Pages = slices.SortedStableFunc(slices.Values(dto.Pages), func(a, b sdk.PageV4) int {
		return cmp.Compare(a.Hrid, b.Hrid)
	})
	dto.PortalNavigation = withResolvedVisibility(dto.PortalNavigation)
	dto.Listeners = withNormalisedListeners(dto.Listeners)
	return dto
}

func withResolvedVisibility(paths []sdk.NavigationPath) []sdk.NavigationPath {
	if len(paths) == 0 {
		return paths
	}
	declared := make(map[string]sdk.PortalVisibility, len(paths))
	for _, entry := range paths {
		if entry.Visibility != nil {
			declared[normalizeNavigationPath(entry.Path)] = *entry.Visibility
		}
	}
	resolved := make([]sdk.NavigationPath, len(paths))
	for i, entry := range paths {
		visibility := resolveNavigationVisibility(entry.Path, declared)
		entry.Visibility = &visibility
		resolved[i] = entry
	}
	return resolved
}

func resolveNavigationVisibility(path string, declared map[string]sdk.PortalVisibility) sdk.PortalVisibility {
	segments := strings.Split(strings.Trim(normalizeNavigationPath(path), "/"), "/")
	for i := len(segments); i > 0; i-- {
		ancestor := "/" + strings.Join(segments[:i], "/")
		if visibility, ok := declared[ancestor]; ok && visibility != "" {
			return visibility
		}
	}
	return sdk.PortalVisibilityPUBLIC
}

func normalizeNavigationPath(path string) string {
	trimmed := strings.TrimRight(path, "/")
	if trimmed == "" {
		return "/"
	}
	return trimmed
}

func withNormalisedListeners(listeners []sdk.Listener) []sdk.Listener {
	if len(listeners) == 0 {
		return listeners
	}
	normalised := make([]sdk.Listener, len(listeners))
	for i, listener := range listeners {
		if listener.Http != nil {
			http := *listener.Http
			http.Paths = withoutOverrideAccess(http.Paths)
			listener.Http = &http
		}
		if listener.Kafka != nil {
			kafka := *listener.Kafka
			kafka.Entrypoints = withoutQos(kafka.Entrypoints)
			listener.Kafka = &kafka
		}
		normalised[i] = listener
	}
	return normalised
}

func withoutOverrideAccess(paths []sdk.PathV4) []sdk.PathV4 {
	if paths == nil {
		return nil
	}
	result := make([]sdk.PathV4, len(paths))
	for i, path := range paths {
		path.OverrideAccess = nil
		result[i] = path
	}
	return result
}

func withoutQos(entrypoints []sdk.Entrypoint) []sdk.Entrypoint {
	if entrypoints == nil {
		return nil
	}
	result := make([]sdk.Entrypoint, len(entrypoints))
	for i, entrypoint := range entrypoints {
		entrypoint.Qos = nil
		result[i] = entrypoint
	}
	return result
}
