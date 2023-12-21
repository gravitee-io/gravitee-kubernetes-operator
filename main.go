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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/extensions"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"

	v1 "k8s.io/api/networking/v1"
	apiExtensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/indexer"
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
	metricsServer "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1beta1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/apidefinition"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/apiresource"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/ingress"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/managementcontext"
	//+kubebuilder:scaffold:imports
)

var (
	scheme = runtime.NewScheme()

	//go:embed helm
	helm embed.FS
)

func init() {
	utilRuntime.Must(clientScheme.AddToScheme(scheme))
	utilRuntime.Must(v1alpha1.AddToScheme(scheme))
	utilRuntime.Must(v1beta1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	if env.Config.HTTPClientSkipCertVerify {
		log.Global.Warn("TLS certificates verification is skipped for APIM HTTP client")
	}

	log.Global.Debugf("enable metrics: %t", env.Config.EnableMetrics)
	metrics := metricsServer.Options{
		BindAddress:   env.GetMetricsAddr(),
		SecureServing: env.Config.SecureMetrics,
		CertDir:       env.Config.MetricsCertDir,
	}

	probeAddr := fmt.Sprintf(":%d", env.Config.ProbePort)

	log.Global.Debug("setting up webhook server")
	var webhookServer webhook.Server
	if env.Config.EnableWebhook {
		if err := extendResources(); err != nil {
			log.Global.Error(err, "unable to apply custom resource definitions")
			os.Exit(1)
		}
		webhookServer = webhook.NewServer(webhook.Options{
			Port: env.Config.WebhookPort,
		})
	}

	log.Global.Debug("creating manager")
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		Metrics:                metrics,
		WebhookServer:          webhookServer,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         env.Config.EnableLeaderElection,
		LeaderElectionID:       "24d975d3.gravitee.io",
		Cache:                  buildCacheOptions(env.Config.WatchNS),
	})

	if err != nil {
		log.Global.Error(err, "unable to create manager")
		os.Exit(1)
	}

	if err = addIndexer(mgr); err != nil {
		log.Global.Error(err, "unable to index required fields")
		os.Exit(1)
	}

	registerControllers(mgr)

	log.Global.Debug("setting up webhook handlers")
	if env.Config.EnableWebhook {
		if err = (&v1beta1.ApiDefinition{}).SetupWebhookWithManager(mgr); err != nil {
			log.Global.Error(err, "unable to setup webhook for API definitions v1beta1")
			os.Exit(1)
		}

		if err = (&v1alpha1.ApiDefinition{}).SetupWebhookWithManager(mgr); err != nil {
			log.Global.Error(err, "unable to setup webhook for API definitions v1alpha1")
			os.Exit(1)
		}
	}

	//+kubebuilder:scaffold:builder

	log.Global.Debug("adding health probe")
	if healthCheckErr := mgr.AddHealthzCheck("healthz", healthz.Ping); healthCheckErr != nil {
		log.Global.Error(healthCheckErr, "unable to set up health check")
		os.Exit(1)
	}

	log.Global.Debug("adding readiness probe")
	if readyCheckErr := mgr.AddReadyzCheck("readyz", healthz.Ping); readyCheckErr != nil {
		log.Global.Error(readyCheckErr, "unable to set up ready check")
		os.Exit(1)
	}

	k8s.RegisterClient(mgr.GetClient())

	log.Global.Info("starting manager")
	ctx := ctrl.SetupSignalHandler()
	if startErr := mgr.Start(ctx); startErr != nil {
		log.Global.Error(startErr, "unable to run manager")
		os.Exit(1)
	}
}

