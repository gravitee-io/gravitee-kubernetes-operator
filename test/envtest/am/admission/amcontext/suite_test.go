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

package amcontext

import (
	"testing"
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/envtest"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/am"
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
	mockServer *am.MockServer
	testEnv    *k8sEnvtest.Environment
)

func TestResources(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AMContext Admission Suite")
}

var _ = SynchronizedBeforeSuite(func() {
	// NOSONAR mandatory noop
}, func() {
	var cfg *rest.Config
	testEnv, cfg = envtest.Start()
	manager.UseConfig(cfg)
	mockServer = am.Start(false)
	fixture.Builder().
		AddSecret(constants.AMContextSecretFile).
		Build().
		Apply()
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
