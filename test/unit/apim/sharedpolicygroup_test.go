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
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/sharedpolicygroups"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	controller "github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/sharedpolicygroups"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	spg "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/sharedpolicygroup"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

const (
	spgPath    = "/automation/organizations/org1/environments/env1/shared-policy-groups"
	spgID      = "5e2a7c1d-0000-0000-0000-000000000001"
	spgCrossID = "5e2a7c1d-0000-0000-0000-000000000002"
	spgState   = `{
		"id": "` + spgID + `",
		"crossId": "` + spgCrossID + `",
		"hrid": "default-limits",
		"name": "Limits",
		"apiType": "PROXY",
		"phase": "REQUEST"
	}`
)

// limits is a shared policy group created through the Automation API: no status yet.
func limits() *v1alpha1.SharedPolicyGroup {
	phase := sharedpolicygroups.FlowPhase("REQUEST")
	return &v1alpha1.SharedPolicyGroup{
		ObjectMeta: metav1.ObjectMeta{Name: "limits", Namespace: "default"},
		Spec: v1alpha1.SharedPolicyGroupSpec{
			Context: &refs.NamespacedName{Name: "dev", Namespace: "default"},
			SharedPolicyGroup: &sharedpolicygroups.SharedPolicyGroup{
				Name:    "Limits",
				ApiType: "PROXY",
				Phase:   &phase,
				Steps:   []*sharedpolicygroups.Step{{Enabled: true, Policy: new("rate-limit")}},
			},
		},
	}
}

// legacyLimits is a shared policy group created before the Automation API: APIM knows it by its cross id.
func legacyLimits() *v1alpha1.SharedPolicyGroup {
	obj := limits()
	obj.Status.CrossID = spgCrossID
	return obj
}

func toSPGDTO(obj *v1alpha1.SharedPolicyGroup) spg.DTO {
	GinkgoHelper()
	dto, err := spg.ToDTO(obj)
	Expect(err).ToNot(HaveOccurred())
	return dto
}

func firstStep(body map[string]any) map[string]any {
	GinkgoHelper()
	steps, ok := body["steps"].([]any)
	Expect(ok).To(BeTrue())
	return steps[0].(map[string]any)
}

// registerSPGK8s registers, and returns, a fake Kubernetes client holding objects.
func registerSPGK8s(objects ...client.Object) client.Client {
	GinkgoHelper()
	scheme := runtime.NewScheme()
	Expect(clientgoscheme.AddToScheme(scheme)).To(Succeed())
	Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
	cli := fake.NewClientBuilder().WithScheme(scheme).WithObjects(objects...).
		WithStatusSubresource(&v1alpha1.SharedPolicyGroup{}).Build()
	k8s.RegisterClient(cli)
	return cli
}

