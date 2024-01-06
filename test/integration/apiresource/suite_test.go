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

package apiresource

import (
	"context"
	"testing"
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1beta1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/integration/internal"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/onsi/gomega/gexec"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var k8sClient client.Client
var ctx context.Context

// Define utility constants for object names and testing timeouts/durations and intervals.
const (
	timeout  = time.Second * 10
	interval = time.Second * 1
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "API Resource Controller Suite")
}

var _ = SynchronizedBeforeSuite(func() {
	By("Setting up the test environment")

	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	Expect(v1alpha1.AddToScheme(scheme.Scheme)).NotTo(HaveOccurred())
	Expect(v1beta1.AddToScheme(scheme.Scheme)).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient = internal.ClusterClient()

}, func() {
	//+kubebuilder:scaffold:scheme
	Expect(v1alpha1.AddToScheme(scheme.Scheme)).NotTo(HaveOccurred())
	Expect(v1beta1.AddToScheme(scheme.Scheme)).NotTo(HaveOccurred())

	k8sClient = internal.ClusterClient()
	ctx = context.Background()
})

var _ = SynchronizedAfterSuite(func() {
	By("Tearing down the test environment")
	gexec.KillAndWait(5 * time.Second)
}, func() {
	// NOSONAR ignore this noop func
})
