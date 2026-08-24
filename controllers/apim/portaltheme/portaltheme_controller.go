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

package portaltheme

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/portaltheme/internal"
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

// +kubebuilder:rbac:groups=gravitee.io,resources=portalthemes,verbs=get;list;watch;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=gravitee.io,resources=portalthemes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gravitee.io,resources=portalthemes/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	thm := &v1alpha1.PortalTheme{}
	if err := r.Client.Get(ctx, req.NamespacedName, thm); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	thm.SetNamespace(req.Namespace)

	events := event.NewRecorder(r.Recorder)

	k8s.ResetConditionsExceptAutomationAPI(thm)
	dc := thm.DeepCopy()

	_, err := util.CreateOrUpdate(ctx, k8s.GetClient(), thm, func() error {
		util.AddFinalizer(thm, core.PortalThemeFinalizer)
		k8s.AddAnnotation(thm, core.LastSpecHashAnnotation, hash.Calculate(&thm.Spec))

		if thm.IsBeingDeleted() {
			if err := template.ReleaseReferences(ctx, thm); err != nil {
				return err
			}
		} else if err := template.Compile(ctx, dc, true); err != nil {
			return err
		}

		var err error
		if thm.IsBeingDeleted() {
			err = events.Record(event.Delete, thm, func() error {
				if err := internal.Delete(ctx, dc); err != nil {
					return err
				}
				util.RemoveFinalizer(thm, core.PortalThemeFinalizer)
				return nil
			})
		} else {
			err = events.Record(event.Update, thm, func() error {
				return internal.CreateOrUpdate(ctx, dc)
			})
		}

		return err
	})

	if err := dc.GetStatus().DeepCopyTo(thm); err != nil {
		return ctrl.Result{}, err
	}

	if err == nil {
		log.InfoEndReconcile(ctx, thm)
		return ctrl.Result{}, internal.UpdateStatusSuccess(ctx, thm)
	}

	if err := internal.UpdateStatusFailure(ctx, thm, err); err != nil {
		return ctrl.Result{}, err
	}

	if errors.IsRecoverable(err) {
		log.ErrorRequeuingReconcile(ctx, err, thm)
		return ctrl.Result{}, err
	}

	log.ErrorAbortingReconcile(ctx, err, thm)
	return ctrl.Result{}, nil
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	newController := ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.PortalTheme{}).
		WithEventFilter(predicate.LastSpecHashPredicate{}).
		Watches(&v1alpha1.ManagementContext{}, r.Watcher.WatchContexts(search.PortalThemeContextField))

	if env.Config.EnableTemplating {
		newController.
			Watches(&corev1.Secret{}, r.Watcher.WatchTemplatingSource("portalthemes")).
			Watches(&corev1.ConfigMap{}, r.Watcher.WatchTemplatingSource("portalthemes"))
	}
	return newController.Complete(r)
}
