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

package mapper

import (
	"fmt"
	"strings"

	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func extractURLRewriteFilter(filters []gwAPIv1.HTTPRouteFilter) *gwAPIv1.HTTPURLRewriteFilter {
	for _, f := range filters {
		if f.URLRewrite != nil {
			return f.URLRewrite
		}
	}
	return nil
}

func getRewrittenPath(path *gwAPIv1.HTTPPathModifier) string {
	if path.ReplaceFullPath != nil {
		return *path.ReplaceFullPath
	}

	prefixMatch := strings.TrimPrefix(*path.ReplacePrefixMatch, "/")
	prefixMatchWithoutTrailingSlash := strings.TrimSuffix(prefixMatch, "/")

	if prefixMatchWithoutTrailingSlash == "" {
		return "{#request.pathInfo}"
	}

	rewrittenPath := fmt.Sprintf(locationPathWhenReplacePrefix, prefixMatchWithoutTrailingSlash)

	return rewrittenPath
}
