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

package portallisting

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/portallisting/internal"
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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type Reconciler struct {
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Watcher  watch.Interface
	Client   client.Client
}

// +kubebuilder:rbac:groups=gravitee.io,resources=portallistings,verbs=get;list;watch;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=gravitee.io,resources=portallistings/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gravitee.io,resources=portallistings/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	listing := &v1alpha1.PortalListing{}
	if err := r.Client.Get(ctx, req.NamespacedName, listing); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	listing.SetNamespace(req.Namespace)

	events := event.NewRecorder(r.Recorder)

	k8s.ResetConditionsExceptAutomationAPI(listing)
	dc := listing.DeepCopy()

	_, err := util.CreateOrUpdate(ctx, k8s.GetClient(), listing, func() error {
		util.AddFinalizer(listing, core.PortalListingFinalizer)
		k8s.AddAnnotation(listing, core.LastSpecHashAnnotation, hash.Calculate(&listing.Spec))

		if listing.IsBeingDeleted() {
			if err := template.ReleaseReferences(ctx, listing); err != nil {
				return err
			}
		} else if err := template.Compile(ctx, dc, true); err != nil {
			return err
		}

		var err error
		if listing.IsBeingDeleted() {
			err = events.Record(event.Delete, listing, func() error {
				if err := internal.Delete(ctx, dc); err != nil {
					return err
				}
				util.RemoveFinalizer(listing, core.PortalListingFinalizer)
				return nil
			})
		} else {
			err = events.Record(event.Update, listing, func() error {
				return internal.CreateOrUpdate(ctx, dc)
			})
		}

		return err
	})

	if err := dc.GetStatus().DeepCopyTo(listing); err != nil {
		return ctrl.Result{}, err
	}

	if err == nil {
		log.InfoEndReconcile(ctx, listing)
		return ctrl.Result{}, internal.UpdateStatusSuccess(ctx, listing)
	}

	if err := internal.UpdateStatusFailure(ctx, listing, err); err != nil {
		return ctrl.Result{}, err
	}

	if errors.IsRecoverable(err) {
		log.ErrorRequeuingReconcile(ctx, err, listing)
		return ctrl.Result{}, err
	}

	log.ErrorAbortingReconcile(ctx, err, listing)
	return ctrl.Result{}, nil
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	newController := ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.PortalListing{}).
		WithEventFilter(predicate.LastSpecHashPredicate{}).
		Watches(&v1alpha1.Portal{}, r.Watcher.WatchPortals(search.PortalListingPortalField)).
		Watches(&v1alpha1.ApiV4Definition{}, r.Watcher.WatchApis(search.PortalListingApiField))

	if env.Config.EnableTemplating {
		newController.
			Watches(&corev1.Secret{}, r.Watcher.WatchTemplatingSource("portallistings")).
			Watches(&corev1.ConfigMap{}, r.Watcher.WatchTemplatingSource("portallistings"))
	}
	return newController.Complete(r)
}
