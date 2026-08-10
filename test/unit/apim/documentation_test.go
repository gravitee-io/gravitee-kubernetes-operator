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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	documentation "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/docs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

func newDocumentation(area *documentation.PageArea) *v1alpha1.Documentation {
	return &v1alpha1.Documentation{
		ObjectMeta: metav1.ObjectMeta{Name: "getting-started", Namespace: "default"},
		Spec: v1alpha1.DocumentationSpec{
			Type: documentation.Type{
				Name:     "Getting Started",
				PageType: documentation.GraviteeMarkdown,
				Content:  "# Getting Started",
				Location: utils.ToReference("/guides"),
				Order:    utils.ToReference(int32(1)),
				Area:     area,
			},
		},
	}
}

var _ = Describe("Documentation spec", func() {
	DescribeTable("area accessor",
		func(given *documentation.PageArea, expected documentation.PageArea) {
			Expect(newDocumentation(given).Spec.GetArea()).To(Equal(expected))
		},
		Entry("with no area", nil, documentation.PageArea("")),
		Entry("with top navbar area", utils.ToReference(documentation.TopNavbar), documentation.TopNavbar),
		Entry("with homepage area", utils.ToReference(documentation.Homepage), documentation.Homepage),
	)

	// The area is part of the spec hash, so that flipping a page to (or away
	// from) the homepage triggers a reconcile instead of being skipped as
	// unchanged by the last-spec-hash predicate.
	It("should be reflected in the spec hash", func() {
		navbar := newDocumentation(utils.ToReference(documentation.TopNavbar))
		homepage := newDocumentation(utils.ToReference(documentation.Homepage))
		unset := newDocumentation(nil)

		Expect(navbar.Spec.Hash()).ToNot(Equal(homepage.Spec.Hash()))
		Expect(unset.Spec.Hash()).ToNot(Equal(navbar.Spec.Hash()))
	})
})
