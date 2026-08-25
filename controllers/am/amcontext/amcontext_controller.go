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

package amcontext

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/am/amcontext/internal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/event"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/predicate"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/template"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/watch"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Reconciler reconciles an AMContext object.
type Reconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Watcher  watch.Interface
}

// +kubebuilder:rbac:groups=gravitee.io,resources=amcontexts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gravitee.io,resources=amcontexts/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gravitee.io,resources=amcontexts/finalizers,verbs=update
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	amContext := &v1alpha1.AMContext{}
	if err := r.Get(ctx, req.NamespacedName, amContext); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	events := event.NewRecorder(r.Recorder)

	amContext.SetConditions([]metav1.Condition{})

	dc := amContext.DeepCopy()

	_, err := util.CreateOrUpdate(ctx, r.Client, amContext, func() error {
		util.AddFinalizer(amContext, core.AMContextFinalizer)
		k8s.AddAnnotation(amContext, core.LastSpecHashAnnotation, hash.Calculate(&amContext.Spec))

		if amContext.IsBeingDeleted() {
			if err := template.ReleaseReferences(ctx, amContext); err != nil {
				return err
			}
		} else if err := template.Compile(ctx, dc, true); err != nil {
			return err
		}

		var err error
		if amContext.IsBeingDeleted() {
			err = events.Record(event.Delete, amContext, func() error {
				if err := internal.Delete(ctx, dc); err != nil {
					return err
				}
				util.RemoveFinalizer(amContext, core.AMContextFinalizer)
				return nil
			})
		} else {
			err = events.Record(event.Update, amContext, func() error {
				if err = internal.CreateOrUpdate(ctx, dc); err != nil {
					return err
				}
				amContext.Annotations[internal.LastSecretReferenceName] =
					dc.Annotations[internal.LastSecretReferenceName]
				return nil
			})
		}

		return err
	})

	if err := dc.GetStatus().DeepCopyTo(amContext); err != nil {
		return ctrl.Result{}, err
	}

	if err == nil {
		log.InfoEndReconcile(ctx, amContext)
		return ctrl.Result{}, internal.UpdateCondition(ctx, amContext, err)
	}

	log.ErrorAbortingReconcile(ctx, err, amContext)
	return ctrl.Result{}, internal.UpdateCondition(ctx, amContext, err)
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	newController := ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.AMContext{}).
		WithEventFilter(predicate.LastSpecHashPredicate{})
	if env.Config.EnableTemplating {
		newController.Watches(&corev1.Secret{}, r.Watcher.WatchTemplatingSource("amcontexts")).
			Watches(&corev1.ConfigMap{}, r.Watcher.WatchTemplatingSource("amcontexts"))
	}
	return newController.Complete(r)
}
