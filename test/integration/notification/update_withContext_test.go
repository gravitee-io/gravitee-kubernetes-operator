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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
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

	When("Updating groups", func() {

		It("should have groups added", func() {

			By("creating a notification with events")
			objects := fixture.Builder().
				WithContext(constants.ContextWithCredentialsFile).
				WithGroup(constants.GroupFile).
				WithNotification(constants.NotificationWithGroupFile).
				Build().
				Apply()

			newGroup := fixture.Builder().
				WithGroup(constants.GroupFile).Build().Apply()

			By("adding a new groups")
			dc := objects.Notification.DeepCopy()
			dc.Spec.Console.GroupRefs = append(dc.Spec.Console.GroupRefs, refs.NamespacedName{
				Name:      newGroup.Group.Name,
				Namespace: constants.Namespace,
			})

			Eventually(func() error {
				return manager.UpdateSafely(ctx, dc)
			}, timeout, interval).Should(Succeed())

			Eventually(func() error {
				err := manager.GetLatest(ctx, objects.Notification)
				if err != nil {
					return err
				}
				return assert.Equals("groups", objects.Notification.Spec.Console.Groups,
					[]string{objects.Group.Status.ID, newGroup.Group.Status.ID})
			}, timeout, interval).Should(Succeed())

		})

		It("should have groups removed", func() {

			By("creating a notification with events")
			objects := fixture.Builder().
				WithContext(constants.ContextWithCredentialsFile).
				WithGroup(constants.GroupFile).
				WithNotification(constants.NotificationWithGroupFile).
				Build().
				Apply()

			By("updating removing groups")
			dc := objects.Notification.DeepCopy()
			dc.Spec.Console.GroupRefs = make([]refs.NamespacedName, 0)

			Eventually(func() error {
				return manager.UpdateSafely(ctx, dc)
			}, timeout, interval).Should(Succeed())
			Eventually(func() error {
				err := manager.GetLatest(ctx, objects.Notification)
				if err != nil {
					return err
				}
				if err := assert.Equals("Group",
					objects.Notification.Spec.Console.Groups,
					[]string{}); err != nil {
					return err
				}
				return nil
			}, timeout, interval).Should(Succeed())
		})
	})

})
