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

package portaltheme

import (
	"context"

	"github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimportaltheme "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portaltheme"
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
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Lifecycle = lifecycle.ResourceLifecycle[*v1alpha1.PortalTheme, sdk.PortalTheme, *apim.SDKClient, apimportaltheme.Response]

type Reconciler struct {
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Watcher  watch.Interface
	Client   client.Client
}

// resourceLifecycle wires the reconcile steps of portal themes.
// ResolveRefs runs the delete guard itself, before the context lookup: a lifecycle DeleteGuard would run
// after it, and so never for a context deleted first.
var resourceLifecycle = lifecycle.NewResourceLifecycle(Lifecycle{
	Finalizer:     core.PortalThemeFinalizer,
	ResolveRefs:   apimportaltheme.ResolveRefs,
	ClientFactory: apimportaltheme.CreateAPIMClient,
	ToDTO:         apimportaltheme.ToDTO,
	Delete:        apimportaltheme.Delete,
	Upsert:        apimportaltheme.Upsert,
	PostUpsert:    apimportaltheme.UpdateStatus,
})

// +kubebuilder:rbac:groups=gravitee.io,resources=portalthemes,verbs=get;list;watch;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=gravitee.io,resources=portalthemes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gravitee.io,resources=portalthemes/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return resourceLifecycle.Reconcile(ctx, r.Client, r.Recorder, req, &v1alpha1.PortalTheme{})
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	newController := ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.PortalTheme{}).
		WithEventFilter(predicate.LastSpecHashPredicate{}).
		Watches(&v1alpha1.ManagementContext{}, r.Watcher.WatchContexts(search.PortalThemeContextField))

	if env.Config.EnableTemplating {
		newController.
			Watches(&corev1.Secret{}, r.Watcher.WatchTemplatingSource("portalthemes")).
			Watches(&corev1.ConfigMap{}, r.Watcher.WatchTemplatingSource("portalthemes"))
	}
	return newController.Complete(r)
}
