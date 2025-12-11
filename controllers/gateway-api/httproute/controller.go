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

package httproute

import (
	"context"
	"errors"
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/gateway"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/gateway-api/httproute/internal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/event"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	kErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/retry"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwAPIv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

type Reconciler struct {
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//nolint:gocognit // acceptable complexity
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	route := &gwAPIv1.HTTPRoute{}
	if err := k8s.GetClient().Get(ctx, req.NamespacedName, route); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	events := event.NewRecorder(r.Recorder)

	dc := route.DeepCopy()

	err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		if err := k8s.GetClient().Get(ctx, req.NamespacedName, dc); client.IgnoreNotFound(err) != nil {
			return err
		}

		return k8s.CreateOrUpdate(ctx, route, func() error {
			util.AddFinalizer(route, core.HTTPRouteFinalizer)

			if !route.DeletionTimestamp.IsZero() {
				return events.Record(event.Delete, route, func() error {
					util.RemoveFinalizer(route, core.HTTPRouteFinalizer)
					return nil
				})
			}

			return events.Record(event.Update, route, func() error {
				internal.Init(dc)

				// Create a shared gateway cache to ensure both Resolve and Accept use the same gateway versions
				cache := make(internal.GatewayCache)

				if err := internal.ResolveWithCache(ctx, dc, cache); err != nil {
					return err
				}

				if err := internal.AcceptWithCache(ctx, dc, cache); err != nil {
					return err
				}

				// TODO: detect with merge conflict resolution
				if err := internal.DetectConflicts(ctx, dc); err != nil {
					return err
				}

				for i := range dc.Status.Parents {
					parent := &dc.Status.Parents[i]
					if k8s.IsConflicted(gateway.WrapRouteParentStatus(parent)) {
						return nil
					}
				}

				if err := internal.Program(ctx, dc); err != nil {
					return err
				}
				return nil
			})
		})
	})

	if err != nil && errors.Is(err, internal.ErrGatewayNotReady) {
		log.Debug(ctx, "Gateway status not ready, requeuing HTTPRoute", "route", req.NamespacedName)
		return ctrl.Result{RequeueAfter: 2 * time.Second}, nil
	}

	if err != nil {
		log.ErrorRequeuingReconcile(ctx, err, route)
		return ctrl.Result{}, err
	}

	if err := k8s.GetClient().Get(ctx, req.NamespacedName, route); err != nil {
		if kErrors.IsNotFound(err) {
			log.Debug(ctx, "Looks like the HTTPRoute was deleted during reconciliation, no need to update status")
			return ctrl.Result{}, nil
		}
		log.ErrorRequeuingReconcile(ctx, err, route)
		return ctrl.Result{}, err
	}

	dc.Status.DeepCopyInto(&route.Status)
	if err := k8s.UpdateStatus(ctx, route); client.IgnoreNotFound(err) != nil {
		log.ErrorRequeuingReconcile(ctx, err, route)
		return ctrl.Result{}, err
	}

	log.InfoEndReconcile(ctx, route)
	return ctrl.Result{}, nil
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gwAPIv1.HTTPRoute{}).
		Watches(&gwAPIv1beta1.ReferenceGrant{}, internal.WatchReferenceGrants()).
		Watches(&gwAPIv1.Gateway{}, internal.WatchGateways()).
		Complete(r)
}
