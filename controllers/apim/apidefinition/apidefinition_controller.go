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

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/indexer"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"

	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1beta1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/apidefinition/internal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/event"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/watch"
)

// Reconciler reconciles a ApiDefinition object.
type Reconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Watcher  watch.Interface
}

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gravitee.io,resources=apidefinitions,verbs=get;list;watch;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=gravitee.io,resources=apidefinitions/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gravitee.io,resources=apidefinitions/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log.InfoInitReconcile(ctx)

	apiDefinition := &v1beta1.ApiDefinition{}

	if err := r.Get(ctx, req.NamespacedName, apiDefinition); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	delegate := internal.NewDelegate(ctx, r.Client)
	if err := delegate.ResolveTemplate(apiDefinition); err != nil {
		return ctrl.Result{}, err
	}

	events := event.NewRecorder(r.Recorder)

	if apiDefinition.GetAnnotations()[keys.IngressTemplateAnnotation] == "true" {
		log.Info(ctx, "Definition is an ingress template")

		if err := delegate.SyncApiDefinitionTemplate(apiDefinition, req.Namespace); err != nil {
			log.Error(ctx, err, "Failed to sync API definition template")
			return ctrl.Result{}, err
		}

		log.Info(ctx, "Ingress template has been reconciled")
		return ctrl.Result{}, nil
	}

	delegate.AddDeletionFinalizer(apiDefinition)

	if apiDefinition.Spec.Context != nil {
		if err := delegate.ResolveContext(apiDefinition); err != nil {
			log.Error(ctx, err, "Unable to resolve context, requeuing reconcile")
			return ctrl.Result{}, err
		}
	}

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
		log.InfoEndReconcile(ctx)
		return ctrl.Result{}, delegate.UpdateStatusSuccess(apiDefinition)
	}

	if err := delegate.UpdateStatusFailure(apiDefinition); err != nil {
		return ctrl.Result{}, err
	}

	if apim.IsRecoverable(reconcileErr) {
		log.ErrorRequeuingReconcile(ctx, reconcileErr)
		return ctrl.Result{}, reconcileErr
	}

	log.ErrorAbortingReconcile(ctx, reconcileErr)
	return ctrl.Result{}, nil
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.ApiDefinition{}).
		Watches(&v1beta1.ManagementContext{}, r.Watcher.WatchContexts(indexer.ContextField)).
		Watches(&v1beta1.ApiResource{}, r.Watcher.WatchResources()).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)
}
