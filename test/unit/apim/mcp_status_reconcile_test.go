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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/mcpproxy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	catalogmcpserverctrl "github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/catalogmcpserver"
	mcpproxyctrl "github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/mcpproxy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

// The reference a policy names a server or a tool by is copied from the state the platform answers
// into status, on every successful apply.
var _ = Describe("MCP status reconcile", func() {
	accepting := &refs.NamespacedName{Name: "accepting-ctx"}

	// setup registers a fake client holding objects. The manager's cache sets the kind on what Get
	// returns, and an Update clears it when it decodes the API server's answer: the fake client does
	// neither.
	setup := func(objects ...client.Object) (*runtime.Scheme, client.Client) {
		scheme := runtime.NewScheme()
		Expect(clientgoscheme.AddToScheme(scheme)).To(Succeed())
		Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
		cli := fake.NewClientBuilder().
			WithScheme(scheme).
			WithObjects(objects...).
			WithStatusSubresource(&v1alpha1.CatalogMcpServer{}, &v1alpha1.McpProxy{}).
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
		return scheme, cli
	}

	It("copies the server's and each tool's authPolicyRef into CatalogMcpServer status", func() {
		srv := newCatalogMcpServer(nil)
		srv.Spec.Context = accepting
		key := types.NamespacedName{Namespace: srv.Namespace, Name: srv.Name}
		scheme, cli := setup(srv)
		r := &catalogmcpserverctrl.Reconciler{Scheme: scheme, Client: cli, Recorder: record.NewFakeRecorder(20)}

		_, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: key})
		Expect(err).ToNot(HaveOccurred())

		stored := &v1alpha1.CatalogMcpServer{}
		Expect(cli.Get(context.Background(), key, stored)).To(Succeed())
		Expect(stored.Status.ID).To(Equal("7a2e"))
		Expect(stored.Status.AuthPolicyRef).To(Equal(`MCPServer::"github"`))
		Expect(stored.Status.Tools).To(HaveLen(1))
		Expect(stored.Status.Tools[0].EntityID).To(Equal("mcp-tool.github.create_issue"))
		Expect(stored.Status.Tools[0].AuthPolicyRef).To(Equal(`MCPTool::"github.create_issue"`))
	})

	It("copies the state's tools list into McpProxy status, an aliased tool keeping its reference", func() {
		proxy := newMcpStudio()
		proxy.Spec.Context = accepting
		synced := func(name, namespace string) *v1alpha1.CatalogMcpServer {
			srv := newCatalogMcpServer(nil)
			srv.ObjectMeta = metav1.ObjectMeta{Name: name, Namespace: namespace}
			srv.Status.ID = name
			return srv
		}
		key := types.NamespacedName{Namespace: proxy.Namespace, Name: proxy.Name}
		scheme, cli := setup(proxy, synced("github-mcp", "gravitee"), synced("jira-mcp", "tools"))
		r := &mcpproxyctrl.Reconciler{Scheme: scheme, Client: cli, Recorder: record.NewFakeRecorder(20)}

		_, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: key})
		Expect(err).ToNot(HaveOccurred())

		stored := &v1alpha1.McpProxy{}
		Expect(cli.Get(context.Background(), key, stored)).To(Succeed())
		Expect(stored.Status.ID).To(Equal("9b3d"))
		Expect(stored.Status.Tools).To(Equal([]mcpproxy.StudioToolStatus{{
			Server:        "gravitee-github-mcp",
			Tool:          "search_issues",
			Alias:         new("gh_search"),
			EntityID:      "mcp-tool.github.search_issues",
			AuthPolicyRef: `MCPTool::"github.search_issues"`,
		}}))
	})
})
