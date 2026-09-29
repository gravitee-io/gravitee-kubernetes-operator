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

package amidentityprovider

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

	It("should update the identity provider in AM", func() {
		idp := fixtures().Apply().AMIdentityProvider

		By("renaming the identity provider")

		Expect(manager.GetLatest(ctx, idp)).To(Succeed())
		idp.Spec.Name = new("Renamed inline users")
		Expect(manager.Client().Update(ctx, idp)).To(Succeed())

		By("expecting the new name in AM")

		Eventually(ctx, func() error {
			status, remote := remoteIdentityProvider(ctx, idp)
			if status != http.StatusOK {
				return fmt.Errorf("status %d", status)
			}
			if remote.Name == nil || *remote.Name != "Renamed inline users" {
				return fmt.Errorf("name %v", remote.Name)
			}
			return nil
		}, constants.EventualTimeout, constants.Interval).Should(Succeed(), idp.Name)
	})
})
