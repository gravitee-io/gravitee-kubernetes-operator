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

package search_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/mcpproxy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
)

// No authz kind targets an AuthzPdp yet. McpProxy and its catalog server index stand in for a
// dependent kind and its pdpRef index: both key a resource by the namespaced name it references.
var _ = Describe("FindAuthzPdpDependents", func() {
	const ns = "orders"

	pdpRef := refs.NamespacedName{Name: "orders-eu", Namespace: ns}
	standIn := search.AuthzPdpDependent{
		Kind:    "McpProxy",
		Field:   search.McpProxyCatalogMcpServerField,
		NewList: func() client.ObjectList { return &v1alpha1.McpProxyList{} },
	}
	targeting := func(name string, target refs.NamespacedName) *v1alpha1.McpProxy {
		return &v1alpha1.McpProxy{
			ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns},
			Spec: v1alpha1.McpProxySpec{Type: mcpproxy.Type{
				Mode: mcpproxy.ModeStudio,
				Studio: &mcpproxy.Studio{Tools: []mcpproxy.StudioTool{
					{ServerRef: target, Tool: "check"},
				}},
			}},
		}
	}
	cluster := func(objects ...client.Object) {
		c := &recordingCache{indexers: map[string]client.IndexerFunc{}}
		Expect(search.InitCache(context.Background(), c)).To(Succeed())
		scheme := runtime.NewScheme()
		Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
		field := search.McpProxyCatalogMcpServerField.String()
		k8s.RegisterClient(fake.NewClientBuilder().
			WithScheme(scheme).
			WithObjects(objects...).
			WithIndex(&v1alpha1.McpProxy{}, field, c.indexers[field]).
			Build())
	}

	It("names the resources that target the PDP, sorted", func() {
		cluster(
			targeting("orders-schema", refs.NamespacedName{Name: "orders-eu"}),
			targeting("billing-schema", refs.NamespacedName{Name: "billing"}),
			targeting("audit-schema", pdpRef),
		)

		found, err := search.FindAuthzPdpDependents(context.Background(), pdpRef,
			[]search.AuthzPdpDependent{standIn})

		Expect(err).ToNot(HaveOccurred())
		Expect(found).To(Equal([]search.Dependents{{
			Kind:  "McpProxy",
			Names: []string{"orders/audit-schema", "orders/orders-schema"},
		}}))
	})

	It("leaves out a kind with no resource targeting the PDP", func() {
		cluster(targeting("billing-schema", refs.NamespacedName{Name: "billing"}))

		found, err := search.FindAuthzPdpDependents(context.Background(), pdpRef,
			[]search.AuthzPdpDependent{standIn})

		Expect(err).ToNot(HaveOccurred())
		Expect(found).To(BeEmpty())
	})

	It("finds nothing while no kind is registered", func() {
		cluster(targeting("orders-schema", pdpRef))

		found, err := search.FindAuthzPdpDependents(context.Background(), pdpRef, search.AuthzPdpDependents)

		Expect(err).ToNot(HaveOccurred())
		Expect(found).To(BeEmpty())
	})
})