var _ = Describe("SharedPolicyGroup through the APIM SDK", func() {
	var (
		ctx    context.Context
		fake   *fakeAutomationAPI
		client *apim.SDKClient
	)

	BeforeEach(func() {
		ctx = context.Background()
		fake = newFakeAutomationAPI()
		fake.response = spgState
		DeferCleanup(fake.Close)

		var err error
		client, err = apim.NewSDKClient(ctx, managementContext(fake.URL, bearer("token")))
		Expect(err).ToNot(HaveOccurred())
	})

	It("reconciles an Automation API managed group by its hrid, reporting Accepted only", func() {
		obj := limits()
		cli := registerSPGK8s(obj)

		// The real lifecycle, with only the management context, which this suite has no cluster for, stubbed.
		l := controller.NewLifecycle()
		l.ClientFactory = func(context.Context, *v1alpha1.SharedPolicyGroup) (*apim.SDKClient, error) { return client, nil }

		req := ctrl.Request{NamespacedName: types.NamespacedName{Namespace: obj.Namespace, Name: obj.Name}}
		_, err := l.Reconcile(ctx, cli, record.NewFakeRecorder(20), req, &v1alpha1.SharedPolicyGroup{})
		Expect(err).ToNot(HaveOccurred())

		sent := fake.lastRequest()
		Expect(sent.Method).To(Equal(http.MethodPut))
		Expect(sent.Path).To(Equal(spgPath))
		Expect(sent.Query.Get("dryRun")).To(Equal("false"))
		Expect(sent.Query.Has("hridContainsUUID")).To(BeFalse())
		Expect(sent.Body).To(HaveKeyWithValue("hrid", refs.NewNamespacedNameFromObject(obj).HRID()))

		got := &v1alpha1.SharedPolicyGroup{}
		Expect(cli.Get(ctx, req.NamespacedName, got)).To(Succeed())
		Expect(got.Status.ID).To(Equal(spgID))
		Expect(got.Status.CrossID).To(Equal(spgCrossID))
		Expect(got.Status.OrgID).To(Equal("org1"))
		Expect(got.Status.EnvID).To(Equal("env1"))
		Expect(got.Status.ProcessingStatus).To(Equal(core.ProcessingStatusCompleted))
		Expect(k8s.IsAccepted(got)).To(BeTrue())
		Expect(k8s.GetCondition(got, k8s.ConditionResolvedRefs)).To(BeNil())
		Expect(k8s.IsAutomationAPIManaged(got)).To(BeTrue())
	})

	It("addresses a legacy group by its cross id on PUT, GET and DELETE", func() {
		dto := toSPGDTO(legacyLimits())

		resp, err := spg.Upsert(ctx, client, dto)
		Expect(err).ToNot(HaveOccurred())
		sent := fake.lastRequest()
		Expect(sent.Path).To(Equal(spgPath))
		Expect(sent.Query.Get("hridContainsUUID")).To(Equal("true"))
		Expect(sent.Body).To(HaveKeyWithValue("hrid", spgCrossID))
		Expect(sent.Body).To(HaveKeyWithValue("crossId", spgCrossID))

		_, err = spg.GetRemote(ctx, client, dto)
		Expect(err).ToNot(HaveOccurred())
		sent = fake.lastRequest()
		Expect(sent.Method).To(Equal(http.MethodGet))
		Expect(sent.Path).To(Equal(spgPath + "/" + spgCrossID))
		Expect(sent.Query.Get("hridContainsUUID")).To(Equal("true"))

		Expect(spg.Delete(ctx, client, dto)).To(Succeed())
		sent = fake.lastRequest()
		Expect(sent.Method).To(Equal(http.MethodDelete))
		Expect(sent.Path).To(Equal(spgPath + "/" + spgCrossID))
		Expect(sent.Query.Get("hridContainsUUID")).To(Equal("true"))

		obj := legacyLimits()
		Expect(spg.UpdateStatus(ctx, obj, resp)).To(Succeed())
		Expect(k8s.IsAutomationAPIManaged(obj)).To(BeFalse())
	})

	It("addresses a group by its hrid once it is Automation API managed, even with a status cross id", func() {
		obj := legacyLimits()
		k8s.AddAutomationAPIManagedCondition(obj)
		dto := toSPGDTO(obj)

		Expect(dto.HridContainsUUID).To(BeFalse())
		Expect(dto.Identity()).To(Equal(refs.NewNamespacedNameFromObject(obj).HRID()))
	})

	It("sends a step message condition", func() {
		obj := limits()
		obj.Spec.Steps[0].MessageCondition = new("{#message.headers['foo'] != null}")

		_, err := spg.Upsert(ctx, client, toSPGDTO(obj))
		Expect(err).ToNot(HaveOccurred())

		Expect(firstStep(fake.lastRequest().Body)).
			To(HaveKeyWithValue("messageCondition", "{#message.headers['foo'] != null}"))
	})

	It("sends a step without a policy with an empty policy, for APIM to refuse", func() {
		obj := limits()
		obj.Spec.Steps[0].Policy = nil

		_, err := spg.Upsert(ctx, client, toSPGDTO(obj))
		Expect(err).ToNot(HaveOccurred())

		Expect(firstStep(fake.lastRequest().Body)).To(HaveKeyWithValue("policy", ""))
	})

	It("counts a group APIM no longer has as deleted", func() {
		fake.status = http.StatusNotFound

		Expect(spg.Delete(ctx, client, toSPGDTO(limits()))).To(Succeed())
	})

	It("reports the severe findings of a dry run, without its warnings", func() {
		fake.response = `{"hrid": "default-limits", "name": "Limits", "apiType": "PROXY", "phase": "REQUEST",
			"errors": {"severe": ["plugin not found"], "warning": ["warning one"]}}`

		errs, err := spg.DryRun(ctx, client, toSPGDTO(limits()))
		Expect(err).ToNot(HaveOccurred())

		Expect(fake.lastRequest().Query.Get("dryRun")).To(Equal("true"))
		Expect(errs.IsSevere()).To(BeTrue())
		Expect(errs.Severe).To(HaveLen(1))
		Expect(errs.Severe[0].Message).To(Equal("plugin not found"))
		Expect(errs.Warning).To(BeEmpty())
	})

	It("reports the warnings of a dry run APIM finds nothing severe in", func() {
		fake.response = `{"hrid": "default-limits", "name": "Limits", "apiType": "PROXY", "phase": "REQUEST",
			"errors": {"severe": [], "warning": ["warning one"]}}`

		errs, err := spg.DryRun(ctx, client, toSPGDTO(limits()))
		Expect(err).ToNot(HaveOccurred())

		Expect(errs.IsSevere()).To(BeFalse())
		Expect(errs.Warning).To(HaveLen(1))
		Expect(errs.Warning[0].Message).To(Equal("warning one"))
	})
})
