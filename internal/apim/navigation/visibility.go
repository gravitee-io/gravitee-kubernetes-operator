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

// Package navigation holds what the portal navigation resources share on the APIM SDK types.
package navigation

import (
	"strings"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
)

// WithResolvedTopNavbar returns a new slice in which every entry carries the visibility APIM will store
// for it, so that drift detection can compare the field exactly instead of ignoring it. The input is
// never written to.
//
// It is meant for the drift path only. The payload sent to APIM keeps an undeclared visibility omitted:
// APIM resolves it against the persisted parent, and sending an explicit PUBLIC under a PRIVATE parent
// would be rejected by its parent-must-be-public-if-child-is-public rule.
func WithResolvedTopNavbar(entries []sdk.PortalNavigationPath) []sdk.PortalNavigationPath {
	if len(entries) == 0 {
		return entries
	}
	declared := make(map[string]sdk.PortalVisibility, len(entries))
	for _, entry := range entries {
		if entry.Visibility != nil {
			declared[NormalizePath(entry.Path)] = *entry.Visibility
		}
	}
	resolved := make([]sdk.PortalNavigationPath, len(entries))
	for i, entry := range entries {
		visibility := ResolveVisibility(entry.Path, declared)
		entry.Visibility = &visibility
		resolved[i] = entry
	}
	return resolved
}

// ResolveVisibility mirrors PortalVisibility.resolve on the APIM side: the entry's own declaration wins,
// then the closest declared ancestor, then PUBLIC. Ancestors the spec does not list are folders APIM
// creates implicitly; they declare nothing, so they inherit too and the walk simply continues upwards.
func ResolveVisibility(path string, declared map[string]sdk.PortalVisibility) sdk.PortalVisibility {
	segments := strings.Split(strings.Trim(NormalizePath(path), "/"), "/")
	for i := len(segments); i > 0; i-- {
		ancestor := "/" + strings.Join(segments[:i], "/")
		if visibility, ok := declared[ancestor]; ok && visibility != "" {
			return visibility
		}
	}
	return sdk.PortalVisibilityPUBLIC
}

// NormalizePath drops a trailing slash so that "/a" and "/a/" describe the same node. The CRD only
// enforces the leading slash.
func NormalizePath(path string) string {
	trimmed := strings.TrimRight(path, "/")
	if trimmed == "" {
		return "/"
	}
	return trimmed
}
