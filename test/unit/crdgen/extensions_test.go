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

package crdgen_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/hack/crdgen/amend"
)

var _ = Describe("ApplyExtensions", func() {
	DescribeTable("runs the custom extensions and removes their keys",
		func(property, want string) {
			doc := thing(property, false)
			Expect(amend.ApplyExtensions(doc)).To(Succeed())

			Expect(value(doc).Description).To(Equal(want))
			for key := range value(doc).Extensions {
				Expect(key == "x-kubebuilder" || strings.HasPrefix(key, "x-gko-")).
					To(BeFalse(), "custom key %s left in the amended OAS", key)
			}
		},
		Entry("x-kubebuilder appends raw markers",
			"{type: integer, description: Doc., x-kubebuilder: ['validation:Minimum=1', 'validation:Maximum=100']}",
			"Doc.\n+kubebuilder:validation:Minimum=1\n+kubebuilder:validation:Maximum=100"),
		Entry("x-gko-default keeps a bool default", "{type: boolean, default: false, x-gko-default: keep}",
			"+kubebuilder:default=false"),
		Entry("x-gko-default keeps a string default", "{type: string, default: AUTO, x-gko-default: keep}",
			`+kubebuilder:default="AUTO"`),
		Entry("standard extensions are left alone", "{type: string, x-go-type: string}", ""),
	)

	DescribeTable("fails on an invalid custom extension",
		func(property, wantErr string) {
			Expect(amend.ApplyExtensions(thing(property, false))).To(MatchError(ContainSubstring(wantErr)))
		},
		Entry("x-gko-default without an OAS default", "{type: boolean, x-gko-default: keep}",
			"Thing/properties/value: x-gko-default: no OAS default to keep"),
		Entry("unknown custom key", "{type: string, x-gko-unknown: true}",
			"Thing/properties/value: x-gko-unknown: unknown extension"),
	)
})
