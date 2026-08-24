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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Portal theme Drift detection", func() {

	DescribeTable("equivalent values",
		func(crd, remote any) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
		},
		Entry("empty struct",
			model.PortalThemeDTO{},
			model.PortalThemeDTO{},
		),
		Entry("equal struct",
			completePortalThemeDTO(),
			completePortalThemeDTO(),
		),
		Entry("equal struct from CRD mapping",
			completePortalThemeDTO(),
			model.ToPortalThemeDTO(completePortalThemeCRD()),
		),
	)

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(completePortalThemeDTO(), completePortalThemeDTO())
		})
	})
})

func completePortalThemeDTO() model.PortalThemeDTO {
	GinkgoHelper()
	return loadFixture[model.PortalThemeDTO]("portal_theme_dto.json")
}

func completePortalThemeCRD() *v1alpha1.PortalTheme {
	GinkgoHelper()
	fixture := loadFixture[v1alpha1.PortalTheme]("portal_theme_crd.json")
	return &fixture
}
