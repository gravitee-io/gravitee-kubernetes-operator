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

// APIM wraps services needed to sync resources with a given environment on a Gravitee.io APIM instance.
type APIM struct {
	APIs *service.APIs

	orgID string
	envID string
}

// EnvID returns the environment ID of the current managed APIM instance.
func (apim *APIM) EnvID() string {
	return apim.envID
}

// OrgID returns the organization ID of the current managed APIM instance.
func (apim *APIM) OrgID() string {
	return apim.orgID
}

// FromContext returns a new APIM instance from a given reconcile context and management context.
func FromContext(ctx context.Context, managementContext model.Context) (*APIM, error) {
	orgID, envID := managementContext.OrgId, managementContext.EnvId
	urls, err := client.NewURLs(managementContext.BaseUrl, orgID, envID)
	if err != nil {
		return nil, err
	}

	client := &client.Client{
		HTTP: http.NewClient(ctx, toHttpAuth(managementContext)),
		URLs: urls,
	}

	return &APIM{
		APIs:  service.NewAPIs(client),
		orgID: orgID,
		envID: envID,
	}, nil
}

func toHttpAuth(management model.Context) *http.Auth {
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
