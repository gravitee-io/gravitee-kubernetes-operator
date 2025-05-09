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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Delete", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval
	ctx := context.Background()

	When("deleting referenced API notification", func() {

		It("should not delete API V4 notification when referenced", func() {

			fixtures := fixture.Builder().
				WithAPIv4(constants.ApiV4).
				WithNotification(constants.NotificationNoGroupFile).
				Build().
				Apply()

			By("deleting Notification")

			Expect(manager.Client().Delete(ctx, fixtures.Notification)).To(Succeed())

			By("expecting to still find Notification")

			checkUntil := constants.ConsistentTimeout
			Consistently(func() error {
				return manager.GetLatest(ctx, fixtures.Notification)
			}, checkUntil, interval).Should(Succeed())

			By("deleting the API")

			Expect(manager.Client().Delete(ctx, fixtures.APIv4)).To(Succeed())

			By("expecting notification to have been deleted")

			Eventually(func() error {
				return assert.Deleted(ctx, "Notification", fixtures.Notification)
			}, timeout, interval).Should(Succeed())
		})

	})

})
