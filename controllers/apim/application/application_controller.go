/*
Copyright 2022 DAVID BRASSELY.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package application

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/event"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/indexer"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/predicate"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/template"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/watch"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/go-logr/logr"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/application/internal"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const requeueAfterTime = time.Second * 5

// Reconciler reconciles a Application object.
type Reconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Watcher  watch.Interface
}

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gravitee.io,resources=applications,verbs=get;list;watch;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=gravitee.io,resources=applications/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gravitee.io,resources=applications/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	application := &v1alpha1.Application{}
	if err := r.Get(ctx, req.NamespacedName, application); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	events := event.NewRecorder(r.Recorder)

	if application.Spec.Context == nil {
		log.ErrorAbortingReconcile(
			ctx,
			fmt.Errorf("no context is provided, no attempt will be made to sync with APIM"),
			application,
		)
		return ctrl.Result{}, nil
	}

	dc := application.DeepCopy()

	_, err := util.CreateOrUpdate(ctx, r.Client, application, func() error {
		util.AddFinalizer(application, core.ApplicationFinalizer)
		k8s.AddAnnotation(application, core.LastSpecHashAnnotation, hash.Calculate(&application.Spec))

		if err := template.Compile(ctx, dc, true); err != nil {
			application.Status.ProcessingStatus = core.ProcessingStatusFailed
			return err
		}

		var err error
		if application.IsBeingDeleted() {
			err = events.Record(event.Delete, application, func() error {
				if err := internal.Delete(ctx, dc); err != nil {
					return err
				}
				util.RemoveFinalizer(application, core.ApplicationFinalizer)
				return nil
			})
		} else {
			err = events.Record(event.Update, application, func() error {
				return internal.CreateOrUpdate(ctx, dc)
			})
		}

		return err
	})

	if err := dc.GetStatus().DeepCopyTo(application); err != nil {
		return ctrl.Result{}, err
	}

	application.SetFinalizers(dc.GetFinalizers())

	if err == nil {
		log.InfoEndReconcile(ctx, application)
		return ctrl.Result{}, internal.UpdateStatusSuccess(ctx, application)
	}

	// An error occurred during the reconcile
	if err := internal.UpdateStatusFailure(ctx, application); err != nil {
		return ctrl.Result{}, err
	}

	if errors.IsRecoverable(err) {
		log.ErrorRequeuingReconcile(ctx, err, application)
		return ctrl.Result{RequeueAfter: requeueAfterTime}, err
	}

	log.ErrorAbortingReconcile(ctx, err, application)
	return ctrl.Result{}, nil
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	newController := ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Application{}).
		WithEventFilter(predicate.LastSpecHashPredicate{}).
		Watches(&v1alpha1.ManagementContext{}, r.Watcher.WatchContexts(indexer.AppContextField))
	if env.Config.EnableTemplating {
		newController.Watches(&corev1.Secret{}, r.Watcher.WatchTemplatingSource("applications")).
			Watches(&corev1.ConfigMap{}, r.Watcher.WatchTemplatingSource("applications"))
	}
	return newController.Complete(r)
}
