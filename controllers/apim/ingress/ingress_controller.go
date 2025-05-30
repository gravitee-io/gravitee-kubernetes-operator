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

	"sigs.k8s.io/controller-runtime/pkg/predicate"

	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/template"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/watch"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"

	e "github.com/gravitee-io/gravitee-kubernetes-operator/internal/event"

	"github.com/go-logr/logr"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/ingress/internal"
	p "github.com/gravitee-io/gravitee-kubernetes-operator/internal/predicate"
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
	ingress := &netV1.Ingress{}
	if err := r.Get(ctx, req.NamespacedName, ingress); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	events := e.NewRecorder(r.Recorder)

	dc := ingress.DeepCopy()

	_, err := util.CreateOrUpdate(ctx, r.Client, ingress, func() error {
		util.AddFinalizer(ingress, core.IngressFinalizer)
		k8s.AddAnnotation(ingress, core.LastSpecHashAnnotation, hash.Calculate(&ingress.Spec))

		err := template.Compile(ctx, dc, true)
		if err != nil {
			return err
		}

		if !ingress.DeletionTimestamp.IsZero() {
			err = events.Record(e.Delete, ingress, func() error {
				if err := internal.Delete(ctx, dc); err != nil {
					return err
				}
				util.RemoveFinalizer(ingress, core.IngressFinalizer)
				return nil
			})
		} else {
			err = events.Record(e.Update, ingress, func() error {
				return internal.CreateOrUpdate(ctx, dc)
			})
		}

		return err
	})

	dc.Status.DeepCopyInto(&ingress.Status)

	if err != nil {
		log.ErrorRequeuingReconcile(
			ctx,
			err,
			ingress,
		)
		return ctrl.Result{}, err
	}

	log.InfoEndReconcile(ctx, ingress)
	return ctrl.Result{}, nil
}

func (r *Reconciler) ingressClassEventFilter() predicate.Predicate {
	reconcilable := func(o runtime.Object) bool {
		switch t := o.(type) {
		case *netV1.Ingress:
			return k8s.IsGraviteeIngress(t)
		case *v1alpha1.ApiDefinition:
			return t.GetAnnotations()[core.IngressTemplateAnnotation] == env.TrueString
		case *corev1.Secret:
			return t.Type == "kubernetes.io/tls"
		default:
			return false
		}
	}

	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			if !reconcilable(e.Object) {
				return false
			}
			return p.LastSpecHashPredicate{}.Create(e)
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			if !reconcilable(e.ObjectNew) {
				return false
			}

			return p.LastSpecHashPredicate{}.Update(e)
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
		Watches(&v1alpha1.ApiDefinition{}, r.Watcher.WatchApiTemplate()).
		Watches(&corev1.Secret{}, r.Watcher.WatchTLSSecret()).
		Watches(&corev1.Secret{}, r.Watcher.WatchTemplatingSource("ingresses")).
		Watches(&corev1.ConfigMap{}, r.Watcher.WatchTemplatingSource("ingresses")).
		WithEventFilter(r.ingressClassEventFilter()).
		Complete(r)
}
