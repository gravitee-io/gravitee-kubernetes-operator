// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package authzpolicy_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

func TestAuthzPolicyAdmission(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AuthzPolicy admission hook tests suite")
}

// recordingCache keeps the index functions InitCache registers, for the fake client to use.
type recordingCache struct {
	cache.Cache
	indexers map[string]client.IndexerFunc
}

func (c *recordingCache) IndexField(_ context.Context, _ client.Object, field string, f client.IndexerFunc) error {
	c.indexers[field] = f
	return nil
}

// platform records what the APIM stub receives and answers what the running spec sets: 500 unless
// a spec calls answer. A GET, drift detection's read, gets what answerGet sets instead, when set,
// and is not recorded.
var platform = &platformStub{}

type platformStub struct {
	mu        sync.Mutex
	status    int
	body      string
	getStatus int
	getBody   string
	calls     []string
}

func (p *platformStub) reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.status, p.body, p.calls = http.StatusInternalServerError, "", nil
	p.getStatus, p.getBody = 0, ""
}

func (p *platformStub) answer(status int, body string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.status, p.body = status, body
}

func (p *platformStub) answerGet(status int, body string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.getStatus, p.getBody = status, body
}

// requests returns the bodies received, one per call.
func (p *platformStub) requests() []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]string(nil), p.calls...)
}

func (p *platformStub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	p.mu.Lock()
	defer p.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet && p.getStatus != 0 {
		w.WriteHeader(p.getStatus)
		_, _ = w.Write([]byte(p.getBody))
		return
	}
	p.calls = append(p.calls, string(body))
	w.WriteHeader(p.status)
	_, _ = w.Write([]byte(p.body))
}

// The dynamic client resolves ManagementContexts: it answers the dev-ctx context of the travel
// namespace, pointing at the APIM stub, and 404 for anything else.
var _ = BeforeSuite(func() {
	drift.Init()

	apim := httptest.NewServer(platform)
	DeferCleanup(apim.Close)

	cluster := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/apis/gravitee.io/v1alpha1/namespaces/travel/managementcontexts/dev-ctx" {
			_, _ = w.Write([]byte(`{"apiVersion":"gravitee.io/v1alpha1","kind":"ManagementContext",` +
				`"metadata":{"name":"dev-ctx","namespace":"travel"},` +
				`"spec":{"baseUrl":"` + apim.URL + `","organizationId":"DEFAULT","environmentId":"DEFAULT",` +
				`"auth":{"bearerToken":"test"}}}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"kind":"Status","apiVersion":"v1","status":"Failure",` +
			`"reason":"NotFound","code":404,"message":"not found"}`))
	}))
	DeferCleanup(cluster.Close)
	dynamic.UseConfig(&rest.Config{Host: cluster.URL})
})
