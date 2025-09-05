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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Create", labels.WithContext, func() {

	timeout := constants.EventualTimeout
	interval := constants.Interval
	ctx := context.Background()

	When("creating a notification with a resolved group", func() {
		It("should be accepted and resolved", func() {
			fixtures := fixture.Builder().
				WithContext(constants.ContextWithCredentialsFile).
				WithGroup(constants.GroupFile).
				WithNotification(constants.NotificationWithGroupFile).
				Build().
				Apply()

			Eventually(func() error {
				if err := manager.GetLatest(ctx, fixtures.Notification); err != nil {
					return err
				}

				if err := assert.HasFinalizer(fixtures.Notification, core.NotificationFinalizer); err != nil {
					return err
				}

				if err := assert.IsAccepted(fixtures.Notification); err != nil {
					return err
				}

				return assert.IsResolved(fixtures.Notification)
			}, timeout, interval).Should(Succeed())
		})
	})
})
