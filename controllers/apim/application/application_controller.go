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

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/indexer"

	"github.com/go-logr/logr"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/application/internal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/event"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/watch"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
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
	logger := log.FromContext(ctx)

	application := &gio.Application{}

	if err := r.Get(ctx, req.NamespacedName, application); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	delegate := internal.NewDelegate(ctx, r.Client, logger)
	events := event.NewRecorder(r.Recorder)

	if application.Spec.Context == nil {
		logger.Error(fmt.Errorf("no context is provided, no attempt will be made to sync with APIM"), "Aborting reconcile")
		return ctrl.Result{}, nil
	}

	if err := delegate.ResolveContext(application); err != nil {
		logger.Error(err, "Unable to resolve context, no attempt will be made to sync with APIM")
		return ctrl.Result{}, err
	}

	if application.IsMissingDeletionFinalizer() {
		err := delegate.AddDeletionFinalizer(application)
		if err != nil {
			logger.Error(err, "Unable to add deletion finalizer to Application")
			return ctrl.Result{}, err
		}
	}

	var reconcileErr error

	if application.IsBeingDeleted() {
		reconcileErr = events.Record(event.Delete, application, func() error {
			return delegate.Delete(application)
		})
	} else {
		reconcileErr = events.Record(event.Update, application, func() error {
			return delegate.CreateOrUpdate(application)
		})
	}

	if reconcileErr == nil {
		logger.Info("Application has been reconciled")
		return ctrl.Result{}, delegate.UpdateStatusSuccess(application)
	}

	// An error occurred during the reconcile
	if err := delegate.UpdateStatusFailure(application); err != nil {
		return ctrl.Result{}, err
	}

	if apim.IsRecoverable(reconcileErr) {
		logger.Error(reconcileErr, "Requeuing reconcile")
		return ctrl.Result{RequeueAfter: requeueAfterTime}, reconcileErr
	}

	logger.Error(reconcileErr, "Aborting reconcile")
	return ctrl.Result{}, nil
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gio.Application{}).
		Watches(&source.Kind{Type: &gio.ManagementContext{}}, r.Watcher.WatchContexts(indexer.AppContextField)).
		Complete(r)
}
