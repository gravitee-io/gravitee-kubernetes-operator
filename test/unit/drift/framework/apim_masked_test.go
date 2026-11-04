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

package framework

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const apimMask = "••••••••••••"

type withAPIMMaskedSecret struct {
	Name   string `json:"name,omitempty"`
	Secret string `json:"secret,omitempty" drift:"apim-masked"`
}

var _ = Describe("apim-masked", func() {
	DescribeTable("no drift",
		func(crd, remote withAPIMMaskedSecret) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
		},
		Entry("the remote answers with the mask",
			withAPIMMaskedSecret{Secret: "s3cret"},
			withAPIMMaskedSecret{Secret: apimMask},
		),
		Entry("neither side holds a value",
			withAPIMMaskedSecret{},
			withAPIMMaskedSecret{},
		),
	)

	It("drifts when the CRD declares a secret the remote does not hold", func() {
		result := drift.DetectWithNamespace(
			withAPIMMaskedSecret{Secret: "s3cret"},
			withAPIMMaskedSecret{},
			"",
		)

		Expect(result.DriftDetected()).To(BeTrue())
	})

	It("drifts when the mask stands against an empty CRD value", func() {
		result := drift.DetectWithNamespace(
			withAPIMMaskedSecret{},
			withAPIMMaskedSecret{Secret: apimMask},
			"",
		)

		Expect(result.DriftDetected()).To(BeTrue())
	})

	Describe("redaction", func() {
		It("names the property and prints neither value", func() {
			result := drift.DetectWithNamespace(
				withAPIMMaskedSecret{Secret: "s3cret"},
				withAPIMMaskedSecret{Secret: "a different secret"},
				"",
			)
			report := result.String()

			Expect(result.DriftDetected()).To(BeTrue())
			Expect(report).To(ContainSubstring("secret"))
			Expect(report).NotTo(ContainSubstring("s3cret"))
			Expect(report).NotTo(ContainSubstring("a different secret"))
		})

		It("leaves untagged fields of the same object printed", func() {
			result := drift.DetectWithNamespace(
				withAPIMMaskedSecret{Name: "from Git", Secret: "s3cret"},
				withAPIMMaskedSecret{Name: "renamed remotely", Secret: "a different secret"},
				"",
			)
			report := result.String()

			Expect(report).To(ContainSubstring("from Git"))
			Expect(report).To(ContainSubstring("renamed remotely"))
			Expect(report).NotTo(ContainSubstring("s3cret"))
		})
	})
})
