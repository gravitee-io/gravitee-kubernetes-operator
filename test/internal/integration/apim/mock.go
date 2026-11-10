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
	"fmt"
	"net/http"
	"time"

	"github.com/gravitee-io/gravitee-automation-sdk/apim-mock-server/server"
	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-automation-sdk/common/pkg/apicontext"
	"github.com/gravitee-io/gravitee-automation-sdk/common/pkg/auth"
)

const (
	// MockPort is where the in-process APIM mock listens; the AM mock uses 30093.
	MockPort = 30094
	// MockBasePath is the Automation API base path the mock serves.
	MockBasePath = "/automation"
	// MockToken is the admin bearer token the mock accepts.
	MockToken = "admin-token"
	// MockOrg and MockEnv are the organization and environment the mock contexts target.
	MockOrg = "DEFAULT"
	MockEnv = "DEFAULT"
)

// MockServer is an in-process APIM Automation API mock, for envtest suites.
type MockServer struct {
	srv *http.Server
}

// StartMock serves the APIM mock on MockPort. With dryRunReject, a dry-run PUT answers its
// errors.severe and persists nothing.
func StartMock(dryRunReject bool) *MockServer {
	reg := auth.NewRegistry(mockAuthConfig(), MockBasePath)
	handler := server.NewWithPath(server.NewMockAPIM(dryRunReject), MockBasePath, reg)
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", MockPort),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go srv.ListenAndServe() //nolint:errcheck
	// give the server a moment to bind
	time.Sleep(100 * time.Millisecond)
	return &MockServer{srv: srv}
}

func (m *MockServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	m.srv.Shutdown(ctx) //nolint:errcheck
}

// NewMockSDKClient is a raw SDK client on the mock, for assertions on what the operator stored.
func NewMockSDKClient() *sdk.APIMClient {
	token := MockToken
	client, err := sdk.NewAPIMClient(sdk.APIContext{APIContext: apicontext.APIContext{
		BaseURL: fmt.Sprintf("http://localhost:%d%s", MockPort, MockBasePath),
		OrgID:   MockOrg,
		EnvID:   MockEnv,
		Auth:    apicontext.Auth{BearerToken: &token},
	}})
	if err != nil {
		panic(fmt.Sprintf("failed to create APIM SDK client: %v", err))
	}
	return client
}

// mockAuthConfig declares one admin user; empty permissions grant every route.
func mockAuthConfig() auth.Config {
	return auth.Config{
		Users: map[string]auth.UserConfig{
			"admin": {Token: MockToken},
		},
	}
}
