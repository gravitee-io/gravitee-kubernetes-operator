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

package test

import (
	"context"
	"testing"
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1beta1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/onsi/gomega/gexec"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var k8sClient client.Client
var ctx context.Context

// Define utility constants for object names and testing timeouts/durations and intervals.
const (
	namespace = "default"
	timeout   = time.Second * 10
	interval  = time.Second * 1

	pemRegistryName = "pem-registry"
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "GKO Controllers Suite")
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
	Expect(k8sClient.Create(ctx, pemRegistry())).Should(Succeed())
})

var _ = SynchronizedAfterSuite(func() {
	By("Tearing down the test environment")
	gexec.KillAndWait(5 * time.Second)
}, func() {
	// NOSONAR ignore this noop func
})

func getEventReasons(obj client.Object) func() []string {
	return func() []string {
		eventsReason := []string{}

		events := &v1.EventList{}

		if err := k8sClient.List(
			ctx,
			events,
			&client.ListOptions{Namespace: obj.GetNamespace()},
			client.MatchingFields{"involvedObject.name": obj.GetName()},
		); err != nil {
			return nil
		}

		for _, event := range events.Items {
			eventsReason = append(eventsReason, event.Reason)
		}
		return eventsReason
	}
}

func pemRegistry() *v1.ConfigMap {
	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pemRegistryName,
			Namespace: namespace,
			Labels: map[string]string{
				keys.GraviteeComponentLabel: keys.GraviteePemRegistryLabel,
				keys.IngressClassAnnotation: keys.IngressClassAnnotationValue,
			},
		},
	}
}
