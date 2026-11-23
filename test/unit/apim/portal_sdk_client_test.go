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
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/portal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/portaltheme"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimportal "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portal"
	apimportaltheme "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portaltheme"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

const automationEnvPath = "/automation/organizations/org1/environments/env1"

// portalState is a state the fake Automation API answers for any portal navigation resource: the ids,
// one finding of each severity, and the fields the payload declares.
func portalState(fields string) string {
	return `{"id": "7a1c", "organizationId": "DEFAULT", "environmentId": "DEV",
		"errors": {"severe": ["severe one"], "warning": ["warning one"]}, ` + fields + `}`
}

var _ = Describe("Portal navigation SDK packages", func() {
	var (
		ctx    context.Context
		fake   *fakeAutomationAPI
		client *apim.SDKClient
	)

	BeforeEach(func() {
		ctx = context.Background()
		fake = newFakeAutomationAPI()
		DeferCleanup(fake.Close)
		var err error
		client, err = apim.NewSDKClient(ctx, managementContext(fake.URL, bearer("token")))
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("PortalTheme", func() {
		newTheme := func() *v1alpha1.PortalTheme {
			return &v1alpha1.PortalTheme{
				ObjectMeta: metav1.ObjectMeta{Name: "brand", Namespace: "default"},
				Spec: v1alpha1.PortalThemeSpec{Type: portaltheme.Type{
					Name:       "Brand",
					Definition: portaltheme.Definition{CustomCss: new(".a {}")},
				}},
			}
		}
		toDTO := func(thm *v1alpha1.PortalTheme) sdk.PortalTheme {
			GinkgoHelper()
			dto, err := apimportaltheme.ToDTO(thm)
			Expect(err).ToNot(HaveOccurred())
			return dto
		}

		BeforeEach(func() {
			fake.response = portalState(`"hrid": "default-brand", "name": "Brand",
				"definition": {"customCss": ".a {}"}`)
		})

		It("upserts by HRID and copies the ids into the status", func() {
			thm := newTheme()
			resp, err := apimportaltheme.Upsert(ctx, client, toDTO(thm))
			Expect(err).ToNot(HaveOccurred())
			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodPut))
			Expect(req.Path).To(Equal(automationEnvPath + "/portal-themes"))
			Expect(req.Query.Get("dryRun")).To(Equal("false"))
			Expect(req.Body).To(HaveKeyWithValue("hrid", "default-brand"))

			Expect(apimportaltheme.UpdateStatus(ctx, thm, resp)).To(Succeed())
			Expect(thm.Status.ID).To(Equal("7a1c"))
			Expect(thm.Status.OrgID).To(Equal("DEFAULT"))
			Expect(thm.Status.EnvID).To(Equal("DEV"))
			Expect(k8s.IsAutomationAPIManaged(thm)).To(BeTrue())
		})

		It("returns the dry run findings", func() {
			errs, err := apimportaltheme.DryRun(ctx, client, toDTO(newTheme()))
			Expect(err).ToNot(HaveOccurred())
			Expect(fake.lastRequest().Query.Get("dryRun")).To(Equal("true"))
			Expect(errs.Severe).To(ConsistOf(MatchError("severe one")))
			Expect(errs.Warning).To(ConsistOf(MatchError("warning one")))
		})

		It("gets the remote as the payload a PUT sends", func() {
			dto := toDTO(newTheme())
			remote, err := apimportaltheme.GetRemote(ctx, client, dto)
			Expect(err).ToNot(HaveOccurred())
			Expect(fake.lastRequest().Path).To(Equal(automationEnvPath + "/portal-themes/default-brand"))
			Expect(remote).To(Equal(dto))
		})

		It("reports a theme missing in APIM as not found, and deletes one already gone", func() {
			fake.status = http.StatusNotFound
			dto := toDTO(newTheme())

			_, err := apimportaltheme.GetRemote(ctx, client, dto)
			Expect(errors.IsNotFound(err)).To(BeTrue())

			Expect(apimportaltheme.Delete(ctx, client, dto)).To(Succeed())
			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodDelete))
			Expect(req.Path).To(Equal(automationEnvPath + "/portal-themes/default-brand"))
		})
	})
	Describe("Portal", func() {
		newPortal := func(theme *refs.NamespacedName) *v1alpha1.Portal {
			return &v1alpha1.Portal{
				ObjectMeta: metav1.ObjectMeta{Name: "dev", Namespace: "default"},
				Spec: v1alpha1.PortalSpec{
					Type: portal.Type{
						Name: "Dev",
						Structure: &portal.NavigationStructure{TopNavbar: []*portal.NavigationEntry{
							{Path: "/alpha"},
						}},
					},
					Theme: theme,
				},
			}
		}
		toDTO := func(prtl *v1alpha1.Portal) sdk.Portal {
			GinkgoHelper()
			dto, err := apimportal.ToDTO(prtl)
			Expect(err).ToNot(HaveOccurred())
			return dto
		}

		BeforeEach(func() {
			fake.response = portalState(`"hrid": "default-dev", "name": "Dev",
				"structure": {"topNavbar": [{"path": "/alpha"}]}`)
		})

		It("upserts by HRID with the theme HRID, and copies the ids into the status", func() {
			prtl := newPortal(&refs.NamespacedName{Name: "brand"})
			resp, err := apimportal.Upsert(ctx, client, toDTO(prtl))
			Expect(err).ToNot(HaveOccurred())
			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodPut))
			Expect(req.Path).To(Equal(automationEnvPath + "/portals"))
			Expect(req.Query.Get("dryRun")).To(Equal("false"))
			Expect(req.Body).To(HaveKeyWithValue("hrid", "default-dev"))
			Expect(req.Body).To(HaveKeyWithValue("activeThemeHrid", "default-brand"))

			Expect(apimportal.UpdateStatus(ctx, prtl, resp)).To(Succeed())
			Expect(prtl.Status.ID).To(Equal("7a1c"))
			Expect(prtl.Status.OrgID).To(Equal("DEFAULT"))
			Expect(prtl.Status.EnvID).To(Equal("DEV"))
			Expect(k8s.IsAutomationAPIManaged(prtl)).To(BeTrue())
		})

		It("sends no active theme without a themeRef, nor an undeclared visibility", func() {
			_, err := apimportal.Upsert(ctx, client, toDTO(newPortal(nil)))
			Expect(err).ToNot(HaveOccurred())
			body := fake.lastRequest().Body
			Expect(body).ToNot(HaveKey("activeThemeHrid"))
			Expect(body["structure"]).To(Equal(map[string]any{"topNavbar": []any{map[string]any{"path": "/alpha"}}}))
		})

		It("returns the dry run findings", func() {
			errs, err := apimportal.DryRun(ctx, client, toDTO(newPortal(nil)))
			Expect(err).ToNot(HaveOccurred())
			Expect(fake.lastRequest().Query.Get("dryRun")).To(Equal("true"))
			Expect(errs.Severe).To(ConsistOf(MatchError("severe one")))
			Expect(errs.Warning).To(ConsistOf(MatchError("warning one")))
		})

		It("gets the remote as the payload a PUT sends", func() {
			dto := toDTO(newPortal(nil))
			remote, err := apimportal.GetRemote(ctx, client, dto)
			Expect(err).ToNot(HaveOccurred())
			Expect(fake.lastRequest().Path).To(Equal(automationEnvPath + "/portals/default-dev"))
			Expect(remote).To(Equal(dto))
		})

		It("reports a portal missing in APIM as not found, and deletes one already gone", func() {
			fake.status = http.StatusNotFound
			dto := toDTO(newPortal(nil))

			_, err := apimportal.GetRemote(ctx, client, dto)
			Expect(errors.IsNotFound(err)).To(BeTrue())

			Expect(apimportal.Delete(ctx, client, dto)).To(Succeed())
			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodDelete))
			Expect(req.Path).To(Equal(automationEnvPath + "/portals/default-dev"))
		})
	})
})
