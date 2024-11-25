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
	"time"

	"github.com/go-logr/logr"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/subscription/internal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/event"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/template"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/predicate"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/watch"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

const requeueAfterTime = time.Second * 5

type Reconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Watcher  watch.Interface
}

// +kubebuilder:rbac:groups=gravitee.io,resources=subscriptions,verbs=get;list;watch;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=gravitee.io,resources=subscriptions/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gravitee.io,resources=subscriptions/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	subscription := &v1alpha1.Subscription{}
	if err := r.Get(ctx, req.NamespacedName, subscription); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	subscription.SetNamespace(req.Namespace)

	events := event.NewRecorder(r.Recorder)

	status := subscription.Status.DeepCopy()
	_, reconcileErr := util.CreateOrUpdate(ctx, r.Client, subscription, func() error {
		util.AddFinalizer(subscription, core.SubscriptionFinalizer)
		k8s.AddAnnotation(subscription, core.LastSpecHashAnnotation, hash.Calculate(&subscription.Spec))

		dc := subscription.DeepCopy()
		if err := template.Compile(ctx, dc); err != nil {
			subscription.Status.ProcessingStatus = core.ProcessingStatusFailed
			return err
		}

		var err error
		if subscription.IsBeingDeleted() {
			err = events.Record(event.Delete, subscription, func() error {
				return internal.Delete(ctx, dc)
			})
		} else {
			err = events.Record(event.Update, subscription, func() error {
				return internal.CreateOrUpdate(ctx, dc)
			})
		}

		dc.Status.DeepCopyInto(status)
		subscription.SetFinalizers(dc.GetFinalizers())

		return err
	})

	status.DeepCopyInto(&subscription.Status)

	if reconcileErr == nil {
		logger.Info("Subscription has been reconciled")
		return ctrl.Result{}, internal.UpdateStatusSuccess(ctx, subscription)
	}

	// An error occurred during the reconcile
	if err := internal.UpdateStatusFailure(ctx, subscription); err != nil {
		return ctrl.Result{}, err
	}

	if errors.IsRecoverable(reconcileErr) {
		logger.Error(reconcileErr, "Requeuing reconcile")
		return ctrl.Result{RequeueAfter: requeueAfterTime}, reconcileErr
	}

	logger.Error(reconcileErr, "Aborting reconcile")
	return ctrl.Result{}, nil
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Subscription{}).
		WithEventFilter(predicate.LastSpecHashPredicate{}).
		Complete(r)
}
