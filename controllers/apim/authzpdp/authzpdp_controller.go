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

package authzpdp

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/authzpdp/internal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/event"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/predicate"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/template"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/watch"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type Reconciler struct {
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Watcher  watch.Interface
	Client   client.Client
}

// +kubebuilder:rbac:groups=gravitee.io,resources=authzpdps,verbs=get;list;watch;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=gravitee.io,resources=authzpdps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gravitee.io,resources=authzpdps/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	pdp := &v1alpha1.AuthzPdp{}
	if err := r.Client.Get(ctx, req.NamespacedName, pdp); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	pdp.SetNamespace(req.Namespace)

	// The finalizer is saved before mutate: mutate writes template references to Secrets and ConfigMaps,
	// and only a deletion reconcile releases them. Without a saved finalizer, there is no deletion reconcile.
	if err := r.addFinalizer(ctx, pdp); err != nil {
		return ctrl.Result{}, err
	}

	k8s.ResetConditionsExceptAutomationAPI(pdp)
	dc := pdp.DeepCopy()

	// Not util.CreateOrUpdate: it reads pdp again from the cache, which may not hold the finalizer saved above yet.
	// pdp is saved only when mutate succeeds: on error, its metadata changes are dropped and only the status is written.
	before := pdp.DeepCopy()
	err := r.mutate(ctx, event.NewRecorder(r.Recorder), pdp, dc)
	if err == nil && !equality.Semantic.DeepEqual(before, pdp) {
		err = r.Client.Update(ctx, pdp)
	}

	if copyErr := dc.GetStatus().DeepCopyTo(pdp); copyErr != nil {
		return ctrl.Result{}, copyErr
	}

	if err == nil {
		log.InfoEndReconcile(ctx, pdp)
		return ctrl.Result{}, internal.UpdateStatusSuccess(ctx, pdp)
	}

	if err := internal.UpdateStatusFailure(ctx, pdp, err); err != nil {
		return ctrl.Result{}, err
	}

	if errors.IsRecoverable(err) {
		log.ErrorRequeuingReconcile(ctx, err, pdp)
		return ctrl.Result{}, err
	}

	log.ErrorAbortingReconcile(ctx, err, pdp)
	return ctrl.Result{}, nil
}

func (r *Reconciler) mutate(ctx context.Context, events *event.Recorder, pdp, dc *v1alpha1.AuthzPdp) error {
	k8s.AddAnnotation(pdp, core.LastSpecHashAnnotation, hash.Calculate(&pdp.Spec))

	if pdp.IsBeingDeleted() {
		// The only place template references are released: the finalizer keeps pdp until this branch has run.
		if err := template.ReleaseReferences(ctx, pdp); err != nil {
			return err
		}
		return events.Record(event.Delete, pdp, func() error {
			if err := internal.Delete(ctx, dc); err != nil {
				return err
			}
			util.RemoveFinalizer(pdp, core.AuthzPdpFinalizer)
			return nil
		})
	}

	// Compile saves a reference to pdp on each Secret and ConfigMap it reads, at once and whatever the apply returns.
	if err := template.Compile(ctx, dc, true); err != nil {
		return err
	}
	return events.Record(event.Update, pdp, func() error {
		return internal.CreateOrUpdate(ctx, dc)
	})
}

// addFinalizer saves the finalizer with its own Update, before the reconcile touches any other object.
// Added inside mutate, it would be lost whenever the apply fails: a PDP the platform never accepted would then
// be deleted at once, leaving its Secrets stuck on finalizers.gravitee.io/templating (GKO-3387).
func (r *Reconciler) addFinalizer(ctx context.Context, pdp *v1alpha1.AuthzPdp) error {
	if pdp.IsBeingDeleted() || !util.AddFinalizer(pdp, core.AuthzPdpFinalizer) {
		return nil
	}
	gvk := pdp.GetObjectKind().GroupVersionKind()
	if err := r.Client.Update(ctx, pdp); err != nil {
		return err
	}
	// Update clears the kind, which keys the references template.Compile records on Secrets and ConfigMaps.
	pdp.GetObjectKind().SetGroupVersionKind(gvk)
	return nil
}

// SetupWithManager filters events per watch rather than with a controller-wide event filter: the
// entity and policy watches (WatchAuthzPdpReleases) need a dependent's status-only update, which
// LastSpecHashPredicate would drop.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	lastSpecHash := builder.WithPredicates(predicate.LastSpecHashPredicate{})

	newController := ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.AuthzPdp{}, lastSpecHash).
		Watches(&v1alpha1.ManagementContext{}, r.Watcher.WatchContexts(search.AuthzPdpContextField), lastSpecHash).
		// A delete the platform refused (409) while an entity still needed the PDP as a target is
		// retried once an entity changes or goes.
		Watches(&v1alpha1.AuthzEntity{}, r.Watcher.WatchAuthzPdpReleases())

	if env.Config.EnableTemplating {
		newController.
			Watches(&corev1.Secret{}, r.Watcher.WatchTemplatingSource("authzpdps"), lastSpecHash).
			Watches(&corev1.ConfigMap{}, r.Watcher.WatchTemplatingSource("authzpdps"), lastSpecHash)
	}
	return newController.Complete(r)
}
