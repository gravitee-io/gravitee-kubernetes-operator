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

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/template"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/predicate"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/managementcontext/internal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/event"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/watch"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
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
	managementContext := &v1alpha1.ManagementContext{}
	if err := r.Get(ctx, req.NamespacedName, managementContext); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	events := event.NewRecorder(r.Recorder)

	dc := managementContext.DeepCopy()

	_, reconcileErr := util.CreateOrUpdate(ctx, r.Client, managementContext, func() error {
		util.AddFinalizer(managementContext, core.ManagementContextFinalizer)
		k8s.AddAnnotation(managementContext, core.LastSpecHashAnnotation, hash.Calculate(&managementContext.Spec))

		if err := template.Compile(ctx, dc); err != nil {
			return err
		}

		var err error
		if managementContext.IsBeingDeleted() {
			err = events.Record(event.Delete, managementContext, func() error {
				if err := internal.Delete(ctx, dc); err != nil {
					return err
				}
				util.RemoveFinalizer(managementContext, core.ManagementContextFinalizer)
				return nil
			})
		} else {
			err = events.Record(event.Update, managementContext, func() error {
				// We don't do anything directly when there is an update on ManagementContext
				return nil
			})
		}

		return err
	})

	if err := dc.GetStatus().DeepCopyTo(managementContext); err != nil {
		return ctrl.Result{}, err
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
		For(&v1alpha1.ManagementContext{}).
		WithEventFilter(predicate.LastSpecHashPredicate{}).
		Complete(r)
}
