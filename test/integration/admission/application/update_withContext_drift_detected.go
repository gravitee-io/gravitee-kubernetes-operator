package application

import (
	"context"

	admission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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

var _ = Describe("Validate update", labels.WithContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := admission.AdmissionCtrl{}

	It("should detect drift", func() {
		fixtures := fixture.
			Builder().
			WithApplication(constants.Application).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("changing the remote application description")

		newApp := fixtures.Application.DeepCopy()
		newApp.PopulateIDs(nil, true)

		newApp.Spec.Description = "remote updated description"

		apim := apim.NewClient(ctx)

		_, err := apim.Applications.CreateOrUpdate(newApp)
		Expect(err).ToNot(HaveOccurred())
		By("changing the CRD application description")

		newApp.Spec.Description = "local CRD description"

		By("checking that application validation returns error")

		Eventually(func() error {
			_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Application, newApp)
			return assert.DriftDetected("description: \"local CRD description\" != \"remote updated description\"", err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})
})
