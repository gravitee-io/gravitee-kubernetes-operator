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

package mcpproxy

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/mcpproxy/internal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/event"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/predicate"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/template"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/watch"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type Reconciler struct {
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Watcher  watch.Interface
	Client   client.Client
}

// +kubebuilder:rbac:groups=gravitee.io,resources=mcpproxies,verbs=get;list;watch;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=gravitee.io,resources=mcpproxies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gravitee.io,resources=mcpproxies/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	proxy := &v1alpha1.McpProxy{}
	if err := r.Client.Get(ctx, req.NamespacedName, proxy); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	proxy.SetNamespace(req.Namespace)

	events := event.NewRecorder(r.Recorder)

	k8s.ResetConditionsExceptAutomationAPI(proxy)
	dc := proxy.DeepCopy()

	_, err := util.CreateOrUpdate(ctx, k8s.GetClient(), proxy, func() error {
		util.AddFinalizer(proxy, core.McpProxyFinalizer)
		k8s.AddAnnotation(proxy, core.LastSpecHashAnnotation, hash.Calculate(&proxy.Spec))

		if proxy.IsBeingDeleted() {
			if err := template.ReleaseReferences(ctx, proxy); err != nil {
				return err
			}
		} else {
			// A studio selecting a server the platform does not know yet would be refused: wait
			// for it (ResolvedRefs=False, requeued) rather than calling the platform. Checked
			// before compiling templates, which registers this proxy on the Secrets it reads: a
			// failed first reconcile does not persist the finalizer that would release them.
			if err := dynamic.AssertCatalogMcpServersSynced(ctx, proxy); err != nil {
				return errors.NewResolveRefError(err)
			}
			if err := template.Compile(ctx, dc, true); err != nil {
				return err
			}
		}

		var err error
		if proxy.IsBeingDeleted() {
			err = events.Record(event.Delete, proxy, func() error {
				if err := internal.Delete(ctx, dc); err != nil {
					return err
				}
				util.RemoveFinalizer(proxy, core.McpProxyFinalizer)
				return nil
			})
		} else {
			err = events.Record(event.Update, proxy, func() error {
				return internal.CreateOrUpdate(ctx, dc)
			})
		}

		return err
	})

	if err := dc.GetStatus().DeepCopyTo(proxy); err != nil {
		return ctrl.Result{}, err
	}

	if err == nil {
		log.InfoEndReconcile(ctx, proxy)
		return ctrl.Result{}, internal.UpdateStatusSuccess(ctx, proxy)
	}

	if err := internal.UpdateStatusFailure(ctx, proxy, err); err != nil {
		return ctrl.Result{}, err
	}

	if errors.IsRecoverable(err) {
		log.ErrorRequeuingReconcile(ctx, err, proxy)
		return ctrl.Result{}, err
	}

	log.ErrorAbortingReconcile(ctx, err, proxy)
	return ctrl.Result{}, nil
}

// templatingSourceKind is the annotation the template compiler writes on a Secret or ConfigMap an
// McpProxy reads: gravitee.io/ plus the lowercase kind plus "s", which is not the plural.
const templatingSourceKind = "mcpproxys"

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	newController := ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.McpProxy{}).
		WithEventFilter(predicate.LastSpecHashPredicate{}).
		Watches(&v1alpha1.ManagementContext{}, r.Watcher.WatchContexts(search.McpProxyContextField)).
		Watches(&v1alpha1.CatalogMcpServer{}, r.Watcher.WatchCatalogMcpServers(search.McpProxyCatalogMcpServerField))

	if env.Config.EnableTemplating {
		newController.
			Watches(&corev1.Secret{}, r.Watcher.WatchTemplatingSource(templatingSourceKind)).
			Watches(&corev1.ConfigMap{}, r.Watcher.WatchTemplatingSource(templatingSourceKind))
	}
	return newController.Complete(r)
}
