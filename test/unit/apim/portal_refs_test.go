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
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/portallisting"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	apimdocs "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/docs"
	apimportal "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portal"
	apimportallink "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portallink"
	apimportaltheme "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portaltheme"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
)

// indexRecorder records the indexers search.InitCache registers, so the fake client indexes as the
// manager's cache does.
type indexRecorder struct {
	cache.Cache
	indexers map[string]client.IndexerFunc
}

func (c *indexRecorder) IndexField(_ context.Context, _ client.Object, field string, f client.IndexerFunc) error {
	c.indexers[field] = f
	return nil
}

// portalCluster registers a fake client holding objects, the portal navigation kinds indexed as the
// manager's cache indexes them.
func portalCluster(objects ...client.Object) {
	GinkgoHelper()
	recorder := &indexRecorder{indexers: map[string]client.IndexerFunc{}}
	Expect(search.InitCache(context.Background(), recorder)).To(Succeed())
	scheme := runtime.NewScheme()
	Expect(clientgoscheme.AddToScheme(scheme)).To(Succeed())
	Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
	builder := fake.NewClientBuilder().WithScheme(scheme).WithObjects(objects...)
	for obj, field := range map[client.Object]search.IndexField{
		&v1alpha1.Portal{}:        search.PortalThemeField,
		&v1alpha1.PortalListing{}: search.PortalListingPortalField,
		&v1alpha1.PortalLink{}:    search.PortalLinkPortalField,
		&v1alpha1.Documentation{}: search.DocumentationPortalField,
	} {
		builder = builder.WithIndex(obj, field.String(), recorder.indexers[field.String()])
	}
	k8s.RegisterClient(builder.Build())
}

func contextRef() *refs.NamespacedName {
	return &refs.NamespacedName{Name: "dev-ctx", Namespace: "default"}
}

func refTheme(withContext bool) *v1alpha1.PortalTheme {
	thm := &v1alpha1.PortalTheme{ObjectMeta: metav1.ObjectMeta{Name: "dark", Namespace: "default"}}
	if withContext {
		thm.Spec.Context = contextRef()
	}
	return thm
}

func refPortal(withContext bool) *v1alpha1.Portal {
	prtl := &v1alpha1.Portal{ObjectMeta: metav1.ObjectMeta{Name: "main", Namespace: "default"}}
	if withContext {
		prtl.Spec.Context = contextRef()
	}
	return prtl
}

func refLink(portal, api *refs.NamespacedName) *v1alpha1.PortalLink {
	link := &v1alpha1.PortalLink{ObjectMeta: metav1.ObjectMeta{Name: "docs-link", Namespace: "default"}}
	link.Spec.Portal = portal
	link.Spec.API = api
	return link
}

func refListing() *v1alpha1.PortalListing {
	listing := &v1alpha1.PortalListing{ObjectMeta: metav1.ObjectMeta{Name: "catalog", Namespace: "default"}}
	listing.Spec.Portal = refs.NamespacedName{Name: "main"}
	listing.Spec.APIs = []portallisting.ApiEntry{{Ref: refs.NamespacedName{Name: "orders"}}}
	return listing
}

func refDoc(portal *refs.NamespacedName) *v1alpha1.Documentation {
	doc := &v1alpha1.Documentation{ObjectMeta: metav1.ObjectMeta{Name: "welcome", Namespace: "default"}}
	doc.Spec.Portal = portal
	return doc
}

func deleting[T client.Object](obj T) T {
	now := metav1.NewTime(time.Now())
	obj.SetDeletionTimestamp(&now)
	return obj
}

func mainPortalRef() *refs.NamespacedName {
	return &refs.NamespacedName{Name: "main"}
}

