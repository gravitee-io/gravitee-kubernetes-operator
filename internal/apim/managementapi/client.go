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

package managementapi

import (
	"context"
	"net/http"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
)

type Client struct {
	ctx     context.Context
	baseUrl string
	orgId   string
	envId   string
	http    http.Client
}

func (client *Client) orgUrl() string {
	return client.baseUrl + "/management/organizations/" + client.orgId
}

func (client *Client) envUrl() string {
	return client.orgUrl() + "/environments/" + client.envId
}

func (client *Client) GetEnvId() string {
	return client.envId
}

func (client *Client) GetOrgId() string {
	return client.orgId
}

type AuthenticatedRoundTripper struct {
	management *model.Management
	transport  http.RoundTripper
}

func newAuthenticatedRoundTripper(
	management *model.Management,
	transport http.RoundTripper,
) *AuthenticatedRoundTripper {
	return &AuthenticatedRoundTripper{
		management, transport,
	}
}

func (t *AuthenticatedRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	t.management.Authenticate(req)
	return t.transport.RoundTrip(req)
}

func NewClient(ctx context.Context, management *model.Management, httpClient http.Client) *Client {
	baseUrl := strings.TrimSuffix(management.BaseUrl, "/")

	authRoundTripper := newAuthenticatedRoundTripper(management, http.DefaultTransport)
	httpClient.Transport = authRoundTripper

	orgId, envId := management.OrgId, management.EnvId

	return &Client{
		ctx, baseUrl, orgId, envId, httpClient,
	}
}
