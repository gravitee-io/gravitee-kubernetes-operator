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
	documentation "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/docs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Documentation Drift detection", func() {

	DescribeTable("equivalent values",
		func(crd, remote any) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
		},
		Entry("empty struct",
			model.DocumentationDTO{},
			model.DocumentationDTO{},
		),
		Entry("equal struct",
			completeDocumentationDTO(),
			completeDocumentationDTO(),
		),
		Entry("equal struct from CRD mapping",
			completeDocumentationDTO(),
			model.ToDocumentationDTO(completeDocumentationCRD()),
		),
		Entry("empty area equivalent to remote TOP_NAVBAR",
			model.DocumentationDTO{},
			model.DocumentationDTO{
				Area: documentation.TopNavbar,
			},
		),
	)

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(completeDocumentationDTO(), completeDocumentationDTO())
		})
	})
})

func completeDocumentationDTO() model.DocumentationDTO {
	GinkgoHelper()
	return loadFixture[model.DocumentationDTO]("documentation_dto.json")
}

func completeDocumentationCRD() *v1alpha1.Documentation {
	GinkgoHelper()
	fixture := loadFixture[v1alpha1.Documentation]("documentation_crd.json")
	return &fixture
}
