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

	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/application"

	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"

	netv1 "k8s.io/api/networking/v1"

	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/ingress"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/indexer"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/watch"
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
		Watcher:  watch.New(context.Background(), k8sManager.GetClient(), &gio.ApiDefinitionList{}),
	}).SetupWithManager(k8sManager)

	Expect(err).ToNot(HaveOccurred())

	err = (&managementcontext.Reconciler{
		Client:   k8sManager.GetClient(),
		Scheme:   k8sManager.GetScheme(),
		Recorder: k8sManager.GetEventRecorderFor("managementcontext_controller"),
		Watcher:  watch.New(context.Background(), k8sManager.GetClient(), &gio.ManagementContextList{}),
	}).SetupWithManager(k8sManager)

	Expect(err).ToNot(HaveOccurred())

	err = (&apiresource.Reconciler{
		Client:   k8sManager.GetClient(),
		Scheme:   k8sManager.GetScheme(),
		Recorder: k8sManager.GetEventRecorderFor("apiresource-controller"),
	}).SetupWithManager(k8sManager)

	Expect(err).ToNot(HaveOccurred())

	err = (&ingress.Reconciler{
		Client:   k8sManager.GetClient(),
		Scheme:   k8sManager.GetScheme(),
		Recorder: k8sManager.GetEventRecorderFor("ingress-controller"),
		Watcher:  watch.New(context.Background(), k8sManager.GetClient(), &netv1.IngressList{}),
	}).SetupWithManager(k8sManager)

	Expect(err).ToNot(HaveOccurred())

	err = (&application.Reconciler{
		Client:   k8sManager.GetClient(),
		Scheme:   k8sManager.GetScheme(),
		Recorder: k8sManager.GetEventRecorderFor("application-controller"),
		Watcher:  watch.New(context.Background(), k8sManager.GetClient(), &gio.ApplicationList{}),
	}).SetupWithManager(k8sManager)

	Expect(err).ToNot(HaveOccurred())

	cache := k8sManager.GetCache()

	contextIndexer := indexer.NewIndexer(indexer.ContextField, indexer.IndexManagementContexts)
	err = cache.IndexField(ctx, &gio.ApiDefinition{}, contextIndexer.Field, contextIndexer.Func)
	Expect(err).ToNot(HaveOccurred())

	contextSecretsIndexer := indexer.NewIndexer(indexer.SecretRefField, indexer.IndexManagementContextSecrets)
	err = cache.IndexField(ctx, &gio.ManagementContext{}, contextSecretsIndexer.Field, contextSecretsIndexer.Func)
	Expect(err).ToNot(HaveOccurred())

	resourceIndexer := indexer.NewIndexer(indexer.ResourceField, indexer.IndexApiResourceRefs)
	err = cache.IndexField(ctx, &gio.ApiDefinition{}, resourceIndexer.Field, resourceIndexer.Func)
	Expect(err).ToNot(HaveOccurred())

	apiTemplateIndexer := indexer.NewIndexer(indexer.ApiTemplateField, indexer.IndexApiTemplate)
	err = cache.IndexField(ctx, &netv1.Ingress{}, apiTemplateIndexer.Field, apiTemplateIndexer.Func)
	Expect(err).ToNot(HaveOccurred())

	tlsSecretIndexer := indexer.NewIndexer(indexer.TLSSecretField, indexer.IndexTLSSecret)
	err = cache.IndexField(ctx, &netv1.Ingress{}, tlsSecretIndexer.Field, tlsSecretIndexer.Func)
	Expect(err).ToNot(HaveOccurred())

	// Set initial values for env variables
	env.Config.CMTemplate404NS = namespace
	env.Config.CMTemplate404Name = "template-404"

	appContextIndexer := indexer.NewIndexer(indexer.AppContextField, indexer.IndexApplicationManagementContexts)
	err = cache.IndexField(ctx, &gio.Application{}, appContextIndexer.Field, appContextIndexer.Func)
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

	Expect(k8sClient.Create(ctx, template404())).Should(Succeed())
})

var _ = SynchronizedAfterSuite(func() {
	By("Tearing down the test environment")

	Expect(k8sClient.DeleteAllOf(
		ctx,
		&netv1.Ingress{},
		client.InNamespace(namespace),
		client.MatchingLabels{keys.IngressLabel: keys.IngressLabelValue}),
	).To(Succeed())
	Consistently(k8sClient.DeleteAllOf(
		ctx,
		&gio.ApiDefinition{},
		client.InNamespace(namespace)), timeout/10, 1*time.Second).Should(Succeed())
	Consistently(k8sClient.DeleteAllOf(
		ctx,
		&gio.Application{},
		client.InNamespace(namespace)), timeout/10, 1*time.Second).Should(Succeed())
	Consistently(k8sClient.DeleteAllOf(
		ctx,
		&gio.ManagementContext{},
		client.InNamespace(namespace)), timeout/10, 1*time.Second).Should(Succeed())
	Expect(k8sClient.DeleteAllOf(ctx, &gio.ApiResource{}, client.InNamespace(namespace))).To(Succeed())
	Expect(k8sClient.Delete(ctx, template404())).Should(Succeed())
	gexec.KillAndWait(5 * time.Second)
}, func() {
	// NOSONAR ignore this noop func
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

func template404() *v1.ConfigMap {
	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "template-404",
			Namespace: namespace,
		},
		Data: map[string]string{
			"content":     `{ "message": "not-found-test" }`,
			"contentType": "application/json",
		},
	}
}
