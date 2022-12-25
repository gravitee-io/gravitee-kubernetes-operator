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

package http

import "net/http"

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
)

type BearerToken string

func (b BearerToken) String() string {
	return BearerPrefix + string(b)
}

type BasicAuth struct {
	Username string
	Password string
}

type Auth struct {
	Basic *BasicAuth
	Token BearerToken
}

type AuthenticatedRoundTripper struct {
	auth      *Auth
	transport http.RoundTripper
}

func (auth *Auth) Authenticate(req *http.Request) {
	bearer := auth.Token
	basic := auth.Basic

	if bearer != "" {
		req.Header.Add(AuthorizationHeader, bearer.String())
	} else if basic != nil {
		req.SetBasicAuth(basic.Username, basic.Password)
	}
}

func (t *AuthenticatedRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	t.auth.Authenticate(req)
	return t.transport.RoundTrip(req)
}

func NewAuthenticatedRoundTripper(
	auth *Auth,
	transport http.RoundTripper,
) *AuthenticatedRoundTripper {
	return &AuthenticatedRoundTripper{
		auth, transport,
	}
}
