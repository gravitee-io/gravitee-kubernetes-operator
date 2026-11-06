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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

var _ = Describe("Authz PDP reconcile", func() {
	key := types.NamespacedName{Namespace: "gravitee", Name: "orders-eu"}
	req := ctrl.Request{NamespacedName: key}
	secretKey := types.NamespacedName{Namespace: "gravitee", Name: "orders-pdp"}

	setup := func(objects ...client.Object) (*authzpdp.Reconciler, client.Client) {
		scheme := runtime.NewScheme()
		Expect(clientgoscheme.AddToScheme(scheme)).To(Succeed())
		Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
		// The manager's cache sets the kind on what Get returns, and an Update clears it when it
		// decodes the API server's answer. The fake client does neither.
		cli := fake.NewClientBuilder().
			WithScheme(scheme).
			WithObjects(objects...).
			WithStatusSubresource(&v1alpha1.AuthzPdp{}).
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
		return &authzpdp.Reconciler{Scheme: scheme, Client: cli, Recorder: record.NewFakeRecorder(20)}, cli
	}

	// The finalizer must be saved even when the apply fails: template.Compile has already recorded
	// the PDP on the Secrets and ConfigMaps it reads, and only a deletion reconcile releases them
	// (GKO-3387). Here the apply fails because the ManagementContext does not exist: the suite's
	// dynamic client answers 404.
	It("saves the finalizer when the apply fails", func() {
		r, cli := setup(newAuthzPdp())

		_, err := r.Reconcile(context.Background(), req)
		Expect(err).To(HaveOccurred())

		stored := &v1alpha1.AuthzPdp{}
		Expect(cli.Get(context.Background(), key, stored)).To(Succeed())
		Expect(stored.Finalizers).To(ContainElement(core.AuthzPdpFinalizer))
		Expect(stored.Status.RoutingScope).To(BeEmpty())
	})

	// The reference recorded on the Secret is keyed by the PDP's kind, which the finalizer Update
	// clears: if it were not restored, the deletion reconcile would look for another key and leave
	// the Secret on its templating finalizer.
	It("releases on delete the template references a failed apply recorded", func() {
		pdp := newAuthzPdp()
		pdp.Spec.Name = `[[ secret "orders-pdp/name" ]]`
		r, cli := setup(pdp, &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: secretKey.Name, Namespace: secretKey.Namespace},
			Data:       map[string][]byte{"name": []byte("Orders PDP (EU)")},
		})

		_, err := r.Reconcile(context.Background(), req)
		Expect(err).To(HaveOccurred())

		secret := &corev1.Secret{}
		Expect(cli.Get(context.Background(), secretKey, secret)).To(Succeed())
		Expect(secret.Finalizers).To(ContainElement(core.TemplatingFinalizer))
		Expect(secret.Annotations).To(HaveKeyWithValue("gravitee.io/authzpdps", ContainSubstring(key.String())))

		stored := &v1alpha1.AuthzPdp{}
		Expect(cli.Get(context.Background(), key, stored)).To(Succeed())
		Expect(cli.Delete(context.Background(), stored)).To(Succeed())

		// The ManagementContext is missing, so there is no platform to delete from: the finalizer
		// goes, and the PDP with it.
		_, err = r.Reconcile(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())
		Expect(apierrors.IsNotFound(cli.Get(context.Background(), key, &v1alpha1.AuthzPdp{}))).To(BeTrue())

		Expect(cli.Get(context.Background(), secretKey, secret)).To(Succeed())
		Expect(secret.Finalizers).ToNot(ContainElement(core.TemplatingFinalizer))
	})
})
