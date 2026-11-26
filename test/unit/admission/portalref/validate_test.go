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

package portalref_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	documentation "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/docs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/portallink"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/portallisting"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	docsadm "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/docs"
	linkadm "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/portallink"
	listingadm "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/portallisting"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

// A Portal applied in the same file as its listing, page or link is written to APIM by its
// reconcile, after their admission. Until then APIM refuses them while the environment holds
// another portal, its default portal included (GKO-3480).
var _ = Describe("Validating a resource of a portal GKO has not synced yet", func() {
	const ns = "travel"
	ctx := context.Background()
	meta := func(name string) metav1.ObjectMeta { return metav1.ObjectMeta{Name: name, Namespace: ns} }
	location := "/travel/apis"

	listing := func(portal string) *v1alpha1.PortalListing {
		return &v1alpha1.PortalListing{
			ObjectMeta: meta("travel-listing"),
			Spec: v1alpha1.PortalListingSpec{Type: portallisting.Type{
				Portal: refs.NamespacedName{Name: portal},
				APIs:   []portallisting.ApiEntry{{Ref: refs.NamespacedName{Name: "orders"}, Location: location}},
			}},
		}
	}
	page := func(portal string) *v1alpha1.Documentation {
		return &v1alpha1.Documentation{
			ObjectMeta: meta("travel-docs"),
			Spec: v1alpha1.DocumentationSpec{Type: documentation.Type{
				Name:     "Getting Started",
				PageType: documentation.GraviteeMarkdown,
				Content:  "# Getting Started",
				Location: &location,
				Portal:   &refs.NamespacedName{Name: portal},
			}},
		}
	}
	link := func(portal string) *v1alpha1.PortalLink {
		return &v1alpha1.PortalLink{
			ObjectMeta: meta("travel-link"),
			Spec: v1alpha1.PortalLinkSpec{Type: portallink.Type{
				Name:     "Status page",
				Href:     "https://status.example.com",
				Location: &location,
				Portal:   &refs.NamespacedName{Name: portal},
			}},
		}
	}

	BeforeEach(func() {
		platformCalls.Store(0)
		scheme := runtime.NewScheme()
		Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
		k8s.RegisterClient(fake.NewClientBuilder().WithScheme(scheme).Build())
	})

	It("warns and skips the dry run of a PortalListing", func() {
		warnings, err := listingadm.AdmissionCtrl{}.ValidateCreate(ctx, listing("unsynced"))

		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(ConsistOf("platform checks deferred (portal not synced: [travel/unsynced])"))
		Expect(platformCalls.Load()).To(BeZero())
	})

	It("warns and skips the dry run of a Documentation", func() {
		warnings, err := docsadm.AdmissionCtrl{}.ValidateCreate(ctx, page("unsynced"))

		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(ConsistOf("platform checks deferred (portal not synced: [travel/unsynced])"))
		Expect(platformCalls.Load()).To(BeZero())
	})

	It("warns and skips the dry run of a PortalLink", func() {
		warnings, err := linkadm.AdmissionCtrl{}.ValidateCreate(ctx, link("unsynced"))

		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(ConsistOf("platform checks deferred (portal not synced: [travel/unsynced])"))
		Expect(platformCalls.Load()).To(BeZero())
	})

	It("dry-runs a PortalListing once its portal has synced", func() {
		_, err := listingadm.AdmissionCtrl{}.ValidateCreate(ctx, listing("synced"))

		Expect(err).To(HaveOccurred())
		Expect(platformCalls.Load()).ToNot(BeZero())
	})
})
