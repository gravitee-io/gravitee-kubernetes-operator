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
	"errors"
	"net/http"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimapp "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

const (
	applicationUUID  = "8f1e2d3c-0000-0000-0000-000000000001"
	applicationState = `{
	"id": "8f1e2d3c-0000-0000-0000-000000000001",
	"hrid": "default-app",
	"name": "app",
	"groups": ["developers"],
	"errors": {"severe": ["severe one"], "warning": ["warning one"]}
}`
	applicationsPath = "/automation/organizations/org1/environments/env1/applications"
)

func newApplication(statusID, specID string, managed bool) *v1alpha1.Application {
	app := &v1alpha1.Application{ObjectMeta: metav1.ObjectMeta{Name: "app", Namespace: "default"}}
	app.Spec.Name = "app"
	app.Spec.ID = specID
	app.Status.ID = statusID
	if managed {
		k8s.AddAutomationAPIManagedCondition(app)
	}
	return app
}

var _ = Describe("Application addressing", func() {
	DescribeTable("ToDTO sends the hrid, or the UUID of an application created before the Automation API",
		func(app *v1alpha1.Application, hrid string, legacy bool) {
			dto, err := apimapp.ToDTO(app)

			Expect(err).ToNot(HaveOccurred())
			Expect(dto.Hrid).To(Equal(hrid))
			Expect(dto.Legacy.HridContainsUUID).To(Equal(legacy))
		},
		Entry("Automation API managed", newApplication(applicationUUID, "", true), "default-app", false),
		Entry("legacy, known by its status id", newApplication(applicationUUID, "other", false), applicationUUID, true),
		Entry("legacy, known only by its spec id", newApplication("", applicationUUID, false), applicationUUID, true),
		Entry("new", newApplication("", "", false), "default-app", false),
	)

	It("flags the ignore-groups annotation", func() {
		app := newApplication("", "", false)
		app.Annotations = map[string]string{core.IgnoreGroupsAnnotation: "true"}

		dto, err := apimapp.ToDTO(app)

		Expect(err).ToNot(HaveOccurred())
		Expect(dto.IgnoreGroups).To(BeTrue())
	})
})

var _ = Describe("Application SDK calls", func() {
	var (
		ctx    context.Context
		fake   *fakeAutomationAPI
		client *apim.SDKClient
	)

	BeforeEach(func() {
		ctx = context.Background()
		fake = newFakeAutomationAPI()
		fake.response = applicationState
		DeferCleanup(fake.Close)

		var err error
		client, err = apim.NewSDKClient(ctx, managementContext(fake.URL, bearer("token")))
		Expect(err).ToNot(HaveOccurred())
	})

	applicationDTO := func(legacy bool, ignoreGroups bool) apimapp.DTO {
		hrid := "default-app"
		if legacy {
			hrid = applicationUUID
		}
		return apimapp.DTO{
			Application:  sdk.Application{Hrid: hrid, Name: "app", Groups: []string{"developers"}},
			Legacy:       apimapp.Legacy{HridContainsUUID: legacy},
			IgnoreGroups: ignoreGroups,
		}
	}

	DescribeTable("Upsert puts the application with its flags",
		func(dto apimapp.DTO, want map[string]string) {
			_, err := apimapp.Upsert(ctx, client, dto)
			Expect(err).ToNot(HaveOccurred())

			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodPut))
			Expect(req.Path).To(Equal(applicationsPath))
			Expect(req.Body).To(HaveKeyWithValue("hrid", dto.Hrid))
			Expect(req.Query).To(HaveLen(len(want)))
			for flag, value := range want {
				Expect(req.Query.Get(flag)).To(Equal(value), flag)
			}
		},
		Entry("hrid-addressed", applicationDTO(false, false), map[string]string{"dryRun": "false"}),
		Entry("created before the Automation API",
			applicationDTO(true, false), map[string]string{"dryRun": "false", "hridContainsUUID": "true"}),
		Entry("with the ignore-groups annotation",
			applicationDTO(false, true), map[string]string{"dryRun": "false", "ignoreGroups": "true"}),
	)

	It("dry runs the application and returns its severe errors, without the warnings", func() {
		findings, err := apimapp.DryRun(ctx, client, applicationDTO(false, false))

		Expect(err).ToNot(HaveOccurred())
		Expect(fake.lastRequest().Query.Get("dryRun")).To(Equal("true"))
		Expect(findings.Severe).To(ConsistOf(MatchError("severe one")))
		Expect(findings.Warning).To(BeEmpty())
	})

	It("dry runs the application and returns its warnings when APIM finds nothing severe", func() {
		fake.response = `{"hrid": "default-app", "name": "app", "errors": {"warning": ["warning one"]}}`

		findings, err := apimapp.DryRun(ctx, client, applicationDTO(false, false))

		Expect(err).ToNot(HaveOccurred())
		Expect(findings.Severe).To(BeEmpty())
		Expect(findings.Warning).To(ConsistOf(MatchError("warning one")))
	})

	DescribeTable("Delete removes the application by its hrid, or a legacy one by its UUID",
		func(legacy bool, path string, query map[string]string) {
			fake.status = http.StatusNoContent
			Expect(apimapp.Delete(ctx, client, applicationDTO(legacy, false))).To(Succeed())

			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodDelete))
			Expect(req.Path).To(Equal(path))
			Expect(req.Query).To(HaveLen(len(query)))
			for flag, value := range query {
				Expect(req.Query.Get(flag)).To(Equal(value), flag)
			}
		},
		Entry("hrid-addressed", false, applicationsPath+"/default-app", map[string]string{}),
		Entry("created before the Automation API",
			true, applicationsPath+"/"+applicationUUID, map[string]string{"hridContainsUUID": "true"}),
	)

	It("counts an application already gone from APIM as deleted", func() {
		fake.status = http.StatusNotFound
		Expect(apimapp.Delete(ctx, client, applicationDTO(false, false))).To(Succeed())
	})

	DescribeTable("GetRemote leaves the remote groups out under ignore-groups, as the CR's are",
		func(ignoreGroups bool, groups []string) {
			remote, err := apimapp.GetRemote(ctx, client, applicationDTO(false, ignoreGroups))

			Expect(err).ToNot(HaveOccurred())
			Expect(fake.lastRequest().Path).To(Equal(applicationsPath + "/default-app"))
			Expect(remote.Groups).To(Equal(groups))
		},
		Entry("groups managed by the CR", false, []string{"developers"}),
		Entry("with the ignore-groups annotation", true, nil),
	)

	It("reports a successful response without a JSON body as unexpected", func() {
		fake.contentType = "text/plain"

		_, err := apimapp.Upsert(ctx, client, applicationDTO(false, false))

		Expect(err).To(MatchError(ContainSubstring(`unexpected response: status 200, content type "text/plain"`)))
	})
})

