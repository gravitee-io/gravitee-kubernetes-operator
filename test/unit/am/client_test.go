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

package am_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	amsdk "github.com/gravitee-io-labs/gravitee-automation-tools/am-sdk/v2/pkg/sdk"
	amodel "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	internal "github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/securitydomain"
)

const automationRoot = "/automation/organizations/DEFAULT/environments/DEFAULT"

// amServer answers every request with status, content type and body, and records the last path.
type amServer struct {
	*httptest.Server
	lastPath string
}

func newAMServer(status int, contentType, body string) *amServer {
	s := &amServer{}
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.lastPath = r.URL.Path
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	DeferCleanup(s.Close)
	return s
}

func newAMContext(baseURL string, path *string, auth *amodel.Auth) *v1alpha1.AMContext {
	return &v1alpha1.AMContext{
		Spec: v1alpha1.AMContextSpec{Context: &amodel.Context{
			BaseUrl: baseURL,
			Path:    path,
			OrgID:   "DEFAULT",
			EnvID:   "DEFAULT",
			Auth:    auth,
		}},
	}
}

func tokenAuth() *amodel.Auth { return &amodel.Auth{BearerToken: "admin-token"} }

func clientFor(s *amServer) *am.Client {
	c, err := am.NewSDKClient(context.Background(), newAMContext(s.URL, nil, tokenAuth()))
	Expect(err).ToNot(HaveOccurred())
	return c
}

var _ = Describe("am.NewSDKClient", func() {
	ctx := context.Background()

	It("joins a base URL ending with a slash and the default path", func() {
		s := newAMServer(http.StatusOK, "application/json", "[]")
		c, err := am.NewSDKClient(ctx, newAMContext(s.URL+"/", nil, tokenAuth()))
		Expect(err).ToNot(HaveOccurred())

		Expect(c.Probe(ctx)).To(Succeed())
		Expect(s.lastPath).To(Equal(automationRoot + "/domains"))
	})

	It("joins a path override given without a leading slash", func() {
		s := newAMServer(http.StatusOK, "application/json", "[]")
		c, err := am.NewSDKClient(ctx, newAMContext(s.URL, new("custom"), tokenAuth()))
		Expect(err).ToNot(HaveOccurred())

		Expect(c.Probe(ctx)).To(Succeed())
		Expect(s.lastPath).To(Equal("/custom/organizations/DEFAULT/environments/DEFAULT/domains"))
	})

	It("returns an error instead of panicking when auth is missing", func() {
		_, err := am.NewSDKClient(ctx, newAMContext("http://am.local", nil, nil))
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("AMSecurityDomain AM calls", func() {
	ctx := context.Background()
	dto := amsdk.Domain{Key: "default-my-domain", Name: "my-domain", Path: "/my-domain"}

	It("reports the AM error message on a rejected upsert", func() {
		s := newAMServer(http.StatusBadRequest, "application/json", `{"message":"path already used"}`)

		_, err := internal.Upsert(ctx, clientFor(s), dto)

		Expect(err).To(MatchError(ContainSubstring("path already used")))
	})

	It("returns an error on an upsert answered without a JSON body", func() {
		s := newAMServer(http.StatusOK, "text/plain", "ok")

		_, err := internal.Upsert(ctx, clientFor(s), dto)

		Expect(err).To(HaveOccurred())
	})

	It("rejects a dry run answered without a JSON body", func() {
		s := newAMServer(http.StatusOK, "text/plain", "ok")

		errs := internal.DryRun(ctx, clientFor(s), dto)

		Expect(errs.IsSevere()).To(BeTrue())
	})

	It("returns an error on a get answered without a JSON body", func() {
		s := newAMServer(http.StatusOK, "text/plain", "ok")

		_, err := internal.GetRemote(ctx, clientFor(s), dto)

		Expect(err).To(HaveOccurred())
	})

	It("treats a domain already gone from AM as deleted", func() {
		s := newAMServer(http.StatusNotFound, "application/json", `{"message":"domain not found"}`)

		Expect(internal.Delete(ctx, clientFor(s), dto)).To(Succeed())
	})
})
