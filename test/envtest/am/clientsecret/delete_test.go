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

package amclientsecret

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Delete", func() {
	ctx := context.Background()

	It("should delete the secret from AM, also the last one, and keep the Secret (S1, AC9)", func() {
		fixtures := fixtures()
		secret := createSecret(ctx, fixtures.AMClientSecret, withValue(fixtures.AMClientSecret, value))
		fixtures.Apply()
		app, cs := fixtures.AMApplication, fixtures.AMClientSecret
		expectAccepted(ctx, cs)
		Expect(remoteSecret(ctx, app, cs)).To(Equal(http.StatusOK))

		deleteAndWait(ctx, "AMClientSecret", cs)

		Expect(remoteSecret(ctx, app, cs)).To(Equal(http.StatusNotFound))
		Expect(remoteSecrets(ctx, app)).To(BeEmpty())
		Expect(manager.GetLatest(ctx, secret)).To(Succeed())
	})

	It("should delete the secret from AM once its Secret is gone, without reading the value (S4)", func() {
		fixtures := fixtures()
		secret := createSecret(ctx, fixtures.AMClientSecret, withValue(fixtures.AMClientSecret, value))
		fixtures.Apply()
		app, cs := fixtures.AMApplication, fixtures.AMClientSecret
		expectAccepted(ctx, cs)

		deleteAndWait(ctx, "Secret", secret)
		deleteAndWait(ctx, "AMClientSecret", cs)

		Expect(remoteSecret(ctx, app, cs)).To(Equal(http.StatusNotFound))
	})

	It("should release the client secret once its application is deleted, without calling AM (S2)", func() {
		fixtures := fixtures()
		secret := createSecret(ctx, fixtures.AMClientSecret, withValue(fixtures.AMClientSecret, value))
		fixtures.Apply()
		app, cs := fixtures.AMApplication, fixtures.AMClientSecret
		expectAccepted(ctx, cs)

		By("deleting the application: AM deletes it with its secrets")

		deleteAndWait(ctx, "AMApplication", app)
		Expect(remoteSecret(ctx, app, cs)).To(Equal(http.StatusNotFound))

		By("stopping the mock AM: a call to AM would fail and keep the finalizer")

		mockServer.Stop()
		DeferCleanup(func() { mockServer = am.Start(false) })

		By("deleting the client secret by hand (envtest has no garbage collector)")

		deleteAndWait(ctx, "AMClientSecret", cs)
		Expect(manager.GetLatest(ctx, secret)).To(Succeed())
	})
})