var _ = Describe("Application UpdateStatus", func() {
	state := func(id *string) sdk.ApplicationState {
		return sdk.ApplicationState{Hrid: "default-app", Id: id, Errors: &sdk.Errors{Warning: []string{"warning one"}}}
	}

	It("copies the id, org, env and errors APIM answers, and marks the application Automation API managed", func() {
		app := newApplication("", "", false)
		resp := apimapp.Response{State: state(new(applicationUUID)), OrgID: "org1", EnvID: "env1"}

		Expect(apimapp.UpdateStatus(context.Background(), app, resp)).To(Succeed())

		Expect(app.Status.ID).To(Equal(applicationUUID))
		Expect(app.Status.OrgID).To(Equal("org1"))
		Expect(app.Status.EnvID).To(Equal("env1"))
		Expect(app.Status.Errors.Warning).To(Equal([]string{"warning one"}))
		Expect(k8s.IsAutomationAPIManaged(app)).To(BeTrue())
	})

	It("leaves an application upserted by its UUID addressed by its UUID", func() {
		app := newApplication(applicationUUID, "", false)
		resp := apimapp.Response{State: state(new(applicationUUID)), Legacy: apimapp.Legacy{HridContainsUUID: true}}

		Expect(apimapp.UpdateStatus(context.Background(), app, resp)).To(Succeed())

		Expect(k8s.IsAutomationAPIManaged(app)).To(BeFalse())
	})

	It("keeps the id it has when APIM answers none", func() {
		app := newApplication(applicationUUID, "", false)
		resp := apimapp.Response{State: state(nil), Legacy: apimapp.Legacy{HridContainsUUID: true}}

		Expect(apimapp.UpdateStatus(context.Background(), app, resp)).To(Succeed())

		Expect(app.Status.ID).To(Equal(applicationUUID))
	})

	It("clears the errors of an earlier sync when APIM answers none", func() {
		app := newApplication(applicationUUID, "", true)
		app.Status.Errors.Severe = []string{"severe one"}
		app.Status.Errors.Warning = []string{"warning one"}
		resp := apimapp.Response{State: sdk.ApplicationState{Hrid: "default-app", Id: new(applicationUUID)}}

		Expect(apimapp.UpdateStatus(context.Background(), app, resp)).To(Succeed())

		Expect(app.Status.Errors.Severe).To(BeEmpty())
		Expect(app.Status.Errors.Warning).To(BeEmpty())
	})
})

var _ = Describe("Application SetProcessingStatus", func() {
	It("writes Completed after a successful reconcile", func() {
		app := newApplication(applicationUUID, "", true)
		app.Status.ProcessingStatus = core.ProcessingStatusFailed

		apimapp.SetProcessingStatus(app, nil)

		Expect(app.Status.ProcessingStatus).To(Equal(core.ProcessingStatusCompleted))
	})

	It("writes Failed after a failed reconcile", func() {
		app := newApplication(applicationUUID, "", true)
		app.Status.ProcessingStatus = core.ProcessingStatusCompleted

		apimapp.SetProcessingStatus(app, errors.New("upsert refused"))

		Expect(app.Status.ProcessingStatus).To(Equal(core.ProcessingStatusFailed))
	})
})
