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

package notification

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/notification"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
)

var _ = Describe("Update", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval
	ctx := context.Background()

	When("updating events", func() {
		It("should have events added", func() {

			By("creating a notification with events")
			objects := fixture.Builder().
				WithNotification(constants.NotificationNoGroupFile).
				Build().
				Apply()

			By("updating adding new event")
			dc := objects.Notification.DeepCopy()
			dc.Spec.Console.APIEvents = append(dc.Spec.Console.APIEvents, "APIKEY_REVOKED")

			Eventually(func() error {
				return manager.UpdateSafely(ctx, dc)
			}, timeout, interval).Should(Succeed())
			Eventually(func() error {
				err := manager.GetLatest(ctx, objects.Notification)
				if err != nil {
					return err
				}
				if err := assert.Equals("APIEvents",
					objects.Notification.Spec.Console.APIEvents,
					[]notification.ApiEvent{"API_STARTED", "API_STOPPED", "APIKEY_REVOKED"}); err != nil {
					return err
				}
				return nil
			}, timeout, interval).Should(Succeed())
		})

		It("should have events removed", func() {

			By("creating a notification with events")
			objects := fixture.Builder().
				WithNotification(constants.NotificationNoGroupFile).
				Build().
				Apply()

			By("updating removing first event")
			dc := objects.Notification.DeepCopy()
			dc.Spec.Console.APIEvents = dc.Spec.Console.APIEvents[1:]

			Eventually(func() error {
				return manager.UpdateSafely(ctx, dc)
			}, timeout, interval).Should(Succeed())
			Eventually(func() error {
				err := manager.GetLatest(ctx, objects.Notification)
				if err != nil {
					return err
				}
				if err := assert.Equals("APIEvents",
					objects.Notification.Spec.Console.APIEvents,
					[]notification.ApiEvent{"API_STOPPED"}); err != nil {
					return err
				}
				return nil
			}, timeout, interval).Should(Succeed())
		})
	})

})
