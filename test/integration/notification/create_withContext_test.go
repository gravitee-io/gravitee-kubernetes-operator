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
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
)

var _ = Describe("Create", labels.WithContext, func() {

	timeout := constants.EventualTimeout
	interval := constants.Interval
	ctx := context.Background()

	When("creating a notification", func() {
		It("should have conditions to true", func() {
			fixtures := fixture.Builder().WithNotification(constants.NotificationNoGroupFile).Build().Apply()
			Eventually(func() error {
				return assert.HasFinalizer(fixtures.Notification, core.NotificationFinalizer)
			})
		})
	})
	When("creating a notification with groups", func() {

		It("should have only have known group ids", func() {
			objects := fixture.Builder().
				WithContext(constants.ContextWithCredentialsFile).
				WithGroup(constants.GroupFile).
				WithNotification(constants.NotificationWithGroupFile).
				// there will be two groups reference because Build() leaves the
				// existing one untouched and add the existing group to the notification
				Build().
				Apply()

			Eventually(func() error {

				err := manager.GetLatest(ctx, objects.Notification)
				if err != nil {
					return err
				}

				if err := assert.Equals("GroupRefs", 2, len(objects.Notification.Spec.Console.GroupRefs)); err != nil {
					return err
				}
				if err := assert.Equals("Groups", 1, len(objects.Notification.Spec.Console.Groups)); err != nil {
					return err
				}
				if err := assert.Equals("Groups[0]",
					objects.Group.Status.ID,
					objects.Notification.Spec.Console.Groups[0]); err != nil {
					return err
				}
				return nil

			}, timeout, interval).Should(Succeed())

		})
	})
})
