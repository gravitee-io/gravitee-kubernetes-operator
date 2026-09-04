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

package apim

import (
	nav "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/navigation"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("PortalLink Drift detection", func() {

	DescribeTable("equivalent values",
		func(crd, remote any) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
		},
		Entry("empty struct",
			model.PortalLinkDTO{},
			model.PortalLinkDTO{},
		),
		Entry("equal struct",
			completePortalLinkDTO(),
			completePortalLinkDTO(),
		),
		Entry("equal struct from CRD mapping",
			completePortalLinkDTO(),
			model.ToPortalLinkDTO(completePortalLinkCRD().Spec.Type, refs.NewNamespacedNameFromObject(completePortalLinkCRD()).HRID()),
		),
		Entry("unset visibility equivalent to any remote visibility resolved by APIM",
			model.PortalLinkDTO{},
			model.PortalLinkDTO{
				Visibility: nav.Private,
			},
		),
	)

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(completePortalLinkDTO(), completePortalLinkDTO())
		})
	})
})

func completePortalLinkDTO() model.PortalLinkDTO {
	GinkgoHelper()
	return loadFixture[model.PortalLinkDTO]("portallink_dto.json")
}

func completePortalLinkCRD() *v1alpha1.PortalLink {
	GinkgoHelper()
	fixture := loadFixture[v1alpha1.PortalLink]("portallink_crd.json")
	return &fixture
}
