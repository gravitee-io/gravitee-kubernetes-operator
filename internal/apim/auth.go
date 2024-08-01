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

package apim

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
)

func toHttpAuth(management core.ContextModel) *http.Auth {
	if !management.HasAuthentication() {
		return nil
	}

	return &http.Auth{
		Basic: toBasicAuth(management.GetAuth()),
		Token: toBearer(management.GetAuth()),
	}
}

func toBasicAuth(auth core.Auth) *http.BasicAuth {
	if auth == nil || !auth.HasCredentials() {
		return nil
	}

	return &http.BasicAuth{
		Username: auth.GetCredentials().GetUsername(),
		Password: auth.GetCredentials().GetPassword(),
	}
}

func toBearer(auth core.Auth) http.BearerToken {
	if auth == nil {
		return ""
	}

	return http.BearerToken(auth.GetBearerToken())
}
