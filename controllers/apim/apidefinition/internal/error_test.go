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

package internal

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi/clienterror"
)

var _ = Describe("Errors", func() {
	DescribeTable("wrap error",
		func(given error, expected bool) {
			Expect(IsRecoverableError(wrapError(given))).To(Equal(expected))
		},
		Entry("With raw error", fmt.Errorf("raw error"), true),
		Entry("With not found error", clienterror.NewCrossIdNotFoundError("cross-id"), true),
		Entry("With illegal state error", clienterror.NewAmbiguousCrossIdError("cross-id", 2), true),
		Entry("With unauthorized api request error", clienterror.NewUnauthorizedApiRequestError("api-id"), false),
		Entry("With unauthorized cross ID request error", clienterror.NewUnauthorizedCrossIdRequestError("api-id"), false),
	)
})
