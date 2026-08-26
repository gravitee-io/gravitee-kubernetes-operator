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
	"net"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	ammodel "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/amctx"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type probeServer struct {
	status int
	auth   string
	path   string
	query  string
	server *httptest.Server
}

func newProbeServer(status int) *probeServer {
	p := &probeServer{status: status}
	p.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p.auth = r.Header.Get("Authorization")
		p.path = r.URL.Path
		p.query = r.URL.RawQuery
		w.WriteHeader(p.status)
	}))
	return p
}

func validAMContext(baseUrl, token string) *v1alpha1.AMContext {
	return &v1alpha1.AMContext{
		ObjectMeta: metav1.ObjectMeta{Name: "test-ctx", Namespace: "default"},
		Spec: v1alpha1.AMContextSpec{
			Context: &ammodel.Context{
				BaseUrl: baseUrl,
				OrgID:   "DEFAULT",
				EnvID:   "DEFAULT",
				Auth:    &ammodel.Auth{BearerToken: token},
			},
		},
	}
}

var _ = Describe("AMContext admission", func() {
	ctrl := amctx.AdmissionCtrl{}
	ctx := context.Background()

	DescribeTable("rejects a missing or invalid required field",
		func(mutate func(*v1alpha1.AMContext), field, extra string) {
			obj := validAMContext("http://am.example", "tok")
			mutate(obj)
			_, err := ctrl.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(field))
			Expect(err.Error()).To(ContainSubstring(extra))
		},
		Entry("nil context", func(o *v1alpha1.AMContext) { o.Spec.Context = nil }, "[baseUrl]", "mandatory"),
		Entry("empty baseUrl", func(o *v1alpha1.AMContext) { o.Spec.BaseUrl = "" }, "[baseUrl]", "mandatory"),
		Entry("blank baseUrl", func(o *v1alpha1.AMContext) { o.Spec.BaseUrl = "  " }, "[baseUrl]", "mandatory"),
		Entry("not a URL", func(o *v1alpha1.AMContext) { o.Spec.BaseUrl = "not-a-url" }, "[baseUrl]", "not a valid URL"),
		Entry("empty org", func(o *v1alpha1.AMContext) { o.Spec.OrgID = "" }, "[orgId]", "mandatory"),
		Entry("empty env", func(o *v1alpha1.AMContext) { o.Spec.EnvID = "" }, "[envId]", "mandatory"),
		Entry("missing auth", func(o *v1alpha1.AMContext) { o.Spec.Auth = nil }, "[auth]", "mandatory"),
	)

	DescribeTable("maps AM HTTP status to admission outcome",
		func(status int, wantErr, wantWarning string) {
			srv := newProbeServer(status)
			DeferCleanup(srv.server.Close)

			warnings, err := ctrl.ValidateCreate(ctx, validAMContext(srv.server.URL, "tok"))
			if wantErr == "" {
				Expect(err).ToNot(HaveOccurred())
			} else {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(wantErr))
			}
			if wantWarning == "" {
				Expect(warnings).To(BeEmpty())
			} else {
				Expect(warnings).To(ContainElement(ContainSubstring(wantWarning)))
			}
		},
		Entry("200 admits", http.StatusOK, "", ""),
		Entry("401 is bad credentials", http.StatusUnauthorized, "bad credentials", ""),
		Entry("400 is unknown org/env", http.StatusBadRequest, "invalid organization or environment", ""),
		Entry("403 is unknown org/env", http.StatusForbidden, "invalid organization or environment", ""),
		Entry("404 is unknown org/env", http.StatusNotFound, "invalid organization or environment", ""),
	)

	It("sends a Bearer token and probes /domains?size=1", func() {
		srv := newProbeServer(http.StatusOK)
		DeferCleanup(srv.server.Close)

		warnings, err := ctrl.ValidateCreate(ctx, validAMContext(srv.server.URL, "s3cret"))
		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(BeEmpty())
		Expect(srv.auth).To(Equal("Bearer s3cret"))
		Expect(srv.path).To(Equal("/automation/organizations/DEFAULT/environments/DEFAULT/domains"))
		Expect(srv.query).To(Equal("size=1"))
	})

	It("uses spec.path instead of /automation", func() {
		srv := newProbeServer(http.StatusOK)
		DeferCleanup(srv.server.Close)

		obj := validAMContext(srv.server.URL, "tok")
		path := "/custom"
		obj.Spec.Path = &path

		_, err := ctrl.ValidateCreate(ctx, obj)
		Expect(err).ToNot(HaveOccurred())
		Expect(srv.path).To(Equal("/custom/organizations/DEFAULT/environments/DEFAULT/domains"))
	})

	It("accepts with a warning when AM is unreachable", func() {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		Expect(err).ToNot(HaveOccurred())
		baseUrl := "http://" + ln.Addr().String()
		Expect(ln.Close()).To(Succeed())

		warnings, err := ctrl.ValidateCreate(ctx, validAMContext(baseUrl, "tok"))
		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(ContainElement(ContainSubstring("unable to reach AM")))
		Expect(warnings).To(ContainElement(ContainSubstring(baseUrl)))
	})

	It("skips validation when the object is being deleted", func() {
		obj := validAMContext("not-a-url", "")
		now := metav1.NewTime(time.Now())
		obj.DeletionTimestamp = &now

		warnings, err := ctrl.ValidateUpdate(ctx, obj, obj)
		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(BeEmpty())
	})
})
