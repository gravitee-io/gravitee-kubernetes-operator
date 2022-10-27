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

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

type Client struct {
	ctx     context.Context
	baseUrl string
	orgUrl  string
	envUrl  string
	http    http.Client
}

type AuthenticatedRoundTripper struct {
	apimCtx   *gio.ManagementContext
	transport http.RoundTripper
}

func newAuthenticatedRoundTripper(
	apimCtx *gio.ManagementContext,
	transport http.RoundTripper,
) *AuthenticatedRoundTripper {
	return &AuthenticatedRoundTripper{
		apimCtx, transport,
	}
}

func (t *AuthenticatedRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	t.apimCtx.Authenticate(req)
	return t.transport.RoundTrip(req)
}

func NewClient(ctx context.Context, apimCtx *gio.ManagementContext, httpCli http.Client) *Client {
	baseUrl := strings.TrimSuffix(apimCtx.Spec.BaseUrl, "/")
	orgUrl := baseUrl + "/management/organizations/" + apimCtx.Spec.OrgId
	envUrl := orgUrl + "/environments/" + apimCtx.Spec.EnvId

	authRoundTripper := newAuthenticatedRoundTripper(apimCtx, http.DefaultTransport)
	httpCli.Transport = authRoundTripper

	return &Client{
		ctx, baseUrl, orgUrl, envUrl, httpCli,
	}
}
