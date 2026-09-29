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

package amidentityprovider

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/amctx"
)

var _ = Describe("Validate AMContext delete (D7)", func() {
	ctx := context.Background()

	It("should reject deleting an AMContext a domain references", func() {
		fixtures := withDomain()
		fixtures.AMIdentityProvider = nil
		fixtures.Apply()

		_, err := amctx.AdmissionCtrl{}.ValidateDelete(ctx, fixtures.AMContext)
		Expect(err).To(MatchError(ContainSubstring("cannot be deleted because 1 AM security domains")))
	})

	It("should accept deleting an AMContext no domain references", func() {
		fixtures := withDomain()
		fixtures.AMSecurityDomain, fixtures.AMIdentityProvider = nil, nil
		fixtures.Apply()

		_, err := amctx.AdmissionCtrl{}.ValidateDelete(ctx, fixtures.AMContext)
		Expect(err).ToNot(HaveOccurred())
	})
})
