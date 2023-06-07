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

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env/template"

	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/managementcontext/internal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/event"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/watch"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

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
	logger := log.FromContext(ctx)
	managementContext := &gio.ManagementContext{}
	if err := r.Get(ctx, req.NamespacedName, managementContext); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if err := template.NewResolver(ctx, r.Client, logger, managementContext).Resolve(); err != nil {
		return ctrl.Result{}, err
	}

	events := event.NewRecorder(r.Recorder)
	var reconcileErr error
	if managementContext.IsBeingDeleted() {
		reconcileErr = events.Record(event.Delete, managementContext, func() error {
			return internal.Delete(ctx, r.Client, managementContext)
		})
	} else {
		reconcileErr = events.Record(event.Update, managementContext, func() error {
			return internal.CreateOrUpdate(ctx, r.Client, managementContext)
		})
	}

	if reconcileErr == nil {
		logger.Info("Management context has been reconciled")
		return ctrl.Result{}, nil
	}

	// There was an error reconciling the Management Context
	return ctrl.Result{}, reconcileErr
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gio.ManagementContext{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)
}
