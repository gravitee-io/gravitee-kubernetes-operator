/*
Copyright 2022 DAVID BRASSELY.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package test

import (
	"context"
	"testing"
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/ingress"

	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	netv1 "k8s.io/api/networking/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/onsi/gomega/gexec"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/apidefinition"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/apiresource"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/managementcontext"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/indexer"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var k8sClient client.Client
var k8sManager ctrl.Manager
var ctx context.Context

// Define utility constants for object names and testing timeouts/durations and intervals.
const (
	metricsAddr = ":10000"
	probeAddr   = ":10001"
	managerPort = 10002

	namespace = "default"
	timeout   = time.Second * 30
	interval  = time.Millisecond * 250
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "GKO Controllers Suite")
}

var _ = SynchronizedBeforeSuite(func() {
	By("Setting up the test environment")

	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	err := gio.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sManager, err = ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme.Scheme,
		Port:                   managerPort,
		MetricsBindAddress:     metricsAddr,
		HealthProbeBindAddress: probeAddr,
	})

	Expect(err).ToNot(HaveOccurred())

	err = addEventIndexes()

	Expect(err).ToNot(HaveOccurred())

	err = (&apidefinition.Reconciler{
		Client:   k8sManager.GetClient(),
		Scheme:   k8sManager.GetScheme(),
		Recorder: k8sManager.GetEventRecorderFor("apidefinition_controller"),
	}).SetupWithManager(k8sManager)

	Expect(err).ToNot(HaveOccurred())

	err = (&managementcontext.Reconciler{
		Client: k8sManager.GetClient(),
		Scheme: k8sManager.GetScheme(),
	}).SetupWithManager(k8sManager)

	Expect(err).ToNot(HaveOccurred())

	err = (&apiresource.Reconciler{
		Client: k8sManager.GetClient(),
		Scheme: k8sManager.GetScheme(),
	}).SetupWithManager(k8sManager)

	Expect(err).ToNot(HaveOccurred())

	err = (&ingress.Reconciler{
		Client:   k8sManager.GetClient(),
		Scheme:   k8sManager.GetScheme(),
		Recorder: k8sManager.GetEventRecorderFor("ingress-controller"),
	}).SetupWithManager(k8sManager)

	Expect(err).ToNot(HaveOccurred())

	cache := k8sManager.GetCache()

	contextIndexer := indexer.NewIndexer(indexer.ContextField, indexer.IndexManagementContexts)
	err = cache.IndexField(ctx, &gio.ApiDefinition{}, contextIndexer.Field, contextIndexer.Func)
	Expect(err).ToNot(HaveOccurred())

	resourceIndexer := indexer.NewIndexer(indexer.ResourceField, indexer.IndexApiResourceRefs)
	err = cache.IndexField(ctx, &gio.ApiDefinition{}, resourceIndexer.Field, resourceIndexer.Func)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		err = k8sManager.Start(ctrl.SetupSignalHandler())
		Expect(err).ToNot(HaveOccurred())
	}()

}, func() {
	//+kubebuilder:scaffold:scheme
	err := gio.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())
	cli, err := client.New(ctrl.GetConfigOrDie(), client.Options{Scheme: scheme.Scheme})
	Expect(err).ToNot(HaveOccurred())
	k8sClient = cli
	ctx = context.Background()
})

var _ = SynchronizedAfterSuite(func() {
	By("Tearing down the test environment")

	Expect(k8sClient.DeleteAllOf(
		ctx,
		&netv1.Ingress{},
		client.InNamespace(namespace),
		client.MatchingLabels{keys.IngressLabel: keys.IngressLabelValue}),
	).To(Succeed())
	Expect(k8sClient.DeleteAllOf(ctx, &gio.ApiDefinition{}, client.InNamespace(namespace))).To(Succeed())
	Expect(k8sClient.DeleteAllOf(ctx, &gio.ManagementContext{}, client.InNamespace(namespace))).To(Succeed())
	Expect(k8sClient.DeleteAllOf(ctx, &gio.ApiResource{}, client.InNamespace(namespace))).To(Succeed())
	gexec.KillAndWait(5 * time.Second)
}, func() {

})

// Add filed indexes for event to be able to filter on it.
func addEventIndexes() error {
	err := k8sManager.GetFieldIndexer().IndexField(
		ctx,
		&v1.Event{},
		"involvedObject.name",
		func(rawObj client.Object) []string {
			event, _ := rawObj.(*v1.Event)
			return []string{event.InvolvedObject.Name}
		},
	)
	return err
}

func getEventsReason(namespace string, name string) []string {
	eventsReason := []string{}

	events := &v1.EventList{}

	err := k8sClient.List(
		ctx,
		events,
		&client.ListOptions{Namespace: namespace},
		client.MatchingFields{"involvedObject.name": name},
	)
	Expect(err).ToNot(HaveOccurred())

	for _, event := range events.Items {
		eventsReason = append(eventsReason, event.Reason)
	}
	return eventsReason
}
