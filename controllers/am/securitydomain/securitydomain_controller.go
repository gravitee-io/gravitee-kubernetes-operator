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

package securitydomain

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	internal "github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/securitydomain"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/lifecycle"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/predicate"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	amsdk "github.com/gravitee-io-labs/gravitee-automation-tools/am-sdk/pkg/sdk/domain"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/watch"
)

type Lifecycle lifecycle.ResourceLifecycle[*v1alpha1.AMSecurityDomain, amsdk.Domain, *am.Client, internal.DomainResponse]

// Reconciler reconciles an AMSecurityDomain object.
type Reconciler struct {
	client.Client
	Lifecycle Lifecycle
	Scheme    *runtime.Scheme
	Recorder  record.EventRecorder
	Watcher   watch.Interface
}

func NewLifecycle() Lifecycle {
	return Lifecycle{
		Finalizer:     core.AMSecurityDomainFinalizer,
		ResolveRefs:   nil,
		ClientFactory: internal.CreateAMClient,
		ToDTO:         internal.ToDomainDTO,
		DeleteGuard:   nil,
		Delete:        internal.Delete,
		Upsert:        internal.Upsert,
		PostUpsert:    internal.UpdateStatus,
	}
}

// +kubebuilder:rbac:groups=gravitee.io,resources=amsecuritydomains,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gravitee.io,resources=amsecuritydomains/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gravitee.io,resources=amsecuritydomains/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	domain := &v1alpha1.AMSecurityDomain{}
	if err := r.Client.Get(ctx, req.NamespacedName, domain); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	return lifecycle.ResourceLifecycle[*v1alpha1.AMSecurityDomain, amsdk.Domain, *am.Client, internal.DomainResponse](r.Lifecycle).
		Reconcile(ctx, r.Client, r.Recorder, req, domain)
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	newController := ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.AMSecurityDomain{}).
		WithEventFilter(predicate.LastSpecHashPredicate{}).
		Watches(&v1alpha1.AMContext{}, r.Watcher.WatchContexts(search.AMSecurityContextField))

	if env.Config.EnableTemplating {
		newController.Watches(&corev1.Secret{}, r.Watcher.WatchTemplatingSource(core.CRDAMSecurityDomainResource)).
			Watches(&corev1.ConfigMap{}, r.Watcher.WatchTemplatingSource(core.CRDAMSecurityDomainResource))
	}
	return newController.Complete(r)
}
