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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/status"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	apimdocs "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/docs"
	apimportal "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portal"
	apimportallink "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portallink"
	apimportallisting "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portallisting"
	apimportaltheme "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portaltheme"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

// UpdateStatus writes the ids field by field: the conditions and the admission findings already in the
// status are kept.
var _ = Describe("Portal navigation UpdateStatus on the APIM SDK types", func() {
	ctx := context.Background()

	accepted := func() []metav1.Condition {
		return []metav1.Condition{{
			Type: "Accepted", Status: metav1.ConditionTrue, Reason: "Accepted", LastTransitionTime: metav1.Now(),
		}}
	}
	findings := status.Errors{Warning: []string{"kept warning"}}

	expectKept := func(conditions []metav1.Condition, errs status.Errors, id, orgID, envID string) {
		GinkgoHelper()
		Expect(id).To(Equal("7a1c"))
		Expect(orgID).To(Equal("DEFAULT"))
		Expect(envID).To(Equal("DEV"))
		Expect(errs).To(Equal(findings))
		Expect(meta.IsStatusConditionTrue(conditions, "Accepted")).To(BeTrue())
	}

	It("keeps a theme's conditions and findings", func() {
		thm := &v1alpha1.PortalTheme{}
		thm.Status.Conditions, thm.Status.Errors = accepted(), findings
		resp := apimportaltheme.Response{State: sdk.PortalThemeState{Id: new("7a1c")}, OrgID: "DEFAULT", EnvID: "DEV"}
		Expect(apimportaltheme.UpdateStatus(ctx, thm, resp)).To(Succeed())
		expectKept(thm.Status.Conditions, thm.Status.Errors, thm.Status.ID, thm.Status.OrgID, thm.Status.EnvID)
		Expect(k8s.IsAutomationAPIManaged(thm)).To(BeTrue())
	})

	It("keeps a portal's conditions and findings", func() {
		prtl := &v1alpha1.Portal{}
		prtl.Status.Conditions, prtl.Status.Errors = accepted(), findings
		resp := apimportal.Response{State: sdk.PortalState{Id: new("7a1c")}, OrgID: "DEFAULT", EnvID: "DEV"}
		Expect(apimportal.UpdateStatus(ctx, prtl, resp)).To(Succeed())
		expectKept(prtl.Status.Conditions, prtl.Status.Errors, prtl.Status.ID, prtl.Status.OrgID, prtl.Status.EnvID)
		Expect(k8s.IsAutomationAPIManaged(prtl)).To(BeTrue())
	})

	It("keeps a link's conditions and findings", func() {
		link := &v1alpha1.PortalLink{}
		link.Status.Conditions, link.Status.Errors = accepted(), findings
		resp := apimportallink.Response{State: sdk.PortalLinkState{Id: new("7a1c")}, OrgID: "DEFAULT", EnvID: "DEV"}
		Expect(apimportallink.UpdateStatus(ctx, link, resp)).To(Succeed())
		expectKept(link.Status.Conditions, link.Status.Errors, link.Status.ID, link.Status.OrgID, link.Status.EnvID)
		Expect(k8s.IsAutomationAPIManaged(link)).To(BeTrue())
	})

	It("keeps a listing's conditions and findings", func() {
		listing := &v1alpha1.PortalListing{}
		listing.Status.Conditions, listing.Status.Errors = accepted(), findings
		resp := apimportallisting.Response{State: sdk.PortalListingState{Id: new("7a1c")}, OrgID: "DEFAULT", EnvID: "DEV"}
		Expect(apimportallisting.UpdateStatus(ctx, listing, resp)).To(Succeed())
		expectKept(listing.Status.Conditions, listing.Status.Errors,
			listing.Status.ID, listing.Status.OrgID, listing.Status.EnvID)
		Expect(k8s.IsAutomationAPIManaged(listing)).To(BeTrue())
	})

	It("keeps a documentation page's conditions and findings", func() {
		doc := &v1alpha1.Documentation{}
		doc.Status.Conditions, doc.Status.Errors = accepted(), findings
		resp := apimdocs.Response{State: sdk.DocumentationState{Id: new("7a1c")}, OrgID: "DEFAULT", EnvID: "DEV"}
		Expect(apimdocs.UpdateStatus(ctx, doc, resp)).To(Succeed())
		expectKept(doc.Status.Conditions, doc.Status.Errors, doc.Status.ID, doc.Status.OrgID, doc.Status.EnvID)
		Expect(k8s.IsAutomationAPIManaged(doc)).To(BeTrue())
	})
})
