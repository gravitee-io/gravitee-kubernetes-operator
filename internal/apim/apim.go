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
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/service"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
)

type APIM struct {
	APIs *service.APIs

	orgID string
	envID string
}

func (apim *APIM) EnvID() string {
	return apim.envID
}

func (apim *APIM) OrgID() string {
	return apim.orgID
}

// FromContext returns a new APIM instance from a given reconcile context.
func FromContext(ctx context.Context, management *model.Management) (*APIM, error) {
	orgID, envID := management.OrgId, management.EnvId
	urls, err := client.NewURLs(management.BaseUrl, orgID, envID)
	if err != nil {
		return nil, err
	}

	client := &client.Client{
		HTTP: http.NewClient(ctx, toHttpAuth(management)),
		URLs: urls,
	}

	return &APIM{
		APIs:  service.NewAPIs(client),
		orgID: orgID,
		envID: envID,
	}, nil
}

func toHttpAuth(management *model.Management) *http.Auth {
	if !management.HasAuthentication() {
		return nil
	}

	return &http.Auth{
		Basic: toBasicAuth(management.Auth),
		Token: toBearer(management.Auth),
	}
}

func toBasicAuth(auth *model.Auth) *http.BasicAuth {
	if auth == nil || auth.Credentials == nil {
		return nil
	}

	return &http.BasicAuth{
		Username: auth.Credentials.Username,
		Password: auth.Credentials.Password,
	}
}

func toBearer(auth *model.Auth) http.BearerToken {
	if auth == nil {
		return ""
	}

	return http.BearerToken(auth.BearerToken)
}
