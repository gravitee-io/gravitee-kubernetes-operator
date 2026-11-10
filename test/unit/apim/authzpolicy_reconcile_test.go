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

package apim_test

import (
	"context"
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/authzpolicy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

var _ = Describe("Authz policy reconcile", func() {
	key := types.NamespacedName{Namespace: "gravitee", Name: "hotel-booking-read"}
	req := ctrl.Request{NamespacedName: key}
	secretKey := types.NamespacedName{Namespace: "gravitee", Name: "hotel-booking"}

	setup := func(objects ...client.Object) (*authzpolicy.Reconciler, client.Client) {
		scheme := runtime.NewScheme()
		Expect(clientgoscheme.AddToScheme(scheme)).To(Succeed())
		Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
		// The manager's cache sets the kind on what Get returns, and an Update clears it when it
		// decodes the API server's answer. The fake client does neither.
		cli := fake.NewClientBuilder().
			WithScheme(scheme).
			WithObjects(objects...).
			WithStatusSubresource(&v1alpha1.AuthzPolicy{}).
			WithInterceptorFuncs(interceptor.Funcs{
				Get: func(ctx context.Context, c client.WithWatch, key client.ObjectKey, obj client.Object,
					opts ...client.GetOption) error {
					if err := c.Get(ctx, key, obj, opts...); err != nil {
						return err
					}
					gvk, err := apiutil.GVKForObject(obj, scheme)
					obj.GetObjectKind().SetGroupVersionKind(gvk)
					return err
				},
				Update: func(ctx context.Context, c client.WithWatch, obj client.Object, opts ...client.UpdateOption) error {
					err := c.Update(ctx, obj, opts...)
					obj.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{})
					return err
				},
			}).
			Build()
		k8s.RegisterClient(cli)
		return &authzpolicy.Reconciler{Scheme: scheme, Client: cli, Recorder: record.NewFakeRecorder(20)}, cli
	}

	stored := func(cli client.Client) *v1alpha1.AuthzPolicy {
		p := &v1alpha1.AuthzPolicy{}
		Expect(cli.Get(context.Background(), key, p)).To(Succeed())
		return p
	}

	condition := func(cli client.Client, conditionType string) *metav1.Condition {
		c := k8s.GetCondition(stored(cli), conditionType)
		Expect(c).ToNot(BeNil())
		return c
	}

	twoPdps := func() *v1alpha1.AuthzPolicy {
		policy := newAuthzPolicy()
		policy.Spec.Targets = []authzpdp.Target{
			{PdpRef: &refs.NamespacedName{Name: "orders-eu"}},
			{Scope: "default"},
			{PdpRef: &refs.NamespacedName{Name: "orders-us"}},
		}
		return policy
	}

	unsynced := func(name string) *v1alpha1.AuthzPdp {
		pdp := newAuthzPdp()
		pdp.Name = name
		return pdp
	}

	// Resolving runs before templates compile: a waiting policy records no reference on the
	// Secrets it reads.
	It("waits on ResolvedRefs=False naming every pending AuthzPdp, recording no template reference", func() {
		policy := twoPdps()
		policy.Spec.Name = `[[ secret "hotel-booking/name" ]]`
		r, cli := setup(policy, unsynced("orders-us"), &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: secretKey.Name, Namespace: secretKey.Namespace},
			Data:       map[string][]byte{"name": []byte("Hotel booking — read access")},
		})

		_, err := r.Reconcile(context.Background(), req)
		Expect(err).To(HaveOccurred())

		p := stored(cli)
		Expect(p.Finalizers).To(ContainElement(core.AuthzPolicyFinalizer))
		resolved := condition(cli, k8s.ConditionResolvedRefs)
		Expect(resolved.Status).To(Equal(metav1.ConditionFalse))
		Expect(resolved.Message).To(ContainSubstring("authz pdp not found: [gravitee/orders-eu]; " +
			"authz pdp not synced: [gravitee/orders-us]"))
		Expect(p.Status.ID).To(BeEmpty())
		Expect(p.Status.PolicyStatus).To(BeEmpty())

		secret := &corev1.Secret{}
		Expect(cli.Get(context.Background(), secretKey, secret)).To(Succeed())
		Expect(secret.Annotations).ToNot(HaveKey("gravitee.io/authzpolicies"))
	})

	// Both contexts are immutable: only a change to the pdpRef fixes it, and that change reconciles
	// the policy again. Requeuing would retry and log the same error forever.
	It("reports an AuthzPdp of another context on Accepted without requeuing, even beside a pending one", func() {
		prod := newAuthzPdp()
		prod.Name = "orders-us"
		prod.Spec.Context = &refs.NamespacedName{Name: "prod-ctx"}
		prod.Status.RoutingScope = "orders@us"
		r, cli := setup(twoPdps(), prod)

		result, err := r.Reconcile(context.Background(), req)

		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(Equal(ctrl.Result{}))
		p := stored(cli)
		Expect(p.Finalizers).To(ContainElement(core.AuthzPolicyFinalizer))
		accepted := condition(cli, k8s.ConditionAccepted)
		Expect(accepted.Status).To(Equal(metav1.ConditionFalse))
		Expect(accepted.Reason).To(Equal("IllegalStateError"))
		Expect(accepted.Message).To(ContainSubstring(
			"authz pdp [gravitee/orders-us] belongs to management context [gravitee/prod-ctx], not [gravitee/dev-ctx]"))
		Expect(k8s.GetCondition(p, k8s.ConditionResolvedRefs).Status).ToNot(Equal(metav1.ConditionFalse))
		Expect(p.Status.ID).To(BeEmpty())
	})

	// Admission had no say, a pdpRef pending at apply time for instance: the reconcile's 400 is
	// where the ownership finding shows, and the policy is retried with backoff.
	It("keeps the platform's refusal in status.errors, naming the holder, and requeues", func() {
		policy := newAuthzPolicy()
		policy.Spec.Context = &refs.NamespacedName{Name: "refusing-ctx"}
		policy.Spec.Targets = []authzpdp.Target{{Scope: "default"}}
		r, cli := setup(policy)

		_, err := r.Reconcile(context.Background(), req)
		Expect(err).To(HaveOccurred())

		p := stored(cli)
		Expect(p.Status.Errors.Severe).To(ConsistOf(policyOwnershipFinding))
		Expect(condition(cli, k8s.ConditionAccepted).Status).To(Equal(metav1.ConditionFalse))
		Expect(p.Status.ID).To(BeEmpty())
		Expect(p.Status.PolicyStatus).To(BeEmpty())
	})

	It("records in status the deployment and the identifiers the platform answers", func() {
		policy := newAuthzPolicy()
		policy.Spec.Context = &refs.NamespacedName{Name: "accepting-ctx"}
		policy.Spec.Targets = []authzpdp.Target{{Scope: "default"}}
		r, cli := setup(policy)

		_, err := r.Reconcile(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())

		p := stored(cli)
		Expect(p.Status.ID).To(Equal("5f1c"))
		Expect(p.Status.OrgID).To(Equal("DEFAULT"))
		Expect(p.Status.EnvID).To(Equal("DEFAULT"))
		Expect(p.Status.HRID).To(Equal("gravitee-hotel-booking-read"))
		Expect(p.Status.PolicyStatus).To(Equal("DEPLOYED"))
		Expect(condition(cli, k8s.ConditionAccepted).Status).To(Equal(metav1.ConditionTrue))
		Expect(k8s.IsAutomationAPIManaged(p)).To(BeTrue())
	})

	// Only the reconcile's working copy is resolved: saving the policy must not replace its pdpRefs
	// with the scopes they resolved to.
	It("sends each pdpRef as its AuthzPdp's routing scope, the stored spec keeping the pdpRef", func() {
		policy := newAuthzPolicy()
		policy.Spec.Context = &refs.NamespacedName{Name: "accepting-ctx"}
		pdp := newAuthzPdp()
		pdp.Spec.Context = &refs.NamespacedName{Name: "accepting-ctx"}
		pdp.Status.RoutingScope = "orders@eu"
		r, cli := setup(policy, pdp)
		lastReceived.Delete("/authz/policies")

		_, err := r.Reconcile(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())

		received, ok := lastReceived.Load("/authz/policies")
		Expect(ok).To(BeTrue())
		body, isString := received.(string)
		Expect(isString).To(BeTrue())
		var sent struct {
			Targets []string `json:"targets"`
		}
		Expect(json.Unmarshal([]byte(body), &sent)).To(Succeed())
		Expect(sent.Targets).To(Equal([]string{"orders@eu", "default"}))
		Expect(stored(cli).Spec.Targets).To(Equal(newAuthzPolicy().Spec.Targets))
	})

	It("deletes a policy still waiting on its AuthzPdp", func() {
		r, cli := setup(twoPdps())

		_, err := r.Reconcile(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(cli.Delete(context.Background(), stored(cli))).To(Succeed())

		// The ManagementContext is missing, so there is no platform to delete from: the finalizer
		// goes, and the policy with it. The delete needs no target.
		_, err = r.Reconcile(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())
		Expect(apierrors.IsNotFound(cli.Get(context.Background(), key, &v1alpha1.AuthzPolicy{}))).To(BeTrue())
	})
})
