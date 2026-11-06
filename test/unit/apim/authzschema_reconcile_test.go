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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/authzschema"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

var _ = Describe("Authz schema reconcile", func() {
	key := types.NamespacedName{Namespace: "gravitee", Name: "travel-schema"}
	req := ctrl.Request{NamespacedName: key}
	secretKey := types.NamespacedName{Namespace: "gravitee", Name: "travel-schema"}

	setup := func(objects ...client.Object) (*authzschema.Reconciler, client.Client) {
		scheme := runtime.NewScheme()
		Expect(clientgoscheme.AddToScheme(scheme)).To(Succeed())
		Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
		// The manager's cache sets the kind on what Get returns, and an Update clears it when it
		// decodes the API server's answer. The fake client does neither.
		cli := fake.NewClientBuilder().
			WithScheme(scheme).
			WithObjects(objects...).
			WithStatusSubresource(&v1alpha1.AuthzSchema{}).
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
		return &authzschema.Reconciler{Scheme: scheme, Client: cli, Recorder: record.NewFakeRecorder(20)}, cli
	}

	stored := func(cli client.Client) *v1alpha1.AuthzSchema {
		s := &v1alpha1.AuthzSchema{}
		Expect(cli.Get(context.Background(), key, s)).To(Succeed())
		return s
	}

	expectWaitingOnPdp := func(cli client.Client, message string) {
		s := stored(cli)
		Expect(s.Finalizers).To(ContainElement(core.AuthzSchemaFinalizer))
		resolved := k8s.GetCondition(s, k8s.ConditionResolvedRefs)
		Expect(resolved).ToNot(BeNil())
		Expect(resolved.Status).To(Equal(metav1.ConditionFalse))
		Expect(resolved.Message).To(ContainSubstring(message))
		Expect(s.Status.ID).To(BeEmpty())
	}

	// The ResolvedRefs error is returned, so the controller requeues it with backoff on top of the
	// AuthzPdp watch. Resolving runs before templates compile: a waiting schema records no reference
	// on the Secrets it reads.
	It("waits on ResolvedRefs=False while the AuthzPdp is missing, recording no template reference", func() {
		s := newAuthzSchema()
		s.Spec.Name = `[[ secret "travel-schema/name" ]]`
		r, cli := setup(s, &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: secretKey.Name, Namespace: secretKey.Namespace},
			Data:       map[string][]byte{"name": []byte("Travel schema")},
		})

		_, err := r.Reconcile(context.Background(), req)
		Expect(err).To(HaveOccurred())

		expectWaitingOnPdp(cli, "authz pdp not found: [gravitee/orders-eu]")
		secret := &corev1.Secret{}
		Expect(cli.Get(context.Background(), secretKey, secret)).To(Succeed())
		Expect(secret.Annotations).ToNot(HaveKey("gravitee.io/authzschemas"))
	})

	It("waits on ResolvedRefs=False while the AuthzPdp has no routing scope", func() {
		pdp := newAuthzPdp()
		r, cli := setup(newAuthzSchema(), pdp)

		_, err := r.Reconcile(context.Background(), req)
		Expect(err).To(HaveOccurred())

		expectWaitingOnPdp(cli, "authz pdp not synced: [gravitee/orders-eu] has no routing scope yet")
	})

	// A scope target needs no AuthzPdp, so the reconcile goes on to the apply, which fails here
	// because the ManagementContext does not exist: the suite's dynamic client answers 404.
	It("applies a scope target without an AuthzPdp, saving the finalizer when the apply fails", func() {
		s := newAuthzSchema()
		s.Spec.Target = &authzpdp.Target{Scope: "default"}
		r, cli := setup(s)

		_, err := r.Reconcile(context.Background(), req)
		Expect(err).To(HaveOccurred())

		// The missing ManagementContext is what fails, reported on ResolvedRefs as for every kind.
		Expect(err.Error()).ToNot(ContainSubstring("authz pdp"))
		st := stored(cli)
		Expect(st.Finalizers).To(ContainElement(core.AuthzSchemaFinalizer))
		resolved := k8s.GetCondition(st, k8s.ConditionResolvedRefs)
		Expect(resolved.Status).To(Equal(metav1.ConditionFalse))
		Expect(resolved.Message).ToNot(ContainSubstring("authz pdp"))
	})

	It("deletes a schema still waiting on its AuthzPdp", func() {
		r, cli := setup(newAuthzSchema())

		_, err := r.Reconcile(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(cli.Delete(context.Background(), stored(cli))).To(Succeed())

		// The ManagementContext is missing, so there is no platform to delete from: the finalizer
		// goes, and the schema with it. The delete needs no target.
		_, err = r.Reconcile(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())
		Expect(apierrors.IsNotFound(cli.Get(context.Background(), key, &v1alpha1.AuthzSchema{}))).To(BeTrue())
	})
})
