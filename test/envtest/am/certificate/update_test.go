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
package amcertificate

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

	It("should update the certificate in AM", func() {
		cert := fixtures().Apply().AMCertificate

		By("changing the key alias")

		Expect(manager.GetLatest(ctx, cert)).To(Succeed())
		cert.Spec.Configuration.Put("alias", "rotated")
		Expect(manager.Client().Update(ctx, cert)).To(Succeed())

		By("expecting the new alias in AM")

		Eventually(ctx, func() error {
			status, remote := remoteCertificate(ctx, cert)
			if status != http.StatusOK {
				return fmt.Errorf("status %d", status)
			}
			if alias := remote.Configuration.Object["alias"]; alias != "rotated" {
				return fmt.Errorf("alias %v", alias)
			}
			return nil
		}, constants.EventualTimeout, constants.Interval).Should(Succeed(), cert.Name)
	})
})
