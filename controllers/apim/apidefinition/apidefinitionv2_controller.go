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

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/predicate"

	"github.com/go-logr/logr"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/indexer"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/watch"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
)

// Reconciler reconciles a ApiDefinition object.
type Reconciler struct {
	client.Client
	Log      logr.Logger
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
	apiDefinition := &v1alpha1.ApiDefinition{}

	if err := r.Get(ctx, req.NamespacedName, apiDefinition); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return Reconcile(ctx, apiDefinition, r.Recorder)
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	newController := ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.ApiDefinition{}).
		Watches(&v1alpha1.ManagementContext{}, r.Watcher.WatchContexts(indexer.ApiContextField)).
		WithEventFilter(predicate.LastSpecHashPredicate{}).
		Watches(&v1alpha1.ApiResource{}, r.Watcher.WatchResources(indexer.ApiResourceField))
	if env.Config.EnableTemplating {
		newController = newController.
			Watches(&corev1.Secret{}, r.Watcher.WatchTemplatingSource("apidefinitions")).
			Watches(&corev1.ConfigMap{}, r.Watcher.WatchTemplatingSource("apidefinitions"))
	}
	return newController.Complete(r)
}
