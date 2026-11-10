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

package apimgroup

import (
	"context"
	"testing"
	"time"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/envtest"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"

	"github.com/onsi/gomega/gexec"
	"k8s.io/client-go/rest"
	k8sEnvtest "sigs.k8s.io/controller-runtime/pkg/envtest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	mockServer *apim.MockServer
	testEnv    *k8sEnvtest.Environment
)

func TestResources(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "APIM Group Suite")
}

var _ = SynchronizedBeforeSuite(func() {
	// NOSONAR mandatory noop
}, func() {
	var cfg *rest.Config
	testEnv, cfg = envtest.Start()
	manager.UseConfig(cfg)
	drift.Init()
	mockServer = apim.StartMock(false)
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

// groupFixtures builds a management context on the APIM mock and the group of file under it.
func groupFixtures(file string) *fixture.Objects {
	return fixture.Builder().
		AddSecret(constants.MockContextSecretFile).
		WithContext(constants.MockContextFile).
		WithGroup(file).
		Build()
}

// remoteGroup fetches the group from the APIM mock by the HRID the operator gives it.
func remoteGroup(ctx context.Context, grp *v1alpha1.Group) (int, *sdk.GroupState) {
	resp, err := apim.NewMockSDKClient().GetGroupWithResponse(ctx, refs.NewNamespacedNameFromObject(grp).HRID(), nil)
	Expect(err).ToNot(HaveOccurred())
	return resp.StatusCode(), resp.JSON200
}
