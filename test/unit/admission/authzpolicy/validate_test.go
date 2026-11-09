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
	"encoding/json"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpolicy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	admission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/authzpolicy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
)

const ownershipFinding = `entityId [policy.hotel-booking-read] is held by policy "Hotel booking" (id: 5f1c), ` +
	`which this manifest did not create. Declare another entityId, or delete that policy first`

const mcpTargetsWarning = `the policy names MCPTool, but its targets leave out default and *: an MCP proxy ` +
	`asks the default engine only, so it never evaluates this policy. Add default to targets`

const policyText = `permit(principal in Group::"agents", action == Action::"read", resource == Hotel::"grand-paris");`

const ns = "travel"

var devCtx = refs.NamespacedName{Name: "dev-ctx"}

func newPolicy(name, entityID string, targets ...authzpdp.Target) *v1alpha1.AuthzPolicy {
	return &v1alpha1.AuthzPolicy{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns},
		Spec: v1alpha1.AuthzPolicySpec{
			Context: &devCtx,
			Type: authzpolicy.Type{
				EntityID: entityID, Name: "Hotel booking — read access", PolicyType: "API",
				PolicyText: policyText, Targets: targets,
			},
		},
	}
}

func byRef(name string) authzpdp.Target {
	return authzpdp.Target{PdpRef: &refs.NamespacedName{Name: name}}
}

func newPdp(name string, contextRef refs.NamespacedName, scope string) *v1alpha1.AuthzPdp {
	return &v1alpha1.AuthzPdp{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns},
		Spec: v1alpha1.AuthzPdpSpec{
			Context: &contextRef,
			Type:    authzpdp.Type{Name: name, TargetPdpID: name},
		},
		Status: v1alpha1.AuthzPdpStatus{Status: authzpdp.Status{RoutingScope: scope}},
	}
}

// cluster registers a fake client holding objects, AuthzPolicy indexed as the manager's cache
// indexes it.
func cluster(objects ...client.Object) {
	c := &recordingCache{indexers: map[string]client.IndexerFunc{}}
	Expect(search.InitCache(context.Background(), c)).To(Succeed())
	scheme := runtime.NewScheme()
	Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
	builder := fake.NewClientBuilder().WithScheme(scheme).WithObjects(objects...)
	for _, field := range []search.IndexField{search.AuthzPolicyContextField, search.AuthzPolicyAuthzPdpField} {
		builder = builder.WithIndex(&v1alpha1.AuthzPolicy{}, field.String(), c.indexers[field.String()])
	}
	k8s.RegisterClient(builder.Build())
}

func mustJSON(v any) string {
	raw, err := json.Marshal(v)
	Expect(err).ToNot(HaveOccurred())
	return string(raw)
}

