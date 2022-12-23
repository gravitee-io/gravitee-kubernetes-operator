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

package apidefinition

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/go-logr/logr"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/apidefinition/internal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/event"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/indexer"
)

const requeueAfterTime = time.Second * 5

// Reconciler reconciles a ApiDefinition object.
type Reconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gravitee.io,resources=apidefinitions,verbs=get;list;watch;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=gravitee.io,resources=apidefinitions/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gravitee.io,resources=apidefinitions/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("namespace", req.Namespace, "name", req.Name)

	apiDefinition := &gio.ApiDefinition{}

	if err := r.Get(ctx, req.NamespacedName, apiDefinition); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	delegate := internal.NewDelegate(ctx, r.Client, log)
	events := event.NewRecorder(r.Recorder)

	if apiDefinition.Spec.Contexts != nil {
		delegate.ResolveContexts(apiDefinition)
	}

	apiDefinition.Status.Initialize()

	var reconcileErr error

	if apiDefinition.IsBeingDeleted() {
		reconcileErr = events.Record(event.Delete, apiDefinition, func() error {
			return delegate.Delete(apiDefinition)
		})
	} else {
		reconcileErr = events.Record(event.Update, apiDefinition, func() error {
			return delegate.CreateOrUpdate(apiDefinition)
		})
	}

	if reconcileErr == nil {
		log.Info("API definition has been reconciled")
		return ctrl.Result{}, delegate.UpdateStatusSuccess(apiDefinition)
	}

	if err := delegate.UpdateStatusFailure(apiDefinition, reconcileErr); err != nil {
		return ctrl.Result{}, err
	}

	if internal.IsRecoverable(reconcileErr) {
		log.Error(reconcileErr, "Requeuing reconcile")
		return ctrl.Result{RequeueAfter: requeueAfterTime}, reconcileErr
	}

	log.Error(reconcileErr, "Aborting reconcile")
	return ctrl.Result{}, nil
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	contextsWatcher := &handler.Funcs{UpdateFunc: r.CreateUpdaterFromLookup(indexer.ContextField.String())}
	resourceWatcher := &handler.Funcs{UpdateFunc: r.CreateUpdaterFromLookup(indexer.ResourceField.String())}

	return ctrl.NewControllerManagedBy(mgr).
		For(&gio.ApiDefinition{}).
		Watches(&source.Kind{Type: &gio.ApiContext{}}, contextsWatcher).
		Watches(&source.Kind{Type: &gio.ApiResource{}}, resourceWatcher).
		WithEventFilter(ApiUpdateFilter{}).
		Complete(r)
}
