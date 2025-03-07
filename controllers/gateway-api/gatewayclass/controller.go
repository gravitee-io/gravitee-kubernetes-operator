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

package gatewayclass

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/gateway"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/gateway-api/gatewayclass/internal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/event"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type Reconciler struct {
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	gwc := gateway.NewGatewayClass(&gwAPIv1.GatewayClass{})

	if err := k8s.GetClient().Get(ctx, req.NamespacedName, gwc.Object); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	controllerName := gwc.Object.Spec.ControllerName
	if controllerName != core.GraviteeGatewayClassController {
		log.Debug(ctx, "unknown controller name", log.KeyValues(gwc.Object, "controllerName", controllerName)...)
		return ctrl.Result{}, nil
	}

	events := event.NewRecorder(r.Recorder)

	dc := gwc.DeepCopy()

	_, err := util.CreateOrUpdate(ctx, k8s.GetClient(), gwc.Object, func() error {
		util.AddFinalizer(gwc.Object, core.GatewayClassFinalizer)

		if !gwc.Object.DeletionTimestamp.IsZero() {
			return events.Record(event.Delete, gwc.Object, func() error {
				util.RemoveFinalizer(gwc.Object, core.GatewayClassFinalizer)
				return nil
			})
		}

		return events.Record(event.Update, gwc.Object, func() error {
			if accepted, err := internal.Accept(ctx, gwc.Object); err != nil {
				return err
			} else {
				k8s.SetCondition(dc, accepted)
				return nil
			}
		})
	})

	if err != nil {
		log.ErrorRequeuingReconcile(ctx, err, gwc.Object)
		return k8s.RequeueError(err)
	}

	dc.Object.Status.DeepCopyInto(&gwc.Object.Status)
	if err := k8s.GetClient().Status().Update(ctx, gwc.Object); err != nil {
		log.ErrorRequeuingReconcile(ctx, err, gwc.Object)
		return k8s.RequeueError(err)
	}

	log.InfoEndReconcile(ctx, gwc.Object)
	return ctrl.Result{}, err
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gwAPIv1.GatewayClass{}).
		Complete(r)
}