// The resolvers read the referenced objects through internal/k8s/dynamic, which talks to the ambient cluster:
// these specs cover what is decided before any such read. The NotFound wording is NotFoundRef's (error_test.go);
// resolution against a cluster is covered by test/integration.
var _ = Describe("Portal navigation reference resolution", func() {
	ctx := context.Background()

	Describe("PortalTheme", func() {
		It("refuses a theme without context with the reconcile's message", func() {
			portalCluster()
			Expect(apimportaltheme.ResolveRefs(ctx, refTheme(false), "default")).
				To(MatchError("theme [dark] has no management context"))
		})

		It("refuses to delete a theme a portal activates before looking at its context", func() {
			prtl := refPortal(true)
			prtl.Spec.Theme = &refs.NamespacedName{Name: "dark"}
			portalCluster(prtl)

			err := apimportaltheme.ResolveRefs(ctx, deleting(refTheme(true)), "default")
			Expect(err).To(MatchError(ContainSubstring("portals are using this theme")))
			Expect(apierrors.IsNotFound(gerrors.NewResolveRefError(err))).To(BeFalse())
		})

		It("deletes a theme without context without resolving anything", func() {
			portalCluster()
			Expect(apimportaltheme.ResolveRefs(ctx, deleting(refTheme(false)), "default")).To(Succeed())
		})
	})

	Describe("Portal", func() {
		It("refuses a portal without context with the reconcile's message", func() {
			portalCluster()
			Expect(apimportal.ResolveRefs(ctx, refPortal(false), "default")).
				To(MatchError("portal [main] has no management context"))
		})

		It("refuses a theme ref of another kind with the webhook's message", func() {
			portalCluster()
			prtl := refPortal(true)
			prtl.Spec.Theme = &refs.NamespacedName{Name: "dark", Kind: "Portal"}

			Expect(apimportal.ResolveRefs(ctx, prtl, "default")).
				To(MatchError("themeRef [dark] is of kind Portal, but only PortalTheme is supported"))
		})

		It("does not look at the theme while deleting", func() {
			portalCluster()
			prtl := deleting(refPortal(false))
			prtl.Spec.Theme = &refs.NamespacedName{Name: "dark", Kind: "Portal"}

			Expect(apimportal.ResolveRefs(ctx, prtl, "default")).To(Succeed())
		})

		DescribeTable("refuses to delete a portal still referenced before looking at its context",
			func(referrer client.Object) {
				portalCluster(referrer)
				err := apimportal.ResolveRefs(ctx, deleting(refPortal(true)), "default")
				Expect(err).To(MatchError(ContainSubstring("cannot be deleted")))
				Expect(apierrors.IsNotFound(gerrors.NewResolveRefError(err))).To(BeFalse())
			},
			Entry("by a listing", refListing()),
			Entry("by a link", refLink(mainPortalRef(), nil)),
			Entry("by a documentation page", refDoc(mainPortalRef())),
		)
	})

	Describe("PortalLink and Documentation", func() {
		DescribeTable("check the parent refs with the webhook's messages",
			func(portal, api *refs.NamespacedName, message string) {
				portalCluster()
				Expect(apimportallink.ResolveRefsForAdmission(ctx, refLink(portal, api), "default")).To(MatchError(message))
				Expect(apimportallink.ResolveRefs(ctx, refLink(portal, api), "default")).To(MatchError(message))
				doc := refDoc(portal)
				doc.Spec.API = api
				Expect(apimdocs.ResolveRefsForAdmission(ctx, doc, "default")).To(MatchError(message))
				Expect(apimdocs.ResolveRefs(ctx, doc, "default")).To(MatchError(message))
			},
			Entry("both", mainPortalRef(), &refs.NamespacedName{Name: "orders"},
				"exactly one of portalRef / apiRef must be set, but both were provided"),
			Entry("neither", nil, nil,
				"exactly one of portalRef / apiRef must be set, but neither was provided"),
			Entry("a v2 API", nil, &refs.NamespacedName{Name: "orders", Kind: "ApiDefinition"},
				"apiRef [orders] must be of kind ApiV4Definition (next-gen portal only supports those)"),
		)
	})
})
