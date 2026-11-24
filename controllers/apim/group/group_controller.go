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

package group

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimgroup "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/group"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/lifecycle"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/predicate"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/watch"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Lifecycle = lifecycle.ResourceLifecycle[*v1alpha1.Group, apimgroup.DTO, *apim.SDKClient, apimgroup.Response]

// Reconciler reconciles a Group object.
type Reconciler struct {
	client.Client
	Lifecycle Lifecycle
	Scheme    *runtime.Scheme
	Recorder  record.EventRecorder
	Watcher   watch.Interface
}

// NewLifecycle wires the Group reconcile steps. Group has no references to resolve, so its status reports
// Accepted only, without a ResolvedRefs condition.
func NewLifecycle() Lifecycle {
	return lifecycle.NewResourceLifecycle(Lifecycle{
		Finalizer:     core.GroupFinalizer,
		ResolveRefs:   nil,
		ClientFactory: apimgroup.CreateAPIMClient,
		ToDTO:         lifecycle.WithoutContext(apimgroup.ToDTO),
		Delete:        apimgroup.Delete,
		Upsert:        apimgroup.Upsert,
		PostUpsert:    apimgroup.UpdateStatus,
	})
}

// +kubebuilder:rbac:groups=gravitee.io,resources=groups,verbs=get;list;watch;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=gravitee.io,resources=groups/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gravitee.io,resources=groups/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return r.Lifecycle.Reconcile(ctx, r.Client, r.Recorder, req, &v1alpha1.Group{})
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Group{}).
		WithEventFilter(predicate.LastSpecHashPredicate{}).
		Watches(&v1alpha1.ManagementContext{}, r.Watcher.WatchContexts(search.GroupContextField)).
		Complete(r)
}
