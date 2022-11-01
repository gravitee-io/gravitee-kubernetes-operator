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

package clienterror

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client Errors", func() {
	DescribeTable("Is Unauthorized",
		func(given error, expected bool) {
			Expect(IsUnauthorized(given)).To(Equal(expected))
		},
		Entry("With raw error", fmt.Errorf("raw error"), false),
		Entry("With nil error", nil, false),
		Entry("With unauthorized API error", NewUnauthorizedApiRequestError("api-id"), true),
		Entry("With unauthorized cross ID error", NewUnauthorizedCrossIdRequestError("cross-id"), true),
	)

	DescribeTable("Is Not Found",
		func(given error, expected bool) {
			Expect(IsNotFound(given)).To(Equal(expected))
		},
		Entry("With raw error", fmt.Errorf("raw error"), false),
		Entry("With nil error", nil, false),
		Entry("With API not found error", NewApiNotFoundError("api-id"), true),
		Entry("With cross ID not found error", NewCrossIdNotFoundError("cross-id"), true),
	)

	DescribeTable("Is Illegal State",
		func(given error, expected bool) {
			Expect(IsIllegalState(given)).To(Equal(expected))
		},
		Entry("With raw error", fmt.Errorf("raw error"), false),
		Entry("With nil error", nil, false),
		Entry("With ambiguous cross ID error", NewAmbiguousCrossIdError("cross-id", 2), true),
	)
})
