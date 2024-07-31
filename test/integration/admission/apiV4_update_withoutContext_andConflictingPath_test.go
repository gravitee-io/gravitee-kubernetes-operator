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

package admission

import (
	"context"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	admission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate update", labels.WithoutContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := admission.AdmissionCtrl{}

	It("should return error on API update with conflicting path", func() {
		fixtures := fixture.
			Builder().
			WithAPI(constants.Api).
			WithAPIv4(constants.ApiV4).
			Build().
			Apply()

		By("checking that API update does not pass validation")

		Eventually(func() error {
			desired := fixtures.APIv4.DeepCopy()
			existingPath := fixtures.API.Spec.Proxy.VirtualHosts[0].Path
			listener, _ := desired.Spec.Listeners[0].ToListener().(*v4.HttpListener)
			listener.Paths[0].Path = existingPath
			_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.APIv4, desired)
			return err
		}, constants.EventualTimeout, interval).ShouldNot(Succeed())
	})
})
