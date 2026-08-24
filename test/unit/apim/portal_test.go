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

	nav "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/navigation"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/portal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

func navPath(path string) *portal.NavigationEntry {
	return &portal.NavigationEntry{Path: path}
}

func legacyPath(path string) *nav.NavigationPath {
	return &nav.NavigationPath{Path: path}
}

func newPortalSpec(structure *portal.NavigationStructure, navigation ...*nav.NavigationPath) *v1alpha1.Portal {
	return &v1alpha1.Portal{
		ObjectMeta: metav1.ObjectMeta{Name: "default-portal", Namespace: "default"},
		Spec: v1alpha1.PortalSpec{
			Type: portal.Type{
				Name:       "Default Portal",
				Structure:  structure,
				Navigation: navigation,
			},
		},
	}
}

func newPortalWithTheme(ref *refs.NamespacedName) *v1alpha1.Portal {
	prtl := newPortalSpec(nil)
	prtl.Spec.Theme = ref
	return prtl
}

var _ = Describe("Portal active theme", func() {
	It("should be reflected in the spec hash", func() {
		a := newPortalWithTheme(&refs.NamespacedName{Name: "theme-a"})
		b := newPortalWithTheme(&refs.NamespacedName{Name: "theme-b"})
		unset := newPortalWithTheme(nil)

		Expect(a.Spec.Hash()).ToNot(Equal(b.Spec.Hash()))
		Expect(unset.Spec.Hash()).ToNot(Equal(a.Spec.Hash()))
	})
})

var _ = Describe("Portal spec", func() {
	DescribeTable("deprecated navigation detection",
		func(prtl *v1alpha1.Portal, expected bool) {
			Expect(prtl.Spec.UsesDeprecatedNavigation()).To(Equal(expected))
		},
		Entry("with no navigation", newPortalSpec(nil), false),
		Entry("with navigation", newPortalSpec(nil, legacyPath("/legacy")), true),
		Entry("with both", newPortalSpec(&portal.NavigationStructure{}, legacyPath("/legacy")), true),
	)

	It("should be reflected in the spec hash", func() {
		structure := newPortalSpec(&portal.NavigationStructure{TopNavbar: []*portal.NavigationEntry{navPath("/new")}})
		legacy := newPortalSpec(nil, legacyPath("/new"))
		empty := newPortalSpec(nil)

		Expect(structure.Spec.Hash()).ToNot(Equal(legacy.Spec.Hash()))
		Expect(empty.Spec.Hash()).ToNot(Equal(structure.Spec.Hash()))
	})
})
