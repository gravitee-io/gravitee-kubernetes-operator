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

package backendtlspolicy

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/gateway-api/backendtlspolicy/internal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type Reconciler struct {
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	policy := &gwAPIv1.BackendTLSPolicy{}
	if err := k8s.GetClient().Get(ctx, req.NamespacedName, policy); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	dc := policy.DeepCopy()

	ancestors, err := internal.Reconcile(ctx, dc)
	if err != nil {
		log.ErrorRequeuingReconcile(ctx, err, policy)
		return ctrl.Result{}, err
	}

	dc.Status.Ancestors = ancestors

	if err := k8s.GetClient().Status().Patch(ctx, dc, client.MergeFrom(policy)); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return ctrl.Result{}, nil
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gwAPIv1.BackendTLSPolicy{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Watches(&coreV1.ConfigMap{}, internal.WatchConfigMaps()).
		Complete(r)
}
