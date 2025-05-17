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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
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

		By("Applying an API with groups and notifications")

		fixtures := fixture.
			Builder().
			WithContext(constants.ContextWithCredentialsFile).
			WithGroup(constants.GroupFile).
			WithNotification(constants.NotificationWithGroupFile).
			WithAPIv4(constants.ApiV4WithNotificationsAndGroups).
			Build().
			Apply()

		client := apim.NewClient(ctx)

		console := getConsoleNotification(client, fixtures)

		By("Checking that the notification has been added to API")
		Expect(console).ToNot(BeNil())
		Expect(console.Origin).To(Equal("KUBERNETES"))
		Expect(console.ReferenceID).To(Equal(fixtures.APIv4.Status.ID))
		Expect(console.Groups).To(ContainElements(fixtures.Group.Status.ID))
		Expect(console.Hooks).To(ContainElements("API_STARTED", "API_STOPPED"))

	})

})

func getConsoleNotification(client *apim.APIM, fixtures *fixture.Objects) *base.ConsoleNotificationConfiguration {
	url := client.APIs.EnvV1Target("apis").WithPath(fixtures.APIv4.GetID(), "notificationsettings")

	notifications := make([]base.ConsoleNotificationConfiguration, 0)
	Expect(client.APIs.HTTP.Get(url.String(), &notifications)).To(Succeed())

	var console *base.ConsoleNotificationConfiguration
	for _, notification := range notifications {
		if notification.ConfigType == "PORTAL" {
			n := notification
			console = &n
			break
		}
	}
	return console
}