func buildCacheOptions(ns string) cache.Options {
	if ns == "" {
		log.Global.Info("manager is watching all cluster namespaces")
		return cache.Options{}
	}

	log.Global.Infof("manager is watching namespace %s", ns)
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
		Watcher:  watch.New(context.Background(), mgr.GetClient(), &v1beta1.ApiDefinitionList{}),
	}).SetupWithManager(mgr); err != nil {
		log.Global.Error(err, "unable to create controller for API definitions")
		os.Exit(1)
	}

	if err := (&managementcontext.Reconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("managementcontext-controller"),
		Watcher:  watch.New(context.Background(), mgr.GetClient(), &v1beta1.ManagementContextList{}),
	}).SetupWithManager(mgr); err != nil {
		log.Global.Error(err, "unable to create controller for management contexts")
		os.Exit(1)
	}
	if err := (&ingress.Reconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("ingress-controller"),
		Watcher:  watch.New(context.Background(), mgr.GetClient(), &v1.IngressList{}),
	}).SetupWithManager(mgr); err != nil {
		log.Global.Error(err, "unable to create controller for ingresses")
		os.Exit(1)
	}
	if err := (&apiresource.Reconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("apiresource-controller"),
	}).SetupWithManager(mgr); err != nil {
		log.Global.Error(err, "unable to create controller for API resources")
		os.Exit(1)
	}
	if err := (&application.Reconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("application-controller"),
		Watcher:  watch.New(context.Background(), mgr.GetClient(), &v1beta1.ApplicationList{}),
	}).SetupWithManager(mgr); err != nil {
		log.Global.Error(err, "unable to create controller for applications")
		os.Exit(1)
	}

	if err := (&secrets.Reconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		log.Global.Error(err, "unable to create controller for secrets")
		os.Exit(1)
	}
}

func addIndexer(mgr manager.Manager) error {
	err := indexApiDefinitionFields(mgr)
	if err != nil {
		return err
	}

	err = indexSecretRefs(mgr)
	if err != nil {
		return err
	}

	err = indexIngressFields(mgr)
	if err != nil {
		return err
	}

	err = indexTLSSecretFields(mgr)
	if err != nil {
		return err
	}

	err = indexApplicationFields(mgr)
	if err != nil {
		return err
	}

	return nil
}

func indexApiDefinitionFields(manager ctrl.Manager) error {
	cache := manager.GetCache()
	ctx := context.Background()

	contextIndexer := indexer.NewIndexer(indexer.ContextField, indexer.IndexManagementContexts)
	err := cache.IndexField(ctx, &v1beta1.ApiDefinition{}, contextIndexer.Field, contextIndexer.Func)
	if err != nil {
		return err
	}

	resourceIndexer := indexer.NewIndexer(indexer.ResourceField, indexer.IndexApiResourceRefs)
	err = cache.IndexField(ctx, &v1beta1.ApiDefinition{}, resourceIndexer.Field, resourceIndexer.Func)
	if err != nil {
		return err
	}

	return nil
}

func indexSecretRefs(manager ctrl.Manager) error {
	cache := manager.GetCache()
	ctx := context.Background()

	secretRefIndexer := indexer.NewIndexer(indexer.SecretRefField, indexer.IndexManagementContextSecrets)
	return cache.IndexField(ctx, &v1beta1.ManagementContext{}, secretRefIndexer.Field, secretRefIndexer.Func)
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
	err := cache.IndexField(ctx, &v1beta1.Application{}, appContextIndexer.Field, appContextIndexer.Func)
	if err != nil {
		return err
	}

	return nil
}

func extendResources() error {
	ctx := context.Background()

	log.Global.Debug("extending scheme for CRD patch")
	sErr := apiExtensions.AddToScheme(scheme)
	if sErr != nil {
		return sErr
	}

	cli, cErr := client.New(ctrl.GetConfigOrDie(), client.Options{Scheme: scheme})
	if cErr != nil {
		return cErr
	}

	vErr := extensions.InitValidatingWebhookConfig(ctx, cli)
	if vErr != nil {
		return vErr
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

		log.Global.Infof("patching custom resource definition %s", unstructured.GetName())

		existing := &apiExtensions.CustomResourceDefinition{}
		if err = cli.Get(ctx, types.NamespacedName{Name: unstructured.GetName()}, existing); err != nil {
			return err
		}

		desired := &apiExtensions.CustomResourceDefinition{}
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured.Object, desired); err != nil {
			return err
		}

		extensions.AddConversionWebhook(desired)
		extensions.InjectCA(desired)

		desired.Spec.DeepCopyInto(&existing.Spec)
		existing.ObjectMeta.Annotations = desired.GetAnnotations()

		if err = cli.Update(ctx, existing); err != nil {
			return err
		}

		return extensions.CreateValidatingWebhooks(ctx, cli, existing)
	})
}
