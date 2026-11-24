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

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/navigation"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/mapper"
)

// DriftDTO applies the normalisations the drift comparison needs, to the CR and the remote side alike. Drift
// path only: the payload sent to APIM keeps the values as they are (AGENTS.md, Drift detection).
//
//   - plans and pages are compared in hrid order (the plan order is compared through PlanOrder);
//   - portal navigation visibility is resolved from the closest declared ancestor, as APIM does;
//   - listener paths lose overrideAccess, which APIM does not store as sent;
//   - KAFKA entrypoints lose their qos, which APIM does not store for a native entrypoint;
//   - listeners are then compared unstructured (DriftListeners), as GKO always compared them;
//   - a boolean left out reads as false, as in GKO's former payload model.
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
	dto.DriftListeners = unstructuredListeners(withNormalisedListeners(dto.Listeners))
	dto.Listeners = nil
	return withFalseWhenUnset(dto)
}

// unstructuredListeners maps the listeners to the free-form maps GKO's former payload model held them in, so
// that drift reports their values as it always did.
func unstructuredListeners(listeners []sdk.Listener) []*model.APIV4GenericListenerDTO {
	if listeners == nil {
		return nil
	}
	generic, err := mapper.MapViaJSON[[]*model.APIV4GenericListenerDTO](listeners)
	if err != nil {
		// unreachable: a listener the SDK decoded or mapped encodes back to JSON
		return nil
	}
	return generic
}

// withFalseWhenUnset sets to false the booleans the SDK leaves nil when absent and GKO's former payload model read
// as false: notifyMembers, flowExecution.matchRequired, and each page's homepage and published.
func withFalseWhenUnset(dto DTO) DTO {
	dto.NotifyMembers = falseWhenUnset(dto.NotifyMembers)
	if dto.FlowExecution != nil {
		flowExecution := *dto.FlowExecution
		flowExecution.MatchRequired = falseWhenUnset(flowExecution.MatchRequired)
		dto.FlowExecution = &flowExecution
	}
	if dto.Pages != nil {
		pages := make([]sdk.PageV4, len(dto.Pages))
		for i, page := range dto.Pages {
			page.Homepage = falseWhenUnset(page.Homepage)
			page.Published = falseWhenUnset(page.Published)
			pages[i] = page
		}
		dto.Pages = pages
	}
	return dto
}

func falseWhenUnset(value *bool) *bool {
	if value != nil {
		return value
	}
	return new(false)
}

// withResolvedVisibility resolves the visibility of each page navigation entry as APIM does, with the rule
// the portal navigation resources share.
func withResolvedVisibility(paths []sdk.NavigationPath) []sdk.NavigationPath {
	if len(paths) == 0 {
		return paths
	}
	declared := make(map[string]sdk.PortalVisibility, len(paths))
	for _, entry := range paths {
		if entry.Visibility != nil {
			declared[navigation.NormalizePath(entry.Path)] = *entry.Visibility
		}
	}
	resolved := make([]sdk.NavigationPath, len(paths))
	for i, entry := range paths {
		visibility := navigation.ResolveVisibility(entry.Path, declared)
		entry.Visibility = &visibility
		resolved[i] = entry
	}
	return resolved
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
