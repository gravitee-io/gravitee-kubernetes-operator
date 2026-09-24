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

package am

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gravitee-io-labs/gravitee-automation-tools/am-mock-server/server"
	pkg "github.com/gravitee-io-labs/gravitee-automation-tools/am-sdk/v2/pkg"
	"github.com/gravitee-io-labs/gravitee-automation-tools/common/pkg/apicontext"
	"github.com/gravitee-io-labs/gravitee-automation-tools/common/pkg/auth"
)

const (
	Port     = 30093
	BasePath = server.BasePath
)

type MockServer struct {
	srv    *http.Server
	MockAM *server.MockAM
}

func Start(dryRunReject bool) *MockServer {
	mockAM := server.NewMockAM()
	reg := auth.NewRegistry(AuthConfig(), server.BasePath)
	handler := server.NewWithPath(mockAM, server.BasePath, reg, dryRunReject)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", Port),
		Handler: handler,
	}
	go srv.ListenAndServe() //nolint:errcheck
	// give the server a moment to bind
	time.Sleep(100 * time.Millisecond)
	return &MockServer{srv: srv, MockAM: mockAM}
}

func (m *MockServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	m.srv.Shutdown(ctx) //nolint:errcheck
}

func NewSDKClient() *pkg.AMClient {
	token := "admin-token"
	ac := apicontext.APIContext{
		BaseURL: fmt.Sprintf("http://localhost:%d%s", Port, BasePath),
		OrgID:   "DEFAULT",
		EnvID:   "DEFAULT",
		Auth: apicontext.Auth{
			BearerToken: &token,
		},
	}
	client, err := pkg.NewClient(ac, 5000)
	if err != nil {
		panic(fmt.Sprintf("failed to create AM SDK client: %v", err))
	}
	return client
}

func AuthConfig() auth.Config {
	return auth.Config{
		Users: map[string]auth.UserConfig{
			"admin": {
				Token:       "admin-token",
				Permissions: []string{},
			},
			"readonly": {
				Token: "readonly-token",
				Permissions: []string{
					"DOMAIN_LIST",
					"DOMAIN_READ",
				},
			},
		},
		Permissions: []auth.RoutePermission{
			{
				Path:   "/organizations/{orgId}/environments/{envId}/domains",
				Get:    "DOMAIN_LIST",
				Put:    "DOMAIN_UPDATE",
			},
			{
				Path:   "/organizations/{orgId}/environments/{envId}/domains/{domainKey}",
				Get:    "DOMAIN_READ",
				Delete: "DOMAIN_DELETE",
			},
		},
	}
}
