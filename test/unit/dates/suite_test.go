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

package dates

import (
	"testing"
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/dates"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDates(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "package dates")
}

var _ = Describe("ParseRFC3339", func() {
	DescribeTable("parses valid timestamps",
		func(value string, expected time.Time) {
			got, err := dates.ParseRFC3339(value)
			Expect(err).NotTo(HaveOccurred())
			Expect(got.Equal(expected)).To(BeTrue())
		},
		Entry("RFC3339 Z",
			"2024-06-15T10:30:00Z",
			time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC),
		),
		Entry("RFC3339 numeric offset",
			"2024-06-15T13:30:00+03:00",
			time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC),
		),
		Entry("RFC3339Nano",
			"2024-06-15T10:30:00.123456789Z",
			time.Date(2024, 6, 15, 10, 30, 0, 123456789, time.UTC),
		),
		Entry("RFC3339Nano without fractional seconds",
			"2024-06-15T10:30:00.000000000Z",
			time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC),
		),
	)

	DescribeTable("rejects invalid timestamps",
		func(value string) {
			_, err := dates.ParseRFC3339(value)
			Expect(err).To(HaveOccurred())
		},
		Entry("empty", ""),
		Entry("not a date", "not-a-date"),
		Entry("missing zone", "2024-06-15T10:30:00"),
	)

	It("treats RFC3339 and RFC3339Nano of the same instant as equal", func() {
		rfc3339, err := dates.ParseRFC3339("2024-06-15T10:30:00Z")
		Expect(err).NotTo(HaveOccurred())
		nano, err := dates.ParseRFC3339("2024-06-15T10:30:00.000000000Z")
		Expect(err).NotTo(HaveOccurred())
		Expect(rfc3339.Equal(nano)).To(BeTrue())
	})
})
