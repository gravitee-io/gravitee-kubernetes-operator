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

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate update", labels.WithContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := v4.NewAdmissionCtrl()

	It("should dry-run the groups when the ignore-groups annotation is removed", func() {
		unknownGroup := random.GetName()

		fixtures := fixture.
			Builder().
			WithAPIv4(constants.ApiV4).
			WithContext(constants.ContextWithCredentialsFile).
			Build()
		fixtures.APIv4.Spec.Groups = []string{unknownGroup}
		fixtures.APIv4.Annotations = map[string]string{core.IgnoreGroupsAnnotation: "true"}
		fixtures = fixtures.Apply()

		By("removing the ignore-groups annotation, so the groups would be sent with the spec unchanged")

		newAPI := fixtures.APIv4.DeepCopy()
		delete(newAPI.Annotations, core.IgnoreGroupsAnnotation)

		Eventually(func() error {
			warnings, err := admissionCtrl.ValidateUpdate(ctx, fixtures.APIv4, newAPI)
			if err != nil {
				return err
			}
			return assert.ContainsInOrder("warnings", []string{
				"Group [" + unknownGroup + "] could not be found in environment [DEFAULT]",
			}, warnings)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})
})
