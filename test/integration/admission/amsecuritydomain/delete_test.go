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

package amsecuritydomain

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/amsecuritydomain"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate delete", labels.WithContext, func() {
	ctx := context.Background()
	admissionCtrl := amsecuritydomain.NewAdmissionCtrl()

	It("should always allow delete", func() {
		fixtures := fixture.Builder().
			WithAMContext(constants.AMContextFile).
			WithAMSecurityDomain(constants.AMSecurityDomainBasicFile).
			Build()
		fixtures.Apply()

		_, err := admissionCtrl.ValidateDelete(ctx, fixtures.AMSecurityDomain)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should allow delete when never created remotely and context does not exist", func() {
		fixtures := fixture.Builder().
			WithAMSecurityDomain(constants.AMSecurityDomainBasicFile).
			Build()

		_, err := admissionCtrl.ValidateDelete(ctx, fixtures.AMSecurityDomain)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should reject delete when created remotely but context is missing", func() {
		fixtures := fixture.Builder().
			WithAMSecurityDomain(constants.AMSecurityDomainBasicFile).
			Build()

		fixtures.AMSecurityDomain.Status.ID = "some-remote-id"

		_, err := admissionCtrl.ValidateDelete(ctx, fixtures.AMSecurityDomain)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("without contextRef"))
	})
})
