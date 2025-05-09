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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create", labels.WithContext, func() {

	ctx := context.Background()

	It("should add notification to API after creation", func() {

		By("Applying context and groups")

		groupFixture := fixture.
			Builder().
			WithContext(constants.ContextWithCredentialsFile).
			WithGroup(constants.GroupFile).
			Build().Apply()

		By("Applying an API and notifications")

		apiFixture := fixture.
			Builder().
			WithAPI(constants.ApiWithNotificationsAndGroups).
			WithNotification(constants.NotificationWithGroupFile).
			Build()
		apiFixture.API.Spec.Context = &refs.NamespacedName{
			Name:      groupFixture.Context.Name,
			Namespace: constants.Namespace,
		}
		apiFixture.API.Spec.GroupRefs = []refs.NamespacedName{
			{
				Name:      groupFixture.Group.Name,
				Namespace: constants.Namespace,
			},
		}
		apiFixture.API.Spec.NotificationsRefs = []refs.NamespacedName{
			{
				Name:      apiFixture.Notification.Name,
				Namespace: constants.Namespace,
			},
		}
		apiFixture.Notification.Spec.Console.GroupRefs = []refs.NamespacedName{
			{
				Name:      groupFixture.Group.Name,
				Namespace: constants.Namespace,
			},
		}
		apiFixture = apiFixture.Apply()

		client := apim.NewClient(ctx)
		var console *base.ConsoleNotificationConfiguration
		Eventually(func() error {
			c, err := client.Notification.GetConsoleNotificationConfiguration(apiFixture.API.GetID())
			console = c
			return err
		}, constants.EventualTimeout, constants.Interval).Should(Succeed())

		Expect(console).NotTo(BeNil())
		Expect(console.Origin).To(Equal("KUBERNETES"))
		Expect(console.ReferenceID).To(Equal(apiFixture.API.Status.ID))
		Expect(console.ReferenceType).To(Equal("API"))
		Expect(console.Groups).To(ContainElements(groupFixture.Group.Spec.Name))
		Expect(console.Hooks).To(ContainElements("API_STARTED", "API_STOPPED"))

	})

})
