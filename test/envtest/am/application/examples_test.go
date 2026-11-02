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
package amapplication

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	internal "github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
)

var _ = Describe("Examples", func() {
	ctx := context.Background()

	DescribeTable("should create the example application in AM",
		func(file string) {
			app := fixtures(file).Apply().AMApplication

			status, remote := remoteApplication(ctx, app)
			Expect(status).To(Equal(http.StatusOK))

			dto, err := internal.ToApplicationDTO(app)
			Expect(err).ToNot(HaveOccurred())
			Expect(remote.Name).To(Equal(dto.Name))
			Expect(remote.Type).To(Equal(dto.Type))
			Expect(remote.Certificate).To(Equal(dto.Certificate))
			Expect(remote.IdentityProviders).To(Equal(dto.IdentityProviders))
			Expect(remote.ClientId).ToNot(BeNil())
			Expect(app.Status.ClientID).To(Equal(*remote.ClientId))
			if app.Spec.ClientId != nil {
				Expect(app.Status.ClientID).To(Equal(*app.Spec.ClientId))
			}
		},
		Entry("minimal", constants.AMApplicationMinimalFile),
		Entry("full", constants.AMApplicationFullFile),
		Entry("keys", constants.AMApplicationKeysFile),
		Entry("client id", constants.AMApplicationClientIDFile),
		Entry("mfa", constants.AMApplicationMFAFile),
	)
})
