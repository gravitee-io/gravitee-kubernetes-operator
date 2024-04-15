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

package manager

import (
	"context"
	"io"

	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/apidefinition"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/apiresource"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/ingress"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/managementcontext"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/secrets"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/indexer"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/watch"

	. "github.com/onsi/ginkgo/v2"

	ctrl "sigs.k8s.io/controller-runtime"

	metrics "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	"sigs.k8s.io/controller-runtime/pkg/client"

	runtimeUtil "k8s.io/apimachinery/pkg/util/runtime"
	clientScheme "k8s.io/client-go/kubernetes/scheme"

	netV1 "k8s.io/api/networking/v1"
)

const (
	metricsAddr = ":0" // disable metrics
	probeAddr   = ":0" // disable probes
	managerPort = 0
)

var mgr ctrl.Manager

func Instance() ctrl.Manager {
	return mgr
}

func init() {
	if _, r := GinkgoConfiguration(); r.Verbose {
		log.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	} else {
		log.SetLogger(zap.New(zap.WriteTo(io.Discard), zap.UseDevMode(true)))
	}

	ctx := context.Background()

	scheme := clientScheme.Scheme
	runtimeUtil.Must(clientScheme.AddToScheme(scheme))
	runtimeUtil.Must(v1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme

	var err error

	mgr, err = ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		WebhookServer: webhook.NewServer(webhook.Options{
			Port: managerPort,
		}),
		Metrics:                metrics.Options{BindAddress: metricsAddr},
		HealthProbeBindAddress: probeAddr,
		Cache:                  cache.Options{},
	})

	runtimeUtil.Must(err)

	cli = mgr.GetClient()

	// register client so that tested code has access to it
	k8s.RegisterClient(cli)

	cache := mgr.GetCache()

	runtimeUtil.Must(indexer.InitCache(ctx, cache))

	// index event for assertions
	runtimeUtil.Must(
		mgr.GetFieldIndexer().IndexField(
			ctx,
			&v1.Event{},
			"involvedObject.name",
			func(rawObj client.Object) []string {
				event, _ := rawObj.(*v1.Event)
				return []string{event.InvolvedObject.Name}
			},
		),
	)

	runtimeUtil.Must(
		(&apidefinition.Reconciler{
			Client:   cli,
			Scheme:   mgr.GetScheme(),
			Recorder: mgr.GetEventRecorderFor("apidefinition_controller"),
			Watcher:  watch.New(context.Background(), cli, &v1alpha1.ApiDefinitionList{}),
		}).SetupWithManager(mgr),
	)

	runtimeUtil.Must(
		(&apidefinition.V4Reconciler{
			Client:   cli,
			Scheme:   mgr.GetScheme(),
			Recorder: mgr.GetEventRecorderFor("apiv4definition-controller"),
			Watcher:  watch.New(context.Background(), cli, &v1alpha1.ApiV4DefinitionList{}),
		}).SetupWithManager(mgr),
	)

	runtimeUtil.Must(
		(&managementcontext.Reconciler{
			Client:   cli,
			Scheme:   mgr.GetScheme(),
			Recorder: mgr.GetEventRecorderFor("managementcontext_controller"),
			Watcher:  watch.New(context.Background(), cli, &v1alpha1.ManagementContextList{}),
		}).SetupWithManager(mgr),
	)

	runtimeUtil.Must(
		(&ingress.Reconciler{
			Client:   mgr.GetClient(),
			Scheme:   mgr.GetScheme(),
			Recorder: mgr.GetEventRecorderFor("ingress-controller"),
			Watcher:  watch.New(context.Background(), mgr.GetClient(), &netV1.IngressList{}),
		}).SetupWithManager(mgr),
	)

	runtimeUtil.Must(
		(&apiresource.Reconciler{
			Client:   mgr.GetClient(),
			Scheme:   mgr.GetScheme(),
			Recorder: mgr.GetEventRecorderFor("apiresource-controller"),
		}).SetupWithManager(mgr),
	)

	runtimeUtil.Must(
		(&application.Reconciler{
			Client:   mgr.GetClient(),
			Scheme:   mgr.GetScheme(),
			Recorder: mgr.GetEventRecorderFor("application-controller"),
			Watcher:  watch.New(context.Background(), mgr.GetClient(), &v1alpha1.ApplicationList{}),
		}).SetupWithManager(mgr),
	)

	runtimeUtil.Must(
		(&secrets.Reconciler{
			Client: mgr.GetClient(),
			Scheme: mgr.GetScheme(),
		}).SetupWithManager(mgr),
	)

	go func() {
		runtimeUtil.Must(Instance().Start(ctrl.SetupSignalHandler()))
	}()
	<-Instance().Elected()
}
