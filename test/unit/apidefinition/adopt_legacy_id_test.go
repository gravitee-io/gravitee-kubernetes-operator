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

package apidefinition_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/management"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/apidefinition"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// fakeAPIM answers every request with the same status and body, and records the requests it got.
type fakeAPIM struct {
	mu       sync.Mutex
	status   int
	body     string
	requests []*http.Request
}

func (f *fakeAPIM) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.requests = append(f.requests, r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(f.status)
	_, _ = w.Write([]byte(f.body))
}

func (f *fakeAPIM) received() []*http.Request {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.requests
}

var _ = Describe("AdoptLegacyID", func() {
	// The id and cross id GKO 4.11 gave default/my-api in the DEFAULT organization and environment.
	const (
		legacyID      = "964236d2-7b02-9a17-a40a-ab686e76433a"
		legacyCrossID = "3d7306da-8a12-9e6b-9749-deaab5bcf711"
	)

	var (
		ctx        context.Context
		fake       *fakeAPIM
		apimClient *apim.APIM
		api        *v1alpha1.ApiV4Definition
	)

	BeforeEach(func() {
		ctx = context.Background()
		fake = &fakeAPIM{}
		server := httptest.NewServer(fake)
		DeferCleanup(server.Close)

		mCtx := &v1alpha1.ManagementContext{
			ObjectMeta: metav1.ObjectMeta{Name: "dev-ctx", Namespace: "default"},
			Spec: v1alpha1.ManagementContextSpec{Context: &management.Context{
				BaseUrl: server.URL,
				OrgID:   "DEFAULT",
				EnvID:   "DEFAULT",
				Auth:    &management.Auth{BearerToken: "token"},
			}},
		}
		var err error
		apimClient, err = apim.FromContext(ctx, mCtx, "default")
		Expect(err).ToNot(HaveOccurred())

		api = &v1alpha1.ApiV4Definition{
			ObjectMeta: metav1.ObjectMeta{Name: "my-api", Namespace: "default"},
			Spec: v1alpha1.ApiV4DefinitionSpec{
				Api: v4.Api{V4BaseApi: &v4.V4BaseApi{ApiBase: &base.ApiBase{}}},
			},
		}
	})

	It("re-attaches the API APIM holds at the 4.11 id with the CR's cross id", func() {
		fake.status = http.StatusOK
		fake.body = `{"id": "` + legacyID + `", "crossId": "` + legacyCrossID + `", "name": "my-api"}`

		Expect(apidefinition.AdoptLegacyID(ctx, apimClient, api)).To(Succeed())

		Expect(api.Status.ID).To(Equal(legacyID))
		Expect(api.Status.CrossID).To(Equal(legacyCrossID))
		Expect(fake.received()).To(HaveLen(1))
		request := fake.received()[0]
		Expect(request.Method).To(Equal(http.MethodGet))
		Expect(request.URL.Path).To(HaveSuffix("/apis/" + legacyID))
		Expect(request.URL.Query().Get("hridContainsUUID")).To(Equal("true"))
	})

	It("leaves the CR to the HRID sync when the API at the 4.11 id has another cross id", func() {
		fake.status = http.StatusOK
		fake.body = `{"id": "` + legacyID + `", "crossId": "another-cross-id", "name": "my-api"}`

		Expect(apidefinition.AdoptLegacyID(ctx, apimClient, api)).To(Succeed())

		Expect(api.Status.ID).To(BeEmpty())
		Expect(api.Status.CrossID).To(BeEmpty())
	})

	DescribeTable("leaves the CR to the HRID sync",
		func(status int) {
			fake.status = status
			fake.body = `{"message": "API not found"}`

			Expect(apidefinition.AdoptLegacyID(ctx, apimClient, api)).To(Succeed())

			Expect(api.Status.ID).To(BeEmpty())
		},
		Entry("on a 404, as an organization admin gets it for an unknown id", http.StatusNotFound),
		Entry("on a 403, as a user without rights gets it for an unknown id", http.StatusForbidden),
		Entry("on a 401, which the PUT that follows gets too", http.StatusUnauthorized),
		Entry("on a 400, which the PUT that follows reports", http.StatusBadRequest),
	)

	DescribeTable("fails the lookup on a status that may not last, so that the reconcile is retried",
		func(status int) {
			fake.status = status
			fake.body = `{"message": "unavailable"}`

			err := apidefinition.AdoptLegacyID(ctx, apimClient, api)

			Expect(err).To(MatchError(ContainSubstring(legacyID)))
			Expect(api.Status.ID).To(BeEmpty())
		},
		Entry("408", http.StatusRequestTimeout),
		Entry("429", http.StatusTooManyRequests),
		Entry("500", http.StatusInternalServerError),
		Entry("503", http.StatusServiceUnavailable),
	)

	It("does not call APIM for a CR that already has a status id", func() {
		api.Status.ID = "an-id"

		Expect(apidefinition.AdoptLegacyID(ctx, apimClient, api)).To(Succeed())

		Expect(fake.received()).To(BeEmpty())
		Expect(api.Status.ID).To(Equal("an-id"))
	})
})
