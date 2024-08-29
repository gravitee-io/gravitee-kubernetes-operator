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

package mctx

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/management"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cloud object context methods", func() {
	type expectations struct {
		hasCloud bool
		enabled  bool
	}
	DescribeTable("possible inputs", func(given *management.Cloud, expected expectations) {
		ctx := management.Context{
			Cloud: given,
		}
		Expect(ctx.HasCloud()).To(Equal(expected.hasCloud))
		if expected.hasCloud {
			Expect(ctx.GetCloud().IsEnabled()).To(Equal(expected.enabled))
		}

	},
		Entry("nil", nil, expectations{hasCloud: false, enabled: false}),
		Entry("not nil", &management.Cloud{}, expectations{hasCloud: true, enabled: false}),
		Entry("token", &management.Cloud{Token: "foo"}, expectations{hasCloud: true, enabled: true}),
		Entry("empty ref",
			&management.Cloud{
				SecretRef: &refs.NamespacedName{},
			}, expectations{hasCloud: true, enabled: false}),
		Entry("ref",
			&management.Cloud{
				SecretRef: &refs.NamespacedName{Name: "foo"},
			}, expectations{hasCloud: true, enabled: true}),
	)
})
