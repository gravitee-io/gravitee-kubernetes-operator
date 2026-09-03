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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ApplicationClientCertificateDTO validity", func() {
	DescribeTable("Expired",
		func(endsAt string, expired bool) {
			Expect(model.ApplicationClientCertificateDTO{EndsAt: endsAt}.Expired()).To(Equal(expired))
		},
		Entry("missing", "", false),
		Entry("past", "2000-01-01T00:00:00Z", true),
		Entry("future", "2099-01-01T00:00:00Z", false),
	)

	DescribeTable("Scheduled",
		func(startsAt string, scheduled bool) {
			Expect(model.ApplicationClientCertificateDTO{StartsAt: startsAt}.Scheduled()).To(Equal(scheduled))
		},
		Entry("missing", "", false),
		Entry("past", "2000-01-01T00:00:00Z", false),
		Entry("future", "2099-01-01T00:00:00Z", true),
	)

	It("does not treat a scheduled cert as expired", func() {
		c := model.ApplicationClientCertificateDTO{
			StartsAt: "2099-01-01T00:00:00Z",
			EndsAt:   "2099-12-31T00:00:00Z",
		}
		Expect(c.Scheduled()).To(BeTrue())
		Expect(c.Expired()).To(BeFalse())
	})
})
