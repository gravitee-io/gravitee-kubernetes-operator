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
	"crypto/tls"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"strings"

	v2Admission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v2"
	v4Admission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v4"
	appAdmission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/application"
	mctxAdmission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/mctx"
	resourceAdmission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/resource"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	wk "github.com/gravitee-io/gravitee-kubernetes-operator/internal/webhook"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/secrets"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	v1 "k8s.io/api/networking/v1"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/indexer"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/logging"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/watch"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.

	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/apidefinition"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/apiresource"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/ingress"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/managementcontext"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricServer "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")

	//go:embed helm
	helm embed.FS
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(v1alpha1.AddToScheme(scheme))
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

	var webhookServer webhook.Server
	if env.Config.EnableWebhook {
		patchAdmissionWebhook()
		webhookServer = webhook.NewServer(webhook.Options{
			CertDir:  "/tmp/webhook-server/certs",
			Port:     env.Config.WebhookPort,
			CertName: wk.CertName,
			KeyName:  wk.KeyName,
			TLSOpts: []func(*tls.Config){func(config *tls.Config) {
				config.InsecureSkipVerify = true
			}},
		})
	}

	opts := zap.Options{
		Development:          env.Config.Development,
		EncoderConfigOptions: logging.NewEncoderConfigOption(),
	}

	opts.BindFlags(flag.CommandLine)

	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	if !env.Config.EnableMetrics {
		metricsAddr = "0" // disables metrics
	}

	if env.Config.HTTPClientInsecureSkipVerify {
		setupLog.Info("TLS verification is skipped for APIM HTTP client")
	}

	metrics := metricServer.Options{BindAddress: metricsAddr}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		Metrics:                metrics,
		WebhookServer:          webhookServer,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "24d975d3.gravitee.io",
		Cache:                  buildCacheOptions(env.Config.NS),
	})

	k8s.RegisterClient(mgr.GetClient())

	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if env.Config.ApplyCRDs {
		if err = applyCRDs(); err != nil {
			setupLog.Error(err, "unable to apply custom resource definitions")
			os.Exit(1)
		}
	}

	if err = indexer.InitCache(context.Background(), mgr.GetCache()); err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	registerControllers(mgr)

	if env.Config.EnableWebhook {
		err = setupAdmissionWebhooks(mgr)
		if err != nil {
			setupLog.Error(err, "unable to start manager")
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

	setupLog.Info("starting manager")
	if startErr := mgr.Start(ctrl.SetupSignalHandler()); startErr != nil {
		setupLog.Error(startErr, "problem running manager")
		os.Exit(1)
	}
}

func buildCacheOptions(ns string) cache.Options {
	if ns == "" {
		return cache.Options{}
	}
	defaultNamespaces := map[string]cache.Config{}
	configNamespaces := strings.Split(env.Config.NS, ",")
	for _, ns := range configNamespaces {
		defaultNamespaces[ns] = cache.Config{}
	}
	return cache.Options{
		DefaultNamespaces: defaultNamespaces,
	}
}

func registerControllers(mgr manager.Manager) {
	const msg = "unable to create controller"
	const controller = "controller"
	if err := (&apidefinition.Reconciler{
		Client:   k8s.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("apidefinitionv2-controller"),
		Watcher:  watch.New(context.Background(), k8s.GetClient(), &v1alpha1.ApiDefinitionList{}),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, msg, controller, "ApiDefinition")
		os.Exit(1)
	}

	if err := (&apidefinition.V4Reconciler{
		Client:   k8s.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("apiv4definition-controller"),
		Watcher:  watch.New(context.Background(), k8s.GetClient(), &v1alpha1.ApiV4DefinitionList{}),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, msg, controller, "ApiV4Definition")
		os.Exit(1)
	}

	if err := (&managementcontext.Reconciler{
		Client:   k8s.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("managementcontext-controller"),
		Watcher:  watch.New(context.Background(), k8s.GetClient(), &v1alpha1.ManagementContextList{}),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, msg, controller, "ManagementContext")
		os.Exit(1)
	}
	if err := (&ingress.Reconciler{
		Client:   k8s.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("ingress-controller"),
		Watcher:  watch.New(context.Background(), k8s.GetClient(), &v1.IngressList{}),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, msg, controller, "Ingress")
		os.Exit(1)
	}
	if err := (&apiresource.Reconciler{
		Client:   k8s.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("apiresource-controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, msg, controller, "ApiResource")
		os.Exit(1)
	}
	if err := (&application.Reconciler{
		Client:   k8s.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("application-controller"),
		Watcher:  watch.New(context.Background(), k8s.GetClient(), &v1alpha1.ApplicationList{}),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, msg, controller, "Application")
		os.Exit(1)
	}

	if err := (&secrets.Reconciler{
		Client: k8s.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, msg, controller, "Secret")
		os.Exit(1)
	}
}

func applyCRDs() error {
	client := dynamic.NewForConfigOrDie(ctrl.GetConfigOrDie())
	ctx := context.Background()

	version := schema.GroupVersionResource{
		Group:    "apiextensions.k8s.io",
		Version:  "v1",
		Resource: "customresourcedefinitions",
	}

	return fs.WalkDir(helm, "helm/gko/crds", func(path string, d fs.DirEntry, walkErr error) error {
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

		opts := metav1.ApplyOptions{Force: true, FieldManager: "gravitee.io/operator"}
		crd := &unstructured.Unstructured{Object: obj}

		if crd, err = client.Resource(version).Apply(ctx, crd.GetName(), crd, opts); err == nil {
			setupLog.Info("applied resource definition", "name", crd.GetName())
		}

		return err
	})
}

func patchAdmissionWebhook() {
	setupLog.Info("setting up Admission Webhook Server")
	webhookPatcher := wk.NewWebhookPatcher()
	svc := env.Config.WebhookService
	ns := env.Config.WebhookNS
	host := strings.Join(
		[]string{
			svc,
			fmt.Sprintf("%s.%s", svc, ns),
			fmt.Sprintf("%s.%s.svc", svc, ns),
			fmt.Sprintf("%s.%s.svc.cluster.local", svc, ns),
		},
		",",
	)
	err := webhookPatcher.CreateSecret(context.Background(), env.Config.WebhookCertSecret, ns, host)
	if err != nil {
		panic(err)
	}

	err = webhookPatcher.UpdateValidationCaBundle(
		context.Background(),
		wk.ValidatingWebhookName,
		env.Config.WebhookCertSecret,
		ns)
	if err != nil {
		setupLog.Error(err, "Can not update CA bundle for GKO validation webhook. GKO can not start")
		panic(err)
	}

	err = webhookPatcher.UpdateMutationCaBundle(
		context.Background(),
		wk.MutatingWebhookName,
		env.Config.WebhookCertSecret,
		ns)
	if err != nil {
		setupLog.Error(err, "Can not update CA bundle for GKO mutation webhook. GKO can not start")
		panic(err)
	}
}

func setupAdmissionWebhooks(mgr manager.Manager) error {
	if err := (resourceAdmission.AdmissionCtrl{}).SetupWithManager(mgr); err != nil {
		return err
	}
	if err := (v2Admission.AdmissionCtrl{}).SetupWithManager(mgr); err != nil {
		return err
	}
	if err := (v4Admission.AdmissionCtrl{}).SetupWithManager(mgr); err != nil {
		return err
	}
	if err := (appAdmission.AdmissionCtrl{}).SetupWithManager(mgr); err != nil {
		return err
	}
	if err := (mctxAdmission.AdmissionCtrl{}).SetupWithManager(mgr); err != nil {
		return err
	}

	return nil
}
