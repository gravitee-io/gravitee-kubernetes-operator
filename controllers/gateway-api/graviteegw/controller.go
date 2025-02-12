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

package graviteegw

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/event"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/predicate"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type Reconciler struct {
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	gw := &v1alpha1.GraviteeGateway{}

	if err := k8s.GetClient().Get(ctx, req.NamespacedName, gw); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	events := event.NewRecorder(r.Recorder)

	_, err := util.CreateOrUpdate(ctx, k8s.GetClient(), gw, func() error {
		util.AddFinalizer(gw, core.GraviteeGatewayFinalizer)
		k8s.AddAnnotation(gw, core.LastSpecHashAnnotation, hash.Calculate(&gw.Spec))

		if !gw.DeletionTimestamp.IsZero() {
			return events.Record(event.Delete, gw, func() error {
				util.RemoveFinalizer(gw, core.GraviteeGatewayFinalizer)
				return nil
			})
		}

		return events.Record(event.Update, gw, func() error {
			return nil
		})
	})

	if err != nil {
		log.ErrorAbortingReconcile(ctx, err, gw)
		return ctrl.Result{}, err
	}

	log.InfoEndReconcile(ctx, gw)
	return ctrl.Result{}, nil
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.GraviteeGateway{}).
		WithEventFilter(predicate.LastSpecHashPredicate{}).
		Complete(r)
}
