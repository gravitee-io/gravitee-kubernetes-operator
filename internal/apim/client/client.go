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

package client

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
)

const (
	orgPath = "/management/organizations/"
	envPath = "/environments/"
)

// Client is the client for a given instance of the Gravitee.io Management API
// The client is created once per reconcile and management context and reused for all the operations
// of a reconcile cycle, using the reconcile context.Context.
type Client struct {
	// HTTP is the HTTP client used to communicate with the Gravitee.io Management API
	HTTP *http.Client
	// URLs contains URLs targeting the organization and environment of the client
	URLs *URLs
}

// URLs contains URLs targeting the organization and environment of the client.
type URLs struct {
	Org *http.URL
	Env *http.URL
}

// EnvTarget returns a new URL with the given path appended to the environment URL.
func (client *Client) EnvTarget(path string) *http.URL {
	return client.URLs.Env.WithPath(path)
}

// OrgTarget returns a new URL with the given path appended to the organization URL.
func (client *Client) OrgTarget(path string) *http.URL {
	return client.URLs.Org.WithPath(path)
}

// NewURLs returns a new URLs instance for the given base URL
// with Org path initialized from the given orgID and Env path initialized from the given envID.
func NewURLs(baseUrl string, orgID, envID string) (*URLs, error) {
	base, err := http.NewURL(baseUrl)
	if err != nil {
		return nil, err
	}

	org := base.WithPath(orgPath, orgID)
	env := org.WithPath(envPath, envID)

	return &URLs{org, env}, nil
}

// NewClient returns a new client for the given management context.
// The client is created once per reconcile and management context and reused for all the operations.
func NewClient(ctx context.Context, management *model.Management) (*Client, error) {
	orgID, envID := management.OrgId, management.EnvId
	urls, err := NewURLs(management.BaseUrl, orgID, envID)
	if err != nil {
		return nil, err
	}

	client := &Client{
		HTTP: http.NewClient(ctx, toHttpAuth(management)),
		URLs: urls,
	}

	return client, nil
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
