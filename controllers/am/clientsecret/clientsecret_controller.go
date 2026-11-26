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

package clientsecret

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/clientsecret"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
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

type Lifecycle = lifecycle.ResourceLifecycle[
	*v1alpha1.AMClientSecret, clientsecret.ClientSecret, *am.Client, clientsecret.Response,
]

// Reconciler reconciles an AMClientSecret object.
type Reconciler struct {
	client.Client
	Lifecycle Lifecycle
	Scheme    *runtime.Scheme
	Recorder  record.EventRecorder
	Watcher   watch.Interface
}

func NewLifecycle() Lifecycle {
	return lifecycle.NewResourceLifecycle(Lifecycle{
		Finalizer:     core.AMClientSecretFinalizer,
		ResolveRefs:   clientsecret.ResolveRefs,
		Owner:         clientsecret.GetOwner,
		ClientFactory: clientsecret.CreateClient,
		ToDTO:         clientsecret.ToClientSecretDTO,
		Delete:        clientsecret.Delete,
		Upsert:        clientsecret.Upsert,
		PostUpsert:    clientsecret.UpdateStatus,
	})
}

// +kubebuilder:rbac:groups=gravitee.io,resources=amclientsecrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gravitee.io,resources=amclientsecrets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gravitee.io,resources=amclientsecrets/finalizers,verbs=update
// +kubebuilder:rbac:groups=gravitee.io,resources=amapplications;amsecuritydomains,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return r.Lifecycle.Reconcile(ctx, r.Client, r.Recorder, req, &v1alpha1.AMClientSecret{})
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	specChanged := builder.WithPredicates(predicate.LastSpecHashPredicate{})
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.AMClientSecret{}, specChanged).
		// Not filtered on the spec hash: an application becoming ready only changes its status.
		Watches(&v1alpha1.AMApplication{}, r.Watcher.WatchApplications(search.AMClientSecretApplicationField)).
		// Not filtered either: a value change in the Secret reaches AM and reports the immutability error at once.
		Watches(&corev1.Secret{}, r.Watcher.WatchResources(search.AMClientSecretSecretField)).
		Complete(r)
}
