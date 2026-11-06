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

	admissionv4 "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate drift for templated APIs", labels.WithContext, func() {
	ctx := context.Background()
	admissionCtrl := admissionv4.AdmissionCtrl{}

	It("should not drift when the CRD changes a templated value", func() {
		fixtures := fixture.Builder().
			WithAPIv4(constants.ApiV4WithTemplatingFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("replacing the templated endpoint target with a plain one")
		newAPI := fixtures.APIv4.DeepCopy()
		newAPI.Spec.EndpointGroups[0].Endpoints[0].Config.Put("target", "https://api.gravitee.io/whattimeisit")

		_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.APIv4, newAPI)
		Expect(err).ToNot(HaveOccurred())
	})
})
