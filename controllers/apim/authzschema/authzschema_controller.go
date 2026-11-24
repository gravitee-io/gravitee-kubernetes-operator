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

package authzschema

import (
	"context"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimauthzschema "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/authzschema"
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
	*v1alpha1.AuthzSchema, sdk.AuthzSchema, *apim.SDKClient, apimauthzschema.Response,
]

type Reconciler struct {
	Scheme    *runtime.Scheme
	Recorder  record.EventRecorder
	Watcher   watch.Interface
	Client    client.Client
	Lifecycle Lifecycle
}

// NewLifecycle wires the AuthzSchema reconcile steps.
func NewLifecycle() Lifecycle {
	return lifecycle.NewResourceLifecycle(Lifecycle{
		Finalizer:     core.AuthzSchemaFinalizer,
		ResolveRefs:   apimauthzschema.ResolveRefs,
		ClientFactory: apimauthzschema.CreateAPIMClient,
		ToDTO:         apimauthzschema.ToDTO,
		Delete:        apimauthzschema.Delete,
		Upsert:        apimauthzschema.Upsert,
		PostUpsert:    apimauthzschema.UpdateStatus,
		OnUpsertError: apimauthzschema.OnUpsertError,
	})
}

// +kubebuilder:rbac:groups=gravitee.io,resources=authzschemas,verbs=get;list;watch;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=gravitee.io,resources=authzschemas/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gravitee.io,resources=authzschemas/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return r.Lifecycle.Reconcile(ctx, r.Client, r.Recorder, req, &v1alpha1.AuthzSchema{})
}

// SetupWithManager filters events per watch rather than with a controller-wide event filter: an
// AuthzPdp gaining its routing scope is a status-only update, which LastSpecHashPredicate would
// drop, leaving a schema waiting on that PDP to its retry backoff.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	lastSpecHash := builder.WithPredicates(predicate.LastSpecHashPredicate{})

	newController := ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.AuthzSchema{}, lastSpecHash).
		Watches(&v1alpha1.ManagementContext{}, r.Watcher.WatchContexts(search.AuthzSchemaContextField), lastSpecHash).
		Watches(
			&v1alpha1.AuthzPdp{},
			r.Watcher.WatchAuthzPdps(search.AuthzSchemaAuthzPdpField),
			builder.WithPredicates(predicate.AuthzPdpRoutedPredicate{}),
		)

	if env.Config.EnableTemplating {
		newController.
			Watches(&corev1.Secret{}, r.Watcher.WatchTemplatingSource("authzschemas"), lastSpecHash).
			Watches(&corev1.ConfigMap{}, r.Watcher.WatchTemplatingSource("authzschemas"), lastSpecHash)
	}
	return newController.Complete(r)
}
