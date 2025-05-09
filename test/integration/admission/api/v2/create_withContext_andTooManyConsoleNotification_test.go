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

package v2

import (
	"context"

	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v2"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("Validate create", labels.WithContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := v2.AdmissionCtrl{}

	It("should return error on API creation when too many notification with target console are used", func() {

		fixture.
			Builder().
			WithContext(constants.ContextWithCredentialsFile).Build().Apply()
		notif1 := fixture.Builder().WithNotification(constants.NotificationNoGroupFile).Build().Apply()
		notif2 := fixture.Builder().WithNotification(constants.NotificationWithGroupFile).Build().Apply()

		fixtures := fixture.
			Builder().
			WithAPI(constants.ApiWithContextFile).Build()

		fixtures.API.Spec.NotificationsRefs = []refs.NamespacedName{
			{
				Name:      notif1.Notification.Name,
				Namespace: notif1.Notification.Namespace,
			}, {
				Name:      notif2.Notification.Name,
				Namespace: notif2.Notification.Namespace,
			},
		}

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, fixtures.API)

			return assert.Equals(
				"severe",
				errors.NewSeveref(
					"api references notification [%s] but there is already another console notification referenced",
					fixtures.API.Spec.NotificationsRefs[1].String(),
				).Error(),
				err.Error(),
			)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})
})
