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

package apidefinition

import (
	"context"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	ginkgotypes "github.com/onsi/ginkgo/v2/types"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	k8sErr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/onsi/gomega/gexec"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var k8sClient client.Client
var k8sManager ctrl.Manager
var ctx = context.Background()

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

var _ = BeforeSuite(func() {
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

	err = (&Reconciler{
		Client:   k8sManager.GetClient(),
		Scheme:   k8sManager.GetScheme(),
		Recorder: k8sManager.GetEventRecorderFor("apidefinition_controller"),
	}).SetupWithManager(k8sManager)

	Expect(err).ToNot(HaveOccurred())

	go func() {
		err = k8sManager.Start(ctrl.SetupSignalHandler())
		Expect(err).ToNot(HaveOccurred())
	}()

	k8sClient = k8sManager.GetClient()

	Expect(k8sClient).ToNot(BeNil())

})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	gexec.KillAndWait(5 * time.Second)
	// err := testEnv.Stop()
	// Expect(err).ToNot(HaveOccurred())
})

var _ = ReportAfterEach(func(specReport ginkgotypes.SpecReport) {
	enableDSmokeExpect := true
	for _, label := range specReport.Labels() {
		if label == "DisableSmokeExpect" {
			enableDSmokeExpect = false
		}
	}

	if enableDSmokeExpect {
		// Smoke test to check there was no unwanted error in the operator's logs. masked by a reconcile for example
		Expect(
			strings.Contains(specReport.CapturedGinkgoWriterOutput, "\tERROR\t"),
		).To(
			BeFalse(), "[Smoke Test] There are errors in the operator logs",
		)
	}
})

func cleanupApiDefinitionAndManagementContext(
	apiDefinition *gio.ApiDefinition,
	managementContext *gio.ManagementContext,
) {
	cleanupApiDefinition(apiDefinition)

	contextLookupKey := types.NamespacedName{Name: managementContext.Name, Namespace: managementContext.Namespace}

	err := k8sClient.Delete(ctx, managementContext)
	if !k8sErr.IsNotFound(err) {
		// wait deleted only if not already deleted
		Eventually(func() error {
			return k8sClient.Get(ctx, contextLookupKey, managementContext)
		}, timeout, interval).ShouldNot(Succeed())
	}
}

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

func cleanupApiDefinition(apiDefinition *gio.ApiDefinition) {
	apiLookupKey := types.NamespacedName{Name: apiDefinition.Name, Namespace: apiDefinition.Namespace}

	err := k8sClient.Delete(ctx, apiDefinition)
	if !k8sErr.IsNotFound(err) {
		// wait deleted only if not already deleted
		Eventually(func() error {
			return k8sClient.Get(ctx, apiLookupKey, apiDefinition)
		}, timeout, interval).ShouldNot(Succeed())
	}
}

func getEventsReason(apiDefinition *gio.ApiDefinition) []string {
	eventsReason := []string{}

	events := &v1.EventList{}

	err := k8sClient.List(
		ctx,
		events,
		&client.ListOptions{Namespace: apiDefinition.GetNamespace()},
		client.MatchingFields{"involvedObject.name": apiDefinition.GetName()},
	)
	Expect(err).ToNot(HaveOccurred())

	for _, event := range events.Items {
		eventsReason = append(eventsReason, event.Reason)
	}
	return eventsReason
}
