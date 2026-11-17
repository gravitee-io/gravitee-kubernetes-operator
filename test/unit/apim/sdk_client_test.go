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

package apim_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/group"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/management"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimgroup "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/group"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

// recordedRequest is what the fake Automation API received.
type recordedRequest struct {
	Method        string
	Path          string
	Query         url.Values
	Authorization string
	Body          map[string]any
}

// fakeAutomationAPI answers group calls with a fixed state, or a fixed status code.
type fakeAutomationAPI struct {
	*httptest.Server
	mu          sync.Mutex
	last        recordedRequest
	status      int
	response    string
	contentType string
}

const groupState = `{
	"id": "4d1c5f0a-0000-0000-0000-000000000001",
	"hrid": "default-developers",
	"name": "Developers",
	"memberCount": 2,
	"errors": {"severe": ["severe one"], "warning": ["warning one"]}
}`

func newFakeAutomationAPI() *fakeAutomationAPI {
	f := &fakeAutomationAPI{status: http.StatusOK, response: groupState, contentType: "application/json"}
	f.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		f.last = recordedRequest{
			Method:        r.Method,
			Path:          r.URL.Path,
			Query:         r.URL.Query(),
			Authorization: r.Header.Get("Authorization"),
		}
		if data, _ := io.ReadAll(r.Body); len(data) > 0 {
			_ = json.Unmarshal(data, &f.last.Body)
		}
		w.Header().Set("Content-Type", f.contentType)
		w.WriteHeader(f.status)
		if f.status == http.StatusOK {
			_, _ = w.Write([]byte(f.response))
		} else if f.status >= http.StatusBadRequest {
			_, _ = w.Write([]byte(`{"message": "HRID [x] cannot be found.", "httpStatus": 404}`))
		}
	}))
	return f
}

func (f *fakeAutomationAPI) lastRequest() recordedRequest {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.last
}

func managementContext(baseURL string, auth *management.Auth) *v1alpha1.ManagementContext {
	return &v1alpha1.ManagementContext{
		ObjectMeta: metav1.ObjectMeta{Name: "dev", Namespace: "default"},
		Spec: v1alpha1.ManagementContextSpec{Context: &management.Context{
			BaseUrl: baseURL,
			OrgID:   "org1",
			EnvID:   "env1",
			Auth:    auth,
		}},
	}
}

func bearer(token string) *management.Auth {
	return &management.Auth{BearerToken: token}
}

func toDTO(grp *v1alpha1.Group) apimgroup.DTO {
	GinkgoHelper()
	dto, err := apimgroup.ToDTO(grp)
	Expect(err).ToNot(HaveOccurred())
	return dto
}

func developers() *v1alpha1.Group {
	return &v1alpha1.Group{
		ObjectMeta: metav1.ObjectMeta{Name: "developers", Namespace: "default"},
		Spec: v1alpha1.GroupSpec{Type: &group.Type{
			HRID:          "default-developers",
			Name:          "Developers",
			NotifyMembers: true,
		}},
	}
}

// cloudToken is an unsigned JWT carrying the claims the SDK reads; the SDK does not verify signatures.
func cloudToken(org string, envs []string, cpg string) string {
	enc := base64.RawURLEncoding
	header := enc.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	claims, _ := json.Marshal(map[string]any{"org": org, "envs": envs, "cpg": cpg})
	return header + "." + enc.EncodeToString(claims) + "." + enc.EncodeToString([]byte("signature"))
}

