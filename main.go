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

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"sigs.k8s.io/controller-runtime/pkg/manager"

	v1 "k8s.io/api/networking/v1"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/indexer"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/logging"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/watch"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/apidefinition"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/apiresource"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/ingress"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/managementcontext"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

const managerPort = 9443

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(gio.AddToScheme(scheme))
	utilruntime.Must(gio.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	opts := zap.Options{
		Development:          env.Config.Development,
		EncoderConfigOptions: logging.NewEncoderConfigOption(),
	}

	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   managerPort,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "24d975d3.gravitee.io",
		Namespace:              env.Config.NS,
	})

	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	err = addIndexer(mgr)
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&apidefinition.Reconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("apidefinition-controller"),
		Watcher:  watch.New(context.Background(), mgr.GetClient(), &gio.ApiDefinitionList{}),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ApiDefinition")
		os.Exit(1)
	}

	if err = (&managementcontext.Reconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("managementcontext-controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ManagementContext")
		os.Exit(1)
	}
	if err = (&ingress.Reconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("ingress-controller"),
		Watcher:  watch.New(context.Background(), mgr.GetClient(), &v1.IngressList{}),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Ingress")
		os.Exit(1)
	}
	if err = (&apiresource.Reconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ApiResource")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if healthCheckErr := mgr.AddHealthzCheck("healthz", healthz.Ping); healthCheckErr != nil {
		setupLog.Error(healthCheckErr, "unable to set up health check")
		os.Exit(1)
	}
	if readyCheckErr := mgr.AddReadyzCheck("readyz", healthz.Ping); readyCheckErr != nil {
		setupLog.Error(readyCheckErr, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if startErr := mgr.Start(ctrl.SetupSignalHandler()); startErr != nil {
		setupLog.Error(startErr, "problem running manager")
		os.Exit(1)
	}
}

func addIndexer(mgr manager.Manager) error {
	err := indexApiDefinitionFields(mgr)
	if err != nil {
		return fmt.Errorf("unable to start manager (Indexing fields in API definition)")
	}

	err = indexIngressFields(mgr)
	if err != nil {
		return fmt.Errorf("unable to start manager (Indexing fields in ingress resources)")
	}

	err = indexTLSSecretFields(mgr)
	if err != nil {
		return fmt.Errorf("unable to start manager (Indexing fields in ingress resources)")
	}

	return nil
}

func indexApiDefinitionFields(manager ctrl.Manager) error {
	cache := manager.GetCache()
	ctx := context.Background()

	contextIndexer := indexer.NewIndexer(indexer.ContextField, indexer.IndexManagementContexts)
	err := cache.IndexField(ctx, &gio.ApiDefinition{}, contextIndexer.Field, contextIndexer.Func)
	if err != nil {
		return err
	}

	resourceIndexer := indexer.NewIndexer(indexer.ResourceField, indexer.IndexApiResourceRefs)
	err = cache.IndexField(ctx, &gio.ApiDefinition{}, resourceIndexer.Field, resourceIndexer.Func)
	if err != nil {
		return err
	}

	return nil
}

func indexIngressFields(manager ctrl.Manager) error {
	cache := manager.GetCache()
	ctx := context.Background()

	apiTemplateIndexer := indexer.NewIndexer(indexer.ApiTemplateField, indexer.IndexApiTemplate)
	err := cache.IndexField(ctx, &v1.Ingress{}, apiTemplateIndexer.Field, apiTemplateIndexer.Func)
	if err != nil {
		return err
	}

	return nil
}

func indexTLSSecretFields(manager ctrl.Manager) error {
	cache := manager.GetCache()
	ctx := context.Background()

	tlsSecretIndexer := indexer.NewIndexer(indexer.TLSSecretField, indexer.IndexTLSSecret)
	err := cache.IndexField(ctx, &v1.Ingress{}, tlsSecretIndexer.Field, tlsSecretIndexer.Func)
	if err != nil {
		return err
	}

	return nil
}
