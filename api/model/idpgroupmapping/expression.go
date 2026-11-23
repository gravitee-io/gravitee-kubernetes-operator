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

package idpgroupmapping

import "regexp"

// groupExpressionPattern matches what the Gravitee expression language engine evaluates
// in an identity provider group mapping value: "{#", "{ #", "{T(" or "{(".
var groupExpressionPattern = regexp.MustCompile(`\{ *[#T(]`)

// IsGroupExpression reports whether a group mapping value is an expression resolving
// to group names at login (e.g. {#jsonPath(#profile, '$.groups')}) rather than a group name.
// Expressions are sent to APIM as-is, names are resolved to group IDs.
func IsGroupExpression(value string) bool {
	return groupExpressionPattern.MatchString(value)
}
