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

package identityprovider

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/identityprovider"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
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
	*v1alpha1.AMIdentityProvider, identityprovider.IdentityProvider, *am.Client, am.DomainSubResourceResponse,
]

// Reconciler reconciles an AMIdentityProvider object.
type Reconciler struct {
	client.Client
	Lifecycle Lifecycle
	Scheme    *runtime.Scheme
	Recorder  record.EventRecorder
	Watcher   watch.Interface
}

func NewLifecycle() Lifecycle {
	return lifecycle.NewResourceLifecycle(Lifecycle{
		Finalizer:     core.AMIdentityProviderFinalizer,
		ResolveRefs:   identityprovider.ResolveDomain,
		Owner:         identityprovider.GetOwner,
		ClientFactory: identityprovider.CreateAMClient,
		ToDTO:         identityprovider.ToIdentityProviderDTO,
		Delete:        identityprovider.Delete,
		Upsert:        identityprovider.Upsert,
		PostUpsert:    identityprovider.UpdateStatus,
	})
}

// +kubebuilder:rbac:groups=gravitee.io,resources=amidentityproviders,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gravitee.io,resources=amidentityproviders/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gravitee.io,resources=amidentityproviders/finalizers,verbs=update
// +kubebuilder:rbac:groups=gravitee.io,resources=amsecuritydomains,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return r.Lifecycle.Reconcile(ctx, r.Client, r.Recorder, req, &v1alpha1.AMIdentityProvider{})
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	specChanged := builder.WithPredicates(predicate.LastSpecHashPredicate{})
	newController := ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.AMIdentityProvider{}, specChanged).
		// Not filtered on the spec hash: a domain becoming ready only changes its status.
		Watches(&v1alpha1.AMSecurityDomain{}, r.Watcher.WatchDomains(search.AMIdentityProviderDomainField))

	if env.Config.EnableTemplating {
		newController.
			Watches(&corev1.Secret{}, r.Watcher.WatchTemplatingSource(core.CRDAMIdentityProviderResource), specChanged).
			Watches(&corev1.ConfigMap{}, r.Watcher.WatchTemplatingSource(core.CRDAMIdentityProviderResource), specChanged)
	}
	return newController.Complete(r)
}
