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

package authzentity

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/authzentity/internal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/event"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
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

// +kubebuilder:rbac:groups=gravitee.io,resources=authzentities,verbs=get;list;watch;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=gravitee.io,resources=authzentities/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gravitee.io,resources=authzentities/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	entity := &v1alpha1.AuthzEntity{}
	if err := r.Client.Get(ctx, req.NamespacedName, entity); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	entity.SetNamespace(req.Namespace)

	// The finalizer is saved before mutate: mutate writes template references to Secrets and ConfigMaps,
	// and only a deletion reconcile releases them. Without a saved finalizer, there is no deletion reconcile.
	if err := r.addFinalizer(ctx, entity); err != nil {
		return ctrl.Result{}, err
	}

	k8s.ResetConditionsExceptAutomationAPI(entity)
	dc := entity.DeepCopy()

	// Not util.CreateOrUpdate: it reads entity again from the cache, which may not hold the finalizer saved above yet.
	// entity is saved only when mutate succeeds: on error, its metadata changes are dropped and only the status is written.
	before := entity.DeepCopy()
	err := r.mutate(ctx, event.NewRecorder(r.Recorder), entity, dc)
	if err == nil && !equality.Semantic.DeepEqual(before, entity) {
		err = r.Client.Update(ctx, entity)
	}

	if copyErr := dc.GetStatus().DeepCopyTo(entity); copyErr != nil {
		return ctrl.Result{}, copyErr
	}

	if err == nil {
		log.InfoEndReconcile(ctx, entity)
		return ctrl.Result{}, internal.UpdateStatusSuccess(ctx, entity)
	}

	if err := internal.UpdateStatusFailure(ctx, entity, err); err != nil {
		return ctrl.Result{}, err
	}

	if errors.IsRecoverable(err) {
		log.ErrorRequeuingReconcile(ctx, err, entity)
		return ctrl.Result{}, err
	}

	log.ErrorAbortingReconcile(ctx, err, entity)
	return ctrl.Result{}, nil
}

func (r *Reconciler) mutate(ctx context.Context, events *event.Recorder, entity, dc *v1alpha1.AuthzEntity) error {
	k8s.AddAnnotation(entity, core.LastSpecHashAnnotation, hash.Calculate(&entity.Spec))

	if entity.IsBeingDeleted() {
		// The only place template references are released: the finalizer keeps entity until this branch has run.
		if err := template.ReleaseReferences(ctx, entity); err != nil {
			return err
		}
		return events.Record(event.Delete, entity, func() error {
			if err := internal.Delete(ctx, dc); err != nil {
				return err
			}
			util.RemoveFinalizer(entity, core.AuthzEntityFinalizer)
			return nil
		})
	}

	// Resolved before templates compile, so an entity waiting on an AuthzPdp records no template
	// reference. A pending AuthzPdp is waited out: the AuthzPdp watch requeues the entity once it
	// syncs, and ResolvedRefs names every pending one. Any other error (an AuthzPdp of another
	// context, a cache read failure) is retried with backoff and reported on Accepted.
	targets, err := dynamic.ResolveAuthzTargets(ctx, entity, entity.Spec.Targets)
	if err != nil {
		if dynamic.IsAuthzPdpPending(err) {
			return errors.NewResolveRefError(err)
		}
		return err
	}

	// Compile saves a reference to entity on each Secret and ConfigMap it reads, at once and whatever the apply returns.
	if err := template.Compile(ctx, dc, true); err != nil {
		return err
	}
	return events.Record(event.Update, entity, func() error {
		return internal.CreateOrUpdate(ctx, dc, targets)
	})
}

// addFinalizer saves the finalizer with its own Update, before the reconcile touches any other object.
// Added inside mutate, it would be lost whenever the apply fails: an entity the platform never accepted would then
// be deleted at once, leaving its Secrets stuck on finalizers.gravitee.io/templating (GKO-3387).
func (r *Reconciler) addFinalizer(ctx context.Context, entity *v1alpha1.AuthzEntity) error {
	if entity.IsBeingDeleted() || !util.AddFinalizer(entity, core.AuthzEntityFinalizer) {
		return nil
	}
	gvk := entity.GetObjectKind().GroupVersionKind()
	if err := r.Client.Update(ctx, entity); err != nil {
		return err
	}
	// Update clears the kind, which keys the references template.Compile records on Secrets and ConfigMaps.
	entity.GetObjectKind().SetGroupVersionKind(gvk)
	return nil
}

// SetupWithManager filters events per watch rather than with a controller-wide event filter: an
// AuthzPdp gaining its routing scope is a status-only update, which LastSpecHashPredicate would
// drop, leaving an entity waiting on that PDP to its retry backoff.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	lastSpecHash := builder.WithPredicates(predicate.LastSpecHashPredicate{})

	newController := ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.AuthzEntity{}, lastSpecHash).
		Watches(&v1alpha1.ManagementContext{}, r.Watcher.WatchContexts(search.AuthzEntityContextField), lastSpecHash).
		Watches(
			&v1alpha1.AuthzPdp{},
			r.Watcher.WatchAuthzPdps(search.AuthzEntityAuthzPdpField),
			builder.WithPredicates(predicate.AuthzPdpRoutedPredicate{}),
		)

	if env.Config.EnableTemplating {
		newController.
			Watches(&corev1.Secret{}, r.Watcher.WatchTemplatingSource("authzentities"), lastSpecHash).
			Watches(&corev1.ConfigMap{}, r.Watcher.WatchTemplatingSource("authzentities"), lastSpecHash)
	}
	return newController.Complete(r)
}
