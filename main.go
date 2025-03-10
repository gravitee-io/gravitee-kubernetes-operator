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
	"fmt"

	"io/fs"
	"os"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"

	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/policygroups"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/gateway-api/gateway"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/gateway-api/gatewayclass"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/gateway-api/httproute"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/gateway-api/parameters"

	v2Admission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v2"
	v4Admission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v4"
	appAdmission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/application"
	groupAdmission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/group"
	mctxAdmission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/mctx"
	spgAdmission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/policygroups"
	resourceAdmission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/resource"
	subAdmission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/subscription"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"

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
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/group"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/notification"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/subscription"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	v1 "k8s.io/api/networking/v1"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/watch"

	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"

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
	runtimeUtil "k8s.io/apimachinery/pkg/util/runtime"
	cliScheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/metrics/filters"
	metricServer "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	//+kubebuilder:scaffold:imports
)

var (
	scheme = runtime.NewScheme()

	//go:embed helm
	helm embed.FS
)

func init() {
	runtimeUtil.Must(cliScheme.AddToScheme(scheme))

	runtimeUtil.Must(v1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
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

	if env.Config.HTTPClientInsecureSkipVerify {
		log.Global.Warn("TLS certificates verification is skipped for APIM HTTP client")
	}

	log.Global.Infof("Probes server listens on %s", env.GetProbesAddr())

	if env.Config.EnableWebhook {
		log.Global.Infof("Webhook server listens on :%d", env.Config.WebhookPort)
	}

	if env.Config.EnableMetrics {
		log.Global.Infof("Metrics server listens on %s", env.GetMetricsAddr())
	}

	metricsOpts := metricServer.Options{
		BindAddress:    env.GetMetricsAddr(),
		SecureServing:  env.Config.SecureMetrics,
		CertDir:        env.Config.MetricsCertDir,
		FilterProvider: filters.WithAuthenticationAndAuthorization,
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		Metrics:                metricsOpts,
		WebhookServer:          webhookServer,
		HealthProbeBindAddress: env.GetProbesAddr(),
		LeaderElection:         env.Config.EnableLeaderElection,
		LeaderElectionID:       "24d975d3.gravitee.io",
		Cache:                  buildCacheOptions(env.Config.NS),
	})

	if err != nil {
		log.Global.Error(err, "Unable to configure manager")
		os.Exit(1)
	}

	k8s.RegisterClient(mgr.GetClient())

	const cannotStartError = "Unable to start manager"

	if env.Config.ApplyCRDs {
		if err = applyCRDs(); err != nil {
			log.Global.Error(err, "Unable to apply custom resource definitions")
			os.Exit(1)
		}
	}

	if err = search.InitCache(context.Background(), mgr.GetCache()); err != nil {
		log.Global.Error(err, cannotStartError)
		os.Exit(1)
	}

	registerControllers(mgr)

	if env.Config.EnableWebhook {
		err = setupAdmissionWebhooks(mgr)
		if err != nil {
			log.Global.Error(err, cannotStartError)
			os.Exit(1)
		}
	}

	//+kubebuilder:scaffold:builder

	if healthCheckErr := mgr.AddHealthzCheck("healthz", healthz.Ping); healthCheckErr != nil {
		log.Global.Error(healthCheckErr, "Unable to set up health check")
		os.Exit(1)
	}

	if readyCheckErr := mgr.AddReadyzCheck("readyz", healthz.Ping); readyCheckErr != nil {
		log.Global.Error(readyCheckErr, "Unable to set up ready check")
		os.Exit(1)
	}

	log.Global.Info("Starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		log.Global.Error(err, cannotStartError)
		os.Exit(1)
	}
}

func buildCacheOptions(ns string) cache.Options {
	if ns == "" {
		log.Global.Info("Listening to all namespaces")
		return cache.Options{}
	}
	defaultNamespaces := map[string]cache.Config{}
	configNamespaces := strings.Split(env.Config.NS, ",")
	for _, ns := range configNamespaces {
		log.Global.Infof("Listening to namespace %s", ns)
		defaultNamespaces[ns] = cache.Config{}
	}
	return cache.Options{
		DefaultNamespaces: defaultNamespaces,
	}
}

