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

package subscription

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimsubscription "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/subscription"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/lifecycle"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/predicate"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/watch"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Lifecycle = lifecycle.ResourceLifecycle[
	*v1alpha1.Subscription, apimsubscription.DTO, *apim.SDKClient, apimsubscription.Response,
]

type Reconciler struct {
	client.Client
	Lifecycle Lifecycle
	Scheme    *runtime.Scheme
	Recorder  record.EventRecorder
	Watcher   watch.Interface
}

// NewLifecycle wires the Subscription reconcile steps. The management context is the subscribed API's.
func NewLifecycle() Lifecycle {
	return lifecycle.NewResourceLifecycle(Lifecycle{
		Finalizer:     core.SubscriptionFinalizer,
		ResolveRefs:   apimsubscription.ResolveRefs,
		ClientFactory: apimsubscription.CreateAPIMClient,
		ToDTO:         apimsubscription.ToResolvedDTO,
		Delete:        apimsubscription.Delete,
		Upsert:        apimsubscription.Upsert,
		PostUpsert:    apimsubscription.UpdateStatusAndRecount,
		PostDelete:    apimsubscription.RecountParents,
	})
}

// +kubebuilder:rbac:groups=gravitee.io,resources=subscriptions,verbs=get;list;watch;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=gravitee.io,resources=subscriptions/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gravitee.io,resources=subscriptions/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return r.Lifecycle.Reconcile(ctx, r.Client, r.Recorder, req, &v1alpha1.Subscription{})
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	newController := ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Subscription{}).
		WithEventFilter(predicate.LastSpecHashPredicate{})

	if env.Config.EnableTemplating {
		newController.Watches(&corev1.Secret{}, r.Watcher.WatchTemplatingSource("subscriptions"))
	}

	return newController.Complete(r)
}
