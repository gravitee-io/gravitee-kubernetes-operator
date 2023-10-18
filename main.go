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
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"os"

	"gopkg.in/yaml.v3"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/secrets"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/extensions"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"

	v1 "k8s.io/api/networking/v1"
	apiExtensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/indexer"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/logging"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/watch"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.

	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	utilRuntime "k8s.io/apimachinery/pkg/util/runtime"
	clientScheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsServer "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	gioV1Alpha1 "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	gioV1Beta1 "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1beta1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/apidefinition"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/apiresource"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/ingress"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/managementcontext"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")

	//go:embed helm
	helm embed.FS
)

func init() {
	utilRuntime.Must(clientScheme.AddToScheme(scheme))
	utilRuntime.Must(gioV1Alpha1.AddToScheme(scheme))
	utilRuntime.Must(gioV1Beta1.AddToScheme(scheme))
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
		Development:          env.Config.DisableJSONLogs,
		EncoderConfigOptions: logging.NewEncoderConfigOption(),
	}

	opts.BindFlags(flag.CommandLine)

	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	if !env.Config.EnableMetrics {
		setupLog.Info("metrics are disabled")
		metricsAddr = "0" // disables metrics
	}

	metrics := metricsServer.Options{BindAddress: metricsAddr}

	if env.Config.InsecureSkipVerify {
		setupLog.Info("TLS certificates verification is skipped for APIM HTTP client")
	}

	var webhookServer webhook.Server
	if env.Config.EnableWebhook {
		if err := patchCRDs(); err != nil {
			setupLog.Error(err, "unable to apply custom resource definitions")
			os.Exit(1)
		}
		webhookServer = webhook.NewServer(webhook.Options{
			Port: env.Config.WebhookPort,
		})
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		Metrics:                metrics,
		WebhookServer:          webhookServer,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "24d975d3.gravitee.io",
		Cache:                  buildCacheOptions(env.Config.WatchNS),
	})

	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = addIndexer(mgr); err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	registerControllers(mgr)

	if env.Config.EnableWebhook {
		if err = (&gioV1Beta1.ApiDefinition{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "ApiDefinition")
			os.Exit(1)
		}
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

	k8s.RegisterClient(mgr.GetClient())

	setupLog.Info("starting manager")
	if startErr := mgr.Start(ctrl.SetupSignalHandler()); startErr != nil {
		setupLog.Error(startErr, "problem running manager")
		os.Exit(1)
	}
}

func buildCacheOptions(ns string) cache.Options {
	if ns == "" {
		setupLog.Info("no namespace defined, watching all namespaces")
		return cache.Options{}
	}

	setupLog.Info("watching on dedicated namespace", "namespace", ns)
	return cache.Options{
		DefaultNamespaces: map[string]cache.Config{
			ns: {},
		},
	}
}

func registerControllers(mgr manager.Manager) {
	if err := (&apidefinition.Reconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("apidefinition-controller"),
		Watcher:  watch.New(context.Background(), mgr.GetClient(), &gioV1Alpha1.ApiDefinitionList{}),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ApiDefinition")
		os.Exit(1)
	}

	if err := (&managementcontext.Reconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("managementcontext-controller"),
		Watcher:  watch.New(context.Background(), mgr.GetClient(), &gioV1Alpha1.ManagementContextList{}),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ManagementContext")
		os.Exit(1)
	}
	if err := (&ingress.Reconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("ingress-controller"),
		Watcher:  watch.New(context.Background(), mgr.GetClient(), &v1.IngressList{}),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Ingress")
		os.Exit(1)
	}
	if err := (&apiresource.Reconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("apiresource-controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ApiResource")
		os.Exit(1)
	}
	if err := (&application.Reconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("application-controller"),
		Watcher:  watch.New(context.Background(), mgr.GetClient(), &gioV1Alpha1.ApplicationList{}),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Application")
		os.Exit(1)
	}

	if err := (&secrets.Reconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Secret")
		os.Exit(1)
	}
}

