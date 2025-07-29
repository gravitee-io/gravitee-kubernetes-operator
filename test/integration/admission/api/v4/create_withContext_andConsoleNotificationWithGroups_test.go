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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
)

var _ = Describe("Validate create", labels.WithContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := v4.AdmissionCtrl{}

	It("should create an API with notification groups and no events without warnings", func() {

		// no group setup
		fixtures := fixture.Builder().
			WithContext(constants.ContextWithCredentialsFile).
			WithGroup(constants.GroupFile).
			WithNotification(constants.NotificationWithGroupFile).
			WithAPIv4(constants.ApiV4WithContextFile).Build()

		// remove the first group ref as we don't want this one
		fixtures.Notification.Spec.Console.GroupRefs = fixtures.Notification.Spec.Console.GroupRefs[1:]

		fixtures.Apply()

		Eventually(func() error {
			warnings, err := admissionCtrl.ValidateCreate(ctx, fixtures.APIv4)

			if err != nil {
				return err
			}

			// as defined in the file
			return assert.Equals(
				"warning",
				admission.Warnings{},
				warnings,
			)
		}, constants.EventualTimeout, interval).Should(Succeed())

	})
})
