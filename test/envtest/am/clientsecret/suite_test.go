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
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	amsdk "github.com/gravitee-io/gravitee-automation-sdk/am-sdk/v2/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/envtest"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"

	"github.com/onsi/gomega/gexec"
	corev1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	k8sEnvtest "sigs.k8s.io/controller-runtime/pkg/envtest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	value      = "jFXuh6wb1_dhEu5h8Yi4Dp0fC3ZboI0I-SDoaGE22eI"
	otherValue = "Zq3v8Kp0Lr2Xw9Ty4Ub7Nc1Md6Fe5Gh0Ji3Ko8Lp2Qs"
)

var (
	mockServer *am.MockServer
	testEnv    *k8sEnvtest.Environment
)

func TestResources(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AM Client Secret Suite")
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

// fixtures builds an AMContext, a domain with its system certificate, the minimal application and a client
// secret of it. The caller creates the Secret holding the value, then applies them.
func fixtures() *fixture.Objects {
	return fixture.Builder().
		AddSecret(constants.AMContextSecretFile).
		WithAMContext(constants.AMContextFile).
		WithAMSecurityDomain(constants.AMSecurityDomainBasicFile).
		WithAMCertificate(constants.AMCertificateSystemFile).
		WithAMApplication(constants.AMApplicationMinimalFile).
		WithAMClientSecret(constants.AMClientSecretFile).
		Build()
}

// createSecret creates the Secret the client secret references, holding data.
func createSecret(ctx context.Context, cs *v1alpha1.AMClientSecret, data map[string]string) *corev1.Secret {
	GinkgoHelper()
	secret := &corev1.Secret{
		ObjectMeta: metaV1.ObjectMeta{Name: cs.Spec.SecretRef.Name, Namespace: cs.Namespace},
		StringData: data,
	}
	Expect(manager.Client().Create(ctx, secret)).To(Succeed())
	return secret
}

// withValue is the data of a Secret holding v under the key the client secret reads.
func withValue(cs *v1alpha1.AMClientSecret, v string) map[string]string {
	return map[string]string{cs.Spec.SecretRef.Key: v}
}

// remoteSecret fetches the client secret from the mock AM, under its application.
func remoteSecret(ctx context.Context, app *v1alpha1.AMApplication, cs *v1alpha1.AMClientSecret) int {
	GinkgoHelper()
	key := refs.NewNamespacedNameFromObject(cs)
	resp, err := am.NewSDKClient().GetClientSecretWithResponse(ctx, app.Status.DomainKey, app.Status.Key, key.HRID())
	Expect(err).ToNot(HaveOccurred())
	return resp.StatusCode()
}

// remoteSecrets lists the client secrets of the application in the mock AM.
func remoteSecrets(ctx context.Context, app *v1alpha1.AMApplication) []amsdk.ClientSecret {
	GinkgoHelper()
	resp, err := am.NewSDKClient().ListClientSecretsWithResponse(ctx, app.Status.DomainKey, app.Status.Key)
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode()).To(Equal(http.StatusOK))
	return *resp.JSON200
}

// expectAccepted waits until the client secret is accepted.
func expectAccepted(ctx context.Context, cs *v1alpha1.AMClientSecret) {
	GinkgoHelper()
	Eventually(ctx, func() error {
		if err := manager.GetLatest(ctx, cs); err != nil {
			return err
		}
		return assert.AMClientSecretAccepted(cs)
	}, constants.EventualTimeout, constants.Interval).Should(Succeed(), cs.Name)
}

// expectCondition waits until the client secret reports condition False with a message containing part.
func expectCondition(ctx context.Context, cs *v1alpha1.AMClientSecret, condition, part string) {
	GinkgoHelper()
	Eventually(ctx, func() error {
		if err := manager.GetLatest(ctx, cs); err != nil {
			return err
		}
		c := k8s.MapConditions(cs.Status.Conditions)[condition]
		if c.Status != metaV1.ConditionFalse || !strings.Contains(c.Message, part) {
			return fmt.Errorf("%s %s: %s", condition, c.Status, c.Message)
		}
		return nil
	}, constants.EventualTimeout, constants.Interval).Should(Succeed(), cs.Name)
}

// deleteAndWait deletes obj and waits until it is gone from the cluster.
func deleteAndWait(ctx context.Context, kind string, obj client.Object) {
	GinkgoHelper()
	cp, ok := obj.DeepCopyObject().(client.Object)
	Expect(ok).To(BeTrue())
	Expect(manager.Client().Delete(ctx, cp)).To(Succeed())
	Eventually(ctx, func() error {
		return assert.Deleted(ctx, kind, obj)
	}, constants.EventualTimeout, constants.Interval).Should(Succeed(), obj.GetName())
}