func addIndexer(mgr manager.Manager) error {
	err := indexApiDefinitionFields(mgr)
	if err != nil {
		return fmt.Errorf("unable to start manager (Indexing fields in API definition)")
	}

	err = indexSecretRefs(mgr)
	if err != nil {
		return fmt.Errorf("unable to start manager (Indexing fields in context resources)")
	}

	err = indexIngressFields(mgr)
	if err != nil {
		return fmt.Errorf("unable to start manager (Indexing fields in ingress resources)")
	}

	err = indexTLSSecretFields(mgr)
	if err != nil {
		return fmt.Errorf("unable to start manager (Indexing fields in ingress resources)")
	}

	err = indexApplicationFields(mgr)
	if err != nil {
		return fmt.Errorf("unable to start manager (Indexing fields in application resources)")
	}

	return nil
}

func indexApiDefinitionFields(manager ctrl.Manager) error {
	cache := manager.GetCache()
	ctx := context.Background()

	contextIndexer := indexer.NewIndexer(indexer.ContextField, indexer.IndexManagementContexts)
	err := cache.IndexField(ctx, &gioV1Alpha1.ApiDefinition{}, contextIndexer.Field, contextIndexer.Func)
	if err != nil {
		return err
	}

	resourceIndexer := indexer.NewIndexer(indexer.ResourceField, indexer.IndexApiResourceRefs)
	err = cache.IndexField(ctx, &gioV1Alpha1.ApiDefinition{}, resourceIndexer.Field, resourceIndexer.Func)
	if err != nil {
		return err
	}

	return nil
}

func indexSecretRefs(manager ctrl.Manager) error {
	cache := manager.GetCache()
	ctx := context.Background()

	secretRefIndexer := indexer.NewIndexer(indexer.SecretRefField, indexer.IndexManagementContextSecrets)
	return cache.IndexField(ctx, &gioV1Alpha1.ManagementContext{}, secretRefIndexer.Field, secretRefIndexer.Func)
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

func indexApplicationFields(manager ctrl.Manager) error {
	cache := manager.GetCache()
	ctx := context.Background()

	appContextIndexer := indexer.NewIndexer(indexer.AppContextField, indexer.IndexApplicationManagementContexts)
	err := cache.IndexField(ctx, &gioV1Alpha1.Application{}, appContextIndexer.Field, appContextIndexer.Func)
	if err != nil {
		return err
	}

	return nil
}

func patchCRDs() error {
	ctx := context.Background()

	setupLog.Info("extending scheme for CRD patch")
	sErr := apiExtensions.AddToScheme(scheme)
	if sErr != nil {
		return sErr
	}

	cli, cErr := client.New(ctrl.GetConfigOrDie(), client.Options{Scheme: scheme})
	if cErr != nil {
		return cErr
	}

	return fs.WalkDir(helm, keys.CRDBase, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if d.IsDir() {
			return nil
		}

		b, err := fs.ReadFile(helm, path)
		if err != nil {
			return err
		}

		obj := make(map[string]interface{})
		if err = yaml.Unmarshal(b, &obj); err != nil {
			return err
		}

		unstructured := &unstructured.Unstructured{Object: obj}
		if err = yaml.Unmarshal(b, unstructured); err != nil {
			return err
		}

		setupLog.Info("patching custom resource definition", "name", unstructured.GetName())

		existing := &apiExtensions.CustomResourceDefinition{}
		if err = cli.Get(ctx, types.NamespacedName{Name: unstructured.GetName()}, existing); err != nil {
			return err
		}

		desired := &apiExtensions.CustomResourceDefinition{}
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured.Object, desired); err != nil {
			return err
		}

		extensions.ExtendWithWebhook(desired)
		extensions.InjectCA(desired)

		desired.Spec.DeepCopyInto(&existing.Spec)
		existing.ObjectMeta.Annotations = desired.GetAnnotations()

		if err = cli.Update(ctx, existing); err != nil {
			return err
		}

		return err
	})
}