var _ = Describe("Validating an AuthzPolicy", func() {
	ctx := context.Background()
	ctrl := admission.AdmissionCtrl{}

	BeforeEach(func() {
		platform.reset()
	})

	It("warns once, naming every pending AuthzPdp, and skips the dry run", func() {
		cluster(newPdp("orders-us", devCtx, ""))

		warnings, err := ctrl.ValidateCreate(ctx, newPolicy("hotel-booking-read", "policy.hotel-booking-read",
			byRef("orders-eu"), authzpdp.Target{Scope: "default"}, byRef("orders-us")))

		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(ConsistOf("authz pdp not found: [travel/orders-eu]; " +
			"authz pdp not synced: [travel/orders-us] has no routing scope yet: " +
			"the policy is checked against the platform once every AuthzPdp it targets has synced"))
		Expect(platform.requests()).To(BeEmpty())
	})

	// The controller's own finalizer and annotation writes go through the webhook: refusing them
	// while an AuthzPdp is pending would stop the controller from saving its finalizer.
	It("admits a metadata-only update while an AuthzPdp is pending", func() {
		cluster()
		oldObj := newPolicy("hotel-booking-read", "policy.hotel-booking-read", byRef("orders-eu"))
		oldObj.Status.ID = "5f1c"
		newObj := oldObj.DeepCopy()
		newObj.Finalizers = []string{"finalizers.gravitee.io/authzpolicies"}

		_, err := ctrl.ValidateUpdate(ctx, oldObj, newObj)

		Expect(err).ToNot(HaveOccurred())
		Expect(platform.requests()).To(BeEmpty())
	})

	// Only a spec change is dry-run: the controller writes its finalizer and annotations on every
	// reconcile.
	It("makes no platform call on a metadata-only update whose targets resolve", func() {
		cluster(newPdp("orders-eu", devCtx, "orders@eu"))
		oldObj := newPolicy("hotel-booking-read", "policy.hotel-booking-read", byRef("orders-eu"))
		oldObj.Status.ID = "5f1c"
		newObj := oldObj.DeepCopy()
		newObj.Annotations = map[string]string{"gravitee.io/last-spec-hash": "5f1c"}

		_, err := ctrl.ValidateUpdate(ctx, oldObj, newObj)

		Expect(err).ToNot(HaveOccurred())
		Expect(platform.requests()).To(BeEmpty())
	})

	It("refuses a pdpRef naming another kind, naming its index", func() {
		cluster()
		mcp := authzpdp.Target{PdpRef: &refs.NamespacedName{Name: "orders-eu", Kind: "McpProxy"}}

		_, err := ctrl.ValidateCreate(ctx,
			newPolicy("hotel-booking-read", "policy.hotel-booking-read", authzpdp.Target{Scope: "default"}, mcp))

		Expect(err).To(MatchError(ContainSubstring(
			"targets[1].pdpRef [travel/orders-eu]: kind must be AuthzPdp, got [McpProxy]")))
		Expect(platform.requests()).To(BeEmpty())
	})

	It("refuses an AuthzPdp of another management context", func() {
		cluster(newPdp("orders-eu", refs.NamespacedName{Name: "prod-ctx"}, "orders@eu"))

		_, err := ctrl.ValidateCreate(ctx, newPolicy("hotel-booking-read", "policy.hotel-booking-read", byRef("orders-eu")))

		Expect(err).To(MatchError(ContainSubstring(
			"authz pdp [travel/orders-eu] belongs to management context [travel/prod-ctx], not [travel/dev-ctx]")))
		Expect(platform.requests()).To(BeEmpty())
	})

	It("dry-runs the resolved targets in declaration order, and the text verbatim", func() {
		cluster(newPdp("orders-eu", devCtx, "orders@eu"), newPdp("orders-us", devCtx, "orders@us"))
		platform.answer(http.StatusOK, `{"hrid":"travel-hotel-booking-read"}`)

		_, err := ctrl.ValidateCreate(ctx, newPolicy("hotel-booking-read", "policy.hotel-booking-read",
			byRef("orders-us"), authzpdp.Target{Scope: "default"}, byRef("orders-eu")))

		Expect(err).ToNot(HaveOccurred())
		Expect(platform.requests()).To(HaveLen(1))
		sent := struct {
			Targets    []string `json:"targets"`
			PolicyText string   `json:"policyText"`
		}{}
		Expect(json.Unmarshal([]byte(platform.requests()[0]), &sent)).To(Succeed())
		Expect(sent.Targets).To(Equal([]string{"orders@us", "default", "orders@eu"}))
		Expect(sent.PolicyText).To(Equal(policyText))
	})

	It("refuses an entityId another policy holds, naming that policy", func() {
		cluster(newPdp("orders-eu", devCtx, "orders@eu"))
		platform.answer(http.StatusOK,
			`{"hrid":"travel-hotel-booking-read","errors":{"severe":[`+mustJSON(ownershipFinding)+`]}}`)

		_, err := ctrl.ValidateCreate(ctx, newPolicy("hotel-booking-read", "policy.hotel-booking-read", byRef("orders-eu")))

		Expect(err).To(MatchError(ContainSubstring(ownershipFinding)))
	})

	It("admits an MCP policy the default engine never sees, with the platform's warning", func() {
		cluster(newPdp("orders-eu", devCtx, "orders@eu"))
		platform.answer(http.StatusOK,
			`{"hrid":"travel-github-create-issue","errors":{"warning":[`+mustJSON(mcpTargetsWarning)+`]}}`)
		policy := newPolicy("github-create-issue", "policy.github.create-issue", byRef("orders-eu"))
		policy.Spec.PolicyType = "MCP"

		warnings, err := ctrl.ValidateCreate(ctx, policy)

		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(ConsistOf(mcpTargetsWarning))
	})

	// Policies have no children: the platform withdraws the policy from its gateways and nothing
	// else goes with it.
	It("admits a delete with no warning and no platform call", func() {
		cluster()

		warnings, err := ctrl.ValidateDelete(ctx, newPolicy("hotel-booking-read", "policy.hotel-booking-read"))

		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(BeEmpty())
		Expect(platform.requests()).To(BeEmpty())
	})
})

