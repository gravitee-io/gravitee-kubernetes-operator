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
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create", labels.WithContext, func() {
	ctx := context.Background()

	It("should sync the API and leave no console notification events", func() {
		fixtures := fixture.Builder().
			WithContext(constants.ContextWithCredentialsFile).
			WithNotification(constants.NotificationNoGroupFile).
			WithAPI(constants.ApiWithContextFile).
			Build()
		fixtures.Notification.Spec.Console.APIEvents = []notification.ApiEvent{}

		fixtures.Apply()

		Expect(assert.ApiCompleted(fixtures.API)).To(Succeed())

		console, err := apim.NewClient(ctx).Notification.GetConsoleNotificationConfiguration(fixtures.API.GetID())
		Expect(err).ToNot(HaveOccurred())
		if console != nil {
			Expect(console.Hooks).To(BeEmpty())
		}
	})
})
