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

	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/ingress/internal"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/go-logr/logr"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	netV1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// Reconciler watches and reconciles Ingress objects.
type Reconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
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

	var operation util.OperationResult
	var createError error
	if ingress.ObjectMeta.DeletionTimestamp.IsZero() {
		operation, createError = d.CreateOrUpdateApiDefintion(ingress, req.Namespace)
		if createError != nil {
			log.Error(createError, "An error occurs while creating or updating the ApiDefinition", "Operation", operation)
			return ctrl.Result{}, createError
		}
	}

	operation, createError = d.CreateOrUpdateIngress(ingress)
	if createError != nil {
		log.Error(createError, "An error occurs while updating the ingress", "Operation", operation)
		return ctrl.Result{}, createError
	}

	// TODO: templating
	// TODO: transform
	log.Info("Sync ingress DONE")

	return ctrl.Result{}, nil
}

func (r *Reconciler) ingressClassEventFilter() predicate.Predicate {
	isGraviteeIngress := func(o runtime.Object) bool {
		switch e := o.(type) {
		case *netV1.Ingress:
			return e.GetAnnotations()[keys.IngressClassAnnotation] == keys.IngressClassAnnotationValue
		default:
			return false
		}
	}

	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return isGraviteeIngress(e.Object)
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			return isGraviteeIngress(e.ObjectNew)
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return isGraviteeIngress(e.Object)
		},
	}
}

// SetupWithManager initializes ingress controller manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&netV1.Ingress{}).
		Owns(&v1alpha1.ApiDefinition{}).
		WithEventFilter(r.ingressClassEventFilter()).
		Complete(r)
}
