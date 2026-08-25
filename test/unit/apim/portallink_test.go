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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/portallink"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

func newPortalAttachedLink() *v1alpha1.PortalLink {
	return &v1alpha1.PortalLink{
		ObjectMeta: metav1.ObjectMeta{Name: "external-docs", Namespace: "default"},
		Spec: v1alpha1.PortalLinkSpec{
			Type: portallink.Type{
				Portal: &refs.NamespacedName{Name: "default-portal"},
				Name:   "External Docs",
				Href:   "https://docs.example.com",
			},
		},
	}
}

func newApiAttachedLink() *v1alpha1.PortalLink {
	return &v1alpha1.PortalLink{
		ObjectMeta: metav1.ObjectMeta{Name: "external-docs", Namespace: "default"},
		Spec: v1alpha1.PortalLinkSpec{
			Type: portallink.Type{
				API:  &refs.NamespacedName{Name: "pets-api"},
				Name: "External Docs",
				Href: "https://docs.example.com",
			},
		},
	}
}

var _ = Describe("PortalLink spec", func() {
	It("should report IsPortalLink/IsApiLink correctly for a portal-attached link", func() {
		link := newPortalAttachedLink()
		Expect(link.IsPortalLink()).To(BeTrue())
		Expect(link.IsApiLink()).To(BeFalse())
		Expect(link.GetPortalRef()).To(Equal(&refs.NamespacedName{Name: "default-portal"}))
		Expect(link.GetApiRef()).To(BeNil())
	})

	It("should report IsPortalLink/IsApiLink correctly for an API-attached link", func() {
		link := newApiAttachedLink()
		Expect(link.IsPortalLink()).To(BeFalse())
		Expect(link.IsApiLink()).To(BeTrue())
		Expect(link.GetPortalRef()).To(BeNil())
		Expect(link.GetApiRef()).To(Equal(&refs.NamespacedName{Name: "pets-api"}))
	})

	It("should reflect a parent-kind change in the spec hash", func() {
		portalLink := newPortalAttachedLink()
		apiLink := newApiAttachedLink()
		Expect(portalLink.Spec.Hash()).ToNot(Equal(apiLink.Spec.Hash()))
	})
})
