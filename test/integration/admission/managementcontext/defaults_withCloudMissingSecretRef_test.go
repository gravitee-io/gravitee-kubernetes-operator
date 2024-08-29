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

package managementcontext

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/mctx"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
)

var _ = Describe("Default create", labels.WithContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := mctx.AdmissionCtrl{}

	It("should fail configuring context with missing secret", func() {
		By("setting a context with wrong secret ref")

		fixtures := fixture.Builder().
			WithContext(constants.ContextCloudWithUnknownSecretRefFile).
			Build().
			Apply()

		By("defaulting the context")

		Consistently(func() error {
			err := admissionCtrl.Default(ctx, fixtures.Context)
			return assert.Equals(
				"error",
				errors.NewSevere(
					"secret [default/unknown] doesn't exist in the cluster"),
				err)
		}, constants.ConsistentTimeout, interval).Should(Succeed())
	})

})
