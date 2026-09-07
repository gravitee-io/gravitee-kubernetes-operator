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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	nav "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/navigation"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
)

func navStructure(entries ...*model.NavigationEntryDTO) *model.NavigationStructureDTO {
	return &model.NavigationStructureDTO{TopNavbar: entries}
}

func navEntry(path string, visibility nav.Visibility) *model.NavigationEntryDTO {
	return &model.NavigationEntryDTO{Path: path, Visibility: visibility}
}

var _ = Describe("PortalDTO.WithResolvedVisibility", func() {
	It("defaults a root entry that declares nothing to PUBLIC", func() {
		resolved := model.PortalDTO{Structure: navStructure(navEntry("/alpha", ""))}.WithResolvedVisibility()
		Expect(resolved.Structure.TopNavbar[0].Visibility).To(Equal(nav.Public))
	})

	It("keeps a declared visibility", func() {
		resolved := model.PortalDTO{Structure: navStructure(navEntry("/alpha", nav.Private))}.WithResolvedVisibility()
		Expect(resolved.Structure.TopNavbar[0].Visibility).To(Equal(nav.Private))
	})

	It("inherits the visibility declared by the direct parent", func() {
		resolved := model.PortalDTO{Structure: navStructure(
			navEntry("/alpha", nav.Private),
			navEntry("/alpha/docs", ""),
		)}.WithResolvedVisibility()
		Expect(resolved.Structure.TopNavbar[1].Visibility).To(Equal(nav.Private))
	})

	It("inherits from the closest declared ancestor across an implicit folder", func() {
		resolved := model.PortalDTO{Structure: navStructure(
			navEntry("/alpha", nav.Private),
			navEntry("/alpha/docs/guides", ""),
		)}.WithResolvedVisibility()
		Expect(resolved.Structure.TopNavbar[1].Visibility).To(Equal(nav.Private))
	})

	It("lets the closest declared ancestor win over a farther one", func() {
		resolved := model.PortalDTO{Structure: navStructure(
			navEntry("/alpha", nav.Private),
			navEntry("/alpha/docs", nav.Public),
			navEntry("/alpha/docs/guides", ""),
		)}.WithResolvedVisibility()
		Expect(resolved.Structure.TopNavbar[2].Visibility).To(Equal(nav.Public))
	})

	It("defaults to PUBLIC when no ancestor declares anything", func() {
		resolved := model.PortalDTO{Structure: navStructure(
			navEntry("/alpha", ""),
			navEntry("/alpha/docs", ""),
		)}.WithResolvedVisibility()
		Expect(resolved.Structure.TopNavbar[1].Visibility).To(Equal(nav.Public))
	})

	It("ignores a sibling prefix that is not an ancestor", func() {
		resolved := model.PortalDTO{Structure: navStructure(
			navEntry("/alpha", nav.Private),
			navEntry("/alphabet", ""),
		)}.WithResolvedVisibility()
		Expect(resolved.Structure.TopNavbar[1].Visibility).To(Equal(nav.Public))
	})

	It("tolerates a trailing slash on either side of the relation", func() {
		resolved := model.PortalDTO{Structure: navStructure(
			navEntry("/alpha/", nav.Private),
			navEntry("/alpha/docs/", ""),
		)}.WithResolvedVisibility()
		Expect(resolved.Structure.TopNavbar[1].Visibility).To(Equal(nav.Private))
	})

	It("does not mutate the receiver", func() {
		original := model.PortalDTO{Structure: navStructure(navEntry("/alpha", ""))}
		original.WithResolvedVisibility()
		Expect(original.Structure.TopNavbar[0].Visibility).To(BeEmpty())
	})

	It("leaves a DTO without a structure untouched", func() {
		resolved := model.PortalDTO{Name: "portal"}.WithResolvedVisibility()
		Expect(resolved.Structure).To(BeNil())
		Expect(resolved.Name).To(Equal("portal"))
	})

	It("leaves the deprecated navigation list untouched", func() {
		original := model.PortalDTO{
			Navigation: []model.NavigationPathDTO{{Path: "/legacy"}},
			Structure:  navStructure(navEntry("/alpha", "")),
		}
		resolved := original.WithResolvedVisibility()
		Expect(resolved.Navigation[0].Visibility).To(BeEmpty())
	})
})

var _ = Describe("APIV4DTO.WithResolvedVisibility", func() {
	It("resolves portal navigation entries the same way", func() {
		resolved := model.APIV4DTO{PortalNavigation: []*model.APIV4NavigationPathDTO{
			{Path: "/alpha", Visibility: nav.Private},
			{Path: "/alpha/docs"},
			{Path: "/beta"},
		}}.WithResolvedVisibility()
		Expect(resolved.PortalNavigation[0].Visibility).To(Equal(nav.Private))
		Expect(resolved.PortalNavigation[1].Visibility).To(Equal(nav.Private))
		Expect(resolved.PortalNavigation[2].Visibility).To(Equal(nav.Public))
	})

	It("leaves a DTO without portal navigation untouched", func() {
		resolved := model.APIV4DTO{}.WithResolvedVisibility()
		Expect(resolved.PortalNavigation).To(BeNil())
	})

	It("does not mutate the receiver", func() {
		original := model.APIV4DTO{PortalNavigation: []*model.APIV4NavigationPathDTO{{Path: "/alpha"}}}
		original.WithResolvedVisibility()
		Expect(original.PortalNavigation[0].Visibility).To(BeEmpty())
	})
})
