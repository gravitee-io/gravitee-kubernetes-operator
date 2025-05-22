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

package notification

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/predicate"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/notification/internal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/event"
)

// Reconciler reconciles a Notification object.
type Reconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=gravitee.io,resources=notification,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=gravitee.io,resources=notification/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=gravitee.io,resources=notification/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	notification := &v1alpha1.Notification{}
	conditions := make([]v1.Condition, 0)
	notification.Status.Conditions = &conditions

	if err := r.Get(ctx, req.NamespacedName, notification); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	events := event.NewRecorder(r.Recorder)

	deepCopy := notification.DeepCopy()

	_, err := util.CreateOrUpdate(ctx, r.Client, notification, func() error {
		util.AddFinalizer(notification, core.NotificationFinalizer)
		k8s.AddAnnotation(notification, core.LastSpecHashAnnotation, hash.Calculate(&notification.Spec))

		var err error
		if notification.IsBeingDeleted() {
			err = events.Record(event.Delete, notification, func() error {
				if err := internal.Delete(ctx, deepCopy); err != nil {
					return err
				}
				util.RemoveFinalizer(notification, core.NotificationFinalizer)
				return nil
			})
		} else {
			err = events.Record(event.Update, notification, func() error {
				return internal.ResolveGroupRefs(ctx, notification.Spec.Type, notification.Namespace)
			})
		}
		return err
	})

	if !notification.IsBeingDeleted() && err != nil {
		if err := internal.SetGroupRefsConditions(ctx, r.Client, err, notification); err != nil {
			return ctrl.Result{}, err
		}
	}

	// in any case of error
	if err != nil {
		return ctrl.Result{}, err
	}

	// update status
	if err := internal.SetAcceptedCondition(ctx, r.Client, notification); err != nil {
		return ctrl.Result{}, err
	}

	// no error, we are done
	log.InfoEndReconcile(ctx, notification)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Notification{}).
		WithEventFilter(predicate.LastSpecHashPredicate{}).
		Complete(r)
}
