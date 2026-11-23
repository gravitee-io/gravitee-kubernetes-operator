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
	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ApplicationClientCertificateDTO validity", func() {
	DescribeTable("Expired",
		func(endsAt string, expired bool) {
			Expect(sdk.ClientCertificate{EndsAt: optional(endsAt)}.Expired()).To(Equal(expired))
		},
		Entry("missing", "", false),
		Entry("past", "2000-01-01T00:00:00Z", true),
		Entry("future", "2099-01-01T00:00:00Z", false),
	)

	DescribeTable("Scheduled",
		func(startsAt string, scheduled bool) {
			Expect(sdk.ClientCertificate{StartsAt: optional(startsAt)}.Scheduled()).To(Equal(scheduled))
		},
		Entry("missing", "", false),
		Entry("past", "2000-01-01T00:00:00Z", false),
		Entry("future", "2099-01-01T00:00:00Z", true),
	)

	It("does not treat a scheduled cert as expired", func() {
		c := sdk.ClientCertificate{
			StartsAt: new("2099-01-01T00:00:00Z"),
			EndsAt:   new("2099-12-31T00:00:00Z"),
		}
		Expect(c.Scheduled()).To(BeTrue())
		Expect(c.Expired()).To(BeFalse())
	})
})

// optional maps a missing date ("") to an unset one.
func optional(date string) *string {
	if date == "" {
		return nil
	}
	return &date
}
