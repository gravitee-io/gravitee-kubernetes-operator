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
	"testing"
	"time"

	amsdk "github.com/gravitee-io-labs/gravitee-automation-tools/am-sdk/v2/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	internal "github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/identityprovider"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/envtest"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"

	"github.com/onsi/gomega/gexec"
	"k8s.io/client-go/rest"
	k8sEnvtest "sigs.k8s.io/controller-runtime/pkg/envtest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	mockServer *am.MockServer
	testEnv    *k8sEnvtest.Environment
)

func TestResources(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AM IdentityProvider Suite")
}

var _ = SynchronizedBeforeSuite(func() {
	// NOSONAR mandatory noop
}, func() {
	var cfg *rest.Config
	testEnv, cfg = envtest.Start()
	manager.UseConfig(cfg)
	mockServer = am.Start(false)
})

var _ = SynchronizedAfterSuite(func() {
	By("Tearing down the test environment")
	if mockServer != nil {
		mockServer.Stop()
	}
	manager.Stop()
	if testEnv != nil {
		envtest.Stop(testEnv)
	}
	gexec.KillAndWait(5 * time.Second)
}, func() {
	// NOSONAR mandatory noop
})

// remoteIdentityProvider fetches the identity provider from the mock AM, under its domain.
func remoteIdentityProvider(ctx context.Context, idp *v1alpha1.AMIdentityProvider) (int, *amsdk.IdentityProvider) {
	dto, err := internal.ToIdentityProviderDTO(idp)
	Expect(err).ToNot(HaveOccurred())
	resp, err := am.NewSDKClient().GetIdentityProviderWithResponse(ctx, dto.DomainKey, dto.Key)
	Expect(err).ToNot(HaveOccurred())
	return resp.StatusCode(), resp.JSON200
}

// fixtures builds an AMContext, a domain and an identity provider under it, applied by the caller.
func fixtures() *fixture.Objects {
	return fixture.Builder().
		AddSecret(constants.AMContextSecretFile).
		WithAMContext(constants.AMContextFile).
		WithAMSecurityDomain(constants.AMSecurityDomainBasicFile).
		WithAMIdentityProvider(constants.AMIdentityProviderFile).
		Build()
}
