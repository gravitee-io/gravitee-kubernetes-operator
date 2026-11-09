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

package v4

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	admissionv4 "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// GKO-3433: what APIM stores differently from what the CRD sends is not drift.
var _ = Describe("Validate drift on values APIM normalises", labels.WithContext, func() {
	ctx := context.Background()
	admissionCtrl := admissionv4.AdmissionCtrl{}

	It("should not drift on a lowercase flow mode APIM stores in uppercase", func() {
		fixtures := fixture.Builder().
			WithAPIv4(constants.ApiV4WithContextFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build()
		fixtures.APIv4.Spec.FlowExecution.Mode = v4.FlowMode("default")
		fixtures.Apply()

		By("changing the API description after the first sync")
		newAPI := fixtures.APIv4.DeepCopy()
		setDescription(newAPI, updatedDescription)

		_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.APIv4, newAPI)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should not drift on a plan order the Automation API does not keep", func() {
		fixtures := fixture.Builder().
			WithAPIv4(constants.ApiV4WithContextFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build()
		(*fixtures.APIv4.Spec.Plans)["KeyLess"].Plan = &base.Plan{Order: new(7)}
		fixtures.Apply()

		By("changing the API description after the first sync")
		newAPI := fixtures.APIv4.DeepCopy()
		setDescription(newAPI, updatedDescription)

		_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.APIv4, newAPI)
		Expect(err).ToNot(HaveOccurred())
	})
})
