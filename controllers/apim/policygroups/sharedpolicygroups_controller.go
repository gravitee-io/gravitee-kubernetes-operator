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

package policygroups

import (
	"context"
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
	corev1 "k8s.io/api/core/v1"

	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/policygroups/internal"

	"github.com/go-logr/logr"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/event"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/template"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

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

// Reconcile will handle v1alpha1.SharedPolicyGroup CRDs
// +kubebuilder:rbac:groups=gravitee.io,resources=sharedpolicygroups,verbs=get;list;watch;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=gravitee.io,resources=sharedpolicygroups/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gravitee.io,resources=sharedpolicygroups/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	spg := &v1alpha1.SharedPolicyGroup{}
	if err := r.Get(ctx, req.NamespacedName, spg); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	spg.SetNamespace(req.Namespace)

	events := event.NewRecorder(r.Recorder)

	dc := spg.DeepCopy()

	_, err := util.CreateOrUpdate(ctx, r.Client, spg, func() error {
		util.AddFinalizer(spg, core.SharedPolicyGroupFinalizer)
		k8s.AddAnnotation(spg, core.LastSpecHashAnnotation, hash.Calculate(&spg.Spec))

		if err := template.Compile(ctx, dc, true); err != nil {
			spg.Status.ProcessingStatus = core.ProcessingStatusFailed
			return err
		}

		var err error
		if spg.IsBeingDeleted() {
			err = events.Record(event.Delete, spg, func() error {
				if err := internal.Delete(ctx, dc); err != nil {
					return err
				}
				util.RemoveFinalizer(spg, core.SharedPolicyGroupFinalizer)
				return nil
			})
		} else {
			err = events.Record(event.Update, spg, func() error {
				return internal.CreateOrUpdate(ctx, dc)
			})
		}

		return err
	})

	if err := dc.GetStatus().DeepCopyTo(spg); err != nil {
		return ctrl.Result{}, err
	}

	if err == nil {
		log.InfoEndReconcile(ctx, spg)
		return ctrl.Result{}, internal.UpdateStatusSuccess(ctx, spg)
	}

	// An error occurred during the reconcile
	if err := internal.UpdateStatusFailure(ctx, spg); err != nil {
		return ctrl.Result{}, err
	}

	if errors.IsRecoverable(err) {
		log.ErrorRequeuingReconcile(ctx, err, spg)
		return ctrl.Result{RequeueAfter: requeueAfterTime}, err
	}

	log.ErrorAbortingReconcile(ctx, err, spg)
	return ctrl.Result{}, nil
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.SharedPolicyGroup{}).
		Watches(&v1alpha1.ManagementContext{}, r.Watcher.WatchContexts(search.SPGContextField)).
		Watches(&corev1.Secret{}, r.Watcher.WatchTemplatingSource("sharedpolicygroups")).
		Watches(&corev1.ConfigMap{}, r.Watcher.WatchTemplatingSource("sharedpolicygroups")).
		WithEventFilter(predicate.LastSpecHashPredicate{}).
		Complete(r)
}