func registerControllers(mgr manager.Manager) {
	if err := (&apidefinition.Reconciler{
		Client:   k8s.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("apidefinitionv2-controller"),
		Watcher:  watch.New(context.Background(), k8s.GetClient(), &v1alpha1.ApiDefinitionList{}),
	}).SetupWithManager(mgr); err != nil {
		log.Global.Error(err, "Unable to create controller for API definitions")
		os.Exit(1)
	}

	if err := (&apidefinition.V4Reconciler{
		Client:   k8s.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("apiv4definition-controller"),
		Watcher:  watch.New(context.Background(), k8s.GetClient(), &v1alpha1.ApiV4DefinitionList{}),
	}).SetupWithManager(mgr); err != nil {
		log.Global.Error(err, "Unable to create controller for API v4 definitions")
		os.Exit(1)
	}

	if err := (&managementcontext.Reconciler{
		Client:   k8s.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("managementcontext-controller"),
		Watcher:  watch.New(context.Background(), k8s.GetClient(), &v1alpha1.ManagementContextList{}),
	}).SetupWithManager(mgr); err != nil {
		log.Global.Error(err, "Unable to create controller for management contexts")
		os.Exit(1)
	}

	if env.Config.EnableIngress {
		if err := (&ingress.Reconciler{
			Client:   k8s.GetClient(),
			Scheme:   mgr.GetScheme(),
			Recorder: mgr.GetEventRecorderFor("ingress-controller"),
			Watcher:  watch.New(context.Background(), k8s.GetClient(), &v1.IngressList{}),
		}).SetupWithManager(mgr); err != nil {
			log.Global.Error(err, "Unable to create controller for ingresses")
			os.Exit(1)
		}
	}

	if err := (&apiresource.Reconciler{
		Client:   k8s.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("apiresource-controller"),
	}).SetupWithManager(mgr); err != nil {
		log.Global.Error(err, "Unable to create controller for API resources")
		os.Exit(1)
	}
	if err := (&application.Reconciler{
		Client:   k8s.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("application-controller"),
		Watcher:  watch.New(context.Background(), k8s.GetClient(), &v1alpha1.ApplicationList{}),
	}).SetupWithManager(mgr); err != nil {
		log.Global.Error(err, "Unable to create controller for applications")
		os.Exit(1)
	}

	if err := (&subscription.Reconciler{
		Scheme:   mgr.GetScheme(),
		Client:   mgr.GetClient(),
		Recorder: mgr.GetEventRecorderFor("subscription-controller"),
	}).SetupWithManager(mgr); err != nil {
		log.Global.Error(err, "Unable to create controller for subscriptions")
		os.Exit(1)
	}

	if err := (&group.Reconciler{
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("group-controller"),
	}).SetupWithManager(mgr); err != nil {
		log.Global.Error(err, "Unable to create controller for groups")
		os.Exit(1)
	}

	if err := (&policygroups.Reconciler{
		Scheme:   mgr.GetScheme(),
		Client:   mgr.GetClient(),
		Recorder: mgr.GetEventRecorderFor("sharedpolicygroups-controller"),
		Watcher:  watch.New(context.Background(), k8s.GetClient(), &v1alpha1.SharedPolicyGroupList{}),
	}).SetupWithManager(mgr); err != nil {
		log.Global.Error(err, "Unable to create controller for shared policy groups")
		os.Exit(1)
	}

	if err := (&notification.Reconciler{
		Scheme:   mgr.GetScheme(),
		Client:   mgr.GetClient(),
		Recorder: mgr.GetEventRecorderFor("notification-controller"),
	}).SetupWithManager(mgr); err != nil {
		log.Global.Error(err, "Unable to create controller for notification")
		os.Exit(1)
	}

	if env.Config.EnableGatewayAPI {
		registerGatewayAPIsControllers(mgr)
	}
}

func registerGatewayAPIsControllers(mgr ctrl.Manager) {
	runtimeUtil.Must(gwAPIv1.SchemeBuilder.AddToScheme(scheme))

	if err := (&parameters.Reconciler{
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("gravitee-gateway-controller"),
	}).SetupWithManager(mgr); err != nil {
		log.Global.Error(err, "Unable to create controller for gravitee gateway")
		os.Exit(1)
	}

	if err := (&gatewayclass.Reconciler{
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("gateway-class-controller"),
	}).SetupWithManager(mgr); err != nil {
		log.Global.Error(err, "Unable to create controller for gateway class")
		os.Exit(1)
	}

	if err := (&gateway.Reconciler{
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("gateway-controller"),
	}).SetupWithManager(mgr); err != nil {
		log.Global.Error(err, "Unable to create controller for gateway")
		os.Exit(1)
	}

	if err := (&httproute.Reconciler{
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("http-route"),
	}).SetupWithManager(mgr); err != nil {
		log.Global.Error(err, "Unable to create controller for HTTP route")
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

		if strings.Contains(path, "gateway-api") && !env.Config.EnableGatewayAPI {
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
			log.Global.Infof("Applied resource definition [%s]", crd.GetName())
		}

		return err
	})
}

func patchAdmissionWebhook() {
	log.Global.Debug("Setting up Admission Webhook Server")
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
		env.Config.WebhookValidatingConfigurationName,
		env.Config.WebhookCertSecret,
		ns)
	if err != nil {
		log.Global.Error(err, "Can not update CA bundle for GKO validation webhook. GKO can not start")
		panic(err)
	}

	err = webhookPatcher.UpdateMutationCaBundle(
		context.Background(),
		env.Config.WebhookMutatingConfigurationName,
		env.Config.WebhookCertSecret,
		ns)
	if err != nil {
		log.Global.Error(err, "Can not update CA bundle for GKO mutation webhook. GKO can not start")
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
	if err := (subAdmission.AdmissionCtrl{}).SetupWithManager(mgr); err != nil {
		return err
	}
	if err := (spgAdmission.AdmissionCtrl{}).SetupWithManager(mgr); err != nil {
		return err
	}
	if err := (groupAdmission.AdmissionCtrl{}).SetupWithManager(mgr); err != nil {
		return err
	}
	return nil
}
