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

package apiresource

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/template"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/predicate"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/apiresource/internal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/event"
)

// Reconciler reconciles a ApiResource object.
type Reconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=gravitee.io,resources=apiresources,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=gravitee.io,resources=apiresources/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=gravitee.io,resources=apiresources/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	apiResource := &v1alpha1.ApiResource{}
	if err := r.Get(ctx, req.NamespacedName, apiResource); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	events := event.NewRecorder(r.Recorder)

	status := apiResource.Status.DeepCopy()
	_, reconcileErr := util.CreateOrUpdate(ctx, r.Client, apiResource, func() error {
		util.AddFinalizer(apiResource, core.ApiResourceFinalizer)
		k8s.AddAnnotation(apiResource, core.LastSpecHashAnnotation, hash.Calculate(&apiResource.Spec))

		dc := apiResource.DeepCopy()

		if err := template.Compile(ctx, dc); err != nil {
			return err
		}

		var err error
		if apiResource.IsBeingDeleted() {
			err = events.Record(event.Delete, apiResource, func() error {
				return internal.Delete(ctx, dc)
			})
		} else {
			err = events.Record(event.Update, apiResource, func() error {
				// We don't do anything directly when there is an update on ApiResource
				return nil
			})
		}

		dc.Status.DeepCopyInto(status)
		apiResource.SetFinalizers(dc.GetFinalizers())

		return err
	})

	status.DeepCopyInto(&apiResource.Status)

	if reconcileErr == nil {
		logger.Info("API Resource has been reconciled")
		return ctrl.Result{}, nil
	}

	// There was an error reconciling the Management Context
	return ctrl.Result{}, reconcileErr
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.ApiResource{}).
		WithEventFilter(predicate.LastSpecHashPredicate{}).
		Complete(r)
}