var _ = Describe("Detecting drift on an AuthzPolicy", func() {
	ctx := context.Background()
	ctrl := admission.AdmissionCtrl{}

	BeforeEach(func() {
		platform.reset()
		platform.answer(http.StatusOK, `{"hrid":"travel-hotel-booking-read"}`)
		saved := env.Config.DriftDetection
		env.Config.DriftDetection.Enabled = true
		env.Config.DriftDetection.Policy = env.DriftPolicyDeny
		env.Config.DriftDetection.OnRemoteMissing = env.DriftPolicyDeny
		DeferCleanup(func() { env.Config.DriftDetection = saved })
	})

	applied := func() *v1alpha1.AuthzPolicy {
		policy := newPolicy("hotel-booking-read", "policy.hotel-booking-read", authzpdp.Target{Scope: "default"})
		policy.Status.ID = "5f1c"
		return policy
	}
	renamed := func(policy *v1alpha1.AuthzPolicy) *v1alpha1.AuthzPolicy {
		edited := policy.DeepCopy()
		edited.Spec.Name = "Hotel booking — read only"
		return edited
	}
	remote := func(text, policyStatus string) string {
		return `{"hrid":"travel-hotel-booking-read","entityId":"policy.hotel-booking-read",` +
			`"name":"Hotel booking — read access","policyType":"API","policyText":` + mustJSON(text) + `,` +
			`"targets":["default"],"id":"5f1c","policyStatus":"` + policyStatus + `"}`
	}

	// A disable in the Console changes policyStatus only: the next apply deploys the policy again,
	// with the platform's warning.
	It("does not take a policy disabled in the Console as drift", func() {
		cluster()
		platform.answerGet(http.StatusOK, remote(policyText, "DISABLED"))

		_, err := ctrl.ValidateUpdate(ctx, applied(), renamed(applied()))

		Expect(err).ToNot(HaveOccurred())
	})

	It("refuses an edit while the platform holds another text", func() {
		cluster()
		platform.answerGet(http.StatusOK, remote(`forbid(principal, action, resource);`, "DEPLOYED"))

		_, err := ctrl.ValidateUpdate(ctx, applied(), renamed(applied()))

		Expect(err).To(MatchError(ContainSubstring("drift detected")))
		Expect(err).To(MatchError(ContainSubstring("forbid(principal, action, resource);")))
	})

	// The platform holds the routing scope the old pdpRef resolves to, so moving the policy to a
	// scope is an edit. Admission resolves the new spec only: drift resolves the old one, without
	// which its pdpRef would compare as an empty scope and both sides would differ from the platform.
	It("admits retargeting a pdpRef while the platform holds its routing scope", func() {
		cluster(newPdp("orders-eu", devCtx, "orders@eu"))
		platform.answerGet(http.StatusOK, strings.Replace(remote(policyText, "DEPLOYED"),
			`"targets":["default"]`, `"targets":["orders@eu"]`, 1))
		retargeted := applied()
		routed := retargeted.DeepCopy()
		routed.Spec.Targets = []authzpdp.Target{byRef("orders-eu")}

		_, err := ctrl.ValidateUpdate(ctx, routed, retargeted)

		Expect(err).ToNot(HaveOccurred())
	})

	// Comparing would apply the remote-missing policy, deny here, and refuse the very edit that
	// unblocks a policy waiting on its AuthzPdp.
	It("compares nothing for a policy the platform never received", func() {
		cluster()
		platform.answerGet(http.StatusNotFound, `{"message":"not found"}`)
		neverApplied := applied()
		neverApplied.Status.ID = ""

		_, err := ctrl.ValidateUpdate(ctx, neverApplied, renamed(neverApplied))

		Expect(err).ToNot(HaveOccurred())
		Expect(platform.requests()).To(HaveLen(1))
	})

	// The usual cause is a deleted AuthzPdp, whose scope the platform has already removed from the
	// policy: comparing the old targets would refuse the edit that retargets it.
	It("compares nothing while an old pdpRef no longer resolves", func() {
		cluster()
		platform.answerGet(http.StatusOK, remote(`forbid(principal, action, resource);`, "DEPLOYED"))
		retargeted := applied()
		orphaned := retargeted.DeepCopy()
		orphaned.Spec.Targets = []authzpdp.Target{byRef("orders-eu")}

		_, err := ctrl.ValidateUpdate(ctx, orphaned, retargeted)

		Expect(err).ToNot(HaveOccurred())
	})
})
