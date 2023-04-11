/*
 * Copyright (C) 2015 The Gravitee team (http://gravitee.io)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package managementcontext

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/managementcontext/internal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/event"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/filter"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/watch"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

// Reconciler reconciles a ManagementContext object.
type Reconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Watcher  watch.Interface
}

// +kubebuilder:rbac:groups=gravitee.io,resources=managementcontexts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gravitee.io,resources=managementcontexts/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gravitee.io,resources=managementcontexts/finalizers,verbs=update
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	instance := &gio.ManagementContext{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	events := event.NewRecorder(r.Recorder)
	var reconcileErr error
	if instance.IsBeingDeleted() {
		reconcileErr = events.Record(event.Delete, instance, func() error {
			return internal.Delete(ctx, r.Client, instance)
		})
	} else {
		reconcileErr = events.Record(event.Update, instance, func() error {
			return internal.CreateOrUpdate(ctx, r.Client, instance)
		})
	}

	if reconcileErr == nil {
		log.Info("Management context has been reconciled")
		return ctrl.Result{}, nil
	}

	// There was an error reconciling the Management Context
	return ctrl.Result{}, reconcileErr
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gio.ManagementContext{}).
		WithEventFilter(filter.NoFinalizerUpdateFilter{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Watches(&source.Kind{Type: &v1.Secret{}}, r.Watcher.WatchContextSecrets()).
		Complete(r)
}
