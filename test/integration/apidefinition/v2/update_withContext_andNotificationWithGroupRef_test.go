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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/notification"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Update", labels.WithContext, func() {

	ctx := context.Background()

	var fixtures *fixture.Objects

	BeforeEach(func() {
		By("Applying an API with groups and notifications")
		fixtures = fixture.
			Builder().
			WithContext(constants.ContextWithCredentialsFile).
			WithGroup(constants.GroupFile).
			WithNotification(constants.NotificationWithGroupFile).
			WithAPI(constants.ApiWithNotificationsAndGroups).
			Build().
			Apply()
	})

	It("should update notification events after creation", func() {

		fixtures.Notification.Spec.Console.APIEvents = []notification.ApiEvent{"API_STARTED", "API_STOPPED", "APIKEY_EXPIRED"}
		Expect(manager.UpdateSafely(ctx, fixtures.Notification)).To(Succeed())

		client := apim.NewClient(ctx)

		Eventually(func() error {

			console, err := client.Notification.GetConsoleNotificationConfiguration(fixtures.API.GetID())
			if err != nil {
				return err
			}

			By("Checking that the notification has been added to API")
			if err := assert.NotNil("console", console); err != nil {
				return err
			}

			if err := assert.Equals("console.Hooks",
				[]notification.ApiEvent{"API_STARTED", "API_STOPPED", "APIKEY_EXPIRED"},
				console.Hooks); err != nil {
				return err
			}

			return nil

		})
	})

	It("should update notification groups after creation", func() {

		fixtures.Notification.Spec.Console.GroupRefs = make([]refs.NamespacedName, 0)

		Expect(manager.UpdateSafely(ctx, fixtures.Notification)).To(Succeed())

		client := apim.NewClient(ctx)

		Eventually(func() error {

			console, err := client.Notification.GetConsoleNotificationConfiguration(fixtures.API.GetID())
			if err != nil {
				return err
			}
			By("Checking that the notification exists for the API")
			if err := assert.NotNil("console", console); err != nil {
				return err
			}

			if err := assert.SliceOfSize("console.Groups", console.Groups, 0); err != nil {
				return err
			}

			return nil

		})
	})

	It("should delete notification after creation", func() {

		fixtures.API.Spec.NotificationsRefs = make([]refs.NamespacedName, 0)

		Expect(manager.UpdateSafely(ctx, fixtures.API)).To(Succeed())

		client := apim.NewClient(ctx)

		Eventually(func() error {

			console, err := client.Notification.GetConsoleNotificationConfiguration(fixtures.API.GetID())
			if err != nil {
				return err
			}

			By("Checking that the notification does not exist for the API")
			if err := assert.Nil("console", console); err != nil {
				return err
			}

			return nil

		})
	})

})
