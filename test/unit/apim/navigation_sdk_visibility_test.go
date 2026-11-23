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

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/navigation"
)

func sdkNavEntry(path string, visibility sdk.PortalVisibility) sdk.PortalNavigationPath {
	entry := sdk.PortalNavigationPath{Path: path}
	if visibility != "" {
		entry.Visibility = &visibility
	}
	return entry
}

func resolvedVisibilities(entries ...sdk.PortalNavigationPath) []sdk.PortalVisibility {
	resolved := navigation.WithResolvedTopNavbar(entries)
	out := make([]sdk.PortalVisibility, len(resolved))
	for i, entry := range resolved {
		Expect(entry.Visibility).ToNot(BeNil())
		out[i] = *entry.Visibility
	}
	return out
}

var _ = Describe("navigation.WithResolvedTopNavbar", func() {
	public, private := sdk.PortalVisibilityPUBLIC, sdk.PortalVisibilityPRIVATE

	It("defaults a root entry that declares nothing to PUBLIC", func() {
		Expect(resolvedVisibilities(sdkNavEntry("/alpha", ""))).To(Equal([]sdk.PortalVisibility{public}))
	})

	It("keeps a declared visibility", func() {
		Expect(resolvedVisibilities(sdkNavEntry("/alpha", private))).To(Equal([]sdk.PortalVisibility{private}))
	})

	It("inherits the visibility declared by the direct parent", func() {
		Expect(resolvedVisibilities(
			sdkNavEntry("/alpha", private),
			sdkNavEntry("/alpha/docs", ""),
		)).To(Equal([]sdk.PortalVisibility{private, private}))
	})

	It("inherits from the closest declared ancestor across an implicit folder", func() {
		Expect(resolvedVisibilities(
			sdkNavEntry("/alpha", private),
			sdkNavEntry("/alpha/docs/guides", ""),
		)).To(Equal([]sdk.PortalVisibility{private, private}))
	})

	It("lets the closest declared ancestor win over a farther one", func() {
		Expect(resolvedVisibilities(
			sdkNavEntry("/alpha", private),
			sdkNavEntry("/alpha/docs", public),
			sdkNavEntry("/alpha/docs/guides", ""),
		)).To(Equal([]sdk.PortalVisibility{private, public, public}))
	})

	It("ignores a sibling prefix that is not an ancestor", func() {
		Expect(resolvedVisibilities(
			sdkNavEntry("/alpha", private),
			sdkNavEntry("/alphabet", ""),
		)).To(Equal([]sdk.PortalVisibility{private, public}))
	})

	It("tolerates a trailing slash on either side of the relation", func() {
		Expect(resolvedVisibilities(
			sdkNavEntry("/alpha/", private),
			sdkNavEntry("/alpha/docs/", ""),
		)).To(Equal([]sdk.PortalVisibility{private, private}))
	})

	It("does not write into the input slice", func() {
		input := []sdk.PortalNavigationPath{sdkNavEntry("/alpha", private), sdkNavEntry("/alpha/docs", "")}
		resolved := navigation.WithResolvedTopNavbar(input)
		Expect(input[1].Visibility).To(BeNil())
		*resolved[0].Visibility = public
		Expect(*input[0].Visibility).To(Equal(private))
	})

	It("returns an empty list as it is", func() {
		Expect(navigation.WithResolvedTopNavbar(nil)).To(BeNil())
	})
})
