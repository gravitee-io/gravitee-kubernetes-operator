// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package amreporter

import (
	"context"
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Update", func() {
	ctx := context.Background()

	It("should update the reporter in AM", func() {
		reporter := fixtures().Apply().AMReporter

		By("changing the retention and disabling the reporter")

		Expect(manager.GetLatest(ctx, reporter)).To(Succeed())
		reporter.Spec.Configuration.Put("retainDays", 30)
		reporter.Spec.Enabled = new(false)
		Expect(manager.Client().Update(ctx, reporter)).To(Succeed())

		By("expecting the new retention and the reporter disabled in AM")

		Eventually(ctx, func() error {
			status, remote := remoteReporter(ctx, reporter)
			if status != http.StatusOK {
				return fmt.Errorf("status %d", status)
			}
			if retainDays := remote.Configuration.Object["retainDays"]; retainDays != float64(30) {
				return fmt.Errorf("retainDays %v", retainDays)
			}
			if remote.Enabled == nil || *remote.Enabled {
				return fmt.Errorf("enabled %v", remote.Enabled)
			}
			return nil
		}, constants.EventualTimeout, constants.Interval).Should(Succeed(), reporter.Name)
	})
})
