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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

var _ = Describe("Authz PDP reconcile", func() {
	// The finalizer must be saved even when the apply fails: template.Compile has already recorded
	// the PDP on the Secrets and ConfigMaps it reads, and only a deletion reconcile releases them
	// (GKO-3387). Here the apply fails because the ManagementContext does not exist.
	It("saves the finalizer when the apply fails", func() {
		scheme := runtime.NewScheme()
		Expect(clientgoscheme.AddToScheme(scheme)).To(Succeed())
		Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
		cli := fake.NewClientBuilder().
			WithScheme(scheme).
			WithObjects(newAuthzPdp()).
			WithStatusSubresource(&v1alpha1.AuthzPdp{}).
			Build()
		k8s.RegisterClient(cli)

		r := &authzpdp.Reconciler{Scheme: scheme, Client: cli, Recorder: record.NewFakeRecorder(20)}
		key := types.NamespacedName{Namespace: "gravitee", Name: "orders-eu"}
		_, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: key})
		Expect(err).To(HaveOccurred())

		stored := &v1alpha1.AuthzPdp{}
		Expect(cli.Get(context.Background(), key, stored)).To(Succeed())
		Expect(stored.Finalizers).To(ContainElement(core.AuthzPdpFinalizer))
		Expect(stored.Status.RoutingScope).To(BeEmpty())
	})
})
