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

package ingress

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/watch"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"sigs.k8s.io/controller-runtime/pkg/source"

	e "github.com/gravitee-io/gravitee-kubernetes-operator/internal/event"

	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/ingress/internal"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/go-logr/logr"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	netV1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Reconciler watches and reconciles Ingress objects.
type Reconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Watcher  watch.Interface
}

// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch

// Reconcile perform reconciliation logic for Ingress resource that is managed
// by the operator.
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	ingress := &netV1.Ingress{}
	if err := r.Get(ctx, req.NamespacedName, ingress); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	d := internal.NewDelegate(ctx, r.Client, log)

	events := e.NewRecorder(r.Recorder)
	var reconcileErr error
	if !ingress.DeletionTimestamp.IsZero() {
		reconcileErr = events.Record(e.Delete, ingress, func() error {
			return d.Delete(ingress)
		})
	} else {
		reconcileErr = events.Record(e.Update, ingress, func() error {
			return d.CreateOrUpdate(ingress)
		})
	}

	if reconcileErr != nil {
		log.Error(reconcileErr, "An error occurs while reconciling the Ingress", "Ingress", ingress)
		return ctrl.Result{}, reconcileErr
	}

	log.Info("Sync ingress DONE")
	return ctrl.Result{}, nil
}

func (r *Reconciler) ingressClassEventFilter() predicate.Predicate {
	reconcilable := func(o runtime.Object) bool {
		switch t := o.(type) {
		case *netV1.Ingress:
			return t.GetAnnotations()[keys.IngressClassAnnotation] == keys.IngressClassAnnotationValue
		case *v1alpha1.ApiDefinition:
			return t.GetLabels()[keys.CrdApiDefinitionTemplate] == "true"
		default:
			return false
		}
	}

	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return reconcilable(e.Object)
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			if !reconcilable(e.ObjectNew) {
				return false
			}
			if e.ObjectOld == nil || e.ObjectNew == nil {
				return false
			}
			if len(e.ObjectOld.GetFinalizers()) != len(e.ObjectNew.GetFinalizers()) {
				return false
			}

			return e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return reconcilable(e.Object)
		},
	}
}

// SetupWithManager initializes ingress controller manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&netV1.Ingress{}).
		Owns(&v1alpha1.ApiDefinition{}).
		Watches(&source.Kind{Type: &v1alpha1.ApiDefinition{}}, r.Watcher.WatchApiTemplate()).
		WithEventFilter(r.ingressClassEventFilter()).
		Complete(r)
}