var _ = Describe("APIM SDK client factory", func() {
	var (
		ctx  context.Context
		fake *fakeAutomationAPI
	)

	BeforeEach(func() {
		ctx = context.Background()
		fake = newFakeAutomationAPI()
		DeferCleanup(fake.Close)
	})

	newClient := func(mctx *v1alpha1.ManagementContext) *apim.SDKClient {
		GinkgoHelper()
		client, err := apim.NewSDKClient(ctx, mctx)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	upsertGroup := func(mctx *v1alpha1.ManagementContext, grp *v1alpha1.Group) (apimgroup.Response, error) {
		GinkgoHelper()
		return apimgroup.Upsert(ctx, newClient(mctx), toDTO(grp))
	}

	It("targets the context's automation URL, organization and environment with its bearer token", func() {
		resp, err := upsertGroup(managementContext(fake.URL, bearer("token")), developers())
		Expect(err).ToNot(HaveOccurred())

		req := fake.lastRequest()
		Expect(req.Method).To(Equal(http.MethodPut))
		Expect(req.Path).To(Equal("/automation/organizations/org1/environments/env1/groups"))
		Expect(req.Query.Get("dryRun")).To(Equal("false"))
		Expect(req.Query.Has("hridContainsUUID")).To(BeFalse())
		Expect(req.Authorization).To(Equal("Bearer token"))
		Expect(req.Body).To(HaveKeyWithValue("hrid", "default-developers"))
		Expect(req.Body).To(HaveKeyWithValue("name", "Developers"))

		Expect(*resp.State.Id).To(Equal("4d1c5f0a-0000-0000-0000-000000000001"))
		Expect(resp.GetOrgID()).To(Equal("org1"))
		Expect(resp.GetEnvID()).To(Equal("env1"))
	})

	It("sends basic credentials when the context has no bearer token", func() {
		auth := &management.Auth{Credentials: &management.BasicAuth{Username: "admin", Password: "secret"}}
		_, err := upsertGroup(managementContext(fake.URL, auth), developers())
		Expect(err).ToNot(HaveOccurred())

		basic := "Basic " + base64.StdEncoding.EncodeToString([]byte("admin:secret"))
		Expect(fake.lastRequest().Authorization).To(Equal(basic))
	})

	It("sends only the bearer token when secret injection also left empty credentials", func() {
		auth := &management.Auth{BearerToken: "token", Credentials: &management.BasicAuth{}}
		_, err := upsertGroup(managementContext(fake.URL, auth), developers())
		Expect(err).ToNot(HaveOccurred())
		Expect(fake.lastRequest().Authorization).To(Equal("Bearer token"))
	})

	It("uses the context path instead of the automation base path", func() {
		mctx := managementContext(fake.URL, bearer("token"))
		mctx.Spec.Path = new("/custom")
		_, err := upsertGroup(mctx, developers())
		Expect(err).ToNot(HaveOccurred())
		Expect(fake.lastRequest().Path).To(Equal("/custom/organizations/org1/environments/env1/groups"))
	})

	It("fails only the SDK client, not the APIM client, for a context with no usable auth", func() {
		mctx := managementContext(fake.URL, &management.Auth{})

		_, err := apim.FromContext(ctx, mctx, "default")
		Expect(err).ToNot(HaveOccurred())

		_, err = apim.NewSDKClient(ctx, mctx)
		Expect(err).To(MatchError(ContainSubstring("APIM SDK client")))
	})

	It("reads org and env from a cloud token and uses the cloud automation path", func() {
		token := cloudToken("cloud-org", []string{"cloud-env"}, "eu")
		mctx := managementContext(fake.URL, bearer(token))
		mctx.Spec.Cloud = &management.Cloud{Token: token}
		mctx.Spec.OrgID, mctx.Spec.EnvID = "cloud-org", "cloud-env"

		_, err := upsertGroup(mctx, developers())
		Expect(err).ToNot(HaveOccurred())

		req := fake.lastRequest()
		Expect(req.Path).To(Equal("/apim/automation/organizations/cloud-org/environments/cloud-env/groups"))
		Expect(req.Authorization).To(Equal("Bearer " + token))
	})

	It("treats a present but disabled cloud block as a plain context", func() {
		auth := &management.Auth{Credentials: &management.BasicAuth{Username: "admin", Password: "secret"}}
		mctx := managementContext(fake.URL, auth)
		mctx.Spec.Cloud = &management.Cloud{}

		_, err := upsertGroup(mctx, developers())
		Expect(err).ToNot(HaveOccurred())
		Expect(fake.lastRequest().Authorization).To(HavePrefix("Basic "))
	})

	It("addresses a group created before the Automation API by its UUID", func() {
		grp := developers()
		// a legacy group carries its UUID in status
		grp.Status.ID = "4d1c5f0a-0000-0000-0000-000000000001"

		_, err := upsertGroup(managementContext(fake.URL, bearer("token")), grp)
		Expect(err).ToNot(HaveOccurred())

		req := fake.lastRequest()
		Expect(req.Query.Get("hridContainsUUID")).To(Equal("true"))
		Expect(req.Body).To(HaveKeyWithValue("hrid", grp.Status.ID))
	})

	It("asks APIM to leave members alone when the ignore-members annotation is set", func() {
		grp := developers()
		grp.Annotations = map[string]string{core.IgnoreMembersAnnotation: "true"}
		_, err := upsertGroup(managementContext(fake.URL, bearer("token")), grp)
		Expect(err).ToNot(HaveOccurred())
		Expect(fake.lastRequest().Query.Get("ignoreMembers")).To(Equal("true"))
	})

	It("returns the dry run errors APIM reports", func() {
		mctx := managementContext(fake.URL, bearer("token"))
		errs, err := apimgroup.DryRun(ctx, newClient(mctx), toDTO(developers()))
		Expect(err).ToNot(HaveOccurred())
		Expect(fake.lastRequest().Query.Get("dryRun")).To(Equal("true"))
		Expect(errs.Severe).To(ConsistOf(MatchError("severe one")))
		Expect(errs.Warning).To(ConsistOf(MatchError("warning one")))
	})

	It("deletes a legacy group by UUID and treats a group already gone as deleted", func() {
		grp := developers()
		grp.Status.ID = "4d1c5f0a-0000-0000-0000-000000000001"
		fake.status = http.StatusNotFound

		Expect(apimgroup.Delete(ctx, newClient(managementContext(fake.URL, bearer("token"))), toDTO(grp))).To(Succeed())

		req := fake.lastRequest()
		Expect(req.Method).To(Equal(http.MethodDelete))
		Expect(req.Path).To(Equal("/automation/organizations/org1/environments/env1/groups/" + grp.Status.ID))
		Expect(req.Query.Get("hridContainsUUID")).To(Equal("true"))
	})

	It("gets a group by HRID, or a legacy one by UUID", func() {
		client := newClient(managementContext(fake.URL, bearer("token")))

		remote, err := apimgroup.GetRemote(ctx, client, toDTO(developers()))
		Expect(err).ToNot(HaveOccurred())
		Expect(remote.Name).To(Equal("Developers"))
		req := fake.lastRequest()
		Expect(req.Method).To(Equal(http.MethodGet))
		Expect(req.Path).To(Equal("/automation/organizations/org1/environments/env1/groups/default-developers"))
		Expect(req.Query.Get("hridContainsUUID")).To(Equal("false"))

		legacy := developers()
		legacy.Status.ID = "4d1c5f0a-0000-0000-0000-000000000001"
		_, err = apimgroup.GetRemote(ctx, client, toDTO(legacy))
		Expect(err).ToNot(HaveOccurred())
		req = fake.lastRequest()
		Expect(req.Path).To(Equal("/automation/organizations/org1/environments/env1/groups/4d1c5f0a-0000-0000-0000-000000000001"))
		Expect(req.Query.Get("hridContainsUUID")).To(Equal("true"))
	})

	It("reports a successful response without a JSON body as unexpected", func() {
		fake.contentType = "text/plain"
		client := newClient(managementContext(fake.URL, bearer("token")))

		_, err := apimgroup.Upsert(ctx, client, toDTO(developers()))
		Expect(err).To(MatchError(ContainSubstring(`unexpected response: status 200, content type "text/plain"`)))

		_, err = apimgroup.GetRemote(ctx, client, toDTO(developers()))
		Expect(err).To(MatchError(ContainSubstring("unexpected response: status 200")))
	})

	It("deletes a group by HRID", func() {
		grp := developers()
		k8s.AddAutomationAPIManagedCondition(grp)
		fake.status = http.StatusNoContent

		client := newClient(managementContext(fake.URL, bearer("token")))
		Expect(apimgroup.Delete(ctx, client, toDTO(grp))).To(Succeed())

		req := fake.lastRequest()
		Expect(req.Path).To(Equal("/automation/organizations/org1/environments/env1/groups/default-developers"))
		Expect(req.Query.Get("hridContainsUUID")).To(Equal("false"))
	})

	Describe("Probe", func() {
		probe := func() error {
			return newClient(managementContext(fake.URL, bearer("token"))).Probe(ctx)
		}

		It("gets the probe group, and treats not found as reachable", func() {
			fake.status = http.StatusNotFound
			Expect(probe()).To(Succeed())

			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodGet))
			Expect(req.Path).To(Equal("/automation/organizations/org1/environments/env1/groups/gko-probe-not-found"))
		})

		It("succeeds when APIM answers the group", func() {
			Expect(probe()).To(Succeed())
		})

		It("fails when APIM refuses the credentials", func() {
			fake.status = http.StatusUnauthorized
			Expect(probe()).ToNot(Succeed())
		})

		It("fails when APIM rejects the org or env", func() {
			fake.status = http.StatusBadRequest
			Expect(probe()).ToNot(Succeed())
		})

		It("fails when APIM answers with a server error", func() {
			fake.status = http.StatusInternalServerError
			Expect(probe()).ToNot(Succeed())
		})
	})
})
