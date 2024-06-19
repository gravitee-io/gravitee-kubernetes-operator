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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
)

const (
	basePath = "/management"
	envPath  = "/environments"
	orgPath  = "/organizations"
)

// Client is the client for a given instance of the Gravitee.io Management API
// The client is created once per reconcile and API context and reused for all the operations
// of a reconcile cycle, using the reconcile context.Context.
type Client struct {
	// HTTP is the HTTP client used to communicate with the Gravitee.io Management API
	HTTP *http.Client
	// URLs contains URLs targeting the organization and environment of the client
	URLs *URLs
}

// URLs contains URLs targeting the organization and environment of the client.
type URLs struct {
	Org   http.URL
	EnvV1 http.URL
	EnvV2 http.URL
}

// EnvV1Target returns a new URL with the given path appended to the environment URL.
func (client *Client) EnvV1Target(path string) http.URL {
	return client.URLs.EnvV1.WithPath(path)
}

// EnvV2Target returns a new URL with the given path appended to the environment URL.
func (client *Client) EnvV2Target(path string) http.URL {
	return client.URLs.EnvV2.WithPath(path)
}

// OrgTarget returns a new URL with the given path appended to the organization URL.
func (client *Client) OrgTarget(path string) http.URL {
	return client.URLs.Org.WithPath(path)
}

// NewURLs returns a new URLs instance for the given base URL
// with Org path initialized from the given orgID and Env path initialized from the given envID.
func NewURLs(baseUrl string, orgID, envID string) (*URLs, error) {
	base, err := http.NewURL(baseUrl)
	if err != nil {
		return nil, err
	}

	// /management
	root := base.WithPath(basePath)

	// /organizations
	org := base.WithPath(orgPath, orgID)

	// /management/organizations/{orgID}/environments/{envID}
	envV1 := root.WithPath(orgPath, orgID).WithPath(envPath, envID)

	// /management/v2/organizations/{orgID}/environments/{envID}
	envV2 := root.WithPath("v2").WithPath(orgPath, orgID).WithPath(envPath, envID)

	return &URLs{org, envV1, envV2}, nil
}
