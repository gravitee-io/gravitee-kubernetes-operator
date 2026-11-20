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

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimmcpproxy "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/mcpproxy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/lifecycle"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/predicate"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/watch"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Lifecycle = lifecycle.ResourceLifecycle[*v1alpha1.McpProxy, sdk.AimMcpProxy, *apim.SDKClient, apimmcpproxy.Response]

type Reconciler struct {
	Scheme    *runtime.Scheme
	Recorder  record.EventRecorder
	Watcher   watch.Interface
	Client    client.Client
	Lifecycle Lifecycle
}

// NewLifecycle wires the McpProxy reconcile steps.
func NewLifecycle() Lifecycle {
	return lifecycle.NewResourceLifecycle(Lifecycle{
		Finalizer:     core.McpProxyFinalizer,
		ResolveRefs:   apimmcpproxy.ResolveRefs,
		ClientFactory: apimmcpproxy.CreateAPIMClient,
		ToDTO:         apimmcpproxy.ToDTO,
		Delete:        apimmcpproxy.Delete,
		Upsert:        apimmcpproxy.Upsert,
		PostUpsert:    apimmcpproxy.UpdateStatus,
		OnUpsertError: apimmcpproxy.OnUpsertError,
	})
}

// +kubebuilder:rbac:groups=gravitee.io,resources=mcpproxies,verbs=get;list;watch;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=gravitee.io,resources=mcpproxies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gravitee.io,resources=mcpproxies/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return r.Lifecycle.Reconcile(ctx, r.Client, r.Recorder, req, &v1alpha1.McpProxy{})
}

// templatingSourceKind is the annotation the template compiler writes on a Secret or ConfigMap an
// McpProxy reads: gravitee.io/ plus the lowercase kind plus "s", which is not the plural.
const templatingSourceKind = "mcpproxys"

// SetupWithManager filters events per watch rather than with a controller-wide event filter: a
// catalog server becoming synced is a status-only update, which LastSpecHashPredicate would drop,
// leaving a studio waiting on that server to its retry backoff.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	lastSpecHash := builder.WithPredicates(predicate.LastSpecHashPredicate{})

	newController := ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.McpProxy{}, lastSpecHash).
		Watches(&v1alpha1.ManagementContext{}, r.Watcher.WatchContexts(search.McpProxyContextField), lastSpecHash).
		Watches(
			&v1alpha1.CatalogMcpServer{},
			r.Watcher.WatchCatalogMcpServers(search.McpProxyCatalogMcpServerField),
			builder.WithPredicates(predicate.CatalogMcpServerSyncedPredicate{}),
		)

	if env.Config.EnableTemplating {
		newController.
			Watches(&corev1.Secret{}, r.Watcher.WatchTemplatingSource(templatingSourceKind), lastSpecHash).
			Watches(&corev1.ConfigMap{}, r.Watcher.WatchTemplatingSource(templatingSourceKind), lastSpecHash)
	}
	return newController.Complete(r)
}
