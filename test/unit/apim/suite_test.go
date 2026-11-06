package apim_test

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

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/rest"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	//+kubebuilder:scaffold:imports
)

func TestAPIM(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "APIM unit tests suite")
}

// The dynamic client resolves ManagementContexts and falls back to the ambient kubeconfig, which
// reaches a developer's cluster and exits where there is none (CI). Every lookup answers 404 here.
var _ = BeforeSuite(func() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"kind":"Status","apiVersion":"v1","status":"Failure",` +
			`"reason":"NotFound","code":404,"message":"not found"}`))
	}))
	DeferCleanup(server.Close)
	dynamic.UseConfig(&rest.Config{Host: server.URL})
})
