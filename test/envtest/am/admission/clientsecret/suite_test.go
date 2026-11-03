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
	"testing"
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/drift"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/envtest"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	"github.com/onsi/gomega/gexec"
	corev1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	k8sEnvtest "sigs.k8s.io/controller-runtime/pkg/envtest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const value = "jFXuh6wb1_dhEu5h8Yi4Dp0fC3ZboI0I-SDoaGE22eI"

var (
	mockServer *am.MockServer
	testEnv    *k8sEnvtest.Environment
)

func TestResources(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AM Client Secret Admission Suite")
}

var _ = SynchronizedBeforeSuite(func() {
	// NOSONAR mandatory noop
}, func() {
	var cfg *rest.Config
	testEnv, cfg = envtest.Start()
	manager.UseConfig(cfg)
	drift.Init()
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

// withApplication builds an AMContext, a domain with its system certificate, the minimal application and a
// client secret of it; the caller applies them.
func withApplication() *fixture.Objects {
	return fixture.Builder().
		WithAMContext(constants.AMContextFile).
		WithAMSecurityDomain(constants.AMSecurityDomainBasicFile).
		WithAMCertificate(constants.AMCertificateSystemFile).
		WithAMApplication(constants.AMApplicationMinimalFile).
		WithAMClientSecret(constants.AMClientSecretFile).
		Build()
}

// applyAllButClientSecret applies the fixtures except the client secret, and returns them all.
func applyAllButClientSecret() *fixture.Objects {
	fixtures := withApplication()
	cs := fixtures.AMClientSecret
	fixtures.AMClientSecret = nil
	fixtures.Apply()
	fixtures.AMClientSecret = cs
	return fixtures
}

// createSecret creates the Secret the client secret references, holding v, and waits until the manager's
// cache, which admission reads, has it.
func createSecret(ctx context.Context, cs *v1alpha1.AMClientSecret, v string) {
	GinkgoHelper()
	secret := &corev1.Secret{
		ObjectMeta: metaV1.ObjectMeta{Name: cs.Spec.SecretRef.Name, Namespace: cs.Namespace},
		StringData: map[string]string{cs.Spec.SecretRef.Key: v},
	}
	Expect(manager.Client().Create(ctx, secret)).To(Succeed())
	Eventually(ctx, func() error {
		return manager.GetLatest(ctx, secret)
	}, constants.EventualTimeout, constants.Interval).Should(Succeed())
}
