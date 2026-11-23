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

package authzentity

import (
	"context"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimauthzentity "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/authzentity"
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

type Lifecycle = lifecycle.ResourceLifecycle[
	*v1alpha1.AuthzEntity, sdk.AuthzEntity, *apim.SDKClient, apimauthzentity.Response,
]

type Reconciler struct {
	Scheme    *runtime.Scheme
	Recorder  record.EventRecorder
	Watcher   watch.Interface
	Client    client.Client
	Lifecycle Lifecycle
}

// NewLifecycle wires the AuthzEntity reconcile steps.
func NewLifecycle() Lifecycle {
	return lifecycle.NewResourceLifecycle(Lifecycle{
		Finalizer:     core.AuthzEntityFinalizer,
		ResolveRefs:   apimauthzentity.ResolveRefs,
		ClientFactory: apimauthzentity.CreateAPIMClient,
		ToDTO:         apimauthzentity.ToDTO,
		Delete:        apimauthzentity.Delete,
		Upsert:        apimauthzentity.Upsert,
		PostUpsert:    apimauthzentity.UpdateStatus,
		OnUpsertError: apimauthzentity.OnUpsertError,
	})
}

// +kubebuilder:rbac:groups=gravitee.io,resources=authzentities,verbs=get;list;watch;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=gravitee.io,resources=authzentities/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gravitee.io,resources=authzentities/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return r.Lifecycle.Reconcile(ctx, r.Client, r.Recorder, req, &v1alpha1.AuthzEntity{})
}

// SetupWithManager filters events per watch rather than with a controller-wide event filter: an
// AuthzPdp gaining its routing scope is a status-only update, which LastSpecHashPredicate would
// drop, leaving an entity waiting on that PDP to its retry backoff.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	lastSpecHash := builder.WithPredicates(predicate.LastSpecHashPredicate{})

	newController := ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.AuthzEntity{}, lastSpecHash).
		Watches(&v1alpha1.ManagementContext{}, r.Watcher.WatchContexts(search.AuthzEntityContextField), lastSpecHash).
		Watches(
			&v1alpha1.AuthzPdp{},
			r.Watcher.WatchAuthzPdps(search.AuthzEntityAuthzPdpField),
			builder.WithPredicates(predicate.AuthzPdpRoutedPredicate{}),
		).
		// A parent's delete takes its RESOURCE children and detaches its PRINCIPAL children on the
		// platform: each child is applied again once the parent is gone.
		Watches(&v1alpha1.AuthzEntity{}, r.Watcher.WatchAuthzEntityParents())

	if env.Config.EnableTemplating {
		newController.
			Watches(&corev1.Secret{}, r.Watcher.WatchTemplatingSource("authzentities"), lastSpecHash).
			Watches(&corev1.ConfigMap{}, r.Watcher.WatchTemplatingSource("authzentities"), lastSpecHash)
	}
	return newController.Complete(r)
}
