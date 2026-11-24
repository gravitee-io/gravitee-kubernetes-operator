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

	"github.com/go-logr/logr"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/apiv4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/lifecycle"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/predicate"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/watch"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// V4Lifecycle is the reconcile of an ApiV4Definition on internal/lifecycle.
type V4Lifecycle = lifecycle.ResourceLifecycle[*v1alpha1.ApiV4Definition, apiv4.DTO, *apiv4.Client, apiv4.Response]

// NewV4Lifecycle wires the ApiV4Definition reconcile steps. The finalizer is the one v2 shares, which existing
// APIs already carry.
func NewV4Lifecycle() V4Lifecycle {
	return lifecycle.NewResourceLifecycle(V4Lifecycle{
		Finalizer:     core.ApiDefinitionFinalizer,
		ResolveRefs:   apiv4.ResolveRefs,
		ClientFactory: apiv4.CreateAPIMClient,
		ToDTO:         lifecycle.WithoutContext(apiv4.ToDTO),
		DeleteGuard:   apiv4.DeleteGuard,
		Delete:        apiv4.Delete,
		Upsert:        apiv4.Upsert,
		PostUpsert:    apiv4.UpdateStatus,
		OnUpsertError: apiv4.KeepAcceptedStatus,
		Status:        apiv4.SetProcessingStatus,
	})
}

// V4Reconciler reconciles a ApiDefinition object.
type V4Reconciler struct {
	client.Client
	Lifecycle V4Lifecycle
	Log       logr.Logger
	Scheme    *runtime.Scheme
	Recorder  record.EventRecorder
	Watcher   watch.Interface
}

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gravitee.io,resources=graviteev4apis,verbs=get;list;watch;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=gravitee.io,resources=graviteev4apis/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gravitee.io,resources=graviteev4apis/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
func (r *V4Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	apiDefinition := &v1alpha1.ApiV4Definition{}

	if err := r.Get(ctx, req.NamespacedName, apiDefinition); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// An ingress template is never sent to APIM: it keeps its own reconcile.
	if apiDefinition.GetAnnotations()[core.IngressTemplateAnnotation] == "true" {
		k8s.ResetConditionsExceptAutomationAPI(apiDefinition)
		return reconcileApiTemplate(ctx, apiDefinition)
	}

	return r.Lifecycle.Reconcile(ctx, r.Client, r.Recorder, req, &v1alpha1.ApiV4Definition{})
}

// SetupWithManager filters events per watch rather than with a controller-wide event filter: a
// group becoming synced to APIM is an update that leaves its spec unchanged, which
// LastSpecHashPredicate would drop, leaving the API imported without that group.
func (r *V4Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	lastSpecHash := builder.WithPredicates(predicate.LastSpecHashPredicate{})

	newController := ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.ApiV4Definition{}, lastSpecHash).
		Watches(&v1alpha1.ManagementContext{}, r.Watcher.WatchContexts(search.ApiV4ContextField), lastSpecHash).
		Watches(&v1alpha1.ApiResource{}, r.Watcher.WatchResources(search.ApiV4ResourceField), lastSpecHash).
		Watches(
			&v1alpha1.SharedPolicyGroup{},
			r.Watcher.WatchSharedPolicyGroups(search.ApiV4SharedPolicyGroupsField),
			lastSpecHash,
		).
		Watches(&v1alpha1.Notification{}, r.Watcher.WatchNotifications(search.ApiV4NotificationRefsField), lastSpecHash).
		Watches(
			&v1alpha1.Group{},
			r.Watcher.WatchGroupsAndDeletes(search.ApiV4GroupField),
			builder.WithPredicates(predicate.GroupSyncedPredicate{}),
		)

	if env.Config.EnableTemplating {
		newController.
			Watches(&corev1.Secret{}, r.Watcher.WatchTemplatingSource("apiv4definitions"), lastSpecHash).
			Watches(&corev1.ConfigMap{}, r.Watcher.WatchTemplatingSource("apiv4definitions"), lastSpecHash)
	}

	return newController.Complete(r)
}
